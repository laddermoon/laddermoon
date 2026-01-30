package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var proposeCmd = &cobra.Command{
	Use:   "propose",
	Short: "Propose improvements and let user decide which become Tasks",
	Long: `Analyze the project to propose improvements, then let user decide which to create Tasks for.

This command:
1. Invokes laddermoon-propose skill to create Proposals
2. Shows each Proposal to user for verification
3. Creates Tasks for approved Proposals

Example:
  lm propose`,
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

	if !meta.SkillsInstalled() {
		printError("LadderMoon skills are not installed.")
		printInfo("Run 'lm init' to reinstall.")
		return fmt.Errorf("skills not installed")
	}

	// Step 1: Invoke propose skill to find suggestions
	printInfo("Step 1: Analyzing project for improvement suggestions...")
	if err := invokeProposeSkill(); err != nil {
		printError("Failed to propose: " + err.Error())
		return err
	}

	// Step 2: Find open proposals and let user verify
	proposals := findOpenProposals()
	if len(proposals) == 0 {
		printSuccess("No proposals found.")
		return nil
	}

	printInfo(fmt.Sprintf("\nFound %d proposal(s). Review each to decide if it should become a Task:\n", len(proposals)))

	for _, proposal := range proposals {
		// Display proposal content
		content, err := meta.ReadFile(proposal)
		if err != nil || content == "" {
			continue
		}

		fmt.Println(strings.Repeat("=", 60))
		fmt.Printf("Proposal: %s\n", proposal)
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println(content)
		fmt.Println(strings.Repeat("-", 60))

		fmt.Println("Options:")
		fmt.Println("  [a] Approve - Create a Task for this proposal")
		fmt.Println("  [r] Reject  - Not a good proposal")
		fmt.Println("  [s] Skip    - Decide later")
		fmt.Println("  [q] Quit    - Stop reviewing")
		fmt.Print("\nYour choice: ")

		var choice string
		fmt.Scanln(&choice)

		switch strings.ToLower(strings.TrimSpace(choice)) {
		case "a", "approve":
			taskFile := createTaskFromProposal(proposal)
			printSuccess(fmt.Sprintf("Task created: %s", taskFile))
		case "r", "reject":
			printInfo("Proposal rejected.")
		case "s", "skip":
			printInfo("Skipped.")
		case "q", "quit":
			printInfo("Stopped reviewing.")
			return nil
		}
		fmt.Println()
	}

	printSuccess("Propose complete!")
	printInfo("Run 'lm tasks' to see created tasks.")
	return nil
}

func invokeProposeSkill() error {
	prompt := "Use the laddermoon-propose skill to propose improvements and create Proposal files."

	cmd := exec.Command("claude", "-p", prompt, "--dangerously-skip-permissions")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func findOpenProposals() []string {
	files, err := meta.GetMetaFileList()
	if err != nil {
		return nil
	}

	var proposals []string
	for _, f := range files {
		if !strings.HasPrefix(f, "Proposals/") && !strings.HasPrefix(f, "Suggestions/") {
			continue
		}
		content, err := meta.ReadFile(f)
		if err != nil || content == "" {
			continue
		}
		if strings.Contains(content, "**Status**: Open") {
			proposals = append(proposals, f)
		}
	}
	return proposals
}

func createTaskFromProposal(proposalFile string) string {
	// Extract proposal ID from filename
	base := proposalFile
	if strings.HasPrefix(base, "Proposals/") {
		base = strings.TrimPrefix(base, "Proposals/")
	} else if strings.HasPrefix(base, "Suggestions/") {
		base = strings.TrimPrefix(base, "Suggestions/")
	}
	taskID := strings.TrimSuffix(base, ".md")
	taskFile := fmt.Sprintf("Tasks/task-from-%s.md", taskID)

	printInfo("Note: Run 'lm workon " + taskFile + "' to start working on this task.")
	return taskFile
}
