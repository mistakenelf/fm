package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the UI.
func (m Model) View() string {
	leftBox := m.filetree[m.activeWorkspace].View()
	rightBox := m.help.View()

	switch m.state {
	case idleState:
		rightBox = m.help.View()
	case showCodeState:
		rightBox = m.code.View()
	case showImageState:
		rightBox = m.image.View()
	case showPdfState:
		rightBox = m.pdf.View()
	case showMarkdownState:
		rightBox = m.markdown.View()
	case showMoveState:
		rightBox = m.secondaryFiletree.View()
	case showCsvState:
		rightBox = m.csv.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
		m.statusbar.View(),
	)
}
