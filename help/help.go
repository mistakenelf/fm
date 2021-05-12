package help

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
	return lipgloss.JoinVertical(lipgloss.Top, m.HeaderText, m.BodyText)
}
