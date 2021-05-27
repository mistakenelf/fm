package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.Ready || m.DirTree.GetTotalFiles() <= 0 {
		return fmt.Sprintf("%s%s", m.Spinner.View(), "loading...")
	}

	panes := lipgloss.JoinHorizontal(lipgloss.Top, m.PrimaryPane.View(), m.SecondaryPane.View())

	return lipgloss.JoinVertical(
		lipgloss.Top,
		panes,
		m.StatusBar.View(),
	)
}
