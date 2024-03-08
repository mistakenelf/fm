// Package code implements a code bubble which renders syntax highlighted
// source code based on a filename.
package code

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/fm/filesystem"
)

type syntaxMsg string
type errorMsg error

// Highlight returns a syntax highlighted string of text.
func Highlight(content, extension, syntaxTheme string) (string, error) {
	buf := new(bytes.Buffer)
	if err := quick.Highlight(buf, content, extension, "terminal256", syntaxTheme); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return buf.String(), nil
}

// readFileContentCmd reads the content of the file.
func readFileContentCmd(fileName, syntaxTheme string) tea.Cmd {
	return func() tea.Msg {
		content, err := filesystem.ReadFileContent(fileName)
		if err != nil {
			return errorMsg(err)
		}

		highlightedContent, err := Highlight(content, filepath.Ext(fileName), syntaxTheme)
		if err != nil {
			return errorMsg(err)
		}

		return syntaxMsg(highlightedContent)
	}
}

// Model represents the properties of a code bubble.
type Model struct {
	Viewport           viewport.Model
	Active             bool
	Filename           string
	HighlightedContent string
	SyntaxTheme        string
}

// New creates a new instance of code.
func New(active bool) Model {
	viewPort := viewport.New(0, 0)

	return Model{
		Viewport:    viewPort,
		Active:      active,
		SyntaxTheme: "dracula",
	}
}

// Init initializes the code bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to highlight.
func (m *Model) SetFileName(filename string) tea.Cmd {
	m.Filename = filename

	return readFileContentCmd(filename, m.SyntaxTheme)
}

// SetIsActive sets if the bubble is currently active.
func (m *Model) SetIsActive(active bool) {
	m.Active = active
}

// SetSyntaxTheme sets the syntax theme of the rendered code.
func (m *Model) SetSyntaxTheme(theme string) {
	m.SyntaxTheme = theme
}

// SetSize sets the size of the bubble.
func (m *Model) SetSize(w, h int) {
	m.Viewport.Width = w
	m.Viewport.Height = h
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// Update handles updating the UI of a code bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case syntaxMsg:
		m.Filename = ""
		m.HighlightedContent = lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render(string(msg))

		m.Viewport.SetContent(m.HighlightedContent)

		return m, nil
	case errorMsg:
		m.Filename = ""
		m.HighlightedContent = lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render("Error: " + msg.Error())

		m.Viewport.SetContent(m.HighlightedContent)

		return m, nil
	}

	if m.Active {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the code bubble.
func (m Model) View() string {
	return m.Viewport.View()
}
