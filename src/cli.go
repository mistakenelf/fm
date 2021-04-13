package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/src/components"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	useHighPerformanceRenderer = false
)

func (m model) Init() tea.Cmd {
	return m.getInitialDirectoryListing
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true

			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case fileStatus:
		m.Files = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Files)-1 {
				m.Cursor++
			}

		case "enter", " ":
			if m.Files[m.Cursor].IsDir() {
				m.Files = getUpdatedDirectoryListing(m.Files[m.Cursor].Name())
				m.Cursor = 0
			} else {
				dat, err := ioutil.ReadFile(m.Files[m.Cursor].Name())

				if err != nil {
					log.Fatal("Error occured reading file")
				}

				m.FileContent = string(dat)
			}
			m.Selected[m.Cursor] = struct{}{}

		case "h", "backspace":
			m.Cursor = 0
			m.Selected[m.Cursor] = struct{}{}
			m.Files = getUpdatedDirectoryListing("..")
		}
	case tea.WindowSizeMsg:
		if !m.ViewportReady {
			m.Viewport = viewport.Model{Width: msg.Width, Height: msg.Height}
			m.Viewport.YPosition = 0
			m.Viewport.HighPerformanceRendering = false
			m.ViewportReady = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}

	default:
		return m, nil
	}

	m.Viewport, cmd = m.Viewport.Update(msg)

	if useHighPerformanceRenderer {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ViewportReady {
		return "Loading..."
	}

	doc := strings.Builder{}
	files := ""

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(50)

	for i, file := range m.Files {
		files += fmt.Sprintf("%s\n", components.FileListing(file.Name(), m.Cursor == i, file.IsDir(), filepath.Ext(file.Name())))
	}

	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		style.Copy().Align(lipgloss.Left).Render(files),
		style.Copy().Align(lipgloss.Left).Render(m.FileContent),
	))

	m.Viewport.SetContent(doc.String())

	return m.Viewport.View()
}
