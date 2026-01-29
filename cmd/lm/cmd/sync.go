package cmd

import (
	"fmt"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize code changes to META",
	Long: `Synchronize the current repository state with the META system.

This command:
1. Detects changes since the last sync
2. Records the git diff and commit log
3. Updates the sync state

After sync, you can run 'lm audit' or 'lm propose' to analyze the changes.`,
	RunE: runSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
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

	printInfo("Synchronizing your intent...")

	// Get current commit
	currentCommit, err := meta.GetCurrentCommitID()
	if err != nil {
		printError("Failed to get current commit: " + err.Error())
		return err
	}

	// Get last synced commit
	lastSyncedCommit, _ := meta.GetSyncedCommitID()

	if lastSyncedCommit == currentCommit {
		printInfo("Already up to date. No changes since last sync.")
		return nil
	}

	// Get diff and log
	var diff, log string
	if lastSyncedCommit != "" {
		printInfo(fmt.Sprintf("Syncing changes from %s to %s...", shortCommit(lastSyncedCommit), shortCommit(currentCommit)))
		diff, _ = meta.GetGitDiff(lastSyncedCommit, currentCommit)
		log, _ = meta.GetGitLog(lastSyncedCommit, currentCommit)
	} else {
		printInfo("First sync - recording initial state...")
		diff, _ = meta.GetGitDiff("", currentCommit)
		log, _ = meta.GetGitLog("", currentCommit)
	}

	// Record sync info to META
	syncEntry := fmt.Sprintf("\n## Sync [%s]\n\n", shortCommit(currentCommit))
	if lastSyncedCommit != "" {
		syncEntry += fmt.Sprintf("Changes from %s to %s:\n\n", shortCommit(lastSyncedCommit), shortCommit(currentCommit))
	} else {
		syncEntry += "Initial sync:\n\n"
	}

	if log != "" {
		syncEntry += "### Commits\n```\n" + log + "```\n\n"
	}

	if diff != "" {
		syncEntry += "### Changed Files\n```\n" + diff + "```\n\n"
	}

	if err := meta.AppendToMetaFile(syncEntry); err != nil {
		printError("Failed to update META.md: " + err.Error())
		return err
	}

	// Update sync state
	if err := meta.SetSyncedCommitID(currentCommit); err != nil {
		printError("Failed to update sync state: " + err.Error())
		return err
	}

	printSuccess("Sync complete!")
	if log != "" {
		fmt.Println("\nRecent commits:")
		fmt.Println(log)
	}

	printInfo("Next: Run 'lm audit' to detect issues or 'lm propose' for suggestions.")

	return nil
}
