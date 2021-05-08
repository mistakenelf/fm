package ui

import (
	"github.com/knipferrc/fm/internal/components"
	"github.com/knipferrc/fm/internal/constants"
	"github.com/knipferrc/fm/internal/utils"

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
		m.ShowCommandBar = false
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
			if !m.ShowCommandBar {
				return m, tea.Quit
			}

		case "left", "h":
			if !m.ShowCommandBar {
				return m, updateDirectoryListing(constants.PreviousDirectory)
			}

		case "down", "j":
			if !m.ShowCommandBar {
				if m.ActivePane == constants.PrimaryPane {
					m.Cursor++
					m.scrollPrimaryViewport()
					m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
				} else {
					m.SecondaryViewport.LineDown(1)
				}
			}

		case "up", "k":
			if !m.ShowCommandBar {
				if m.ActivePane == constants.PrimaryPane {
					m.Cursor--
					m.scrollPrimaryViewport()
					m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
				} else {
					m.SecondaryViewport.LineUp(1)
				}
			}

		case "right", "l":
			if !m.ShowCommandBar {
				if m.ActivePane == constants.PrimaryPane {
					if m.Files[m.Cursor].IsDir() && !m.Textinput.Focused() {
						return m, updateDirectoryListing(m.Files[m.Cursor].Name())
					} else {
						return m, readFileContent(m.Files[m.Cursor].Name())
					}
				}
			}

		case "enter":
			cmd, value := utils.ParseCommand(m.Textinput.Value())

			if cmd == "" {
				return m, nil
			}

			switch cmd {
			case "mkdir":
				return m, createDir(value)

			case "touch":
				return m, createFile(value)

			case "mv":
				return m, renameFileOrDir(m.Files[m.Cursor].Name(), value)

			case "cp":
				if m.Files[m.Cursor].IsDir() {
					return m, moveDir(m.Files[m.Cursor].Name(), value)
				} else {
					return m, moveFile(m.Files[m.Cursor].Name(), value)
				}

			case "rm":
				if m.Files[m.Cursor].IsDir() {
					return m, deleteDir(m.Files[m.Cursor].Name())
				} else {
					return m, deleteFile(m.Files[m.Cursor].Name())
				}

			default:
				return m, nil

			}

		case ":":
			m.ShowCommandBar = true
			m.Textinput.Placeholder = "enter command"
			m.Textinput.Focus()

			return m, nil

		case "tab":
			if !m.ShowCommandBar {
				if m.ActivePane == constants.PrimaryPane {
					m.ActivePane = constants.SecondaryPane
				} else {
					m.ActivePane = constants.PrimaryPane
				}
			}

		case "esc":
			m.ShowCommandBar = false
			m.ActivePane = constants.PrimaryPane
			m.Textinput.Blur()
			m.Textinput.Reset()
			m.SecondaryViewport.GotoTop()
			m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
			m.SecondaryViewport.SetContent(components.Instructions())
		}
	}

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
