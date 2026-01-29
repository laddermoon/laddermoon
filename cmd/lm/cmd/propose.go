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
	Short: "Propose improvements for the project",
	Long: `Analyze the project using AI to suggest improvements.

This command:
1. Checks if META is synced with the latest code
2. Invokes the Suggester AI role to analyze the project
3. Creates Suggestion files in the Suggestions/ directory

Each suggestion is saved as a separate file for tracking.`,
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

	printInfo("Analyzing project for improvement suggestions...")

	// Read META.md for context
	metaContent, err := meta.ReadMetaFile()
	if err != nil {
		printError("Failed to read META.md: " + err.Error())
		return err
	}

	// Build the prompt for Claude
	prompt := buildProposePrompt(metaContent)

	// Invoke Claude Code
	printInfo("Invoking AI Suggester...")
	result, err := invokeClaude(prompt)
	if err != nil {
		printError("Failed to invoke AI: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	if strings.TrimSpace(result) == "" {
		printSuccess("No suggestions at this time!")
		return nil
	}

	printSuccess("Analysis complete!")
	fmt.Println("\n" + result)

	return nil
}

func buildProposePrompt(metaContent string) string {
	return fmt.Sprintf(`You are the Suggester role in the LadderMoon system.

Your task: Analyze this project and suggest IMPROVEMENTS (enhancements, optimizations, new features).

Project META information:
---
%s
---

Instructions:
1. Review the project context from META
2. Identify valuable improvement opportunities
3. For each suggestion, output in this format:

## Suggestion: [Brief Title]
**Impact**: High/Medium/Low
**Effort**: High/Medium/Low
**Category**: Feature/Optimization/Refactoring/Testing/DevOps
**Description**: [Detailed description of the improvement]
**Benefit**: [What value this brings]
**Implementation Notes**: [How to implement it]

If no improvements are needed, simply respond: "Project looks great! No suggestions at this time."

Focus on:
- Feature enhancements
- Performance optimizations
- Code refactoring opportunities
- Testing improvements
- Developer experience
- Maintainability

Be constructive and prioritize high-impact, low-effort improvements.`, metaContent)
}

func invokeClaude(prompt string) (string, error) {
	cmd := exec.Command("claude", "-p", prompt)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("claude command failed: %w", err)
	}
	return string(output), nil
}
