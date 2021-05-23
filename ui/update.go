package ui

import (
	"log"
	"os"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/dirtree"
	"github.com/knipferrc/fm/pane"
	"github.com/knipferrc/fm/statusbar"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) scrollPrimaryPane() {
	top := m.PrimaryPane.Viewport.YOffset
	bottom := m.PrimaryPane.Height + m.PrimaryPane.Viewport.YOffset - 1

	if m.Cursor < top {
		m.PrimaryPane.LineUp(1)
	} else if m.Cursor > bottom {
		m.PrimaryPane.LineDown(1)
	}

	if m.Cursor > len(m.Files)-1 {
		m.Cursor = 0
		m.PrimaryPane.GotoTop()
	} else if m.Cursor < 0 {
		m.Cursor = len(m.Files) - 1
		m.PrimaryPane.GotoBottom()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case directoryMsg:
		m.Files = msg
		m.Cursor = 0
		m.DirTree.SetContent(m.Files, m.Cursor)
		m.PrimaryPane.SetContent(m.DirTree.View())
		m.SecondaryPane.SetContent(m.Text.View())

	case fileContentMsg:
		cfg := config.GetConfig()
		border := lipgloss.NormalBorder()

		if cfg.Settings.RoundedPanes {
			border = lipgloss.RoundedBorder()
		}

		halfScreenWidth := m.ScreenWidth / 2
		borderWidth := lipgloss.Width(border.Left + border.Right)
		m.SecondaryPane.SetContent(lipgloss.NewStyle().Width(halfScreenWidth - borderWidth).Render(utils.ConverTabsToSpaces(string(msg))))

	case tea.WindowSizeMsg:
		cfg := config.GetConfig()
		border := lipgloss.NormalBorder()

		if cfg.Settings.RoundedPanes {
			border = lipgloss.RoundedBorder()
		}

		paneBorderWidth := lipgloss.Width(border.Left + border.Right)
		verticalMargin := lipgloss.Width(border.Bottom) + constants.StatusBarHeight

		if !m.Ready {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height

			m.PrimaryPane = pane.NewModel(
				(msg.Width/2)-paneBorderWidth,
				msg.Height-verticalMargin,
				true,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)
			m.DirTree = dirtree.NewModel(m.Files, m.Cursor)
			m.PrimaryPane.SetContent(m.DirTree.View())

			m.SecondaryPane = pane.NewModel(
				(msg.Width/2)-paneBorderWidth,
				msg.Height-verticalMargin,
				false,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)
			m.SecondaryPane.SetContent(m.Text.View())

			m.StatusBar = statusbar.NewModel(msg.Width, m.Cursor, len(m.Files), m.Files[m.Cursor], m.ShowCommandBar, m.Textinput.View())

			m.Ready = true
		} else {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height
			m.PrimaryPane.SetSize((msg.Width/2)-paneBorderWidth, msg.Height-verticalMargin)
			m.SecondaryPane.SetSize((msg.Width/2)-paneBorderWidth, msg.Height-verticalMargin)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if !m.ShowCommandBar {
				return m, tea.Quit
			}

		case "left", "h":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					previousPath, err := os.Getwd()

					if err != nil {
						log.Fatal("error getting working directory")
					}

					m.PreviousDirectory = previousPath
					return m, updateDirectoryListing(constants.PreviousDirectory)
				}
			}

		case "down", "j":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					m.Cursor++
					m.scrollPrimaryPane()
					m.DirTree.SetContent(m.Files, m.Cursor)
					m.PrimaryPane.SetContent(m.DirTree.View())
				} else {
					m.SecondaryPane.LineDown(1)
				}
			}

		case "up", "k":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					m.Cursor--
					m.scrollPrimaryPane()
					m.DirTree.SetContent(m.Files, m.Cursor)
					m.PrimaryPane.SetContent(m.DirTree.View())
				} else {
					m.SecondaryPane.LineUp(1)
				}
			}

		case "right", "l":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					if m.Files[m.Cursor].IsDir() && !m.Textinput.Focused() {
						return m, updateDirectoryListing(m.Files[m.Cursor].Name())
					} else {
						m.SecondaryPane.GotoTop()
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

		case "~":
			if !m.ShowCommandBar {
				home, err := os.UserHomeDir()
				if err != nil {
					log.Fatal(err)
				}

				return m, updateDirectoryListing(home)
			}

		case "-":
			if !m.ShowCommandBar && m.PreviousDirectory != "" {
				return m, updateDirectoryListing(m.PreviousDirectory)
			}

		case "tab":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					m.PrimaryPane.IsActive = false
					m.SecondaryPane.IsActive = true
				} else {
					m.PrimaryPane.IsActive = true
					m.SecondaryPane.IsActive = false
				}
			}

		case "esc":
			m.ShowCommandBar = false
			m.Textinput.Blur()
			m.Textinput.Reset()
			m.SecondaryPane.GotoTop()
			m.DirTree.SetContent(m.Files, m.Cursor)
			m.PrimaryPane.SetContent(m.DirTree.View())
			m.SecondaryPane.SetContent(m.Text.View())
		}
	}

	m.StatusBar.Update(m.ScreenWidth, m.Cursor, len(m.Files), m.Files[m.Cursor], m.ShowCommandBar, m.Textinput.View())

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
