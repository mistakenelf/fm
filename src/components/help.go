package components

import "github.com/charmbracelet/lipgloss"

func Help() string {
	header := lipgloss.NewStyle().Bold(true).Render("Welcome to FM!")

	return header
}
