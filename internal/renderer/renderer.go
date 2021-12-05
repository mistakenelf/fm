package renderer

import (
	"image"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/fm/internal/commands"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/strfmt"
)

// Model is a struct that contains all the properties of renderer.
type Model struct {
	Image               image.Image
	Style               lipgloss.Style
	Viewport            viewport.Model
	ActiveBorderColor   lipgloss.AdaptiveColor
	InactiveBorderColor lipgloss.AdaptiveColor
	Borderless          bool
	IsActive            bool
	SyntaxTheme         string
	Content             string
	Height              int
	Width               int
}

// NewModel creates a new instance of a renderer.
func NewModel(
	borderless, isActive bool,
	activeBorderColor, inactiveBorderColor lipgloss.AdaptiveColor,
) Model {
	border := lipgloss.NormalBorder()
	padding := 1

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	style := lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border)

	return Model{
		Borderless:          borderless,
		IsActive:            isActive,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		Style:               style,
	}
}

// SetSize sets the size of the renderer.
func (m *Model) SetSize(width, height int) {
	m.Width = (width / 2) - m.Style.GetHorizontalBorderSize()
	m.Height = height - m.Style.GetVerticalBorderSize() - statusbar.StatusbarHeight

	m.Viewport.Width = m.Width - m.Style.GetHorizontalFrameSize() - m.Style.GetHorizontalPadding()
	m.Viewport.Height = m.Height
}

// SetContent sets the content of the renderer.
func (m *Model) SetContent(content string) {
	m.Content = strfmt.ConvertTabsToSpaces(content)
	m.Viewport.SetContent(m.Content)
}

// SetImage sets the image of the renderer.
func (m *Model) SetImage(img image.Image) {
	m.Image = img
}

// GetIsActive returns the active state of the renderer.
func (m Model) GetIsActive() bool {
	return m.IsActive
}

// SetIsActive sets the active state of the renderer.
func (m *Model) SetIsActive(isActive bool) {
	m.IsActive = isActive
}

// Update updates the statusbar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.ReadFileContentMsg:
		switch {
		case msg.Code != "":
			m.Viewport.GotoTop()
			m.SetContent(msg.Code)
			m.SetImage(nil)
		case msg.PDFContent != "":
			m.Viewport.GotoTop()
			m.SetContent(msg.PDFContent)
			m.SetImage(nil)
		case msg.Markdown != "":
			m.Viewport.GotoTop()
			m.SetContent(msg.Markdown)
			m.SetImage(nil)
		case msg.Image != nil:
			m.Viewport.GotoTop()
			m.SetContent(msg.ImageString)
			m.SetImage(msg.Image)
		default:
			m.Viewport.GotoTop()
			m.SetContent(msg.RawContent)
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.SetContent(m.Content)
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			if m.IsActive {
				m.Viewport.LineDown(1)
			}
		case "up", "k":
			if m.IsActive {
				m.Viewport.LineUp(1)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of a renderer.
func (m Model) View() string {
	borderColor := m.InactiveBorderColor

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return m.Style.Copy().
		BorderForeground(borderColor).
		Width(m.Width).
		Height(m.Height).
		Render(m.Viewport.View())
}
