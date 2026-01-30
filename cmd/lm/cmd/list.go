package cmd

import (
	"fmt"
	"strings"

	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks [id]",
	Short: "Show tasks",
	Long: `List all tasks or show a specific task.

Example:
  lm tasks                    # List all tasks
  lm tasks task-from-issue-001  # Show specific task`,
	RunE: runTasks,
}

var issuesCmd = &cobra.Command{
	Use:   "issues [id]",
	Short: "Show issues",
	Long: `List all issues or show a specific issue.

Example:
  lm issues              # List all issues
  lm issues issue-001    # Show specific issue`,
	RunE: runIssues,
}

var proposalsCmd = &cobra.Command{
	Use:   "proposals [id]",
	Short: "Show proposals",
	Long: `List all proposals or show a specific proposal.

Example:
  lm proposals                # List all proposals
  lm proposals proposal-001   # Show specific proposal`,
	RunE: runProposals,
}

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Show META.md content",
	Long:  `Display the content of META.md from the META branch.`,
	RunE:  runMeta,
}

var userlogCmd = &cobra.Command{
	Use:   "userlog",
	Short: "Show UserFeed.log content",
	Long:  `Display the content of UserFeed.log from the META branch.`,
	RunE:  runUserlog,
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	rootCmd.AddCommand(issuesCmd)
	rootCmd.AddCommand(proposalsCmd)
	rootCmd.AddCommand(metaCmd)
	rootCmd.AddCommand(userlogCmd)
}

func runTasks(cmd *cobra.Command, args []string) error {
	return showItems("Tasks", args)
}

func runIssues(cmd *cobra.Command, args []string) error {
	return showItems("Issues", args)
}

func runProposals(cmd *cobra.Command, args []string) error {
	return showItems("Proposals", args)
}

func showItems(directory string, args []string) error {
	if _, err := meta.GetGitRoot(); err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized. Run 'lm init' first.")
		return meta.ErrNotInitialized
	}

	// If specific ID provided, show that item
	if len(args) > 0 {
		itemID := args[0]
		// Try to find the file
		filePath := fmt.Sprintf("%s/%s.md", directory, itemID)
		content, err := meta.ReadFile(filePath)
		if err != nil || content == "" {
			// Try without .md extension in the ID
			filePath = fmt.Sprintf("%s/%s", directory, itemID)
			content, err = meta.ReadFile(filePath)
			if err != nil || content == "" {
				printError(fmt.Sprintf("Item not found: %s", itemID))
				return fmt.Errorf("item not found")
			}
		}
		fmt.Println(content)
		return nil
	}

	// List all items
	files, err := meta.GetMetaFileList()
	if err != nil {
		printError("Failed to get file list: " + err.Error())
		return err
	}

	items := []string{}
	for _, f := range files {
		if strings.HasPrefix(f, directory+"/") {
			items = append(items, f)
		}
	}

	if len(items) == 0 {
		printInfo(fmt.Sprintf("No %s found.", strings.ToLower(directory)))
		return nil
	}

	fmt.Printf("%s (%d):\n", directory, len(items))
	for _, item := range items {
		// Get status from content
		content, _ := meta.ReadFile(item)
		status := "Unknown"
		if strings.Contains(content, "**Status**: Open") {
			status = "Open"
		} else if strings.Contains(content, "**Status**: Resolved") {
			status = "Resolved"
		} else if strings.Contains(content, "**Status**: Approved") {
			status = "Approved"
		} else if strings.Contains(content, "**Status**: Rejected") {
			status = "Rejected"
		}

		// Extract title from first heading
		title := extractTitle(content)
		name := strings.TrimPrefix(item, directory+"/")
		fmt.Printf("  [%s] %s - %s\n", status, name, title)
	}

	return nil
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

func runMeta(cmd *cobra.Command, args []string) error {
	if _, err := meta.GetGitRoot(); err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized. Run 'lm init' first.")
		return meta.ErrNotInitialized
	}

	content, err := meta.ReadMetaFile()
	if err != nil {
		printError("Failed to read META.md: " + err.Error())
		return err
	}

	if content == "" {
		printInfo("META.md is empty.")
		return nil
	}

	fmt.Println(content)
	return nil
}

func runUserlog(cmd *cobra.Command, args []string) error {
	if _, err := meta.GetGitRoot(); err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	if !meta.IsInitialized() {
		printError("LadderMoon is not initialized. Run 'lm init' first.")
		return meta.ErrNotInitialized
	}

	content, err := meta.ReadFile("UserFeed.log")
	if err != nil {
		printError("Failed to read UserFeed.log: " + err.Error())
		return err
	}

	if content == "" {
		printInfo("UserFeed.log is empty.")
		return nil
	}

	fmt.Println(content)
	return nil
}
