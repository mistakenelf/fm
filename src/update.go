package main

import (
	"github.com/knipferrc/fm/src/components"
	"github.com/knipferrc/fm/src/filesystem"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) fixViewport(moveCursor bool) {
	top := m.viewport.YOffset
	bottom := m.viewport.Height + m.viewport.YOffset - 1

	if moveCursor {
		if m.cursor < top {
			m.cursor = top
		} else if m.cursor > bottom {
			m.cursor = bottom
		}
		return
	}

	if m.cursor < top {
		m.viewport.LineUp(1)
	} else if m.cursor > bottom {
		m.viewport.LineDown(1)
	}
}

func (m *model) fixCursor() {
	if m.cursor > len(m.files)-1 {
		m.cursor = 0
	} else if m.cursor < 0 {
		m.cursor = len(m.files) - 1
	}
}

func (m model) handleKeyUp() (tea.Model, tea.Cmd) {
	if !m.textinput.Focused() && !m.showhelp {
		m.cursor--
		m.fixCursor()
		m.fixViewport(false)
	}

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) handleKeyDown() (tea.Model, tea.Cmd) {
	if !m.textinput.Focused() && !m.showhelp {
		m.cursor++
		m.fixCursor()
		m.fixViewport(false)
	}

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.showhelp {
		return m, nil
	}

	if m.files[m.cursor].IsDir() && !m.textinput.Focused() {
		m.files = filesystem.GetDirectoryListing(m.files[m.cursor].Name())
		m.cursor = 0
	} else if m.rename {
		filesystem.RenameDirOrFile(m.files[m.cursor].Name(), m.textinput.Value())
		m.files = filesystem.GetDirectoryListing("./")
		m.textinput.Blur()
		m.rename = false
	} else if m.move {
		if m.files[m.cursor].IsDir() {
			filesystem.CopyDir(m.files[m.cursor].Name(), m.textinput.Value(), true)
			m.files = filesystem.GetDirectoryListing("./")
			m.textinput.Blur()
			m.move = false
		} else {
			filesystem.CopyFile(m.files[m.cursor].Name(), m.textinput.Value(), true)
			m.files = filesystem.GetDirectoryListing("./")
			m.textinput.Blur()
			m.move = false
		}
	} else if m.delete {
		if m.files[m.cursor].IsDir() {
			if m.textinput.Value() == "y" {
				filesystem.DeleteDirectory(m.files[m.cursor].Name())
				m.files = filesystem.GetDirectoryListing("./")
				m.textinput.Blur()
				m.delete = false
			} else {
				m.files = filesystem.GetDirectoryListing("./")
				m.textinput.Blur()
				m.delete = false
			}
		} else {
			if m.textinput.Value() == "y" {
				filesystem.DeleteFile(m.files[m.cursor].Name())
				m.files = filesystem.GetDirectoryListing("./")
				m.textinput.Blur()
				m.delete = false
			} else {
				m.files = filesystem.GetDirectoryListing("./")
				m.textinput.Blur()
				m.delete = false
			}
		}
	} else {
		m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

		return m, nil
	}

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) handleBackKey() (tea.Model, tea.Cmd) {
	if !m.textinput.Focused() && !m.showhelp {
		m.cursor = 0
		m.files = filesystem.GetDirectoryListing("..")
	}

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) handleMoveKey() (tea.Model, tea.Cmd) {
	if !m.textinput.Focused() && !m.showhelp {
		m.move = true
		m.textinput.Placeholder = "/usr/share/"
		m.textinput.Focus()
	}

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) handleRenameKey() (tea.Model, tea.Cmd) {
	if !m.textinput.Focused() && !m.showhelp {
		m.rename = true
		m.textinput.Placeholder = "newfilename.ex"
		m.textinput.Focus()
	}

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) handleDeleteKey() (tea.Model, tea.Cmd) {
	if !m.textinput.Focused() && !m.showhelp {
		m.delete = true
		m.textinput.Placeholder = "[y/n]"
		m.textinput.Focus()
	}

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) handleHelpKey() (tea.Model, tea.Cmd) {
	m.showhelp = true
	m.viewport.SetContent(components.Help())

	return m, nil
}

func (m model) handleEscKey() (tea.Model, tea.Cmd) {
	m.move = false
	m.rename = false
	m.delete = false
	m.showhelp = false
	m.textinput.Blur()

	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.screenwidth = msg.Width
			m.screenheight = msg.Height
			m.viewport = viewport.Model{
				Width:  msg.Width,
				Height: msg.Height - 1,
			}
			m.viewport.YPosition = 0
			m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))
			m.ready = true
		} else {
			m.screenwidth = msg.Width
			m.screenheight = msg.Height
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 1
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			return m.handleKeyUp()
		case "down", "j":
			return m.handleKeyDown()
		case "enter", " ":
			return m.handleEnterKey()
		case "h":
			return m.handleBackKey()
		case "m":
			return m.handleMoveKey()
		case "r":
			return m.handleRenameKey()
		case "d":
			return m.handleDeleteKey()
		case "i":
			return m.handleHelpKey()
		case "esc":
			return m.handleEscKey()
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
