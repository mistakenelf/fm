package pane

import (
	"fmt"

	"github.com/knipferrc/fm/internal/renderer"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model is a struct to represent the properties of a pane.
type Model struct {
	Viewport            viewport.Model
	Style               lipgloss.Style
	IsActive            bool
	Borderless          bool
	AlternateBorder     bool
	ShowLoading         bool
	ActiveBorderColor   lipgloss.AdaptiveColor
	InactiveBorderColor lipgloss.AdaptiveColor
	Spinner             spinner.Model
}

// NewModel creates an instance of a pane.
func NewModel(isActive, borderless bool, activeBorderColor, inactiveBorderColor lipgloss.AdaptiveColor) Model {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot

	return Model{
		IsActive:            isActive,
		Borderless:          borderless,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		Spinner:             s,
	}
}

// SetSize sets the size of the pane and its viewport, useful when resizing the terminal.
func (m *Model) SetSize(width, height int) {
	border := lipgloss.NormalBorder()
	padding := 1

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	// Set the style so that the frame size is able to be determined from other components.
	m.Style = lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border)

	m.Viewport.Width = width - m.Style.GetHorizontalBorderSize()
	m.Viewport.Height = height - m.Style.GetVerticalBorderSize()
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
func (m *Model) ShowAlternateBorder(show bool) {
	m.AlternateBorder = show
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

// ShowSpinner determines wether to show the spinner or not.
func (m *Model) ShowSpinner(show bool) {
	m.ShowLoading = show
}

// Update updates the pane.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the pane.
func (m Model) View() string {
	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()
	padding := 1
	content := m.Viewport.View()
	alternateBorder := lipgloss.Border{
		Top:         "-",
		Bottom:      "-",
		Left:        "|",
		Right:       "|",
		TopLeft:     "*",
		TopRight:    "*",
		BottomLeft:  "*",
		BottomRight: "*",
	}

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	if m.AlternateBorder {
		border = alternateBorder
	}

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	if m.ShowLoading {
		content = fmt.Sprintf("%s%s", m.Spinner.View(), "loading...")
	}

	return m.Style.Copy().
		BorderForeground(borderColor).
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		Width(m.Viewport.Width).
		Height(m.Viewport.Height).
		Render(renderer.ConvertTabsToSpaces(content))
}
