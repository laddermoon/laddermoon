package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
	Use:   "solve [filename]",
	Short: "Solve an Issue or implement a Suggestion",
	Long: `Invoke the laddermoon-solve skill to solve a specific Issue or implement a Suggestion.

This command invokes the Coder AI role to:
1. Read the specified Issue or Suggestion file
2. Understand the project context from META
3. Implement a high-quality solution

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
	if _, err := meta.GetGitRoot(); err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized. Run 'lm init' first.")
		return meta.ErrNotInitialized
	}

	// Check if skills are installed
	if !meta.SkillsInstalled() {
		printError("LadderMoon skills are not installed.")
		printInfo("Run 'lm init' to reinstall.")
		return fmt.Errorf("skills not installed")
	}

	printInfo(fmt.Sprintf("Solving: %s", filename))

	// Invoke Claude Code with the laddermoon-solve skill
	if err := invokeSolveSkill(filename); err != nil {
		printError("Failed to solve: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	return nil
}

func invokeSolveSkill(filename string) error {
	prompt := fmt.Sprintf("Use the laddermoon-solve skill to solve this Issue/Suggestion: %s", filename)

	cmd := exec.Command("claude", "-p", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
