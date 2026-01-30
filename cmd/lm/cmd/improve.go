package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var improveCmd = &cobra.Command{
	Use:   "improve",
	Short: "Analyze and improve the META system quality",
	Long: `Analyze the quality of the META system and suggest improvements.

This command invokes the laddermoon-improve skill to:
- Analyze META.md structure quality
- Check traceability of information
- Examine workflow health (Issues/Suggestions/Questions)
- Check information freshness
- Generate improvement recommendations

Example:
  lm improve`,
	RunE: runImprove,
}

func init() {
	rootCmd.AddCommand(improveCmd)
}

func runImprove(cmd *cobra.Command, args []string) error {
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

	printInfo("Analyzing META system health...")

	// Invoke Claude Code with the laddermoon-improve skill
	if err := invokeImproveSkill(); err != nil {
		printError("Failed to run improve: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	return nil
}

func invokeImproveSkill() error {
	prompt := "Use the laddermoon-improve skill to analyze META system health and suggest improvements."

	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
