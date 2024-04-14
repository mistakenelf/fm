// Package csv implements a csv bubble which renders a table with
// the content of the csv.
package csv

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/mistakenelf/fm/polish"
)

type statusMessageTimeoutMsg struct{}
type errorMsg string
type csvMsg struct {
	headers []string
	records [][]string
}

const (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")
)

var (
	HeaderStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	CellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
	OddRowStyle  = CellStyle.Copy().Foreground(gray)
	EvenRowStyle = CellStyle.Copy().Foreground(lightGray)
)

// Model represents the properties of a code bubble.
type Model struct {
	Viewport              viewport.Model
	Filename              string
	Table                 *table.Table
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
	ViewportDisabled      bool
	Headers               []string
	Records               [][]string
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

	return func() tea.Msg {
		file, err := os.Open(m.Filename)
		if err != nil {
			log.Fatal("Error while reading the file", err)
		}

		defer file.Close()

		reader := csv.NewReader(file)

		headers, err := reader.Read()
		if err != nil {
			fmt.Println(err)
		}

		records, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error reading records")
		}

		return csvMsg{headers, records}
	}
}

// New creates a new instance of code.
func New() Model {
	viewPort := viewport.New(0, 0)
	table := table.New()

	return Model{
		Viewport:              viewPort,
		StatusMessage:         "",
		StatusMessageLifetime: time.Second,
		Table:                 table,
	}
}

// Init initializes the code bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetSize sets the size of the bubble.
func (m *Model) SetSizeCmd(w, h int) tea.Cmd {
	m.Viewport.Width = w
	m.Viewport.Height = h

	if m.Filename != "" {
		return m.SetFileNameCmd(m.Filename)
	}

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

// Update handles updating the UI of the code bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case statusMessageTimeoutMsg:
		m.StatusMessage = ""

		return m, nil
	case csvMsg:
		m.Headers = msg.headers
		m.Records = msg.records

		m.Table = table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == 0:
					return HeaderStyle
				case row%2 == 0:
					return EvenRowStyle
				default:
					return OddRowStyle
				}
			}).
			Headers(m.Headers...).
			Width(m.Viewport.Width).
			Rows(m.Records...)

		m.Viewport.SetContent(m.Table.String())

		return m, nil
	case errorMsg:
		m.Filename = ""
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

// View returns a string representation of the csv bubble.
func (m Model) View() string {
	return m.Viewport.View()
}
