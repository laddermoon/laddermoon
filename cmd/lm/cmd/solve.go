package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve [filename]",
	Short: "Solve an Issue or implement a Suggestion",
	Long: `Invoke the Coder AI role to solve a specific Issue or implement a Suggestion.

This command:
1. Reads the specified Issue or Suggestion file
2. Invokes the Coder AI role to implement the solution
3. Creates a new branch for the changes

Examples:
  lm solve Issues/issue-001.md
  lm solve Suggestions/suggest-refactor-auth.md`,
	Args: cobra.ExactArgs(1),
	RunE: runSolve,
}

func init() {
	rootCmd.AddCommand(solveCmd)
}

func runSolve(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Check prerequisites
	gitRoot, err := meta.GetGitRoot()
	if err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized. Run 'lm init' first.")
		return meta.ErrNotInitialized
	}

	// Read the issue/suggestion file from META branch
	content, err := meta.ReadFile(filename)
	if err != nil {
		printError("Failed to read file: " + err.Error())
		return err
	}

	if content == "" {
		// Try reading from local filesystem as fallback
		localPath := filepath.Join(gitRoot, filename)
		data, err := os.ReadFile(localPath)
		if err != nil {
			printError(fmt.Sprintf("File not found: %s", filename))
			printInfo("Provide a path to an Issue or Suggestion file.")
			return err
		}
		content = string(data)
	}

	printInfo(fmt.Sprintf("Solving: %s", filename))

	// Read META.md for context
	metaContent, _ := meta.ReadMetaFile()

	// Build the prompt
	prompt := buildSolvePrompt(filename, content, metaContent)

	// Invoke Claude Code
	printInfo("Invoking AI Coder...")
	result, err := invokeClaudeCoder(prompt)
	if err != nil {
		printError("Failed to invoke AI: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	printSuccess("Coder analysis complete!")
	fmt.Println("\n" + result)

	return nil
}

func buildSolvePrompt(filename, content, metaContent string) string {
	taskType := "Issue"
	if strings.Contains(strings.ToLower(filename), "suggest") {
		taskType = "Suggestion"
	}

	return fmt.Sprintf(`You are the Coder role in the LadderMoon system.

Your task: Implement a solution for this %s.

%s Content:
---
%s
---

Project META context:
---
%s
---

Instructions:
1. Analyze the %s carefully
2. Understand the project context from META
3. Propose a concrete implementation plan
4. If code changes are needed, provide the specific changes

Output format:

## Implementation Plan

### Understanding
[Your understanding of the problem/request]

### Approach
[Step-by-step approach to solve this]

### Code Changes
[Specific code changes needed, with file paths and code snippets]

### Testing
[How to verify the solution works]

### Risks
[Any potential risks or side effects]

Be precise and actionable. Provide actual code when possible.`, taskType, taskType, content, metaContent, taskType)
}

func invokeClaudeCoder(prompt string) (string, error) {
	cmd := exec.Command("claude", "-p", prompt)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("claude command failed: %w", err)
	}
	return string(output), nil
}
