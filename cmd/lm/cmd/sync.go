package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize code changes to META",
	Long: `Synchronize the current repository state with the META system.

This command invokes the laddermoon-sync skill via Claude Code to:
1. Detect changes since the last sync
2. Analyze the changes and update META.md appropriately
3. Update the sync state

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

	// Check if skills are installed
	if !meta.SkillsInstalled() {
		printError("LadderMoon skills are not installed.")
		printInfo("Run 'lm init' to reinstall.")
		return fmt.Errorf("skills not installed")
	}

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

	printInfo("Synchronizing codebase changes with AI...")
	if lastSyncedCommit != "" {
		printInfo(fmt.Sprintf("Changes: %s → %s", shortCommit(lastSyncedCommit), shortCommit(currentCommit)))
	} else {
		printInfo("First sync - analyzing current state...")
	}

	// Invoke Claude Code with the laddermoon-sync skill
	if err := invokeSyncSkill(); err != nil {
		printError("Failed to sync: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	printInfo("Next: Run 'lm audit' to detect issues or 'lm propose' for suggestions.")

	return nil
}

func invokeSyncSkill() error {
	prompt := "Use the laddermoon-sync skill to synchronize the codebase changes to META."

	cmd := exec.Command("claude", "-p", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
