package cmd

import (
	"github.com/laddermoon/laddermoon/pkg/meta"
	"github.com/laddermoon/laddermoon/skills"
	"github.com/spf13/cobra"
)

var reinstallSkills bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize LadderMoon in the current Git repository",
	Long: `Initialize LadderMoon by creating the shadow branch 'laddermoon-meta'
and setting up the META structure.

This command must be run in a Git-managed repository and can only be
executed once per repository.

Use --reinstall-skills to reinstall skills without re-initializing META.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVar(&reinstallSkills, "reinstall-skills", false, "Reinstall skills only (use if already initialized)")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	gitRoot, err := meta.GetGitRoot()
	if err != nil {
		printError("This command must be run inside a Git repository.")
		return err
	}

	// Handle --reinstall-skills flag
	if reinstallSkills {
		if !meta.IsInitialized() {
			printError("LadderMoon is not initialized. Run 'lm init' first.")
			return meta.ErrNotInitialized
		}

		printInfo("Reinstalling LadderMoon skills...")

		if err := meta.InstallSkills(skills.SkillsFS, skills.SkillNames); err != nil {
			printError("Failed to install skills: " + err.Error())
			return err
		}

		printSuccess("Skills reinstalled successfully!")
		printInfo("Installed skills:")
		for _, name := range skills.SkillNames {
			printInfo("  - " + name)
		}
		return nil
	}

	// Check if already initialized for this branch
	branchInitialized, err := meta.BranchMetaDirExists()
	if err != nil {
		printError("Failed to check initialization status: " + err.Error())
		return err
	}

	if branchInitialized {
		printError("LadderMoon is already initialized for this branch.")
		printInfo("Use --reinstall-skills to reinstall skills.")
		return meta.ErrAlreadyInit
	}

	printInfo("Initializing LadderMoon...")
	printInfo("Git root: " + gitRoot)

	// Create META structure
	if err := meta.InitMetaStructure(); err != nil {
		printError("Failed to initialize: " + err.Error())
		return err
	}

	printInfo("Installing LadderMoon skills...")

	// Install skills to .claude/skills
	if err := meta.InstallSkills(skills.SkillsFS, skills.SkillNames); err != nil {
		printError("Failed to install skills: " + err.Error())
		return err
	}

	// Get current branch for display
	currentBranch, _ := meta.GetCurrentBranch()

	printSuccess("LadderMoon initialized successfully!")
	printInfo("Branch: " + currentBranch)
	printInfo("Shadow branch: laddermoon-meta")
	printInfo("META structure:")
	printInfo("  - META.md (empty)")
	printInfo("  - Questions/")
	printInfo("  - Issues/")
	printInfo("  - Suggestions/")
	printInfo("")
	printInfo("Installed skills:")
	for _, name := range skills.SkillNames {
		printInfo("  - " + name)
	}
	printInfo("")
	printInfo("Next steps:")
	printInfo("  - Run 'lm feed <info>' to add project information")
	printInfo("  - Run 'lm sync' to synchronize with code changes")
	printInfo("  - Run 'lm status' to check current state")

	return nil
}
