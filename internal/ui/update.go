package ui

import (
	"fmt"

	"github.com/knipferrc/fm/internal/constants"
	"github.com/knipferrc/fm/internal/helpers"
	"github.com/knipferrc/fm/internal/statusbar"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// scrollPrimaryPane handles the scrolling of the primary pane which will handle
// infinite scroll on the dirtree and the scrolling of the viewport.
func (m *Model) scrollPrimaryPane() {
	top := m.primaryPane.GetYOffset()
	bottom := m.primaryPane.GetHeight() + m.primaryPane.GetYOffset() - 1

	// If the cursor is above the top of the viewport scroll up on the viewport
	// else were at the bottom and need to scroll the viewport down.
	if m.dirTree.GetCursor() < top {
		m.primaryPane.LineUp(1)
	} else if m.dirTree.GetCursor() > bottom {
		m.primaryPane.LineDown(1)
	}

	// If the cursor of the dirtree is at the bottom of the files
	// set the cursor to 0 to go to the top of the dirtree and
	// scroll the pane to the top else, were at the top of the dirtree and pane so
	// scroll the pane to the bottom and set the cursor to the bottom.
	if m.dirTree.GetCursor() > m.dirTree.GetTotalFiles()-1 {
		m.dirTree.GotoTop()
		m.primaryPane.GotoTop()
	} else if m.dirTree.GetCursor() < top {
		m.dirTree.GotoBottom()
		m.primaryPane.GotoBottom()
	}
}

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
			m.colorimage.SetImage(nil)
			m.markdown.SetContent("")
			m.text.SetContent(msg.code)
			m.secondaryPane.SetContent(m.text.View())
		} else if msg.markdown != "" {
			m.secondaryPane.GotoTop()
			m.colorimage.SetImage(nil)
			m.text.SetContent("")
			m.markdown.SetContent(msg.markdown)
			m.secondaryPane.SetContent(m.markdown.View())
		} else if msg.image != nil {
			m.secondaryPane.GotoTop()
			m.markdown.SetContent("")
			m.text.SetContent("")
			m.colorimage.SetImage(msg.image)
			m.colorimage.SetContent(msg.asciiImage)
			m.secondaryPane.SetContent(m.colorimage.View())
		} else {
			m.secondaryPane.GotoTop()
			m.secondaryPane.SetContent(msg.rawContent)
		}

		return m, nil

	// convertImageToStringMsg is received when an image is to be converted to a string.
	case convertImageToStringMsg:
		m.colorimage.SetContent(string(msg))
		m.secondaryPane.SetContent(m.colorimage.View())

		return m, nil

	// An errorMsg is received any time something in a command goes wrong
	// we receive that error and show it in the secondary pane with red text.
	case errorMsg:
		m.secondaryPane.SetContent(
			lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(constants.Colors.Red)).
				Render(string(msg)),
		)

		return m, nil

	// Any time the window is resized this is called, including when the app
	// is first started.
	case tea.WindowSizeMsg:
		m.primaryPane.SetSize(msg.Width/2, msg.Height-constants.Dimensions.StatusBarHeight)
		m.secondaryPane.SetSize(msg.Width/2, msg.Height-constants.Dimensions.StatusBarHeight)
		m.dirTree.SetSize(m.primaryPane.GetWidth())
		m.primaryPane.SetContent(m.dirTree.View())
		m.statusBar.SetSize(msg.Width, constants.Dimensions.StatusBarHeight)
		m.text.SetSize(m.secondaryPane.GetWidth() - constants.Dimensions.PanePadding)
		m.markdown.SetSize(m.secondaryPane.GetWidth() - constants.Dimensions.PanePadding)
		m.help.Width = msg.Width

		if m.colorimage.Image != nil {
			resizeImageCmd := m.redrawImage(m.secondaryPane.GetWidth()-constants.Dimensions.PanePadding, m.secondaryPane.GetHeight())
			cmds = append(cmds, resizeImageCmd)
		}

		if m.markdown.Content != "" {
			m.secondaryPane.SetContent(m.markdown.View())
		}

		if m.text.Content != "" {
			m.secondaryPane.SetContent(m.text.View())
		}

		if !m.ready {
			m.ready = true
			m.secondaryPane.SetContent(m.help.View(m.keys))
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

		switch {
		// Exit FM.
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit

		// Exit FM if the command bar is not open.
		case key.Matches(msg, m.keys.Quit):
			if !m.showCommandBar {
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.Left):
			// If the command bar is not shown and the primary pane is active
			// set the previous directory to the current directory,
			// and update the directory listing to go back one directory.
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.previousDirectory, _ = helpers.GetWorkingDirectory()

				return m, m.updateDirectoryListing(fmt.Sprintf("%s/%s", m.previousDirectory, constants.Directories.PreviousDirectory))
			}

		case key.Matches(msg, m.keys.Down):
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

		case key.Matches(msg, m.keys.Up):
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

		case key.Matches(msg, m.keys.Right):
			// Open directory or read file content.
			if !m.showCommandBar && m.primaryPane.IsActive {
				if m.dirTree.GetSelectedFile().IsDir() && !m.textInput.Focused() {
					currentDir, _ := helpers.GetWorkingDirectory()
					return m, m.updateDirectoryListing(fmt.Sprintf("%s/%s", currentDir, m.dirTree.GetSelectedFile().Name()))
				}

				return m, m.readFileContent(m.dirTree.GetSelectedFile(), m.secondaryPane.GetWidth()-constants.Dimensions.PanePadding, m.secondaryPane.GetHeight())
			}

		case key.Matches(msg, m.keys.GotoBottom):
			// Go to the bottom of the pane.
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.GotoBottom()
				m.primaryPane.GotoBottom()
				m.primaryPane.SetContent(m.dirTree.View())
			} else {
				// Secondary pane is active so go to the bottom of it.
				m.secondaryPane.GotoBottom()
			}

		case key.Matches(msg, m.keys.Enter):
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
			}

			// Parse the commands from the command bar, command is the name
			// of the command and value is if the command requires input to it
			// get its value, for example (rename test.txt) text.txt is the value.
			command, value := statusbar.ParseCommand(m.textInput.Value())

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
					m.updateDirectoryListing(constants.Directories.CurrentDirectory),
				)

			// Create a new file based on the value passed.
			case "touch":
				return m, tea.Sequentially(
					m.createFile(value),
					m.updateDirectoryListing(constants.Directories.CurrentDirectory),
				)

			// Rename the currently selected file or folder based on the value passed.
			case "mv", "rename":
				return m, tea.Sequentially(
					m.renameFileOrDir(m.dirTree.GetSelectedFile().Name(), value),
					m.updateDirectoryListing(constants.Directories.CurrentDirectory),
				)

			// Delete the currently selected item.
			case "rm", "delete":
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, tea.Sequentially(
						m.deleteDir(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListing(constants.Directories.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.deleteFile(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(constants.Directories.CurrentDirectory),
				)

			default:
				return m, nil
			}

		case key.Matches(msg, m.keys.OpenCommandBar):
			// If move mode is not active, activate the command bar.
			if !m.inMoveMode {
				m.showCommandBar = true
				m.textInput.Placeholder = "enter command"
				m.textInput.Focus()
			}

			return m, cmd

		// Shortcut to get back to the home directory if the
		// command bar is not curently open.
		case key.Matches(msg, m.keys.OpenHomeDirectory):
			if !m.showCommandBar {
				homeDir, _ := helpers.GetHomeDirectory()
				return m, m.updateDirectoryListing(homeDir)
			}

		// Shortcut to go back to the previous directory.
		case key.Matches(msg, m.keys.OpenPreviousDirectory):
			if !m.showCommandBar && m.previousDirectory != "" {
				return m, m.updateDirectoryListing(m.previousDirectory)
			}

		// Toggle hidden files and folders.
		case key.Matches(msg, m.keys.ToggleHidden):
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.dirTree.ToggleHidden()
				return m, m.updateDirectoryListing(constants.Directories.CurrentDirectory)
			}

		// Toggle between the two panes if the command bar is not currently active.
		case key.Matches(msg, m.keys.Tab):
			if !m.showCommandBar {
				m.primaryPane.IsActive = !m.primaryPane.IsActive
				m.secondaryPane.IsActive = !m.secondaryPane.IsActive
			}

		// Enter move mode.
		case key.Matches(msg, m.keys.EnterMoveMode):
			if !m.showCommandBar && m.primaryPane.IsActive {
				m.inMoveMode = true
				m.primaryPane.SetActiveBorderColor(constants.Colors.Blue)
				m.initialMoveDirectory, _ = helpers.GetWorkingDirectory()
				m.itemToMove = m.dirTree.GetSelectedFile()
			}

		// Zip up the currently selected item.
		case key.Matches(msg, m.keys.Zip):
			if !m.showCommandBar && m.primaryPane.IsActive {
				return m, tea.Sequentially(
					m.zipDirectory(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(constants.Directories.CurrentDirectory),
				)
			}

		// Unzip the currently selected zip file.
		case key.Matches(msg, m.keys.Unzip):
			if !m.showCommandBar && m.primaryPane.IsActive {
				return m, tea.Sequentially(
					m.unzipDirectory(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(constants.Directories.CurrentDirectory),
				)
			}

		// Copy the currently selected item.
		case key.Matches(msg, m.keys.Copy):
			if !m.showCommandBar && m.primaryPane.IsActive {
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, tea.Sequentially(
						m.copyDirectory(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListing(constants.Directories.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.copyFile(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(constants.Directories.CurrentDirectory),
				)
			}

		// Reset FM to its initial state.
		case key.Matches(msg, m.keys.Escape):
			m.showCommandBar = false
			m.inMoveMode = false
			m.itemToMove = nil
			m.initialMoveDirectory = ""
			m.primaryPane.IsActive = true
			m.secondaryPane.IsActive = false
			m.help.ShowAll = true
			m.textInput.Blur()
			m.textInput.Reset()
			m.secondaryPane.GotoTop()
			m.primaryPane.SetActiveBorderColor(m.appConfig.Colors.Pane.ActiveBorderColor)
			m.secondaryPane.SetContent(m.help.View(m.keys))
			m.colorimage.SetImage(nil)
			m.markdown.SetContent("")
		}

		// Capture the previous key so that we can capture
		// when two keys are pressed.
		m.previousKey = msg
	}

	if m.ready && m.dirTree.GetTotalFiles() > 0 {
		m.statusBar.SetContent(
			m.dirTree.GetTotalFiles(),
			m.dirTree.GetCursor(),
			m.textInput.View(),
			m.appConfig.Settings.ShowIcons,
			m.showCommandBar,
			m.inMoveMode,
			m.dirTree.GetSelectedFile(),
			m.itemToMove,
		)
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
