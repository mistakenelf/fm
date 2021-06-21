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
	IsPadded            bool
}

// Create a new instance of a pane
func NewModel(isActive, rounded, isPadded bool, activeBorderColor, inactiveBorderColor string) Model {
	m := Model{
		IsActive:            isActive,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		Rounded:             rounded,
		IsPadded:            isPadded,
	}

	// Set the offset of the pane to be the same as the viewport
	m.YOffset = m.Viewport.YOffset

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

	// Set widths of both the pane and viewport taking into account borders
	m.Width = width - lipgloss.Width(border.Right+border.Top)
	m.Height = height - lipgloss.Width(border.Bottom)
	m.Viewport.Width = width - lipgloss.Width(border.Right+border.Top)
	m.Viewport.Height = height - lipgloss.Height(border.Bottom)
}

// Set the content of what to be displayed inside it
func (m *Model) SetContent(content string) {
	padding := 0

	// If the pane requires padding, add it
	if m.IsPadded {
		padding = 1
	}

	// Set the content inside a viewport
	m.Viewport.SetContent(
		lipgloss.NewStyle().
			PaddingLeft(padding).
			Render(content),
	)
}

// Scroll the viewport up a specified number of lines
func (m *Model) LineUp(lines int) {
	m.Viewport.LineUp(lines)
}

// Scroll the viewport down a specified number of lines
func (m *Model) LineDown(lines int) {
	m.Viewport.LineDown(lines)
}

// Go to the top of the viewport
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// Go to the bottom of the viewport
func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

// Set the color of the active border
func (m *Model) SetActiveBorderColor(color string) {
	m.ActiveBorderColor = color
}

// Return the pane and all of its content
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
		Width(m.Width).
		Height(m.Height).
		Render(m.Viewport.View())
}
