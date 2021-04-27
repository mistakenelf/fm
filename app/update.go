package app

import (
	"github.com/knipferrc/fm/components"
	"github.com/knipferrc/fm/constants"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"
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

func (m Model) notPerformingAction() bool {
	return !m.ShowRenamePrompt && !m.ShowDeletePrompt && !m.ShowMovePrompt
}

func (m Model) handleKeyDown() (tea.Model, tea.Cmd) {
	if m.ActivePane == constants.PrimaryPane {
		m.Cursor++
		m.scrollPrimaryViewport()
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	} else {
		m.SecondaryViewport.LineDown(1)
	}

	return m, nil
}

func (m Model) handleKeyUp() (tea.Model, tea.Cmd) {
	if m.ActivePane == constants.PrimaryPane {
		m.Cursor--
		m.scrollPrimaryViewport()
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	} else {
		m.SecondaryViewport.LineUp(1)
	}

	return m, nil
}

func (m Model) handleRightKey() (tea.Model, tea.Cmd) {
	if m.ActivePane == constants.PrimaryPane {
		if m.Files[m.Cursor].IsDir() && !m.Textinput.Focused() {
			return m, updateDirectoryListing(m.Files[m.Cursor].Name())
		} else {
			return m, readFileContent(m.Files[m.Cursor].Name())
		}
	}

	return m, nil
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.ShowRenamePrompt {
		return m, renameFileOrDir(m.Files[m.Cursor].Name(), m.Textinput.Value())
	} else if m.ShowMovePrompt {
		if m.Files[m.Cursor].IsDir() {
			return m, moveDir(m.Files[m.Cursor].Name(), m.Textinput.Value())
		} else {
			return m, moveFile(m.Files[m.Cursor].Name(), m.Textinput.Value())
		}
	} else if m.ShowDeletePrompt {
		if m.Files[m.Cursor].IsDir() {
			if m.Textinput.Value() == "y" {
				return m, deleteDir(m.Files[m.Cursor].Name())
			} else {
				m.Textinput.Blur()
				m.Textinput.Reset()
				m.ShowDeletePrompt = false
			}
		} else {
			if m.Textinput.Value() == "y" {
				return m, deleteFile(m.Files[m.Cursor].Name())
			} else {
				m.Textinput.Blur()
				m.Textinput.Reset()
				m.ShowDeletePrompt = false
			}
		}
	} else {
		return m, nil
	}

	return m, nil
}

func (m Model) handleMoveKey() (tea.Model, tea.Cmd) {
	if m.ActivePane == constants.PrimaryPane {
		m.ActivePane = constants.SecondaryPane
		m.ShowMovePrompt = true
		m.Textinput.Placeholder = "/new/dir/name"
		m.Textinput.Focus()

	}

	return m, nil
}

func (m Model) handleRenameKey() (tea.Model, tea.Cmd) {
	if m.ActivePane == constants.PrimaryPane {
		m.ActivePane = constants.SecondaryPane
		m.ShowRenamePrompt = true
		m.Textinput.Placeholder = "new_name"
		m.Textinput.Focus()
	}

	return m, nil
}

func (m Model) handleDeleteKey() (tea.Model, tea.Cmd) {
	if m.ActivePane == constants.PrimaryPane {
		m.ActivePane = constants.SecondaryPane
		m.ShowDeletePrompt = true
		m.Textinput.Placeholder = "[y/n]"
		m.Textinput.Focus()
	}

	return m, nil
}

func (m Model) handleTabKey() (tea.Model, tea.Cmd) {
	if m.ActivePane == constants.PrimaryPane {
		m.ActivePane = constants.SecondaryPane
	} else {
		m.ActivePane = constants.PrimaryPane
	}

	return m, nil
}

func (m Model) handleEscKey() (tea.Model, tea.Cmd) {
	m.ShowMovePrompt = false
	m.ShowRenamePrompt = false
	m.ShowDeletePrompt = false
	m.ActivePane = constants.PrimaryPane
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

	case actionMsg:
		m.Files = msg
		m.Cursor = 0
		m.ShowRenamePrompt = false
		m.ShowMovePrompt = false
		m.ShowDeletePrompt = false
		m.ActivePane = constants.PrimaryPane
		m.Textinput.Blur()
		m.Textinput.Reset()
		m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
		m.SecondaryViewport.SetContent(components.Instructions())

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
			if m.notPerformingAction() {
				return m, tea.Quit
			}

		case "h":
			if m.ActivePane == constants.PrimaryPane {
				return m, updateDirectoryListing(constants.PreviousDirectory)
			}

		case "down", "j":
			if m.notPerformingAction() {
				return m.handleKeyDown()
			}

		case "up", "k":
			if m.notPerformingAction() {
				return m.handleKeyUp()
			}

		case "l":
			return m.handleRightKey()

		case "enter":
			return m.handleEnterKey()

		case "m":
			if m.notPerformingAction() {
				return m.handleMoveKey()
			}

		case "r":
			if m.notPerformingAction() {
				return m.handleRenameKey()
			}

		case "d":
			if m.notPerformingAction() {
				return m.handleDeleteKey()
			}

		case "tab":
			return m.handleTabKey()

		case "esc":
			return m.handleEscKey()
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
