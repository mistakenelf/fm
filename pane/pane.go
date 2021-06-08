package pane

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width               int
	Height              int
	YOffset             int
	IsActive            bool
	Viewport            viewport.Model
	ActiveBorderColor   string
	InactiveBorderColor string
	Rounded             bool
}

func NewModel(width, height int, isActive, rounded bool, activeBorderColor, inactiveBorderColor string) Model {
	border := lipgloss.NormalBorder()

	if rounded {
		border = lipgloss.RoundedBorder()
	}

	m := Model{
		Width:               width - lipgloss.Width(border.Right+border.Top),
		Height:              height - lipgloss.Height(border.Bottom),
		IsActive:            isActive,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		Rounded:             rounded,
	}

	m.Viewport.Width = width - lipgloss.Width(border.Right+border.Top)
	m.Viewport.Height = height - lipgloss.Height(border.Bottom)
	m.YOffset = m.Viewport.YOffset

	return m
}

func (m *Model) SetSize(width, height int) {
	border := lipgloss.NormalBorder()

	if m.Rounded {
		border = lipgloss.RoundedBorder()
	}

	m.Width = width - lipgloss.Width(border.Right+border.Top)
	m.Height = height - lipgloss.Width(border.Bottom)
	m.Viewport.Width = width - lipgloss.Width(border.Right+border.Top)
	m.Viewport.Height = height - lipgloss.Height(border.Bottom)
}

func (m *Model) SetContent(content string) {
	border := lipgloss.NormalBorder()

	if m.Rounded {
		border = lipgloss.RoundedBorder()
	}

	m.Viewport.SetContent(
		lipgloss.NewStyle().
			Width(m.Width - lipgloss.Width(border.Top+border.Left)).
			Render(content),
	)
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

func (m Model) View() string {
	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()

	if m.Rounded {
		border = lipgloss.RoundedBorder()
	}

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		Border(border).
		Width(m.Width).
		Height(m.Height).
		Render(m.Viewport.View())
}
