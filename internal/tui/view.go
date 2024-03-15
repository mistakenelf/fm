package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the UI.
func (m model) View() string {
	leftBox := m.filetree.View()
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
	}

	if m.filetree.GetSelectedItem().Name != "" {
		statusMessage := m.filetree.GetSelectedItem().CurrentDirectory

		if m.filetree.StatusMessage != "" {
			statusMessage = m.filetree.StatusMessage
		}

		if m.code.StatusMessage != "" {
			statusMessage = m.code.StatusMessage
		}

		m.statusbar.SetContent(
			m.filetree.GetSelectedItem().Name,
			statusMessage,
			fmt.Sprintf("%d/%d", m.filetree.Cursor, m.filetree.GetTotalItems()),
			fmt.Sprintf("%s %s", "🗀", "FM"),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
		m.statusbar.View(),
	)
}
