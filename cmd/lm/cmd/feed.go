package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var feedCmd = &cobra.Command{
	Use:   "feed [user input text]",
	Short: "Add project information to META",
	Long: `Process and integrate user-provided information into the META system.
	
This command invokes the laddermoon-feed skill via Claude Code to
intelligently integrate your input into META.md.

Example:
  lm feed "This project uses PostgreSQL for data storage"
  lm feed "The API follows RESTful conventions"`,
	Args: cobra.MinimumNArgs(1),
	RunE: runFeed,
}

func init() {
	rootCmd.AddCommand(feedCmd)
}

func runFeed(cmd *cobra.Command, args []string) error {
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

	// Join all args as the feed content
	content := strings.Join(args, " ")
	if strings.TrimSpace(content) == "" {
		printError("Feed content cannot be empty.")
		return fmt.Errorf("empty feed content")
	}

	printInfo("Processing your input with AI...")
	printInfo("Content: " + truncateString(content, 60))

	// Invoke Claude Code with the laddermoon-feed skill
	if err := invokeSkill("laddermoon-feed", content); err != nil {
		printError("Failed to process feed: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// invokeSkill invokes a Claude Code skill with the given input
func invokeSkill(skillName, input string) error {
	// Claude Code uses /skillname:action format or just mention the skill
	prompt := fmt.Sprintf("Use the %s skill to process this input:\n\n%s", skillName, input)

	cmd := exec.Command("claude", "-p", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
