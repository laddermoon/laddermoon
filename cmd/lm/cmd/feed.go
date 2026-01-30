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

	// Acquire lock for serialized META operations
	printInfo("Acquiring META lock...")
	lock, err := meta.AcquireMetaLock()
	if err != nil {
		printError("Failed to acquire lock: " + err.Error())
		return err
	}
	defer lock.Release()

	// Get next feed ID
	feedID, err := meta.GetNextFeedID()
	if err != nil {
		printError("Failed to get feed ID: " + err.Error())
		return err
	}

	printInfo(fmt.Sprintf("Recording Feed #%d...", feedID))
	printInfo("Content: " + truncateString(content, 60))

	// Record to UserFeed.log
	if err := meta.RecordUserFeed(feedID, content); err != nil {
		printError("Failed to record feed: " + err.Error())
		return err
	}

	// Increment feed ID for next use
	if err := meta.IncrementFeedID(feedID); err != nil {
		printError("Failed to increment feed ID: " + err.Error())
		return err
	}

	printInfo("Processing with AI...")

	// Invoke Claude Code with the laddermoon-feed skill, passing feed ID
	if err := invokeFeedSkill(feedID, content); err != nil {
		printError("Failed to process feed: " + err.Error())
		printInfo("Make sure 'claude' CLI is installed and configured.")
		return err
	}

	printSuccess(fmt.Sprintf("Feed #%d recorded and processed!", feedID))
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// invokeFeedSkill invokes the laddermoon-feed skill with feed ID and content
func invokeFeedSkill(feedID int, content string) error {
	prompt := fmt.Sprintf("Use the laddermoon-feed skill to process Feed #%d:\n\n%s", feedID, content)

	// Use interactive mode (not -p) because the skill needs to modify files
	cmd := exec.Command("claude", prompt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
