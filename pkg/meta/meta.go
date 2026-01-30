package meta

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	// BranchName is the name of the shadow branch for META storage
	BranchName = "laddermoon-meta"
	// MetaFileName is the main META file name
	MetaFileName = "META.md"
	// FeedIDFile stores the next feed ID
	FeedIDFile = ".next_feed_id"
	// UserFeedLog is the file for recording raw user input
	UserFeedLog = "UserFeed.log"
	// LockFile is used for serializing META write operations
	LockFile = ".lm.lock"
)

var (
	ErrNotGitRepo     = errors.New("not a git repository")
	ErrAlreadyInit    = errors.New("laddermoon already initialized (branch laddermoon-meta exists)")
	ErrNotInitialized = errors.New("laddermoon not initialized, run 'lm init' first")
)

// GetGitRoot returns the root directory of the git repository
func GetGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", ErrNotGitRepo
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentCommitID returns the current HEAD commit ID
func GetCurrentCommitID() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current commit: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// BranchExists checks if a branch exists
func BranchExists(branchName string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", fmt.Sprintf("refs/heads/%s", branchName))
	return cmd.Run() == nil
}

// IsInitialized checks if laddermoon is initialized in the current repo
func IsInitialized() bool {
	return BranchExists(BranchName)
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getBranchMetaDir returns the META directory path for a specific branch
func getBranchMetaDir(branch string) string {
	// Sanitize branch name for filesystem (replace / with _)
	safeBranch := strings.ReplaceAll(branch, "/", "_")
	return safeBranch
}

// BranchMetaDirExists checks if the META directory exists for current branch
func BranchMetaDirExists() (bool, error) {
	if !BranchExists(BranchName) {
		return false, nil
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return false, err
	}

	branchDir := getBranchMetaDir(branch)
	// Check if the branch directory exists in the shadow branch
	cmd := exec.Command("git", "ls-tree", "-d", BranchName, branchDir)
	output, err := cmd.Output()
	if err != nil {
		return false, nil
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// InitMetaStructure initializes the META structure on the shadow branch
// Uses git worktree to avoid disrupting the main working directory
func InitMetaStructure() error {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return err
	}

	// Get current branch name
	branch, err := GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	branchDir := getBranchMetaDir(branch)

	// Create temp directory IN project root (Claude Code may not have write access elsewhere)
	tmpDir := filepath.Join(gitRoot, fmt.Sprintf(".lm-tmp-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Check if shadow branch exists
	shadowExists := BranchExists(BranchName)

	if shadowExists {
		// Shadow branch exists, checkout and add branch directory
		cmd := exec.Command("git", "worktree", "add", tmpDir, BranchName)
		cmd.Dir = gitRoot
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}
		defer func() {
			rmCmd := exec.Command("git", "worktree", "remove", "--force", tmpDir)
			rmCmd.Dir = gitRoot
			rmCmd.Run()
		}()
	} else {
		// Create orphan branch using worktree
		cmd := exec.Command("git", "worktree", "add", "--detach", tmpDir)
		cmd.Dir = gitRoot
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}
		defer func() {
			rmCmd := exec.Command("git", "worktree", "remove", "--force", tmpDir)
			rmCmd.Dir = gitRoot
			rmCmd.Run()
		}()

		// In the worktree, create orphan branch
		cmd = exec.Command("git", "checkout", "--orphan", BranchName)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create orphan branch: %w", err)
		}

		// Remove all files from index
		cmd = exec.Command("git", "rm", "-rf", "--cached", ".")
		cmd.Dir = tmpDir
		cmd.Run() // Ignore error if nothing to remove

		// Clean the worktree directory (remove everything except .git)
		entries, _ := os.ReadDir(tmpDir)
		for _, entry := range entries {
			if entry.Name() != ".git" {
				os.RemoveAll(filepath.Join(tmpDir, entry.Name()))
			}
		}
	}

	// Create branch-specific directory
	branchPath := filepath.Join(tmpDir, branchDir)
	if err := os.MkdirAll(branchPath, 0755); err != nil {
		return fmt.Errorf("failed to create branch directory: %w", err)
	}

	// Create META.md (empty file) in branch directory
	metaPath := filepath.Join(branchPath, MetaFileName)
	if err := os.WriteFile(metaPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create META.md: %w", err)
	}

	// Create directories with .gitkeep in branch directory
	dirs := []string{"Questions", "Issues", "Suggestions"}
	for _, dir := range dirs {
		dirPath := filepath.Join(branchPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		gitkeepPath := filepath.Join(dirPath, ".gitkeep")
		if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create .gitkeep in %s: %w", dir, err)
		}
	}

	// Add all files
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Commit
	commitMsg := fmt.Sprintf("Initialize LadderMoon META for branch: %s", branch)
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// GetMetaBranchCommitID returns the latest commit ID of the META branch
func GetMetaBranchCommitID() (string, error) {
	if !IsInitialized() {
		return "", ErrNotInitialized
	}
	cmd := exec.Command("git", "rev-parse", BranchName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get META branch commit: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// ReadMetaFile reads the content of META.md from the shadow branch for current branch
func ReadMetaFile() (string, error) {
	if !IsInitialized() {
		return "", ErrNotInitialized
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return "", err
	}

	branchDir := getBranchMetaDir(branch)
	filePath := filepath.Join(branchDir, MetaFileName)

	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", BranchName, filePath))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read META.md: %w", err)
	}
	return string(output), nil
}

// GetMetaFileList returns a list of files in the META branch for current branch
func GetMetaFileList() ([]string, error) {
	if !IsInitialized() {
		return nil, ErrNotInitialized
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	branchDir := getBranchMetaDir(branch)

	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", BranchName, branchDir)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list META files: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Remove branch directory prefix from paths
	var result []string
	for _, line := range lines {
		if line != "" {
			// Remove the branch directory prefix
			relativePath := strings.TrimPrefix(line, branchDir+"/")
			result = append(result, relativePath)
		}
	}

	return result, nil
}

// ReadFile reads content of a file from the shadow branch for current branch
func ReadFile(filename string) (string, error) {
	if !IsInitialized() {
		return "", ErrNotInitialized
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return "", err
	}

	branchDir := getBranchMetaDir(branch)
	filePath := filepath.Join(branchDir, filename)

	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", BranchName, filePath))
	output, err := cmd.Output()
	if err != nil {
		return "", nil // File doesn't exist, return empty
	}
	return string(output), nil
}

// withWorktree executes a function with a temporary worktree checked out to the META branch
// Creates worktree in project directory for Claude Code compatibility
func withWorktree(fn func(tmpDir string) error) error {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return err
	}

	// Create temp directory IN project root (Claude Code may not have write access elsewhere)
	tmpDir := filepath.Join(gitRoot, fmt.Sprintf(".lm-tmp-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create worktree for META branch
	cmd := exec.Command("git", "worktree", "add", tmpDir, BranchName)
	cmd.Dir = gitRoot
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	defer func() {
		rmCmd := exec.Command("git", "worktree", "remove", "--force", tmpDir)
		rmCmd.Dir = gitRoot
		rmCmd.Run()
	}()

	return fn(tmpDir)
}

// AppendToMetaFile appends content to META.md on the shadow branch
func AppendToMetaFile(content string) error {
	return AppendToFile(MetaFileName, content)
}

// AppendToFile appends content to a file on the shadow branch for current branch
func AppendToFile(filename, content string) error {
	if !IsInitialized() {
		return ErrNotInitialized
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	branchDir := getBranchMetaDir(branch)

	return withWorktree(func(tmpDir string) error {
		// Use branch-specific path
		branchPath := filepath.Join(tmpDir, branchDir)
		filePath := filepath.Join(branchPath, filename)

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Append to file
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		if _, err := f.WriteString(content); err != nil {
			return fmt.Errorf("failed to write content: %w", err)
		}

		// Git add and commit (use full path from tmpDir root)
		relPath := filepath.Join(branchDir, filename)
		cmd := exec.Command("git", "add", relPath)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}

		cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update %s for branch %s", filename, branch))
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		return nil
	})
}

// WriteFile writes content to a file on the shadow branch (overwrites existing) for current branch
func WriteFile(filename, content string) error {
	if !IsInitialized() {
		return ErrNotInitialized
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		return err
	}

	branchDir := getBranchMetaDir(branch)

	return withWorktree(func(tmpDir string) error {
		// Use branch-specific path
		branchPath := filepath.Join(tmpDir, branchDir)
		filePath := filepath.Join(branchPath, filename)

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		// Git add and commit (use full path from tmpDir root)
		relPath := filepath.Join(branchDir, filename)
		cmd := exec.Command("git", "add", relPath)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}

		cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update %s for branch %s", filename, branch))
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		return nil
	})
}

// GetSyncedCommitID reads the last synced commit ID from META branch
func GetSyncedCommitID() (string, error) {
	content, err := ReadFile(".sync_state")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(content), nil
}

// SetSyncedCommitID saves the synced commit ID to META branch
func SetSyncedCommitID(commitID string) error {
	return WriteFile(".sync_state", commitID+"\n")
}

// GetGitDiff returns the diff between two commits
func GetGitDiff(fromCommit, toCommit string) (string, error) {
	var cmd *exec.Cmd
	if fromCommit == "" {
		// If no from commit, get the full diff of the to commit
		cmd = exec.Command("git", "show", "--stat", toCommit)
	} else {
		cmd = exec.Command("git", "diff", "--stat", fromCommit, toCommit)
	}
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff: %w", err)
	}
	return string(output), nil
}

// GetGitLog returns commit log between two commits
func GetGitLog(fromCommit, toCommit string) (string, error) {
	var cmd *exec.Cmd
	if fromCommit == "" {
		cmd = exec.Command("git", "log", "--oneline", "-20", toCommit)
	} else {
		cmd = exec.Command("git", "log", "--oneline", fmt.Sprintf("%s..%s", fromCommit, toCommit))
	}
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get log: %w", err)
	}
	return string(output), nil
}

// InstallSkills installs LadderMoon skills to the project's .claude/skills directory
func InstallSkills(skillsFS embed.FS, skillNames []string) error {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return err
	}

	skillsDir := filepath.Join(gitRoot, ".claude", "skills")

	for _, name := range skillNames {
		skillDir := filepath.Join(skillsDir, name)

		// Create skill directory
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			return fmt.Errorf("failed to create skill directory %s: %w", name, err)
		}

		// Read skill file from embedded FS
		content, err := skillsFS.ReadFile(filepath.Join(name, "SKILL.md"))
		if err != nil {
			return fmt.Errorf("failed to read embedded skill %s: %w", name, err)
		}

		// Write skill file
		skillFile := filepath.Join(skillDir, "SKILL.md")
		if err := os.WriteFile(skillFile, content, 0644); err != nil {
			return fmt.Errorf("failed to write skill file %s: %w", name, err)
		}
	}

	return nil
}

// SkillsInstalled checks if LadderMoon skills are installed
func SkillsInstalled() bool {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return false
	}

	// Check if at least the feed skill exists
	skillFile := filepath.Join(gitRoot, ".claude", "skills", "laddermoon-feed", "SKILL.md")
	_, err = os.Stat(skillFile)
	return err == nil
}

// MetaLock represents a lock file for serializing META operations
type MetaLock struct {
	file *os.File
	path string
}

// AcquireMetaLock acquires an exclusive lock for META operations
// The lock file is created in the project root directory
func AcquireMetaLock() (*MetaLock, error) {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return nil, err
	}

	lockPath := filepath.Join(gitRoot, LockFile)

	// Create or open lock file
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create lock file: %w", err)
	}

	// Try to acquire exclusive lock with timeout
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return &MetaLock{file: file, path: lockPath}, nil
		}

		select {
		case <-timeout:
			file.Close()
			return nil, fmt.Errorf("timeout waiting for META lock")
		case <-ticker.C:
			continue
		}
	}
}

// Release releases the lock
func (l *MetaLock) Release() error {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.path)
	}
	return nil
}

// GetNextFeedID reads the next feed ID from the META branch
func GetNextFeedID() (int, error) {
	content, err := ReadFile(FeedIDFile)
	if err != nil {
		return 1, nil // Default to 1 if file doesn't exist
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return 1, nil
	}
	id, err := strconv.Atoi(content)
	if err != nil {
		return 1, nil
	}
	return id, nil
}

// IncrementFeedID increments the feed ID and saves it
func IncrementFeedID(currentID int) error {
	nextID := currentID + 1
	return WriteFile(FeedIDFile, strconv.Itoa(nextID)+"\n")
}

// RecordUserFeed records the user input to UserFeed.log
func RecordUserFeed(feedID int, content string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("\n=== Feed #%d ===\nDate: %s\nContent:\n%s\n===\n", feedID, timestamp, content)
	return AppendToFile(UserFeedLog, entry)
}
