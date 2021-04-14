package main

import (
	"os"

	"github.com/knipferrc/fm/src/config"
	"github.com/knipferrc/fm/src/directory"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		viewportCmd  tea.Cmd
		textInputCmd tea.Cmd
		cmds         []tea.Cmd
	)

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		if k == "q" || k == "ctrl+c" {
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
			if m.Cursor > 0 && !m.TextInput.Focused() {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Files)-1 && !m.TextInput.Focused() {
				m.Cursor++
			}

		case "enter", " ":
			if m.Files[m.Cursor].IsDir() && !m.TextInput.Focused() {
				m.Files = directory.GetDirectoryListing(m.Files[m.Cursor].Name())
				m.Cursor = 0
			} else {
				os.Rename(m.Files[m.Cursor].Name(), m.TextInput.Value())
				m.Files = directory.GetDirectoryListing("./")
				m.TextInput.Blur()
				m.Rename = false
			}

		case "h", "backspace":
			if !m.TextInput.Focused() {
				m.Cursor = 0
				m.Files = directory.GetDirectoryListing("..")
			}

		case "m":
			m.Move = true
			m.TextInput.Placeholder = "/usr/share/"
			m.TextInput.Focus()

		case "r":
			m.Rename = true
			m.TextInput.Placeholder = "newfilename.ex"
			m.TextInput.Focus()

		case "esc":
			m.Move = false
			m.TextInput.Blur()
			m.Rename = false
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

	m.Viewport, viewportCmd = m.Viewport.Update(msg)
	m.TextInput, textInputCmd = m.TextInput.Update(msg)

	if config.UseHighPerformanceRenderer {
		cmds = append(cmds, viewportCmd)
	}

	cmds = append(cmds, textInputCmd)

	return m, tea.Batch(cmds...)
}
