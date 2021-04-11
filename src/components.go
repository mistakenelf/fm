package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func checkbox(label string, checked bool) string {
	var style = lipgloss.NewStyle().
		Align(lipgloss.Center)

	if checked {
		return style.Render("[x] " + label)
	}

	return fmt.Sprintf(style.Render("[ ] %s"), label)
}
