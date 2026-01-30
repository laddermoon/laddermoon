package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var clarifyCmd = &cobra.Command{
	Use:   "clarify",
	Short: "Analyze META clarity and resolve questions iteratively",
	Long: `Iteratively analyze META for clarity issues and resolve them.

This command loops through:
1. Criticize: Analyze META for unclear/incomplete areas, file Questions
2. User review: Show questions, let user decide which to address
3. Clarify: Resolve approved questions by analyzing code or asking user
4. Repeat until no more questions

Example:
  lm clarify`,
	RunE: runClarify,
}

func init() {
	rootCmd.AddCommand(clarifyCmd)
}

func runClarify(cmd *cobra.Command, args []string) error {
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

	iteration := 1
	for {
		printInfo(fmt.Sprintf("\n=== Clarify Iteration %d ===", iteration))

		// Step 1: Criticize META to find issues
		printInfo("Step 1: Analyzing META for clarity issues...")
		if err := invokeCriticizeSkill(); err != nil {
			printError("Criticize failed: " + err.Error())
			return err
		}

		// Step 2: Check if there are open questions
		questions := findOpenQuestions()
		if len(questions) == 0 {
			printSuccess("META is clear! No more questions to resolve.")
			break
		}

		printInfo(fmt.Sprintf("Found %d open question(s):", len(questions)))
		for i, q := range questions {
			fmt.Printf("  %d. %s\n", i+1, q)
		}

		// Step 3: Let user choose which to address
		fmt.Print("\nEnter question number to clarify (or 'q' to quit, 'a' for all): ")
		var choice string
		fmt.Scanln(&choice)

		if strings.ToLower(choice) == "q" {
			printInfo("Exiting clarify loop.")
			break
		}

		if strings.ToLower(choice) == "a" {
			// Clarify all questions
			for _, q := range questions {
				printInfo("Clarifying: " + q)
				if err := invokeClarifySkillForQuestion(q); err != nil {
					printError("Failed to clarify: " + err.Error())
				}
			}
		} else {
			// Clarify selected question
			var idx int
			fmt.Sscanf(choice, "%d", &idx)
			if idx >= 1 && idx <= len(questions) {
				printInfo("Clarifying: " + questions[idx-1])
				if err := invokeClarifySkillForQuestion(questions[idx-1]); err != nil {
					printError("Failed to clarify: " + err.Error())
				}
			} else {
				printError("Invalid choice")
			}
		}

		iteration++
		if iteration > 10 {
			printInfo("Max iterations reached. Exiting.")
			break
		}
	}

	return nil
}

func invokeCriticizeSkill() error {
	prompt := "Use the laddermoon-criticize skill to analyze META for clarity and completeness, then file Questions for areas that need clarification."

	cmd := exec.Command("claude", "-p", prompt, "--dangerously-skip-permissions")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func findOpenQuestions() []string {
	files, err := meta.GetMetaFileList()
	if err != nil {
		return nil
	}

	var questions []string
	for _, f := range files {
		if !strings.HasPrefix(f, "Questions/") {
			continue
		}
		content, err := meta.ReadFile(f)
		if err != nil || content == "" {
			continue
		}
		if strings.Contains(content, "**Status**: Open") {
			questions = append(questions, f)
		}
	}
	return questions
}

func invokeClarifySkillForQuestion(questionFile string) error {
	prompt := fmt.Sprintf("Use the laddermoon-clarify skill to resolve this question: %s\n\nAnalyze the codebase first. Only ask me if you cannot find the answer in the code.", questionFile)

	// This skill may need user interaction
	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
