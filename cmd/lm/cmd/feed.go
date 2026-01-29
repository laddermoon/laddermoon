package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var feedCmd = &cobra.Command{
	Use:   "feed [user input text]",
	Short: "Add project information to META",
	Long: `Record user-provided information to the META system.
	
This command appends the provided text to META.md and logs
it to UserFeed.log for tracking.

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

	// Join all args as the feed content
	content := strings.Join(args, " ")
	if strings.TrimSpace(content) == "" {
		printError("Feed content cannot be empty.")
		return fmt.Errorf("empty feed content")
	}

	printInfo("Recording your input...")

	// Create feed entry with timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	feedEntry := fmt.Sprintf("\n## User Feed [%s]\n\n%s\n", timestamp, content)

	// Append to META.md
	if err := meta.AppendToMetaFile(feedEntry); err != nil {
		printError("Failed to update META.md: " + err.Error())
		return err
	}

	// Log to UserFeed.log
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, content)
	if err := meta.AppendToFile("UserFeed.log", logEntry); err != nil {
		printError("Failed to log feed: " + err.Error())
		return err
	}

	printSuccess("Information recorded successfully!")
	fmt.Printf("  Content: %s\n", truncateString(content, 60))
	fmt.Printf("  Time: %s\n", timestamp)

	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
