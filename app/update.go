package app

import (
	"github.com/knipferrc/fm/components"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) scrollPrimaryViewport() {
	top := m.Viewport.YOffset
	bottom := m.Viewport.Height + m.Viewport.YOffset - 1

	if m.Cursor < top {
		m.Viewport.LineUp(1)
	} else if m.Cursor > bottom {
		m.Viewport.LineDown(1)
	}

	if m.Cursor > len(m.Files)-1 {
		m.Cursor = 0
		m.Viewport.GotoTop()
	} else if m.Cursor < 0 {
		m.Cursor = len(m.Files) - 1
		m.Viewport.GotoBottom()
	}
}

func (m Model) handleKeyDown() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Cursor++
		m.scrollPrimaryViewport()
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	}

	return m, nil
}

func (m Model) handleKeyUp() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Cursor--
		m.scrollPrimaryViewport()
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	}

	return m, nil
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.Rename {
		return m, renameFileOrDir(m.Files[m.Cursor].Name(), m.Textinput.Value())
	} else if m.Move {
		if m.Files[m.Cursor].IsDir() {
			return m, moveDir(m.Files[m.Cursor].Name(), m.Textinput.Value())
		} else {
			return m, moveFile(m.Files[m.Cursor].Name(), m.Textinput.Value())
		}
	} else if m.Delete {
		if m.Files[m.Cursor].IsDir() {
			if m.Textinput.Value() == "y" {
				return m, deleteDir(m.Files[m.Cursor].Name())
			} else {
				m.Textinput.Blur()
				m.Textinput.Reset()
				m.Delete = false
			}
		} else {
			if m.Textinput.Value() == "y" {
				return m, deleteFile(m.Files[m.Cursor].Name())
			} else {
				m.Textinput.Blur()
				m.Textinput.Reset()
				m.Delete = false
			}
		}
	} else {
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

		return m, nil
	}

	m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

	return m, nil
}

func (m Model) handleMoveKey() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Move = true
		m.Textinput.Placeholder = "/usr/share/"
		m.Textinput.Focus()
	}

	m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

	return m, nil
}

func (m Model) handleRenameKey() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Rename = true
		m.Textinput.Placeholder = "newfilename.ex"
		m.Textinput.Focus()
	}

	m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

	return m, nil
}

func (m Model) handleDeleteKey() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Delete = true
		m.Textinput.Placeholder = "[y/n]"
		m.Textinput.Focus()
	}

	m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

	return m, nil
}

func (m Model) handleEscKey() (tea.Model, tea.Cmd) {
	m.Move = false
	m.Rename = false
	m.Delete = false
	m.Textinput.Blur()
	m.Textinput.Reset()

	m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case updateDirMsg:
		m.Files = msg
		m.Cursor = 0
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	case renameMsg:
		m.Files = msg
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
		m.Textinput.Blur()
		m.Textinput.Reset()
		m.Rename = false
	case moveMsg:
		m.Files = msg
		m.Cursor = 0
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
		m.Textinput.Blur()
		m.Textinput.Reset()
		m.Move = false
	case deleteMsg:
		m.Files = msg
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
		m.Textinput.Blur()
		m.Textinput.Reset()
		m.Delete = false
	case fileContentMsg:
		m.SecondaryViewport.SetContent(string(msg))
	case tea.WindowSizeMsg:
		if !m.Ready {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height

			m.Viewport = viewport.Model{
				Width:  msg.Width / 2,
				Height: msg.Height - 3,
			}
			m.SecondaryViewport = viewport.Model{
				Width:  (msg.Width / 2) - 3,
				Height: msg.Height - 3,
			}

			m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
			m.SecondaryViewport.SetContent(components.Help())

			m.Ready = true
		} else {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height
			m.Viewport.Width = msg.Width / 2
			m.Viewport.Height = msg.Height - 3
			m.SecondaryViewport.Width = (msg.Width / 2) - 3
			m.SecondaryViewport.Height = msg.Height - 3
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.Rename && !m.Delete && !m.Move {
				return m, tea.Quit
			}
		case "h":
			if !m.Rename && !m.Delete && !m.Move {
				return m, updateDirectoryListing("..")
			}
		case "down", "j":
			if !m.Rename && !m.Delete && !m.Move {
				return m.handleKeyDown()
			}
		case "up", "k":
			if !m.Rename && !m.Delete && !m.Move {
				return m.handleKeyUp()
			}
		case "l":
			if !m.Rename && !m.Delete && !m.Move {
				if m.Files[m.Cursor].IsDir() && !m.Textinput.Focused() {
					return m, updateDirectoryListing(m.Files[m.Cursor].Name())
				} else {
					return m, readFileContent(m.Files[m.Cursor].Name())
				}
			}
		case "enter":
			return m.handleEnterKey()
		case "m":
			if !m.Rename && !m.Delete && !m.Move {
				return m.handleMoveKey()
			}
		case "r":
			if !m.Rename && !m.Delete && !m.Move {
				return m.handleRenameKey()
			}
		case "d":
			if !m.Rename && !m.Delete && !m.Move {
				return m.handleDeleteKey()
			}
		case "esc":
			return m.handleEscKey()
		}
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
