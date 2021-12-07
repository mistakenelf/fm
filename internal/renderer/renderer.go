package renderer

import (
	"image"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/fm/internal/commands"
	"github.com/knipferrc/fm/internal/statusbar"
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
	m.Viewport.Width = (width / 2) - m.Style.GetHorizontalBorderSize()
	m.Viewport.Height = height - m.Style.GetVerticalBorderSize() - statusbar.StatusbarHeight
}

// SetContent sets the content of the renderer.
func (m *Model) SetContent(content string) {
	curContent := lipgloss.NewStyle().
		Width(m.Viewport.Width - m.Style.GetHorizontalPadding()).
		Height(m.Viewport.Height - m.Style.GetVerticalPadding()).
		Render(content)

	m.Content = content
	m.Viewport.SetContent(curContent)
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
	case commands.ErrorMsg:
		m.SetContent(string(msg))
	case commands.CopyToClipboardMsg:
		m.SetContent(string(msg))
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
			m.SetImage(nil)
		}
	case commands.ConvertImageToStringMsg:
		m.SetContent(string(msg))
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.SetContent(m.Content)

		if m.Image != nil {
			cmds = append(cmds, commands.RedrawImageCmd(msg.Width, m.Image))
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if m.IsActive {
				m.Viewport.LineUp(1)
			}
		case tea.MouseWheelDown:
			if m.IsActive {
				m.Viewport.LineDown(1)
			}
		}
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
		case "ctrl+g":
			if m.IsActive {
				m.Viewport.GotoTop()
			}
		case "G":
			if m.IsActive {
				m.Viewport.GotoBottom()
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
		Width(m.Viewport.Width).
		Height(m.Viewport.Height).
		Render(m.Viewport.View())
}
