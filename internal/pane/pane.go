package pane

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	IsActive            bool
	Viewport            viewport.Model
	ActiveBorderColor   string
	InactiveBorderColor string
	Rounded             bool
	IsPadded            bool
}

func NewModel(isActive, rounded, isPadded bool, activeBorderColor, inactiveBorderColor string) Model {
	m := Model{
		IsActive:            isActive,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		Rounded:             rounded,
		IsPadded:            isPadded,
	}

	return m
}

// Set the size of the pane, useful when changing screen sizes
func (m *Model) SetSize(width, height int) {
	// Get the border so that when setting the width of a pane,
	// the border is also taken into account
	border := lipgloss.NormalBorder()

	// If borders are rounded, use the rounded border
	if m.Rounded {
		border = lipgloss.RoundedBorder()
	}

	// Set width of the panes viewport
	m.Viewport.Width = width - lipgloss.Width(border.Right+border.Top)
	m.Viewport.Height = height - lipgloss.Width(border.Bottom+border.Top)
}

// Set content of the pane
func (m *Model) SetContent(content string) {
	padding := 0

	// If the pane requires padding, add it
	if m.IsPadded {
		padding = 1
	}

	// Place the pane content in a viewport
	m.Viewport.SetContent(
		lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			PaddingLeft(padding).
			Render(content),
	)
}

// Scroll pane up the number of specified lines
func (m *Model) LineUp(lines int) {
	m.Viewport.LineUp(lines)
}

// Scroll pane down the specified number of lines
func (m *Model) LineDown(lines int) {
	m.Viewport.LineDown(lines)
}

// Go to the top of the pane
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// Go to the bottom of the pane
func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

// Set active border color
func (m *Model) SetActiveBorderColor(color string) {
	m.ActiveBorderColor = color
}

// Display the pane
func (m Model) View() string {
	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()

	// If rounding is enabled on borders, use the round border
	if m.Rounded {
		border = lipgloss.RoundedBorder()
	}

	// If the pane is active, use the active border color
	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		Border(border).
		Width(m.Viewport.Width).
		Height(m.Viewport.Height).
		Render(m.Viewport.View())
}
