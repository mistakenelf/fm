package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/pane"
	"github.com/knipferrc/fm/statusbar"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) scrollPrimaryPane() {
	top := m.primaryPane.Viewport.YOffset
	bottom := m.primaryPane.Height + m.primaryPane.YOffset - 1

	if m.dirTree.GetCursor() < top {
		m.primaryPane.LineUp(1)
	} else if m.dirTree.GetCursor() > bottom {
		m.primaryPane.LineDown(1)
	}

	if m.dirTree.GetCursor() > m.dirTree.GetTotalFiles()-1 {
		m.dirTree.GotoTop()
		m.primaryPane.GotoTop()
	} else if m.dirTree.GetCursor() < top {
		m.dirTree.GotoBottom()
		m.primaryPane.GotoBottom()
	}
}

func (m model) getStatusBarContent() (string, string, string, string) {
	cfg := config.GetConfig()
	currentPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	logo := ""
	if cfg.Settings.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	status := fmt.Sprintf("%s %s %s",
		utils.ConvertBytesToSizeString(m.dirTree.GetSelectedFile().Size()),
		m.dirTree.GetSelectedFile().Mode().String(),
		currentPath,
	)

	if m.showCommandBar {
		status = m.textInput.View()
	}

	return m.dirTree.GetSelectedFile().Name(), status, fmt.Sprintf("%d/%d", m.dirTree.GetCursor()+1, m.dirTree.GetTotalFiles()), logo
}

func (m model) renderMarkdown(str string) string {
	bg := "light"

	if lipgloss.HasDarkBackground() {
		bg = "dark"
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(m.secondaryPane.Width),
		glamour.WithStandardStyle(bg),
	)

	out, err := r.Render(str)
	if err != nil {
		// FIXME: show an error in the UI
		log.Fatal(err)
	}

	return out
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case directoryMsg:
		m.dirTree.SetContent(msg)
		m.dirTree.GotoTop()
		m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
		m.showCommandBar = false
		m.textInput.Blur()
		m.textInput.Reset()
		selectedFile, status, fileTotals, logo := m.getStatusBarContent()
		m.statusBar.SetContent(selectedFile, status, fileTotals, logo)

		return m, cmd

	case fileContentMsg:
		cfg := config.GetConfig()
		content := string(msg)

		if filepath.Ext(m.dirTree.GetSelectedFile().Name()) == ".md" && cfg.Settings.PrettyMarkdown {
			m.activeMarkdownSource = string(msg)
			content = m.renderMarkdown(m.activeMarkdownSource)
		} else {
			m.activeMarkdownSource = ""
		}

		m.secondaryPane.SetContent(utils.ConverTabsToSpaces(content))

		return m, cmd

	case tea.WindowSizeMsg:
		cfg := config.GetConfig()

		if !m.ready {
			m.screenWidth = msg.Width
			m.screenHeight = msg.Height

			m.primaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
				true,
				cfg.Settings.RoundedPanes,
				cfg.Colors.Pane.ActiveBorderColor,
				cfg.Colors.Pane.InactiveBorderColor,
			)
			m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))

			m.secondaryPane = pane.NewModel(
				msg.Width/2,
				msg.Height-constants.StatusBarHeight,
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

			m.statusBar.SetContent(selectedFile, status, fileTotals, logo)

			m.ready = true
		} else {
			m.screenHeight = msg.Width
			m.screenHeight = msg.Height
			m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.statusBar.SetSize(msg.Width)
		}

		if m.activeMarkdownSource != "" {
			m.secondaryPane.SetContent(m.renderMarkdown(m.activeMarkdownSource))
		}

		return m, cmd

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoUp()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.LineUp(3)
				}
			}

			return m, cmd

		case tea.MouseWheelDown:
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoDown()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.LineDown(3)
				}
			}

			return m, cmd
		}

	case tea.KeyMsg:
		if msg.String() == "g" && m.previousKey.String() == "g" {
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GotoTop()
					m.primaryPane.GotoTop()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.GotoTop()
				}
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
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					previousPath, err := os.Getwd()

					if err != nil {
						log.Fatal("error getting working directory")
					}

					m.previousDirectory = previousPath

					return m, updateDirectoryListing(constants.PreviousDirectory, m.dirTree.ShowHidden)
				}
			}

			return m, cmd

		case "down", "j":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoDown()
					m.scrollPrimaryPane()
					selectedFile, status, fileTotals, logo := m.getStatusBarContent()
					m.statusBar.SetContent(selectedFile, status, fileTotals, logo)
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.LineDown(1)
				}
			}

			return m, cmd

		case "up", "k":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoUp()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
					selectedFile, status, fileTotals, logo := m.getStatusBarContent()
					m.statusBar.SetContent(selectedFile, status, fileTotals, logo)
				} else {
					m.secondaryPane.LineUp(1)
				}
			}

			return m, cmd

		case "G":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GotoBottom()
					m.primaryPane.GotoBottom()
					m.primaryPane.SetContent(lipgloss.NewStyle().PaddingLeft(1).Render(m.dirTree.View()))
				} else {
					m.secondaryPane.GotoBottom()
				}
			}

			return m, cmd

		case "right", "l":
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					if m.dirTree.GetSelectedFile().IsDir() && !m.textInput.Focused() {
						return m, updateDirectoryListing(m.dirTree.GetSelectedFile().Name(), m.dirTree.ShowHidden)
					} else {
						m.secondaryPane.GotoTop()

						return m, readFileContent(m.dirTree.GetSelectedFile().Name())
					}
				}
			}

			return m, cmd

		case "enter":
			command, value := utils.ParseCommand(m.textInput.Value())

			if command == "" {
				return m, nil
			}

			switch command {
			case "mkdir":
				return m, createDir(value, m.dirTree.ShowHidden)

			case "touch":
				return m, createFile(value, m.dirTree.ShowHidden)

			case "mv":
				return m, renameFileOrDir(m.dirTree.GetSelectedFile().Name(), value, m.dirTree.ShowHidden)

			case "cp":
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, moveDir(m.dirTree.GetSelectedFile().Name(), value, m.dirTree.ShowHidden)
				} else {
					return m, moveFile(m.dirTree.GetSelectedFile().Name(), value, m.dirTree.ShowHidden)
				}

			case "rm":
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, deleteDir(m.dirTree.GetSelectedFile().Name(), m.dirTree.ShowHidden)
				} else {
					return m, deleteFile(m.dirTree.GetSelectedFile().Name(), m.dirTree.ShowHidden)
				}

			default:
				return m, nil
			}

		case ":":
			m.showCommandBar = true
			m.textInput.Placeholder = "enter command"
			m.textInput.Focus()

			return m, cmd

		case "~":
			if !m.showCommandBar {
				return m, updateDirectoryListing(utils.GetHomeDirectory(), m.dirTree.ShowHidden)
			}

			return m, cmd

		case "-":
			if !m.showCommandBar && m.previousDirectory != "" {
				return m, updateDirectoryListing(m.previousDirectory, m.dirTree.ShowHidden)
			}

			return m, cmd

		case ".":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.ToggleHidden()

				return m, updateDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)
			}

			return m, cmd

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

			return m, cmd

		case "esc":
			m.showCommandBar = false
			m.textInput.Blur()
			m.textInput.Reset()
			m.secondaryPane.GotoTop()
			m.activeMarkdownSource = constants.HelpText
			m.secondaryPane.SetContent(m.renderMarkdown(m.activeMarkdownSource))
			m.primaryPane.IsActive = true
			m.secondaryPane.IsActive = false
			selectedFile, status, fileTotals, logo := m.getStatusBarContent()
			m.statusBar.SetContent(selectedFile, status, fileTotals, logo)

			return m, cmd
		}

		m.previousKey = msg
	}

	if m.showCommandBar {
		selectedFile, status, fileTotals, logo := m.getStatusBarContent()
		m.statusBar.SetContent(selectedFile, status, fileTotals, logo)
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
