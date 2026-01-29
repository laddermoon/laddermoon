package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print LadderMoon version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("LadderMoon v%s\n", Version)
		fmt.Println("AI-driven project management tool")
		fmt.Println("https://github.com/laddermoon/laddermoon")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
