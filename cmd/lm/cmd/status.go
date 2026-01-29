package cmd

import (
	"fmt"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show LadderMoon META status",
	Long: `Display the current status of the LadderMoon META system including:
- Initialization status
- Current main branch commit ID
- META branch commit ID
- Sync status
- Pending items count`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	_, err := meta.GetGitRoot()
	if err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	// Check if initialized
	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized.")
		printInfo("Run 'lm init' to initialize.")
		return meta.ErrNotInitialized
	}

	printInfo("Checking status...")
	fmt.Println()

	// Get current commit ID
	currentCommit, err := meta.GetCurrentCommitID()
	if err != nil {
		printError("Failed to get current commit: " + err.Error())
		return err
	}

	// Get META branch commit ID
	metaCommit, err := meta.GetMetaBranchCommitID()
	if err != nil {
		printError("Failed to get META branch commit: " + err.Error())
		return err
	}

	// Get META file list
	files, err := meta.GetMetaFileList()
	if err != nil {
		printError("Failed to list META files: " + err.Error())
		return err
	}

	// Read META.md content
	metaContent, err := meta.ReadMetaFile()
	if err != nil {
		printError("Failed to read META.md: " + err.Error())
		return err
	}

	// Display status
	fmt.Println("╭─────────────────────────────────────────╮")
	fmt.Println("│         LadderMoon Status               │")
	fmt.Println("╰─────────────────────────────────────────╯")
	fmt.Println()

	fmt.Printf("  %-20s %s\n", "Initialized:", "✓ Yes")
	fmt.Printf("  %-20s %s\n", "META Branch:", meta.BranchName)
	fmt.Printf("  %-20s %s\n", "Main Commit:", shortCommit(currentCommit))
	fmt.Printf("  %-20s %s\n", "META Commit:", shortCommit(metaCommit))
	fmt.Println()

	fmt.Println("  META Files:")
	for _, f := range files {
		fmt.Printf("    - %s\n", f)
	}
	fmt.Println()

	if len(metaContent) == 0 {
		fmt.Println("  META.md: (empty)")
		fmt.Println("  Hint: Run 'lm feed <info>' to add project information")
	} else {
		fmt.Printf("  META.md: %d bytes\n", len(metaContent))
	}

	fmt.Println()

	return nil
}

func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}
