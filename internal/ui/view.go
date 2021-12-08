package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the entire application UI.
func (m Bubble) View() string {
	currentPane := m.renderer.View()

	if m.appState.ShowPreview {
		currentPane = m.previewer.View()
	}
	horizontalView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.fileTree.View(),
		currentPane,
	)

	if m.appConfig.Settings.SimpleMode {
		horizontalView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.fileTree.View(),
		)
	}

	// Return the UI with the two panes side by side and
	// the status bar at the bottom of the screen.
	return lipgloss.JoinVertical(
		lipgloss.Top,
		horizontalView,
		m.statusBar.View(),
	)
}
