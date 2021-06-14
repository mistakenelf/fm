package ui

import (
	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/pane"
	"github.com/knipferrc/fm/statusbar"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case directoryMsg:
		if len(msg) == 0 {
			m.primaryPane.SetContent("Directory is empty")
		} else {
			m.dirTree.SetContent(msg)
			m.dirTree.GotoTop()
			m.primaryPane.SetContent(m.dirTree.View())
		}

		m.showCommandBar = false
		m.textInput.Blur()
		m.textInput.Reset()
		m.statusBar.SetContent(m.getStatusBarContent())

		return m, cmd

	case moveFileMsg:
		cfg := config.GetConfig()

		m.primaryPane.SetActiveBorderColor(cfg.Colors.Pane.ActiveBorderColor)
		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())
		m.inMoveMode = false
		m.initialMoveDirectory = ""
		m.itemToMove = nil
		m.statusBar.SetContent(m.getStatusBarContent())

		return m, cmd

	case fileContentMsg:
		m.activeMarkdownSource = string(msg.markdownContent)
		m.secondaryPane.SetContent(utils.ConverTabsToSpaces(string(msg.fileContent)))

		return m, cmd

	case markdownMsg:
		m.secondaryPane.SetContent(utils.ConverTabsToSpaces(string(msg)))

		return m, cmd

	case errorMsg:
		m.secondaryPane.SetContent(
			lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(constants.Red)).
				Render(string(msg)),
		)

		return m, cmd

	case tea.WindowSizeMsg:
		cfg := config.GetConfig()

		if !m.ready {
			m.primaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
				true,
				cfg.Settings.RoundedPanes,
				true,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)
			m.primaryPane.SetContent(m.dirTree.View())

			m.secondaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
				false,
				false,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)

			selectedFile, status, fileTotals, logo := m.getStatusBarContent()
			m.statusBar = statusbar.NewModel(
				msg.Width,
				selectedFile,
				status,
				fileTotals,
				logo,
				statusbar.Color{
					Background: cfg.Colors.StatusBar.SelectedFile.Background,
					Foreground: cfg.Colors.StatusBar.SelectedFile.Foreground,
				},
				statusbar.Color{
					Background: cfg.Colors.StatusBar.Bar.Background,
					Foreground: cfg.Colors.StatusBar.Bar.Foreground,
				},
				statusbar.Color{
					Background: cfg.Colors.StatusBar.TotalFiles.Background,
					Foreground: cfg.Colors.StatusBar.TotalFiles.Foreground,
				},
				statusbar.Color{
					Background: cfg.Colors.StatusBar.Logo.Background,
					Foreground: cfg.Colors.StatusBar.Logo.Foreground,
				},
			)

			m.ready = true
		} else {
			m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.statusBar.SetSize(msg.Width)
		}

		if m.activeMarkdownSource != "" {
			return m, renderMarkdownContent(m.secondaryPane.Width, m.activeMarkdownSource)
		}

		return m, cmd

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				m.secondaryPane.LineUp(3)
			}

			return m, cmd

		case tea.MouseWheelDown:
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				m.secondaryPane.LineDown(3)
			}

			return m, cmd
		}

	case tea.KeyMsg:
		if msg.String() == "g" && m.previousKey.String() == "g" {
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.previousKey = tea.KeyMsg{}
				m.dirTree.GotoTop()
				m.primaryPane.GotoTop()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				m.secondaryPane.GotoTop()
			}

			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if !m.showCommandBar {
				return m, tea.Quit
			}

		case "left", "h":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.previousDirectory, _ = utils.GetWorkingDirectory()

				return m, m.updateDirectoryListing(constants.PreviousDirectory)
			}

		case "down", "j":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.statusBar.SetContent(m.getStatusBarContent())
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				m.secondaryPane.LineDown(1)
			}

		case "up", "k":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetContent(m.getStatusBarContent())
			} else {
				m.secondaryPane.LineUp(1)
			}

		case "right", "l":
			if !m.showCommandBar && m.primaryPane.IsActive {
				if m.dirTree.GetSelectedFile().IsDir() && !m.textInput.Focused() {
					return m, m.updateDirectoryListing(m.dirTree.GetSelectedFile().Name())
				} else {
					m.secondaryPane.GotoTop()
					return m, m.readFileContent(m.dirTree.GetSelectedFile())
				}
			}

		case "G":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GotoBottom()
				m.primaryPane.GotoBottom()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				m.secondaryPane.GotoBottom()
			}

		case "enter":
			if m.inMoveMode {
				if m.itemToMove.IsDir() {
					return m, m.moveDir(m.itemToMove.Name())
				} else {
					return m, m.moveFile(m.itemToMove.Name())
				}
			} else {
				command, value := utils.ParseCommand(m.textInput.Value())

				if command == "" {
					return m, nil
				}

				switch command {
				case "mkdir":
					return m, m.createDir(value)

				case "touch":
					return m, m.createFile(value)

				case "mv", "rename":
					return m, m.renameFileOrDir(m.dirTree.GetSelectedFile().Name(), value)

				case "rm", "delete":
					if m.dirTree.GetSelectedFile().IsDir() {
						return m, m.deleteDir(m.dirTree.GetSelectedFile().Name())
					} else {
						return m, m.deleteFile(m.dirTree.GetSelectedFile().Name())
					}

				default:
					return m, nil
				}
			}

		case ":":
			if !m.inMoveMode {
				m.showCommandBar = true
				m.textInput.Placeholder = "enter command"
				m.textInput.Focus()
			}

			return m, cmd

		case "~":
			if !m.showCommandBar {
				homeDir, _ := utils.GetHomeDirectory()
				return m, m.updateDirectoryListing(homeDir)
			}

		case "-":
			if !m.showCommandBar && m.previousDirectory != "" {
				return m, m.updateDirectoryListing(m.previousDirectory)
			}

		case ".":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.ToggleHidden()
				return m, m.updateDirectoryListing(constants.CurrentDirectory)
			}

		case "tab":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.primaryPane.IsActive = false
					m.secondaryPane.IsActive = true
				} else {
					m.primaryPane.IsActive = true
					m.secondaryPane.IsActive = false
				}
			}

		case "m":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.inMoveMode = true
				m.primaryPane.SetActiveBorderColor(constants.Blue)
				m.initialMoveDirectory, _ = utils.GetWorkingDirectory()
				m.itemToMove = m.dirTree.GetSelectedFile()
				m.statusBar.SetContent(m.getStatusBarContent())
			}

		case "esc":
			cfg := config.GetConfig()

			m.showCommandBar = false
			m.inMoveMode = false
			m.itemToMove = nil
			m.initialMoveDirectory = ""
			m.textInput.Blur()
			m.textInput.Reset()
			m.secondaryPane.GotoTop()
			m.primaryPane.IsActive = true
			m.secondaryPane.IsActive = false
			m.primaryPane.SetActiveBorderColor(cfg.Colors.Pane.ActiveBorderColor)
			m.statusBar.SetContent(m.getStatusBarContent())

			return m, renderMarkdownContent(m.secondaryPane.Width, constants.HelpText)
		}

		m.previousKey = msg
	}

	if m.showCommandBar {
		m.statusBar.SetContent(m.getStatusBarContent())
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
