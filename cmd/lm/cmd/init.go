package cmd

import (
	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize LadderMoon in the current Git repository",
	Long: `Initialize LadderMoon by creating the shadow branch 'laddermoon-meta'
and setting up the META structure.

This command must be run in a Git-managed repository and can only be
executed once per repository.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	gitRoot, err := meta.GetGitRoot()
	if err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	// Check if already initialized
	if meta.IsInitialized() {
		printError("LadderMoon is already initialized in this repository.")
		printInfo("The 'laddermoon-meta' branch already exists.")
		return meta.ErrAlreadyInit
	}

	printInfo("Initializing LadderMoon...")
	printInfo("Git root: " + gitRoot)

	// Create META structure
	if err := meta.InitMetaStructure(); err != nil {
		printError("Failed to initialize: " + err.Error())
		return err
	}

	printSuccess("LadderMoon initialized successfully!")
	printInfo("Created shadow branch: laddermoon-meta")
	printInfo("META structure:")
	printInfo("  - META.md (empty)")
	printInfo("  - Questions/")
	printInfo("  - Issues/")
	printInfo("  - Suggestions/")
	printInfo("")
	printInfo("Next steps:")
	printInfo("  - Run 'lm feed <info>' to add project information")
	printInfo("  - Run 'lm sync' to synchronize with code changes")
	printInfo("  - Run 'lm status' to check current state")

	return nil
}
