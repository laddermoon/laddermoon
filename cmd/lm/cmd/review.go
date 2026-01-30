package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review <issue-or-suggestion-file>",
	Short: "Review changes made to resolve an Issue or Suggestion",
	Long: `Review and verify that changes correctly resolve an Issue or implement a Suggestion.

This command invokes the laddermoon-review skill to:
- Read the Issue/Suggestion that was addressed
- Examine the changes made
- Verify the solution meets requirements
- Approve, Request Changes, or Reject

Example:
  lm review Issues/issue-001-bug-fix.md
  lm review Suggestions/suggest-002-performance.md`,
	Args: cobra.ExactArgs(1),
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}

func runReview(cmd *cobra.Command, args []string) error {
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

	filename := args[0]

	printInfo("Reviewing: " + filename)

	// Invoke Claude Code with the laddermoon-review skill
	if err := invokeReviewSkill(filename); err != nil {
		printError("Failed to run review: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	return nil
}

func invokeReviewSkill(filename string) error {
	prompt := fmt.Sprintf("Use the laddermoon-review skill to review changes for: %s", filename)

	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
