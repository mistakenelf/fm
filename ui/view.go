package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if !m.ready || m.dirTree.GetTotalFiles() <= 0 {
		return fmt.Sprintf("%s%s", m.loader.View(), "loading...")
	}

	panes := lipgloss.JoinHorizontal(lipgloss.Top, m.primaryPane.View(), m.secondaryPane.View())

	return lipgloss.JoinVertical(
		lipgloss.Top,
		panes,
		m.statusBar.View(),
	)
}
