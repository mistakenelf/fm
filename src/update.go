package main

import (
	"github.com/knipferrc/fm/src/config"
	"github.com/knipferrc/fm/src/directory"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case fileStatus:
		m.Files = msg
		m.CurrentlyHighlighted = m.Files[m.Cursor].Name()

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
				m.CurrentlyHighlighted = m.Files[m.Cursor].Name()
			}

		case "down", "j":
			if m.Cursor < len(m.Files)-1 {
				m.Cursor++
				m.CurrentlyHighlighted = m.Files[m.Cursor].Name()
			}

		case "enter", " ":
			if m.Files[m.Cursor].IsDir() {
				m.Files = directory.GetDirectoryListing(m.Files[m.Cursor].Name())
				m.Cursor = 0
			}

		case "h", "backspace":
			m.Cursor = 0
			m.Files = directory.GetDirectoryListing("..")

		case "m":
			m.Move = true
		}

	case tea.WindowSizeMsg:
		m.ScreenWidth = msg.Width
		if !m.ViewportReady {
			m.Viewport = viewport.Model{Width: msg.Width, Height: msg.Height}
			m.Viewport.YPosition = 0
			m.Viewport.HighPerformanceRendering = config.UseHighPerformanceRenderer
			m.ViewportReady = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height
		}

		if config.UseHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}

	default:
		return m, nil
	}

	m.Viewport, cmd = m.Viewport.Update(msg)

	if config.UseHighPerformanceRenderer {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
