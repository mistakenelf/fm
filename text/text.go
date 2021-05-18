package text

import (
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	HeaderText string
	BodyText   string
}

func NewModel() Model {
	return Model{
		HeaderText: "Help Header",
		BodyText:   "Help Body Text",
	}
}

func (m Model) View() string {
	text := ""

	if m.HeaderText != "" {
		text = lipgloss.JoinVertical(lipgloss.Top, m.HeaderText, m.BodyText)
	} else {
		text = lipgloss.JoinVertical(lipgloss.Top, m.BodyText)
	}

	return text
}
