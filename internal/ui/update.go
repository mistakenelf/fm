package ui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/knipferrc/fm/dirfs"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// updateStatusBarContent updates the content of the statusbar.
func (m *Model) updateStatusBarContent() {
	m.statusBar.SetContent(
		m.dirTree.GetTotalFiles(),
		m.dirTree.GetCursor(),
		m.showCommandBar,
		m.inMoveMode,
		m.dirTree.GetSelectedFile(),
		m.itemToMove,
	)
}

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
		m.showCommandBar = false
		m.inCreateFileMode = false
		m.inCreateDirectoryMode = false
		m.inRenameMode = false

		m.dirTree.GotoTop()
		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())
		m.primaryPane.GotoTop()
		m.statusBar.BlurCommandBar()
		m.statusBar.ResetCommandBar()
		m.updateStatusBarContent()

		return m, m.getDirectoryItemSize(m.dirTree.GetSelectedFile().Name())

	// A moveDirItemMsg is received any time a file or directory has been moved.
	case moveDirItemMsg:
		m.inMoveMode = false
		m.initialMoveDirectory = ""
		m.itemToMove = nil

		m.primaryPane.ShowAlternateBorder(false)
		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())
		m.updateStatusBarContent()

		return m, nil

	// A fileContentMsg is received anytime a file is read from returning its content
	// along with the markdown content to be rendered by glamour.
	case readFileContentMsg:
		switch {
		case msg.code != "":
			m.secondaryPane.GotoTop()
			m.colorimage.SetImage(nil)
			m.markdown.SetContent("")
			m.sourcecode.SetContent(msg.code)
			m.secondaryPane.SetContent(m.sourcecode.View())
		case msg.markdown != "":
			m.secondaryPane.GotoTop()
			m.colorimage.SetImage(nil)
			m.sourcecode.SetContent("")
			m.markdown.SetContent(msg.markdown)
			m.secondaryPane.SetContent(m.markdown.View())
		case msg.image != nil:
			m.secondaryPane.GotoTop()
			m.markdown.SetContent("")
			m.sourcecode.SetContent("")
			m.colorimage.SetImage(msg.image)
			m.colorimage.SetContent(msg.imageString)
			m.secondaryPane.SetContent(m.colorimage.View())
		default:
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
				Foreground(m.theme.ErrorColor).
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize()).
				Render(string(msg)),
		)

		return m, nil

	// A directoryItemSizeMsg is received anytime a new file is selected
	// in the dirtree returning the file size as a string.
	case directoryItemSizeMsg:
		m.statusBar.SetItemSize(string(msg))

		return m, nil

	// Any time the window is resized this is called, including when the app
	// is first started.
	case tea.WindowSizeMsg:
		m.primaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
		m.secondaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
		m.dirTree.SetSize(m.primaryPane.GetWidth())
		m.sourcecode.SetSize(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize())
		m.markdown.SetSize(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize())
		m.primaryPane.SetContent(m.dirTree.View())
		m.help.Width = msg.Width

		switch {
		case m.colorimage.GetImage() != nil:
			return m, m.redrawImage(
				m.secondaryPane.GetWidth()-m.secondaryPane.GetHorizontalFrameSize(),
				m.secondaryPane.GetHeight(),
			)
		case m.markdown.GetContent() != "":
			m.secondaryPane.SetContent(m.markdown.View())
		case m.sourcecode.GetContent() != "":
			m.secondaryPane.SetContent(m.sourcecode.View())
		case m.sourcecode.GetContent() == "" && m.markdown.GetContent() == "" && m.colorimage.GetImage() == nil:
			m.secondaryPane.SetContent(lipgloss.NewStyle().
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize()).
				Render(m.help.View(m.keys)),
			)
		default:
			return m, nil
		}

		if !m.ready {
			m.ready = true
			m.secondaryPane.SetContent(lipgloss.NewStyle().
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.Style.GetHorizontalFrameSize()).
				Render(m.help.View(m.keys)),
			)
		}

	// Any time a mouse event is received, we get this message.
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			// The command bar is not open and the primary pane is active
			// so scroll the dirtree up and update the primary panes content.
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")

				return m, m.getDirectoryItemSize(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineUp(3)

			return m, nil

		case tea.MouseWheelDown:
			// Command bar is not shown and the primary pane is active
			// so scroll the dirtree down and update the primary panes content.
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")

				return m, m.getDirectoryItemSize(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineDown(3)

			return m, nil
		}

	case tea.KeyMsg:
		// If gg is pressed.
		if msg.String() == "g" && m.previousKey.String() == "g" && !m.showCommandBar {
			// If the command bar is not shown and the primary pane is active,
			// reset the previous key, go to the top of the dirtree and pane.
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.previousKey = tea.KeyMsg{}
				m.dirTree.GotoTop()
				m.primaryPane.GotoTop()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")

				return m, m.getDirectoryItemSize(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.GotoTop()

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

		// If the command bar is not shown and the primary pane is active
		// set the previous directory to the current directory,
		// and update the directory listing to go back one directory.
		case key.Matches(msg, m.keys.Left):
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.statusBar.SetItemSize("")
				m.previousDirectory, _ = dirfs.GetWorkingDirectory()

				return m, m.updateDirectoryListing(
					fmt.Sprintf("%s/%s", m.previousDirectory, dirfs.PreviousDirectory),
				)
			}

			// Scroll pane down.
		case key.Matches(msg, m.keys.Down):
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")

				return m, m.getDirectoryItemSize(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineDown(1)

		// Scroll pane up.
		case key.Matches(msg, m.keys.Up):
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")

				return m, m.getDirectoryItemSize(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineUp(1)

		// Open directory or read file content and display in secondary pane.
		case key.Matches(msg, m.keys.Right):
			if !m.showCommandBar && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				switch {
				case m.dirTree.GetSelectedFile().IsDir() && !m.statusBar.CommandBarFocused():
					m.statusBar.SetItemSize("")
					currentDir, _ := dirfs.GetWorkingDirectory()

					return m, m.updateDirectoryListing(
						fmt.Sprintf("%s/%s", currentDir, m.dirTree.GetSelectedFile().Name()),
					)
				case m.dirTree.GetSelectedFile().Mode()&os.ModeSymlink == os.ModeSymlink:
					m.statusBar.SetItemSize("")
					symlinkFile, _ := os.Readlink(m.dirTree.GetSelectedFile().Name())

					return m, m.updateDirectoryListing(symlinkFile)
				default:
					return m, m.readFileContent(
						m.dirTree.GetSelectedFile(),
						m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(),
						m.secondaryPane.GetHeight(),
					)
				}
			}

		// Jump to the bottom of a pane.
		case key.Matches(msg, m.keys.GotoBottom):
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.dirTree.GotoBottom()
				m.primaryPane.GotoBottom()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")

				return m, m.getDirectoryItemSize(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.GotoBottom()

		// If in move mode, place the selected item to move into
		// the current directory. If the commandbar is open, process the command.
		case key.Matches(msg, m.keys.Enter):
			switch {
			case m.inMoveMode:
				return m, m.moveDirectoryItem(m.itemToMove.Name())
			case m.inCreateFileMode:
				return m, tea.Sequentially(
					m.createFile(m.statusBar.CommandBarValue()),
					m.updateDirectoryListing(dirfs.CurrentDirectory),
				)
			case m.inCreateDirectoryMode:
				return m, tea.Sequentially(
					m.createDir(m.statusBar.CommandBarValue()),
					m.updateDirectoryListing(dirfs.CurrentDirectory),
				)
			case m.inRenameMode:
				return m, tea.Sequentially(
					m.renameFileOrDir(m.dirTree.GetSelectedFile().Name(), m.statusBar.CommandBarValue()),
					m.updateDirectoryListing(dirfs.CurrentDirectory),
				)
			default:
				return m, nil
			}

		// Delete the currently selected item.
		case key.Matches(msg, m.keys.Delete):
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, tea.Sequentially(
						m.deleteDir(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListing(dirfs.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.deleteFile(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(dirfs.CurrentDirectory),
				)
			}

		// Enter create file mode.
		case key.Matches(msg, m.keys.CreateFile):
			if !m.inMoveMode && !m.inCreateDirectoryMode && !m.showCommandBar {
				m.inCreateFileMode = true
				m.showCommandBar = true
				m.statusBar.FocusCommandBar()
				m.updateStatusBarContent()

				return m, nil
			}

		// Enter create directory mode.
		case key.Matches(msg, m.keys.CreateDirectory):
			if !m.inMoveMode && !m.inCreateFileMode && !m.showCommandBar {
				m.inCreateDirectoryMode = true
				m.showCommandBar = true
				m.statusBar.FocusCommandBar()
				m.updateStatusBarContent()

				return m, nil
			}

		// Enter create directory mode.
		case key.Matches(msg, m.keys.Rename):
			if !m.inMoveMode && !m.inCreateFileMode && !m.inCreateDirectoryMode && !m.showCommandBar {
				m.inRenameMode = true
				m.showCommandBar = true
				m.statusBar.FocusCommandBar()
				m.updateStatusBarContent()

				return m, nil
			}

		// Shortcut to get back to the home directory if the
		// command bar is not curently open.
		case key.Matches(msg, m.keys.OpenHomeDirectory):
			if !m.showCommandBar {
				homeDir, _ := dirfs.GetHomeDirectory()
				return m, m.updateDirectoryListing(homeDir)
			}

		// Shortcut to go back to the previous directory.
		case key.Matches(msg, m.keys.OpenPreviousDirectory):
			if !m.showCommandBar && m.previousDirectory != "" {
				return m, m.updateDirectoryListing(m.previousDirectory)
			}

		// Shortcut to go back to the root directory.
		case key.Matches(msg, m.keys.OpenRootDirectory):
			if !m.showCommandBar {
				return m, m.updateDirectoryListing(dirfs.RootDirectory)
			}

		// Toggle hidden files and folders.
		case key.Matches(msg, m.keys.ToggleHidden):
			if !m.showCommandBar && m.primaryPane.GetIsActive() {
				m.dirTree.ToggleHidden()
				return m, m.updateDirectoryListing(dirfs.CurrentDirectory)
			}

		// Toggle between the two panes if the command bar is not currently active.
		case key.Matches(msg, m.keys.Tab):
			if !m.showCommandBar {
				m.primaryPane.SetActive(!m.primaryPane.GetIsActive())
				m.secondaryPane.SetActive(!m.secondaryPane.GetIsActive())
			}

		// Enter move mode.
		case key.Matches(msg, m.keys.EnterMoveMode):
			if !m.showCommandBar && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				m.inMoveMode = true
				m.primaryPane.ShowAlternateBorder(true)
				m.initialMoveDirectory, _ = dirfs.GetWorkingDirectory()
				m.itemToMove = m.dirTree.GetSelectedFile()
				m.updateStatusBarContent()
			}

		// Zip up the currently selected item.
		case key.Matches(msg, m.keys.Zip):
			if !m.showCommandBar && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				return m, tea.Sequentially(
					m.zipDirectory(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(dirfs.CurrentDirectory),
				)
			}

		// Unzip the currently selected zip file.
		case key.Matches(msg, m.keys.Unzip):
			if !m.showCommandBar && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				return m, tea.Sequentially(
					m.unzipDirectory(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(dirfs.CurrentDirectory),
				)
			}

		// Copy the currently selected item.
		case key.Matches(msg, m.keys.Copy):
			if !m.showCommandBar && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, tea.Sequentially(
						m.copyDirectory(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListing(dirfs.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.copyFile(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListing(dirfs.CurrentDirectory),
				)
			}

		// Edit the currently selected file.
		case key.Matches(msg, m.keys.EditFile):
			if !m.showCommandBar && m.primaryPane.GetIsActive() && !m.dirTree.GetSelectedFile().IsDir() {
				vimCmd := exec.Command("vim", m.dirTree.GetSelectedFile().Name())
				vimCmd.Stdin = os.Stdin
				vimCmd.Stdout = os.Stdout
				vimCmd.Stderr = os.Stderr
				err := vimCmd.Start()
				if err != nil {
					return m, nil
				}
				err = vimCmd.Wait()
				if err != nil {
					return m, nil
				}
			}

			return m, m.updateDirectoryListing(dirfs.CurrentDirectory)

		// Reset FM to its initial state.
		case key.Matches(msg, m.keys.Escape):
			m.showCommandBar = false
			m.inMoveMode = false
			m.itemToMove = nil
			m.initialMoveDirectory = ""
			m.help.ShowAll = true
			m.inCreateFileMode = false
			m.inCreateDirectoryMode = false
			m.inRenameMode = false
			m.primaryPane.SetActive(true)
			m.secondaryPane.SetActive(false)
			m.statusBar.BlurCommandBar()
			m.statusBar.ResetCommandBar()
			m.secondaryPane.GotoTop()
			m.primaryPane.ShowAlternateBorder(false)
			m.secondaryPane.SetContent(lipgloss.NewStyle().
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.Style.GetHorizontalFrameSize()).
				Render(m.help.View(m.keys)),
			)
			m.colorimage.SetImage(nil)
			m.markdown.SetContent("")
			m.sourcecode.SetContent("")
			m.updateStatusBarContent()
		}

		// Capture the previous key so that we can capture
		// when two keys are pressed.
		m.previousKey = msg
	}

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
