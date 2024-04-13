// Package pdf provides an pdf bubble which can render
// pdf files as strings.
package pdf

import (
	"bytes"
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ledongthuc/pdf"
	"github.com/mistakenelf/fm/polish"
)

type renderPDFMsg string
type errorMsg string
type statusMessageTimeoutMsg struct{}

// Model represents the properties of a pdf bubble.
type Model struct {
	Viewport              viewport.Model
	ViewportDisabled      bool
	FileName              string
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
}

// ReadPDF reads the content of a PDF and returns it as a string.
func ReadPDF(name string) (string, error) {
	file, reader, err := pdf.Open(name)
	if err != nil {
		return "", errors.Unwrap(err)
	}

	defer func() {
		if e := file.Close(); e != nil {
			err = e
		}
	}()

	buf := new(bytes.Buffer)
	buffer, err := reader.GetPlainText()

	if err != nil {
		return "", errors.Unwrap(err)
	}

	_, err = buf.ReadFrom(buffer)
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return buf.String(), nil
}

func renderPDFCmd(filename string) tea.Cmd {
	return func() tea.Msg {
		pdfContent, err := ReadPDF(filename)
		if err != nil {
			return errorMsg(err.Error())
		}

		return renderPDFMsg(pdfContent)
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
// returns a cmd which will render the pdf.
func (m *Model) SetFileNameCmd(filename string) tea.Cmd {
	m.FileName = filename

	return renderPDFCmd(filename)
}

// New creates a new instance of a PDF.
func New() Model {
	viewPort := viewport.New(0, 0)

	return Model{
		Viewport:              viewPort,
		ViewportDisabled:      false,
		StatusMessageLifetime: time.Second,
	}
}

// Init initializes the PDF bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetSize sets the size of the bubble.
func (m *Model) SetSize(w, h int) {
	m.Viewport.Width = w
	m.Viewport.Height = h
}

// SetViewportDisabled toggles the state of the viewport.
func (m *Model) SetViewportDisabled(disabled bool) {
	m.ViewportDisabled = disabled
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// GotoBottom jumps to the bottom of the viewport.
func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

// Update handles updating the UI of the pdf bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case renderPDFMsg:
		pdfContent := lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render(string(msg))

		m.Viewport.SetContent(pdfContent)

		return m, nil
	case errorMsg:
		m.FileName = ""
		return m, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Foreground(polish.Colors.Red600).
				Bold(true).
				Render(string(msg)),
		)
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
