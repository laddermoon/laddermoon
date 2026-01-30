package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var answerCmd = &cobra.Command{
	Use:   "answer <question-file> [answer text]",
	Short: "Answer a Question and update META",
	Long: `Process user's answer to a filed Question and update META accordingly.

This command invokes the laddermoon-answer skill to:
- Read the Question file
- Process your answer
- Update META.md to reflect the resolution
- Mark the Question as resolved

Example:
  lm answer Questions/question-001-priority.md "Minimal dependencies is more important"
  lm answer Questions/question-002-env.md "Target is Linux servers only"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAnswer,
}

func init() {
	rootCmd.AddCommand(answerCmd)
}

func runAnswer(cmd *cobra.Command, args []string) error {
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

	questionFile := args[0]
	answer := ""
	if len(args) > 1 {
		answer = strings.Join(args[1:], " ")
	}

	if answer == "" {
		printError("Please provide an answer.")
		printInfo("Usage: lm answer <question-file> \"your answer\"")
		return fmt.Errorf("no answer provided")
	}

	printInfo("Processing answer for: " + questionFile)
	printInfo("Answer: " + truncateString(answer, 60))

	// Invoke Claude Code with the laddermoon-answer skill
	if err := invokeAnswerSkill(questionFile, answer); err != nil {
		printError("Failed to process answer: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	return nil
}

func invokeAnswerSkill(questionFile, answer string) error {
	prompt := fmt.Sprintf("Use the laddermoon-answer skill to process this answer.\n\nQuestion file: %s\n\nUser's answer: %s", questionFile, answer)

	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
