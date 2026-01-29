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
	Short: "Detect potential issues in the project",
	Long: `Analyze the project using AI to detect potential issues.

This command:
1. Checks if META is synced with the latest code
2. Invokes the Issuer AI role to analyze the project
3. Creates Issue files in the Issues/ directory

Each detected issue is saved as a separate file for tracking.`,
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

	printInfo("Analyzing project for potential issues...")

	// Read META.md for context
	metaContent, err := meta.ReadMetaFile()
	if err != nil {
		printError("Failed to read META.md: " + err.Error())
		return err
	}

	// Build the prompt for Claude
	prompt := buildAuditPrompt(metaContent)

	// Invoke Claude Code
	printInfo("Invoking AI Issuer...")
	result, err := invokeClaudeCode(prompt)
	if err != nil {
		printError("Failed to invoke AI: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	if strings.TrimSpace(result) == "" {
		printSuccess("No issues detected!")
		return nil
	}

	printSuccess("Audit complete!")
	fmt.Println("\n" + result)

	return nil
}

func buildAuditPrompt(metaContent string) string {
	return fmt.Sprintf(`You are the Issuer role in the LadderMoon system.

Your task: Analyze this project and identify potential ISSUES (problems, bugs, risks).

Project META information:
---
%s
---

Instructions:
1. Review the project context from META
2. Identify concrete, actionable issues
3. For each issue, output in this format:

## Issue: [Brief Title]
**Severity**: High/Medium/Low
**Category**: Bug/Security/Performance/Architecture/Documentation
**Description**: [Detailed description of the issue]
**Recommendation**: [How to fix it]

If no issues are found, simply respond: "No issues detected."

Focus on:
- Code quality problems
- Potential bugs
- Security vulnerabilities
- Performance concerns
- Architectural issues
- Missing documentation

Be specific and actionable. Don't be overly critical - focus on genuine problems.`, metaContent)
}

func invokeClaudeCode(prompt string) (string, error) {
	cmd := exec.Command("claude", "-p", prompt)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("claude command failed: %w", err)
	}
	return string(output), nil
}
