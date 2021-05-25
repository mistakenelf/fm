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
		HeaderText: "Header",
		BodyText:   "Body Text",
	}
}

func (m Model) View() string {
	if m.HeaderText != "" {
		return lipgloss.JoinVertical(lipgloss.Top, m.HeaderText, m.BodyText)
	} else {
		return m.BodyText
	}
}
