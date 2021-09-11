package pane

import (
	"github.com/knipferrc/fm/formatter"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// Model struct represents property of a pane.
type Model struct {
	Viewport            viewport.Model
	Style               lipgloss.Style
	IsActive            bool
	Borderless          bool
	ActiveBorderColor   string
	InactiveBorderColor string
}

// NewModel creates a new instance of a pane.
func NewModel(isActive, borderless bool, activeBorderColor, inactiveBorderColor string) Model {
	return Model{
		IsActive:            isActive,
		Borderless:          borderless,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
	}
}

// SetSize sets the size of the pane and its viewport, useful when resizing the terminal.
func (m *Model) SetSize(width, height int) {
	// Get the border so that when setting the width of a pane,
	// the border is also taken into account.
	border := lipgloss.NormalBorder()

	m.Viewport.Width = width - lipgloss.Width(border.Right+border.Top)
	m.Viewport.Height = height - lipgloss.Width(border.Bottom+border.Top)
}

// GetHorizontalFrameSize returns the horizontal frame size of the pane.
func (m Model) GetHorizontalFrameSize() int {
	return m.Style.GetHorizontalFrameSize()
}

// GetIsActive returns the active state of the pane.
func (m Model) GetIsActive() bool {
	return m.IsActive
}

// SetActive sets the active state of the pane.
func (m *Model) SetActive(isActive bool) {
	m.IsActive = isActive
}

// SetContent sets the content of the pane.
func (m *Model) SetContent(content string) {
	borderColor := m.InactiveBorderColor
	border := lipgloss.ThickBorder()
	padding := 1

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	// If the pane is active, use the active border color.
	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	// Set the style so that the frame size is able to be determined from other components.
	m.Style = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		Width(m.Viewport.Width)

	m.Viewport.SetContent(content)
}

// LineUp scrolls the pane up the specified number of lines.
func (m *Model) LineUp(lines int) {
	m.Viewport.LineUp(lines)
}

// LineDown scrolls the pane down the specified number of lines.
func (m *Model) LineDown(lines int) {
	m.Viewport.LineDown(lines)
}

// GotoTop goes to the top of the pane.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// GotoBottom goes to the bottom of the pane.
func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

// SetActiveBorderColors sets the active border colors.
func (m *Model) SetActiveBorderColor(color string) {
	m.ActiveBorderColor = color
}

// GetWidth returns the width of the pane.
func (m Model) GetWidth() int {
	return m.Viewport.Width
}

// GetHeight returns the height of the pane.
func (m Model) GetHeight() int {
	return m.Viewport.Height
}

// GetYOffset returns the y offset of the pane.
func (m Model) GetYOffset() int {
	return m.Viewport.YOffset
}

// View returns a string representation of the pane.
func (m Model) View() string {
	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()
	padding := 1

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	// If the pane is active, use the active border color.
	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return m.Style.Copy().
		BorderForeground(lipgloss.Color(borderColor)).
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		Width(m.Viewport.Width).
		Render(formatter.ConvertTabsToSpaces(m.Viewport.View()))
}
