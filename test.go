package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func main() {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingTop(1).
		PaddingLeft(4).
		Width(40).
		Align(lipgloss.Center).
		Render("LADDERMOON")

	fmt.Println(style)
}
