package meta

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// BranchName is the name of the shadow branch for META storage
	BranchName = "laddermoon-meta"
	// MetaFileName is the main META file name
	MetaFileName = "META.md"
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

// getCurrentBranch returns the current branch name
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// InitMetaStructure initializes the META structure on the shadow branch
// Uses git worktree to avoid disrupting the main working directory
func InitMetaStructure() error {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return err
	}

	// Create a temporary directory for the worktree
	tmpDir, err := os.MkdirTemp("", "laddermoon-init-")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create orphan branch using worktree
	cmd := exec.Command("git", "worktree", "add", "--detach", tmpDir)
	cmd.Dir = gitRoot
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	defer func() {
		// Clean up worktree
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

	// Create META.md (empty file)
	metaPath := filepath.Join(tmpDir, MetaFileName)
	if err := os.WriteFile(metaPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create META.md: %w", err)
	}

	// Create directories with .gitkeep
	dirs := []string{"Questions", "Issues", "Suggestions"}
	for _, dir := range dirs {
		dirPath := filepath.Join(tmpDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		gitkeepPath := filepath.Join(dirPath, ".gitkeep")
		if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to create .gitkeep in %s: %w", dir, err)
		}
	}

	// Add all files
	cmd = exec.Command("git", "add", "-A")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Commit
	cmd = exec.Command("git", "commit", "-m", "Initialize LadderMoon META structure")
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

// ReadMetaFile reads the content of META.md from the shadow branch
func ReadMetaFile() (string, error) {
	if !IsInitialized() {
		return "", ErrNotInitialized
	}
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", BranchName, MetaFileName))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read META.md: %w", err)
	}
	return string(output), nil
}

// GetMetaFileList returns a list of files in the META branch
func GetMetaFileList() ([]string, error) {
	if !IsInitialized() {
		return nil, ErrNotInitialized
	}
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", BranchName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list META files: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	return lines, nil
}

// ReadFile reads content of a file from the shadow branch
func ReadFile(filename string) (string, error) {
	if !IsInitialized() {
		return "", ErrNotInitialized
	}
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", BranchName, filename))
	output, err := cmd.Output()
	if err != nil {
		return "", nil // File doesn't exist, return empty
	}
	return string(output), nil
}

// withWorktree executes a function with a temporary worktree checked out to the META branch
func withWorktree(fn func(tmpDir string) error) error {
	gitRoot, err := GetGitRoot()
	if err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "laddermoon-work-")
	if err != nil {
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

// AppendToFile appends content to a file on the shadow branch
func AppendToFile(filename, content string) error {
	if !IsInitialized() {
		return ErrNotInitialized
	}

	return withWorktree(func(tmpDir string) error {
		filePath := filepath.Join(tmpDir, filename)

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

		// Git add and commit
		cmd := exec.Command("git", "add", filename)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}

		cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update %s", filename))
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		return nil
	})
}

// WriteFile writes content to a file on the shadow branch (overwrites existing)
func WriteFile(filename, content string) error {
	if !IsInitialized() {
		return ErrNotInitialized
	}

	return withWorktree(func(tmpDir string) error {
		filePath := filepath.Join(tmpDir, filename)

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Write file
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		// Git add and commit
		cmd := exec.Command("git", "add", filename)
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}

		cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Update %s", filename))
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
