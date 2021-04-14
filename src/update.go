package main

import (
	"github.com/knipferrc/fm/src/config"
	"github.com/knipferrc/fm/src/filesystem"

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
				m.Files = filesystem.GetDirectoryListing(m.Files[m.Cursor].Name())
				m.Cursor = 0
			} else if m.Rename {
				filesystem.RenameDirOrFile(m.Files[m.Cursor].Name(), m.TextInput.Value())
				m.Files = filesystem.GetDirectoryListing("./")
				m.TextInput.Blur()
				m.Rename = false
			} else if m.Move {
				if m.Files[m.Cursor].IsDir() {
					filesystem.MoveDir(m.Files[m.Cursor].Name(), m.TextInput.Value())
					m.Files = filesystem.GetDirectoryListing("./")
					m.TextInput.Blur()
					m.Move = false
				} else {
					filesystem.MoveFile(m.Files[m.Cursor].Name(), m.TextInput.Value())
					m.Files = filesystem.GetDirectoryListing("./")
					m.TextInput.Blur()
					m.Move = false
				}
			} else if m.Delete {
				if m.Files[m.Cursor].IsDir() {
					if m.TextInput.Value() == "y" {
						filesystem.DeleteDirectory(m.Files[m.Cursor].Name())
						m.Files = filesystem.GetDirectoryListing("./")
						m.TextInput.Blur()
						m.Delete = false
					}
				} else {
					if m.TextInput.Value() == "y" {
						filesystem.DeleteFile(m.Files[m.Cursor].Name())
						m.Files = filesystem.GetDirectoryListing("./")
						m.TextInput.Blur()
						m.Delete = false
					}
				}
			}

		case "h", "backspace":
			if !m.TextInput.Focused() {
				m.Cursor = 0
				m.Files = filesystem.GetDirectoryListing("..")
			}

		case "m":
			if !m.TextInput.Focused() {
				m.Move = true
				m.TextInput.Placeholder = "/usr/share/"
				m.TextInput.Focus()
			}

		case "r":
			if !m.TextInput.Focused() {
				m.Rename = true
				m.TextInput.Placeholder = "newfilename.ex"
				m.TextInput.Focus()
			}

		case "d":
			if !m.TextInput.Focused() {
				m.Delete = true
				m.TextInput.Placeholder = "[y/n]"
				m.TextInput.Focus()
			}

		case "esc":
			m.Move = false
			m.Rename = false
			m.Delete = false
			m.TextInput.Blur()
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
