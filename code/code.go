// Package code implements a code bubble which renders syntax highlighted
// source code based on a filename.
package code

import (
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mistakenelf/fm/filesystem"
)

type syntaxMsg string
type statusMessageTimeoutMsg struct{}
type errorMsg string

// Model represents the properties of a code bubble.
type Model struct {
	Viewport              viewport.Model
	Filename              string
	Content               string
	SyntaxTheme           string
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
	ViewportDisabled      bool
}

// Highlight returns a syntax highlighted string of text.
func Highlight(content, extension, syntaxTheme string) (string, error) {
	buf := new(bytes.Buffer)
	if err := quick.Highlight(buf, content, extension, "terminal256", syntaxTheme); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return buf.String(), nil
}

func readFileContentCmd(fileName, syntaxTheme string) tea.Cmd {
	return func() tea.Msg {
		content, err := filesystem.ReadFileContent(fileName)
		if err != nil {
			return errorMsg(err.Error())
		}

		highlightedContent, err := Highlight(content, filepath.Ext(fileName), syntaxTheme)
		if err != nil {
			return errorMsg(err.Error())
		}

		return syntaxMsg(highlightedContent)
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

	// Wait for timeout
	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}

// SetFileName sets current file to highlight.
func (m *Model) SetFileNameCmd(filename string) tea.Cmd {
	m.Filename = filename

	return readFileContentCmd(filename, m.SyntaxTheme)
}

// New creates a new instance of code.
func New() Model {
	viewPort := viewport.New(0, 0)

	return Model{
		Viewport:              viewPort,
		SyntaxTheme:           "dracula",
		StatusMessage:         "",
		StatusMessageLifetime: time.Second,
	}
}

// Init initializes the code bubble.
func (m Model) Init() tea.Cmd {
	return nil
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

// GotoBottom jumps to the bottom of the viewport.
func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

// SetViewportDisabled toggles the state of the viewport.
func (m *Model) SetViewportDisabled(disabled bool) {
	m.ViewportDisabled = disabled
}

// Update handles updating the UI of the code bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case syntaxMsg:
		m.Filename = ""
		m.Content = lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render(string(msg))

		m.Viewport.SetContent(m.Content)

	case statusMessageTimeoutMsg:
		m.StatusMessage = ""
	case errorMsg:
		m.Filename = ""
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cc241d")).
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

// View returns a string representation of the code bubble.
func (m Model) View() string {
	return m.Viewport.View()
}
