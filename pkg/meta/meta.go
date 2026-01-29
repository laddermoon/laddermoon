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
