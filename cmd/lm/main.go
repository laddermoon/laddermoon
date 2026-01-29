package main

import (
	"os"

	"github.com/laddermoon/laddermoon/cmd/lm/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
