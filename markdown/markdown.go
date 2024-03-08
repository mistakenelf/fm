// Package markdown provides an markdown bubble which can render
// markdown in a pretty manor.
package markdown

import (
	"errors"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/fm/filesystem"
)

type renderMarkdownMsg string
type errorMsg error

// Model represents the properties of a code bubble.
type Model struct {
	Viewport viewport.Model
	Active   bool
	FileName string
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

// renderMarkdownCmd renders text as pretty markdown.
func renderMarkdownCmd(width int, filename string) tea.Cmd {
	return func() tea.Msg {
		content, err := filesystem.ReadFileContent(filename)
		if err != nil {
			return errorMsg(err)
		}

		markdownContent, err := RenderMarkdown(width, content)
		if err != nil {
			return errorMsg(err)
		}

		return renderMarkdownMsg(markdownContent)
	}
}

// New creates a new instance of markdown.
func New(active bool) Model {
	viewPort := viewport.New(0, 0)

	return Model{
		Viewport: viewPort,
		Active:   active,
	}
}

// Init initializes the code bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to render, this
// returns a cmd which will render the text.
func (m *Model) SetFileName(filename string) tea.Cmd {
	m.FileName = filename

	return renderMarkdownCmd(m.Viewport.Width, filename)
}

// SetSize sets the size of the bubble.
func (m *Model) SetSize(w, h int) tea.Cmd {
	m.Viewport.Width = w
	m.Viewport.Height = h

	if m.FileName != "" {
		return renderMarkdownCmd(m.Viewport.Width, m.FileName)
	}

	return nil
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// SetIsActive sets if the bubble is currently active.
func (m *Model) SetIsActive(active bool) {
	m.Active = active
}

// Update handles updating the UI of a code bubble.
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
		m.Viewport.SetContent(msg.Error())

		return m, nil
	}

	if m.Active {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the markdown bubble.
func (m Model) View() string {
	return m.Viewport.View()
}
