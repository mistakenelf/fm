package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Help() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFDF5")).
		MarginBottom(1).
		Render("Welcome to FM!")

	helpText := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
		"h - go back a directory",
		"j - move cursor down",
		"k - move cursor up",
		"l - open selected folder",
		"m - move file or folder to another directory",
		"d - delete a file or directory",
		"r - rename a file or directory",
	)

	return lipgloss.JoinVertical(lipgloss.Top, header, helpText)
}
