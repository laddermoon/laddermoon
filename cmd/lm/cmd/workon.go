package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var workonCmd = &cobra.Command{
	Use:     "workon <task>",
	Aliases: []string{"solve"},
	Short:   "Work on a Task: code, review, and merge",
	Long: `Complete a Task through the full development cycle.

This command orchestrates:
1. Code: Write code to implement the task (laddermoon-code skill)
2. Review: Review the changes (laddermoon-review skill)
3. Apply: Merge the feature branch (laddermoon-apply skill)

Example:
  lm workon Tasks/task-from-issue-001.md
  lm workon "Add user authentication"
  lm solve Tasks/task-001.md`,
	Args: cobra.MinimumNArgs(1),
	RunE: runWorkon,
}

func init() {
	rootCmd.AddCommand(workonCmd)
}

func runWorkon(cmd *cobra.Command, args []string) error {
	// Check prerequisites
	if _, err := meta.GetGitRoot(); err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized. Run 'lm init' first.")
		return meta.ErrNotInitialized
	}

	if !meta.SkillsInstalled() {
		printError("LadderMoon skills are not installed.")
		printInfo("Run 'lm init' to reinstall.")
		return fmt.Errorf("skills not installed")
	}

	taskInput := strings.Join(args, " ")

	// Step 1: Code - implement the task
	printInfo("\n=== Step 1: Code ===")
	printInfo("Task: " + taskInput)
	printInfo("Starting implementation...")

	if err := invokeCodeSkill(taskInput); err != nil {
		printError("Code step failed: " + err.Error())
		return err
	}

	// Ask user if they want to continue to review
	fmt.Print("\nContinue to Review? [y/n]: ")
	var choice string
	fmt.Scanln(&choice)
	if strings.ToLower(choice) != "y" {
		printInfo("Stopped after Code step. Run 'lm workon' again to continue.")
		return nil
	}

	// Step 2: Review - review the changes
	printInfo("\n=== Step 2: Review ===")
	printInfo("Reviewing changes...")

	if err := invokeReviewSkill(); err != nil {
		printError("Review step failed: " + err.Error())
		return err
	}

	// Ask user if they want to continue to apply
	fmt.Print("\nReview passed. Apply changes (merge)? [y/n]: ")
	fmt.Scanln(&choice)
	if strings.ToLower(choice) != "y" {
		printInfo("Stopped after Review step. Changes are in the feature branch.")
		return nil
	}

	// Step 3: Apply - merge the feature branch
	printInfo("\n=== Step 3: Apply ===")
	printInfo("Merging changes...")

	if err := invokeApplySkill(); err != nil {
		printError("Apply step failed: " + err.Error())
		printInfo("You may need to resolve merge conflicts manually.")
		return err
	}

	printSuccess("\nTask completed successfully!")
	printInfo("Run 'lm sync' to update META with the changes.")
	return nil
}

func invokeCodeSkill(taskInput string) error {
	prompt := fmt.Sprintf("Use the laddermoon-code skill to implement this task: %s", taskInput)

	// Coding needs interaction
	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func invokeReviewSkill() error {
	prompt := "Use the laddermoon-review skill to review the changes in the current feature branch."

	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func invokeApplySkill() error {
	prompt := "Use the laddermoon-apply skill to merge the current feature branch into main. If there are conflicts, try to resolve them."

	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
