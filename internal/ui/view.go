package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the entire application UI.
func (m Model) View() string {
	// If the viewport on the panes is not ready or we dont have any files to display
	// show the spinner.
	if !m.ready {
		return fmt.Sprintf("%s%s", m.loader.View(), "loading...")
	}

	// Return the UI with the two panes side by side and
	// the status bar at the bottom of the screen.
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.primaryPane.View(),
			m.secondaryPane.View(),
		),
		m.statusBar.View(),
	)
}
