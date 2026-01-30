package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var questionCmd = &cobra.Command{
	Use:   "question [focus area]",
	Short: "Proactively identify and file Questions",
	Long: `Analyze the project and identify things that need user clarification.

This command invokes the laddermoon-question skill to find gaps in 
understanding and file Questions that need user input.

Focus areas:
- General questioning (no argument)
- Specific: goals, architecture, decisions, etc.

Example:
  lm question
  lm question goals
  lm question architecture`,
	RunE: runQuestion,
}

func init() {
	rootCmd.AddCommand(questionCmd)
}

func runQuestion(cmd *cobra.Command, args []string) error {
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

	focusArea := "general"
	if len(args) > 0 {
		focusArea = args[0]
	}

	printInfo("Analyzing project for questions...")
	printInfo("Focus: " + focusArea)

	// Invoke Claude Code with the laddermoon-question skill
	if err := invokeQuestionSkill(focusArea); err != nil {
		printError("Failed to run question skill: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	return nil
}

func invokeQuestionSkill(focusArea string) error {
	var prompt string
	if focusArea == "general" {
		prompt = "Use the laddermoon-question skill to identify and file Questions that need user clarification."
	} else {
		prompt = fmt.Sprintf("Use the laddermoon-question skill to identify Questions with focus on: %s", focusArea)
	}

	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
