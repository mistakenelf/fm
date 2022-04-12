package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// View returns a string representation of the UI.
func (b Bubble) View() string {
	leftBox := b.filetree.View()
	rightBox := b.help.View()

	switch b.state {
	case idleState:
		rightBox = b.help.View()
	case showCodeState:
		rightBox = b.code.View()
	case showImageState:
		rightBox = b.image.View()
	case showPdfState:
		rightBox = b.pdf.View()
	case showMarkdownState:
		rightBox = b.markdown.View()
	}

	return lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, leftBox, rightBox),
		b.statusbar.View(),
	)
}
