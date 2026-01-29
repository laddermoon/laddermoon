package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lm",
	Short: "LadderMoon - AI-driven project management",
	Long: `LadderMoon (lm) is an AI-driven project management tool.
	
Core concept: "AI AS ME" - Let AI become your shadow self,
learning your architectural preferences and decision patterns.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func printSuccess(msg string) {
	fmt.Fprintf(os.Stdout, "[LadderMoon] %s\n", msg)
}

func printError(msg string) {
	fmt.Fprintf(os.Stderr, "[LadderMoon] Error: %s\n", msg)
}

func printInfo(msg string) {
	fmt.Fprintf(os.Stdout, "[LadderMoon] %s\n", msg)
}
