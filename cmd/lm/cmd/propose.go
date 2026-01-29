package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var proposeCmd = &cobra.Command{
	Use:   "propose [focus area]",
	Short: "Propose improvements for the project",
	Long: `Analyze the project using the laddermoon-propose skill to suggest improvements.

This command invokes the Suggester AI role to:
1. Analyze the project through the lens of META
2. Identify valuable improvement opportunities
3. Create Suggestion files in the Suggestions/ directory

Example:
  lm propose              # General suggestions
  lm propose DX           # Focus on developer experience
  lm propose performance  # Focus on performance optimizations`,
	RunE: runPropose,
}

func init() {
	rootCmd.AddCommand(proposeCmd)
}

func runPropose(cmd *cobra.Command, args []string) error {
	// Check prerequisites
	if _, err := meta.GetGitRoot(); err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized. Run 'lm init' first.")
		return meta.ErrNotInitialized
	}

	// Check if skills are installed
	if !meta.SkillsInstalled() {
		printError("LadderMoon skills are not installed.")
		printInfo("Run 'lm init' to reinstall.")
		return fmt.Errorf("skills not installed")
	}

	// Check sync status
	currentCommit, _ := meta.GetCurrentCommitID()
	syncedCommit, _ := meta.GetSyncedCommitID()

	if syncedCommit == "" {
		printError("No sync found. Run 'lm sync' first to sync your codebase.")
		return fmt.Errorf("not synced")
	}

	if currentCommit != syncedCommit {
		printError("META is out of sync with the codebase.")
		printInfo(fmt.Sprintf("Current: %s, Synced: %s", shortCommit(currentCommit), shortCommit(syncedCommit)))
		printInfo("Run 'lm sync' first to update.")
		return fmt.Errorf("out of sync")
	}

	// Determine focus area
	focusArea := "general"
	if len(args) > 0 {
		focusArea = args[0]
	}

	printInfo("Analyzing project for improvement suggestions...")
	if focusArea != "general" {
		printInfo("Focus area: " + focusArea)
	}

	// Invoke Claude Code with the laddermoon-propose skill
	if err := invokeProposeSkill(focusArea); err != nil {
		printError("Failed to propose: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	return nil
}

func invokeProposeSkill(focusArea string) error {
	var prompt string
	if focusArea == "general" {
		prompt = "Use the laddermoon-propose skill to suggest general improvements for the project."
	} else {
		prompt = fmt.Sprintf("Use the laddermoon-propose skill to suggest improvements with focus on: %s", focusArea)
	}

	cmd := exec.Command("claude", "-p", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
