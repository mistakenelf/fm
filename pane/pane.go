package pane

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width       int
	Height      int
	IsActive    bool
	Viewport    viewport.Model
	ActiveBorderColor string
	InactiveBorderColor string
	BorderType  lipgloss.Border
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
	m.Viewport.Width = width
	m.Viewport.Height = height
}

func (m Model) View() string {
	borderColor := m.InactiveBorderColor

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	} 

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		Border(m.BorderType).
		Width(m.Width).
		Height(m.Height).
		Render(m.Viewport.View())
}

func (m *Model) SetContent(content string) {
	m.Viewport.SetContent(content)
}

func (m *Model) LineUp(lines int) {
	m.Viewport.LineUp(lines)
}

func (m *Model) LineDown(lines int) {
	m.Viewport.LineDown(lines)
}

func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}
