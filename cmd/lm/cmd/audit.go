package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Detect issues and let user decide which become Tasks",
	Long: `Analyze the project to detect potential issues, then let user decide which to create Tasks for.

This command:
1. Invokes laddermoon-audit skill to find and file Issues
2. Shows each Issue to user for verification
3. Creates Tasks for approved Issues

Example:
  lm audit`,
	RunE: runAudit,
}

func init() {
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) error {
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

	// Step 1: Invoke audit skill to find issues
	printInfo("Step 1: Analyzing project for issues...")
	if err := invokeAuditSkill(); err != nil {
		printError("Failed to audit: " + err.Error())
		return err
	}

	// Step 2: Find open issues and let user verify
	issues := findOpenIssues()
	if len(issues) == 0 {
		printSuccess("No issues found.")
		return nil
	}

	printInfo(fmt.Sprintf("\nFound %d issue(s). Review each to decide if it should become a Task:\n", len(issues)))

	for _, issue := range issues {
		// Display issue content
		content, err := meta.ReadFile(issue)
		if err != nil || content == "" {
			continue
		}

		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("Issue: %s\n", issue)
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println(content)
		fmt.Println(strings.Repeat("-", 60))

		fmt.Println("Options:")
		fmt.Println("  [a] Approve - Create a Task for this issue")
		fmt.Println("  [r] Reject  - Not a valid issue")
		fmt.Println("  [s] Skip    - Decide later")
		fmt.Println("  [q] Quit    - Stop reviewing")
		fmt.Print("\nYour choice: ")

		var choice string
		fmt.Scanln(&choice)

		switch strings.ToLower(strings.TrimSpace(choice)) {
		case "a", "approve":
			taskFile := createTaskFromIssue(issue)
			printSuccess(fmt.Sprintf("Task created: %s", taskFile))
		case "r", "reject":
			printInfo("Issue rejected.")
		case "s", "skip":
			printInfo("Skipped.")
		case "q", "quit":
			printInfo("Stopped reviewing.")
			return nil
		}
		fmt.Println()
	}

	printSuccess("Audit complete!")
	printInfo("Run 'lm tasks' to see created tasks.")
	return nil
}

func invokeAuditSkill() error {
	prompt := "Use the laddermoon-audit skill to detect potential issues and create Issue files."

	cmd := exec.Command("claude", "-p", prompt, "--dangerously-skip-permissions")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func findOpenIssues() []string {
	files, err := meta.GetMetaFileList()
	if err != nil {
		return nil
	}

	var issues []string
	for _, f := range files {
		if !strings.HasPrefix(f, "Issues/") {
			continue
		}
		content, err := meta.ReadFile(f)
		if err != nil || content == "" {
			continue
		}
		if strings.Contains(content, "**Status**: Open") {
			issues = append(issues, f)
		}
	}
	return issues
}

func createTaskFromIssue(issueFile string) string {
	// Extract issue ID from filename
	base := strings.TrimPrefix(issueFile, "Issues/")
	taskID := strings.TrimSuffix(base, ".md")
	taskFile := fmt.Sprintf("Tasks/task-from-%s.md", taskID)

	// Note: Actual task creation via worktree would be done by a skill
	// For now just indicate the task file path
	printInfo("Note: Run 'lm workon " + taskFile + "' to start working on this task.")
	return taskFile
}
