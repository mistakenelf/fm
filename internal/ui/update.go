package ui

import (
	"github.com/knipferrc/fm/internal/constants"
	"github.com/knipferrc/fm/internal/helpers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/reflow/wrap"
)

// Update handles all UI interactions and events for updating the screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// A directoryUpdateMsg returns an updated listing of files to display
	// in the UI. Any time an action is performed, this is called
	// for example, changing directories, or performing most
	// file operations.
	case updateDirectoryListingMsg:
		if len(msg) == 0 {
			m.primaryPane.SetContent("Directory is empty")
		} else {
			m.dirTree.GotoTop()
			m.dirTree.SetContent(msg)
			m.primaryPane.SetContent(m.dirTree.View())
			m.primaryPane.GotoTop()
		}

		m.showCommandBar = false
		m.textInput.Blur()
		m.textInput.Reset()

		return m, nil

	// A moveDirItemMsg is received any time a file or directory has been moved.
	case moveDirItemMsg:
		// Set active color back to default.
		m.primaryPane.SetActiveBorderColor(m.appConfig.Colors.Pane.ActiveBorderColor)

		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())

		// Set move mode back to false, set the initial moving directory to empty,
		// the item that was moving back to nil, and update the status bars content.
		m.inMoveMode = false
		m.initialMoveDirectory = ""
		m.itemToMove = nil

		return m, nil

	// A fileContentMsg is received anytime a file is read from returning its content
	// along with the markdown content to be rendered by glamour.
	case readFileContentMsg:
		if msg.code != "" {
			m.secondaryPane.GotoTop()
			m.secondaryPane.SetContent(helpers.ConvertTabsToSpaces(msg.code))
		} else if msg.markdown != "" {
			m.secondaryPane.GotoTop()
			m.secondaryPane.SetContent(helpers.ConvertTabsToSpaces(msg.markdown))
		} else if msg.image != nil {
			m.secondaryPane.GotoTop()
			m.asciiImage.SetImage(msg.image)
			m.asciiImage.SetContent(msg.asciiImage)
			m.secondaryPane.SetContent(m.asciiImage.View())
		} else {
			m.secondaryPane.GotoTop()
			m.secondaryPane.SetContent(helpers.ConvertTabsToSpaces(msg.rawContent))
		}

		return m, nil

	// convertImageMsg is received when an image is to be converted to ASCII.
	case convertImageToAsciiMsg:
		m.asciiImage.SetContent(string(msg))
		m.secondaryPane.SetContent(m.asciiImage.View())

		return m, nil

	// An errorMsg is received any time something in a command goes wrong
	// we receive that error and show it in the secondary pane with red text.
	case errorMsg:
		m.secondaryPane.SetContent(
			lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(constants.Red)).
				Render(string(msg)),
		)

		return m, nil

	// Any time the window is resized this is called, including when the app
	// is first started.
	case tea.WindowSizeMsg:
		if !m.ready {
			m.secondaryPane.SetContent(wrap.String(wordwrap.String(constants.IntroText, msg.Width/2), msg.Width/2))
			m.dirTree.SetSize(msg.Width / 2)
			m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.statusBar.SetSize(msg.Width, constants.StatusBarHeight)
			m.ready = true
		} else {
			m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.dirTree.SetSize(msg.Width / 2)
			m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.statusBar.SetSize(msg.Width, constants.StatusBarHeight)
		}

		if m.asciiImage.Image != nil {
			resizeImageCmd := m.redrawImage(m.secondaryPane.GetWidth()-2, m.secondaryPane.GetHeight())
			cmds = append(cmds, resizeImageCmd)
		}

		return m, tea.Batch(cmds...)

	// Any time a mouse event is received, we get this message.
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			// The command bar is not open and the primary pane is active
			// so scroll the dirtree up and update the primary panes content.
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so scroll its content up.
				m.secondaryPane.LineUp(3)
			}

			return m, nil

		case tea.MouseWheelDown:
			// Command bar is not shown and the primary pane is active
			// so scroll the dirtree down and update the primary panes content.
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so scroll its content down.
				m.secondaryPane.LineDown(3)
			}

			return m, nil
		}

	case tea.KeyMsg:
		// If gg is pressed.
		if msg.String() == "g" && m.previousKey.String() == "g" {
			// If the command bar is not shown and the primary pane is active,
			// reset the previous key, go to the top of the dirtree and pane.
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.previousKey = tea.KeyMsg{}
				m.dirTree.GotoTop()
				m.primaryPane.GotoTop()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so go to the top.
				m.secondaryPane.GotoTop()
			}

			return m, nil
		}

		switch msg.String() {
		// Exit FM.
		case "ctrl+c":
			return m, tea.Quit

		// Exit FM if the command bar is not open.
		case "q":
			if !m.showCommandBar {
				return m, tea.Quit
			}

		case "left", "h":
			// If the command bar is not shown and the primary pane is active
			// set the previous directory to the current directory,
			// and update the directory listing to go back one directory.
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.previousDirectory, _ = helpers.GetWorkingDirectory()

				return m, m.updateDirectoryListing(constants.PreviousDirectory)
			}

		case "down", "j":
			// Scroll down in pane.
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoDown()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(m.dirTree.View())
				} else {
					m.secondaryPane.LineDown(1)
				}
			}

		case "up", "k":
			// Scroll up in pane.
			if !m.showCommandBar {
				if m.primaryPane.IsActive {
					m.dirTree.GoUp()
					m.scrollPrimaryPane()
					m.primaryPane.SetContent(m.dirTree.View())
				} else {
					m.secondaryPane.LineUp(1)
				}
			}

		case "right", "l":
			// Open directory or read file content.
			if !m.showCommandBar && m.primaryPane.IsActive {
				if m.dirTree.GetSelectedFile().IsDir() && !m.textInput.Focused() {
					return m, m.updateDirectoryListing(m.dirTree.GetSelectedFile().Name())
				}

				return m, m.readFileContent(m.dirTree.GetSelectedFile(), m.secondaryPane.GetWidth()-2, m.secondaryPane.GetHeight())

			}

		case "G":
			// Go to the bottom of the pane.
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GotoBottom()
				m.primaryPane.GotoBottom()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so go to the bottom of it.
				m.secondaryPane.GotoBottom()
			}

		case "enter":
			// If pressing enter while in move mode.
			if m.inMoveMode {
				if m.itemToMove.IsDir() {
					return m, tea.Sequentially(
						m.moveDir(m.itemToMove.Name()),
						m.updateDirectoryListing(m.initialMoveDirectory),
					)
				}

				return m, tea.Sequentially(
					m.moveFile(m.itemToMove.Name()),
					m.updateDirectoryListing(m.initialMoveDirectory),
				)

			} else {
				// Parse the commands from the command bar, command is the name
				// of the command and value is if the command requires input to it
				// get its value, for example (rename test.txt) text.txt is the value.
				command, value := helpers.ParseCommand(m.textInput.Value())

				// Nothing was input for a command.
				if command == "" {
					return m, nil
				}

				switch command {
				// Exit FM.
				case "exit", "q", "quit":
					return m, tea.Quit

				// Create a new directory based on the value passed.
				case "mkdir":
					return m, tea.Sequentially(
						m.createDir(value),
						m.updateDirectoryListing(constants.CurrentDirectory),
					)

				// Create a new file based on the value passed.
				case "touch":
					return m, tea.Sequentially(
						m.createFile(value),
						m.updateDirectoryListing(constants.CurrentDirectory),
					)

				// Rename the currently selected file or folder based on the value passed.
				case "mv", "rename":
					return m, tea.Sequentially(
						m.renameFileOrDir(m.dirTree.GetSelectedFile().Name(), value),
						m.updateDirectoryListing(constants.CurrentDirectory),
					)

				// Delete the currently selected item.
				case "rm", "delete":
					if m.dirTree.GetSelectedFile().IsDir() {
						return m, tea.Sequentially(
							m.deleteDir(m.dirTree.GetSelectedFile().Name()),
							m.updateDirectoryListing(constants.CurrentDirectory),
						)
					}

					return m, tea.Sequentially(
						m.deleteFile(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListing(constants.CurrentDirectory),
					)

				default:
					return m, nil
				}
			}

		case ":":
			// If move mode is not active, activate the command bar.
			if !m.inMoveMode {
				m.showCommandBar = true
				m.textInput.Placeholder = "enter command"
				m.textInput.Focus()
			}

			return m, cmd

		// Shortcut to get back to the home directory if the
		// command bar is not curently open.
		case "~":
			if !m.showCommandBar {
				homeDir, _ := helpers.GetHomeDirectory()
				return m, m.updateDirectoryListing(homeDir)
			}

		// Shortcut to go back to the previous directory.
		case "-":
			if !m.showCommandBar && m.previousDirectory != "" {
				return m, m.updateDirectoryListing(m.previousDirectory)
			}

		// Toggle hidden files and folders.
		case ".":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.ToggleHidden()
				return m, m.updateDirectoryListing(constants.CurrentDirectory)
			}

		// Toggle between the two panes if the command bar is not currently active.
		case "tab":
			if !m.showCommandBar {
				m.primaryPane.IsActive = !m.primaryPane.IsActive
				m.secondaryPane.IsActive = !m.secondaryPane.IsActive
			}

		// Enter move mode.
		case "m":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.inMoveMode = true
				m.primaryPane.SetActiveBorderColor(constants.Blue)
				m.initialMoveDirectory, _ = helpers.GetWorkingDirectory()
				m.itemToMove = m.dirTree.GetSelectedFile()
			}

		// Zip up the currently selected item.
		case "z":
			if !m.showCommandBar && m.primaryPane.IsActive {
				return m, tea.Sequentially(
					m.zipDirectory(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(constants.CurrentDirectory),
				)
			}

		// Unzip the currently selected zip file.
		case "u":
			if !m.showCommandBar && m.primaryPane.IsActive {
				return m, tea.Sequentially(
					m.unzipDirectory(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(constants.CurrentDirectory),
				)
			}

		// Copy the currently selected item.
		case "c":
			if !m.showCommandBar && m.primaryPane.IsActive {
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, tea.Sequentially(
						m.copyDirectory(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListing(constants.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.copyFile(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(constants.CurrentDirectory),
				)
			}

		// Reset FM to its initial state.
		case "esc":
			m.showCommandBar = false
			m.inMoveMode = false
			m.itemToMove = nil
			m.initialMoveDirectory = ""
			m.primaryPane.IsActive = true
			m.secondaryPane.IsActive = false
			m.textInput.Blur()
			m.textInput.Reset()
			m.secondaryPane.GotoTop()
			m.primaryPane.SetActiveBorderColor(m.appConfig.Colors.Pane.ActiveBorderColor)
			m.secondaryPane.SetContent(helpers.ConvertTabsToSpaces(constants.IntroText))
		}

		// Capture the previous key so that we can capture
		// when two keys are pressed.
		m.previousKey = msg
	}

	// Keep status bar content updated.
	m.statusBar.SetContent(m.getStatusBarContent())

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
