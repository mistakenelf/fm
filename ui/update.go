package ui

import (
	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	// A directoryMsg returns an updated listing of files to display
	// in the UI. Any time an action is performed, this is called
	// for example, changing directories, or performing most
	// file operations
	case directoryMsg:
		// Directory is empty so lets display a message in the pane
		// to indicate nothing is in that directory
		if len(msg) == 0 {
			m.primaryPane.SetContent("Directory is empty")
		} else {
			// Update the dirtree with new content, scroll to the top
			// of the dirtree and set the primary panes content to the dirtree
			m.dirTree.SetContent(msg)
			m.dirTree.GotoTop()
			m.primaryPane.SetContent(m.dirTree.View())
		}

		// Hide the command bar, reset the textinput, blur its focus
		// and update the status bars content
		m.showCommandBar = false
		m.textInput.Blur()
		m.textInput.Reset()
		m.statusBar.SetContent(m.getStatusBarContent())

		return m, cmd

	// A moveMsg is received any time a move operation has been performed
	// returning an updated listing of files
	case moveMsg:
		cfg := config.GetConfig()

		// Set active color back to the config default
		m.primaryPane.SetActiveBorderColor(cfg.Colors.Pane.ActiveBorderColor)

		// Set the dirtrees content to the new file listing and
		// display it in the primary pane
		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())

		// Set move mode back to false, set the initial moving directory to empty
		// the item that was moving back to nil and update the status bars content
		m.inMoveMode = false
		m.initialMoveDirectory = ""
		m.itemToMove = nil
		m.statusBar.SetContent(m.getStatusBarContent())

		return m, cmd

	// A fileContentMsg is received anytime a file is read from returning its content
	// along with the markdown content to be rendered by glamour
	case fileContentMsg:
		// Update the active markdown source to the markdownContent
		m.activeMarkdownSource = string(msg.markdownContent)

		// Set the content of the secondary pane to the file content removing any tabs
		// and converting them to spaces
		m.secondaryPane.SetContent(utils.ConverTabsToSpaces(string(msg.fileContent)))

		return m, cmd

	// Anytime markdown is being rendered, this message is received
	case markdownMsg:
		// Set the content of the secondary pane to the markdown content
		// converting any tabs into spaces
		m.secondaryPane.SetContent(utils.ConverTabsToSpaces(string(msg)))

		return m, cmd

	// An errorMsg is received any time something in a command goes wrong
	// we receive that error and show it in the secondary pane with red text
	case errorMsg:
		m.secondaryPane.SetContent(
			lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(constants.Red)).
				Render(string(msg)),
		)

		return m, cmd

	// Any time the window is resized this is called, including when the app
	// is first started
	case tea.WindowSizeMsg:
		if !m.ready {
			m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.primaryPane.SetContent(m.dirTree.View())
			m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.statusBar.SetContent(m.getStatusBarContent())
			m.statusBar.SetSize(msg.Width)
			m.ready = true
		} else {
			m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.StatusBarHeight)
			m.statusBar.SetSize(msg.Width)
		}

		// If we have some active markdown source to render, re render its content
		// when the window is resized so that glamour knows how to wrap its text
		if m.activeMarkdownSource != "" {
			return m, renderMarkdownContent(m.secondaryPane.Width, m.activeMarkdownSource)
		}

		return m, cmd

	// Any time a mouse event is received, we get this message
	case tea.MouseMsg:
		switch msg.Type {
		// Scroll up on the mouse wheel
		case tea.MouseWheelUp:
			// The command bar is not open and the primary pane is active
			// so scroll the dirtree up and update the primary panes content
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so scroll its content up
				m.secondaryPane.LineUp(3)
			}

			return m, cmd

		// Scroll down on the mouse wheel
		case tea.MouseWheelDown:
			// Command bar is not shown and the primary pane is active
			// so scroll the dirtree down and update the primary panes content
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so scroll its content down
				m.secondaryPane.LineDown(3)
			}

			return m, cmd
		}

	case tea.KeyMsg:
		// If gg is pressed
		if msg.String() == "g" && m.previousKey.String() == "g" {
			// If the command bar is not shown and the primary pane is active,
			// reset the previous key, go to the top of the dirtree, and the
			// primary pane and set the content of the primary pane to the contents
			// of the dirtree
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.previousKey = tea.KeyMsg{}
				m.dirTree.GotoTop()
				m.primaryPane.GotoTop()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so go to the top
				m.secondaryPane.GotoTop()
			}

			return m, cmd
		}

		switch msg.String() {
		// Exit FM
		case "ctrl+c":
			return m, tea.Quit

		// Exit FM is the command bar is not open
		case "q":
			if !m.showCommandBar {
				return m, tea.Quit
			}

		// Left arrow of h key is pressed
		case "left", "h":
			// If the command bar is not shown and the primary pane is active
			// set the previous directory to the current directory,
			// and update the directory listing to go back one directory
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.previousDirectory, _ = utils.GetWorkingDirectory()

				return m, m.updateDirectoryListing(constants.PreviousDirectory)
			}

		// Down arrow or j is pressed
		case "down", "j":
			// If the command bar is not shown and the primary pane is active
			// go down in the dirtree, update the status bar content and set the content
			// of the primary pane to the dirtree
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.statusBar.SetContent(m.getStatusBarContent())
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so scroll its content down
				m.secondaryPane.LineDown(1)
			}

		// Up arrow or k is pressed
		case "up", "k":
			// If the command bar is not shown and the primary pane is active
			// go up in the dirtree, update the status bar content and set the content
			// of the primary pane to the dirtree
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetContent(m.getStatusBarContent())
			} else {
				// Secondary pane is active so scroll its content up
				m.secondaryPane.LineUp(1)
			}

		// Right arrow or l is pressed
		case "right", "l":
			// Command bar is not shown and the primary pane is active
			if !m.showCommandBar && m.primaryPane.IsActive {
				// If the selected file is a directory and the textinput is not focused,
				// get an updated directory listing based on the currently selected file
				if m.dirTree.GetSelectedFile().IsDir() && !m.textInput.Focused() {
					return m, m.updateDirectoryListing(m.dirTree.GetSelectedFile().Name())
				} else {
					// The currently selected item is a file so scroll the secondary pane to the top
					// and get the file content to display
					m.secondaryPane.GotoTop()
					return m, m.readFileContent(m.dirTree.GetSelectedFile())
				}
			}

		case "G":
			// If the command bar is not shown and the primary pane is active
			// go to the bottom of the dirtree, and the bottom of the pane and set its
			// content to the dirtree
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GotoBottom()
				m.primaryPane.GotoBottom()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so go to the bottom of it
				m.secondaryPane.GotoBottom()
			}

		case "enter":
			// If a file or folder is currently being moved
			if m.inMoveMode {
				// The item to move is a directory
				if m.itemToMove.IsDir() {
					return m, m.moveDir(m.itemToMove.Name())
				} else {
					// The item to move is a file
					return m, m.moveFile(m.itemToMove.Name())
				}
			} else {
				// Parse the commands from the command bar, command is the name
				// of the command and value is if the command requires input to it
				// get its value, for example (rename test.txt) text.txt is the value
				command, value := utils.ParseCommand(m.textInput.Value())

				// Nothing was input for a command
				if command == "" {
					return m, nil
				}

				switch command {
				// Create a new directory based on the value passed
				case "mkdir":
					return m, m.createDir(value)

				// Create a new file based on the value passed
				case "touch":
					return m, m.createFile(value)

				// Rename the currently selected file or folder based on the value passed
				case "mv", "rename":
					return m, m.renameFileOrDir(m.dirTree.GetSelectedFile().Name(), value)

				// Delete the currently selected item
				case "rm", "delete":
					// If a directory is being deleted
					if m.dirTree.GetSelectedFile().IsDir() {
						return m, m.deleteDir(m.dirTree.GetSelectedFile().Name())
					} else {
						// The currently highlighted item is a file
						return m, m.deleteFile(m.dirTree.GetSelectedFile().Name())
					}

				default:
					return m, nil
				}
			}

		case ":":
			// If move mode is not active, activate the command bar,
			// focus the text input and give it a placeholder
			if !m.inMoveMode {
				m.showCommandBar = true
				m.textInput.Placeholder = "enter command"
				m.textInput.Focus()
			}

			return m, cmd

		// Shortcut to get back to the home directory if the
		// command bar is not curently open
		case "~":
			if !m.showCommandBar {
				homeDir, _ := utils.GetHomeDirectory()
				return m, m.updateDirectoryListing(homeDir)
			}

		// Shortcut to go back to the previous directory when your switching
		// between different directories, will only go back once to the previous
		// directory
		case "-":
			if !m.showCommandBar && m.previousDirectory != "" {
				return m, m.updateDirectoryListing(m.previousDirectory)
			}

		// Toggle weather or not to show hidden files and folders
		case ".":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.ToggleHidden()
				return m, m.updateDirectoryListing(constants.CurrentDirectory)
			}

		// Toggle between the two panes if the command bar is not currently active
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

		// If the command bar is not active and the primary pane is active
		// enter move mode, setting the active border color of the primary pane
		// to blue to indicate move mode. Get the directory in which the move was initiated
		// from, along with the selected file. Update the content of the status bar to indicate
		// a move is in progress
		case "m":
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.inMoveMode = true
				m.primaryPane.SetActiveBorderColor(constants.Blue)
				m.initialMoveDirectory, _ = utils.GetWorkingDirectory()
				m.itemToMove = m.dirTree.GetSelectedFile()
				m.statusBar.SetContent(m.getStatusBarContent())
			}

		// Zip up the currently selected item
		case "z":
			if !m.showCommandBar && m.primaryPane.IsActive {
				return m, m.zipDirectory(m.dirTree.GetSelectedFile().Name())
			}

		// Unzip the currently selected zip file
		case "u":
			if !m.showCommandBar && m.primaryPane.IsActive {
				return m, m.unzipDirectory(m.dirTree.GetSelectedFile().Name())
			}

		// Reset FM to its initial state
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

		// Capture the previous key so that we can capture
		// when two keys are pressed
		m.previousKey = msg
	}

	// If the command bar is shown, make sure to keep it updated
	// for when the user is typing into the command bar
	if m.showCommandBar {
		m.statusBar.SetContent(m.getStatusBarContent())
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
