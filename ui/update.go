package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/pane"
	"github.com/knipferrc/fm/statusbar"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) scrollPrimaryPane() {
	top := m.PrimaryPane.Viewport.YOffset
	bottom := m.PrimaryPane.Height + m.PrimaryPane.YOffset - 1

	if m.DirTree.GetCursor() < top {
		m.PrimaryPane.LineUp(1)
	} else if m.DirTree.GetCursor() > bottom {
		m.PrimaryPane.LineDown(1)
	}

	if m.DirTree.GetCursor() > m.DirTree.GetTotalFiles()-1 {
		m.DirTree.GotoTop()
		m.PrimaryPane.GotoTop()
	} else if m.DirTree.GetCursor() < top {
		m.DirTree.GotoBottom()
		m.PrimaryPane.GotoBottom()
	}
}

func (m Model) getStatusBarContent() (string, string, string, string) {
	cfg := config.GetConfig()
	currentPath, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	logo := ""
	if cfg.Settings.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	status := fmt.Sprintf("%s %s %s",
		utils.ConvertBytesToSizeString(m.DirTree.GetSelectedFile().Size()),
		m.DirTree.GetSelectedFile().Mode().String(),
		currentPath,
	)

	if m.ShowCommandBar {
		status = m.Textinput.View()
	}

	return m.DirTree.GetSelectedFile().Name(), status, fmt.Sprintf("%d/%d", m.DirTree.GetCursor()+1, m.DirTree.GetTotalFiles()), logo
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case directoryMsg:
		m.DirTree.SetContent(msg)
		m.DirTree.GotoTop()
		m.PrimaryPane.SetContent(m.DirTree.View())
		m.ShowCommandBar = false
		m.Textinput.Blur()
		m.Textinput.Reset()
		selectedFile, status, fileTotals, logo := m.getStatusBarContent()
		m.StatusBar.SetContent(selectedFile, status, fileTotals, logo)

	case fileContentMsg:
		m.SecondaryPane.SetContent(utils.ConverTabsToSpaces(string(msg)))

	case tea.WindowSizeMsg:
		cfg := config.GetConfig()

		if !m.Ready {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height

			m.PrimaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
				true,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)
			m.PrimaryPane.SetContent(m.DirTree.View())

			m.SecondaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
				false,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)
			m.SecondaryPane.SetContent(m.HelpText.View())

			selectedFile, status, fileTotals, logo := m.getStatusBarContent()
			m.StatusBar = statusbar.NewModel(
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

			m.StatusBar.SetContent(selectedFile, status, fileTotals, logo)

			m.Ready = true
		} else {
			m.ScreenWidth = msg.Width
			m.ScreenHeight = msg.Height
			m.PrimaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.SecondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.StatusBar.SetSize(msg.Width)
		}

	case tea.KeyMsg:
		if msg.String() == "g" && m.LastKey.String() == "g" {
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					m.DirTree.GotoTop()
					m.PrimaryPane.GotoTop()
					m.PrimaryPane.SetContent(m.DirTree.View())
				} else {
					m.SecondaryPane.GotoTop()
				}
			}
		}

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
					return m, updateDirectoryListing(constants.PreviousDirectory, m.DirTree.ShowHidden)
				}
			}

		case "down", "j":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					m.DirTree.GoDown()
					m.scrollPrimaryPane()
					selectedFile, status, fileTotals, logo := m.getStatusBarContent()
					m.StatusBar.SetContent(selectedFile, status, fileTotals, logo)
					m.PrimaryPane.SetContent(m.DirTree.View())
				} else {
					m.SecondaryPane.LineDown(1)
				}
			}

		case "up", "k":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					m.DirTree.GoUp()
					m.scrollPrimaryPane()
					m.PrimaryPane.SetContent(m.DirTree.View())
					selectedFile, status, fileTotals, logo := m.getStatusBarContent()
					m.StatusBar.SetContent(selectedFile, status, fileTotals, logo)
				} else {
					m.SecondaryPane.LineUp(1)
				}
			}

		case "G":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					m.DirTree.GotoBottom()
					m.PrimaryPane.GotoBottom()
					m.PrimaryPane.SetContent(m.DirTree.View())
				} else {
					m.SecondaryPane.GotoBottom()
				}
			}

		case "right", "l":
			if !m.ShowCommandBar {
				if m.PrimaryPane.IsActive {
					if m.DirTree.GetSelectedFile().IsDir() && !m.Textinput.Focused() {
						return m, updateDirectoryListing(m.DirTree.GetSelectedFile().Name(), m.DirTree.ShowHidden)
					} else {
						m.SecondaryPane.GotoTop()
						return m, readFileContent(m.DirTree.GetSelectedFile().Name())
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
				return m, createDir(value, m.DirTree.ShowHidden)

			case "touch":
				return m, createFile(value, m.DirTree.ShowHidden)

			case "mv":
				return m, renameFileOrDir(m.DirTree.GetSelectedFile().Name(), value, m.DirTree.ShowHidden)

			case "cp":
				if m.DirTree.GetSelectedFile().IsDir() {
					return m, moveDir(m.DirTree.GetSelectedFile().Name(), value, m.DirTree.ShowHidden)
				} else {
					return m, moveFile(m.DirTree.GetSelectedFile().Name(), value, m.DirTree.ShowHidden)
				}

			case "rm":
				if m.DirTree.GetSelectedFile().IsDir() {
					return m, deleteDir(m.DirTree.GetSelectedFile().Name(), m.DirTree.ShowHidden)
				} else {
					return m, deleteFile(m.DirTree.GetSelectedFile().Name(), m.DirTree.ShowHidden)
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
				return m, updateDirectoryListing(utils.GetHomeDirectory(), m.DirTree.ShowHidden)
			}

		case "-":
			if !m.ShowCommandBar && m.PreviousDirectory != "" {
				return m, updateDirectoryListing(m.PreviousDirectory, m.DirTree.ShowHidden)
			}

		case ".":
			if !m.ShowCommandBar && m.PrimaryPane.IsActive {
				m.DirTree.ToggleHidden()

				return m, updateDirectoryListing(constants.CurrentDirectory, m.DirTree.ShowHidden)
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
			m.SecondaryPane.SetContent(m.HelpText.View())
			m.PrimaryPane.IsActive = true
			m.SecondaryPane.IsActive = false
			selectedFile, status, fileTotals, logo := m.getStatusBarContent()
			m.StatusBar.SetContent(selectedFile, status, fileTotals, logo)
		}

		m.LastKey = msg
	}

	if m.ShowCommandBar {
		selectedFile, status, fileTotals, logo := m.getStatusBarContent()
		m.StatusBar.SetContent(selectedFile, status, fileTotals, logo)
	}

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
