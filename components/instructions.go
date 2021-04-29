package components

import (
	"fmt"

	"github.com/knipferrc/fm/constants"

	"github.com/charmbracelet/lipgloss"
)

func Instructions() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(constants.White)).
		MarginBottom(1).
		Render("FM (File Manager)")

	instructionText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.White)).
		Render(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
			"h - go back a directory",
			"j - move cursor down",
			"k - move cursor up",
			"l - open selected folder / view file",
			"m - move file or folder to another directory",
			"d - delete a file or directory",
			"r - rename a file or directory",
			"tab - toggle between panes"),
		)

	return lipgloss.JoinVertical(lipgloss.Top, header, instructionText)
}
