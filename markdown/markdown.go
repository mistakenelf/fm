// Package markdown provides an markdown bubble which can render
// markdown in a pretty manor.
package markdown

import (
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/mistakenelf/fm/filesystem"
	"github.com/mistakenelf/fm/polish"
)

type renderMarkdownMsg string
type errorMsg string
type statusMessageTimeoutMsg struct{}

// Model represents the properties of a markdown bubble.
type Model struct {
	Viewport              viewport.Model
	ViewportDisabled      bool
	FileName              string
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
}

// RenderMarkdown renders the markdown content with glamour.
func RenderMarkdown(width int, content string) (string, error) {
	background := "light"

	if lipgloss.HasDarkBackground() {
		background = "dark"
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithStandardStyle(background),
	)

	out, err := r.Render(content)
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return out, nil
}

func renderMarkdownCmd(width int, filename string) tea.Cmd {
	return func() tea.Msg {
		content, err := filesystem.ReadFileContent(filename)
		if err != nil {
			return errorMsg(err.Error())
		}

		markdownContent, err := RenderMarkdown(width, content)
		if err != nil {
			return errorMsg(err.Error())
		}

		return renderMarkdownMsg(markdownContent)
	}
}

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time.
func (m *Model) NewStatusMessageCmd(s string) tea.Cmd {
	m.StatusMessage = s

	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}

	m.statusMessageTimer = time.NewTimer(m.StatusMessageLifetime)

	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}

// SetFileName sets current file to render, this
// returns a cmd which will render the text.
func (m *Model) SetFileNameCmd(filename string) tea.Cmd {
	m.FileName = filename

	return renderMarkdownCmd(m.Viewport.Width, filename)
}

// SetSize sets the size of the bubble.
func (m *Model) SetSizeCmd(w, h int) tea.Cmd {
	m.Viewport.Width = w
	m.Viewport.Height = h

	if m.FileName != "" {
		return renderMarkdownCmd(m.Viewport.Width, m.FileName)
	}

	return nil
}

// New creates a new instance of markdown.
func New() Model {
	viewPort := viewport.New(0, 0)

	return Model{
		Viewport:              viewPort,
		ViewportDisabled:      false,
		StatusMessageLifetime: time.Second,
	}
}

// Init initializes the code bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// GotoBottom jumps to the bottom of the viewport.
func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

// SetViewportDisabled toggles the state of the viewport.
func (m *Model) SetViewportDisabled(disabled bool) {
	m.ViewportDisabled = disabled
}

// Update handles updating the UI of a markdown bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case renderMarkdownMsg:
		content := lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render(string(msg))

		m.Viewport.SetContent(content)

		return m, nil
	case errorMsg:
		m.FileName = ""
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Foreground(polish.Colors.Red600).
				Bold(true).
				Render(string(msg)),
		))
	}

	if !m.ViewportDisabled {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the markdown bubble.
func (m Model) View() string {
	return m.Viewport.View()
}
