package app

import (
	"github.com/knipferrc/fm/components"
	"github.com/knipferrc/fm/constants"
	"github.com/muesli/reflow/wrap"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) scrollPrimaryViewport() {
	top := m.PrimaryViewport.YOffset
	bottom := m.PrimaryViewport.Height + m.PrimaryViewport.YOffset - 1

	if m.Cursor < top {
		m.PrimaryViewport.LineUp(1)
	} else if m.Cursor > bottom {
		m.PrimaryViewport.LineDown(1)
	}

	if m.Cursor > len(m.Files)-1 {
		m.Cursor = 0
		m.PrimaryViewport.GotoTop()
	} else if m.Cursor < 0 {
		m.Cursor = len(m.Files) - 1
		m.PrimaryViewport.GotoBottom()
	}
}

func (m Model) handleKeyDown() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() && m.ActivePane == constants.PrimaryPane {
		m.Cursor++
		m.scrollPrimaryViewport()
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	} else {
		m.SecondaryViewport.LineDown(1)
	}

	return m, nil
}

func (m Model) handleKeyUp() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() && m.ActivePane == constants.PrimaryPane {
		m.Cursor--
		m.scrollPrimaryViewport()
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	} else {
		m.SecondaryViewport.LineUp(1)
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
		return m, nil
	}

	return m, nil
}

func (m Model) handleMoveKey() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Move = true
		m.Textinput.Placeholder = "/usr/share/"
		m.Textinput.Focus()
	}

	return m, nil
}

func (m Model) handleRenameKey() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Rename = true
		m.Textinput.Placeholder = "new_name"
		m.Textinput.Focus()
	}

	return m, nil
}

func (m Model) handleDeleteKey() (tea.Model, tea.Cmd) {
	if !m.Textinput.Focused() {
		m.Delete = true
		m.Textinput.Placeholder = "[y/n]"
		m.Textinput.Focus()
	}

	return m, nil
}

func (m Model) handleEscKey() (tea.Model, tea.Cmd) {
	m.Move = false
	m.Rename = false
	m.Delete = false
	m.Textinput.Blur()
	m.Textinput.Reset()

	m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	m.SecondaryViewport.SetContent(components.Instructions())

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
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

	case renameMsg:
		m.Files = msg
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
		m.Textinput.Blur()
		m.Textinput.Reset()
		m.Rename = false

	case moveMsg:
		m.Files = msg
		m.Cursor = 0
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
		m.Textinput.Blur()
		m.Textinput.Reset()
		m.Move = false

	case deleteMsg:
		m.Files = msg
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
		m.Textinput.Blur()
		m.Textinput.Reset()
		m.Delete = false

	case fileContentMsg:
		border := lipgloss.NormalBorder()
		halfScreenWidth := m.ScreenWidth / 2
		borderWidth := lipgloss.Width(border.Left + border.Right + border.Top + border.Bottom)
		m.SecondaryViewport.SetContent(wrap.String(string(msg), halfScreenWidth-borderWidth))

	case tea.WindowSizeMsg:
		borderWidth := lipgloss.Width(lipgloss.NormalBorder().Top)
		statusBarHeight := 2
		verticalMargin := borderWidth + statusBarHeight

		if !m.Ready {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height

			m.PrimaryViewport = viewport.Model{
				Width:  msg.Width,
				Height: msg.Height - verticalMargin,
			}
			m.SecondaryViewport = viewport.Model{
				Width:  msg.Width,
				Height: msg.Height - verticalMargin,
			}

			m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
			m.SecondaryViewport.SetContent(components.Instructions())

			m.Ready = true
		} else {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height
			m.PrimaryViewport.Width = msg.Width
			m.PrimaryViewport.Height = msg.Height - verticalMargin
			m.SecondaryViewport.Width = msg.Width
			m.SecondaryViewport.Height = msg.Height - verticalMargin
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.Rename && !m.Delete && !m.Move {
				return m, tea.Quit
			}
		case "h":
			if !m.Rename && !m.Delete && !m.Move && m.ActivePane == constants.PrimaryPane {
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
			if !m.Rename && !m.Delete && !m.Move && m.ActivePane == constants.PrimaryPane {
				if m.Files[m.Cursor].IsDir() && !m.Textinput.Focused() {
					return m, updateDirectoryListing(m.Files[m.Cursor].Name())
				} else {
					return m, readFileContent(m.Files[m.Cursor].Name())
				}
			}
		case "enter":
			if m.ActivePane == constants.PrimaryPane {
				return m.handleEnterKey()
			}
		case "m":
			if !m.Rename && !m.Delete && !m.Move && m.ActivePane == constants.PrimaryPane {
				return m.handleMoveKey()
			}
		case "r":
			if !m.Rename && !m.Delete && !m.Move && m.ActivePane == constants.PrimaryPane {
				return m.handleRenameKey()
			}
		case "d":
			if !m.Rename && !m.Delete && !m.Move && m.ActivePane == constants.PrimaryPane {
				return m.handleDeleteKey()
			}
		case "tab":
			if m.ActivePane == constants.PrimaryPane {
				m.ActivePane = constants.SecondaryPane
			} else {
				m.ActivePane = constants.PrimaryPane
			}
		case "esc":
			if m.ActivePane == constants.PrimaryPane {
				return m.handleEscKey()
			}
		}
	}

	m.PrimaryViewport, cmd = m.PrimaryViewport.Update(msg)
	cmds = append(cmds, cmd)

	m.SecondaryViewport, cmd = m.SecondaryViewport.Update(msg)
	cmds = append(cmds, cmd)

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
