package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Main app view
func (m model) View() string {
	// If the viewport on the panes is not ready or we dont have any files to display
	// show the spinner
	if !m.ready || m.dirTree.GetTotalFiles() <= 0 {
		return fmt.Sprintf("%s%s", m.loader.View(), "loading...")
	}

	// Join the two panes horizontally and top aligned
	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.primaryPane.View(),
		m.secondaryPane.View(),
	)

	// Return the UI with the two panes side by side and
	// the status bar at the bottom of the screen
	return lipgloss.JoinVertical(
		lipgloss.Top,
		panes,
		m.statusBar.View(),
	)
}
