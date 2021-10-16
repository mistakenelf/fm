package ui

import (
	"errors"
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
		m.showCommandInput,
		m.moveMode,
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
	// updateDirectoryListingMsg is received any time a directory is read from.
	case updateDirectoryListingMsg:
		m.showCommandInput = false
		m.createFileMode = false
		m.createDirectoryMode = false
		m.renameMode = false

		m.dirTree.GotoTop()
		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())
		m.primaryPane.GotoTop()
		m.statusBar.BlurCommandBar()
		m.statusBar.ResetCommandBar()
		m.updateStatusBarContent()

		return m, m.getDirectoryItemSizeCmd(m.dirTree.GetSelectedFile().Name())

	// previewDirectoryListingMsg is received any time a directory is previewed.
	case previewDirectoryListingMsg:
		m.showCommandInput = false
		m.createFileMode = false
		m.createDirectoryMode = false
		m.renameMode = false

		m.markdown.SetContent("")
		m.sourcecode.SetContent("")
		m.colorimage.SetImage(nil)
		m.dirTreePreview.GotoTop()
		m.dirTreePreview.SetContent(msg)
		m.secondaryPane.SetContent(m.dirTreePreview.View())
		m.secondaryPane.GotoTop()
		m.statusBar.BlurCommandBar()
		m.statusBar.ResetCommandBar()

		return m, nil

	// moveDirItemMsg is received any time a file or directory has been moved.
	case moveDirectoryItemMsg:
		m.moveMode = false
		m.initialMoveDirectory = ""
		m.itemToMove = nil

		m.primaryPane.ShowAlternateBorder(false)
		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())
		m.updateStatusBarContent()

		return m, nil

	// readFileContentMsg is received any time a file is read from.
	case readFileContentMsg:
		switch {
		case msg.code != "":
			m.secondaryPane.GotoTop()
			m.colorimage.SetImage(nil)
			m.markdown.SetContent("")
			m.dirTreePreview.SetContent(nil)
			m.sourcecode.SetContent(msg.code)
			m.secondaryPane.SetContent(m.sourcecode.View())
		case msg.markdown != "":
			m.secondaryPane.GotoTop()
			m.colorimage.SetImage(nil)
			m.sourcecode.SetContent("")
			m.dirTreePreview.SetContent(nil)
			m.markdown.SetContent(msg.markdown)
			m.secondaryPane.SetContent(m.markdown.View())
		case msg.image != nil:
			m.secondaryPane.GotoTop()
			m.markdown.SetContent("")
			m.sourcecode.SetContent("")
			m.dirTreePreview.SetContent(nil)
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

	// errorMsg is received any time something goes wrong.
	case errorMsg:
		m.secondaryPane.SetContent(
			lipgloss.NewStyle().
				Bold(true).
				Foreground(m.theme.ErrorColor).
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize()).
				Render(string(msg)),
		)

		return m, nil

	// directoryItemSizeMsg is received any time an updated directory size is needed.
	case directoryItemSizeMsg:
		m.statusBar.SetItemSize(string(msg))

		return m, nil

	// tea.WindowSizeMsg is received anytime the terminal is resized.
	case tea.WindowSizeMsg:
		m.primaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
		m.secondaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
		m.dirTree.SetSize(m.primaryPane.GetWidth())
		m.dirTreePreview.SetSize(m.secondaryPane.GetWidth())
		m.sourcecode.SetSize(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize())
		m.markdown.SetSize(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize())
		m.primaryPane.SetContent(m.dirTree.View())
		m.help.Width = msg.Width

		switch {
		case m.colorimage.GetImage() != nil:
			return m, m.redrawImageCmd(
				m.secondaryPane.GetWidth()-m.secondaryPane.GetHorizontalFrameSize(),
				m.secondaryPane.GetHeight(),
			)
		case m.markdown.GetContent() != "":
			m.secondaryPane.SetContent(m.markdown.View())
		case m.sourcecode.GetContent() != "":
			m.secondaryPane.SetContent(m.sourcecode.View())
		case m.dirTreePreview.GetTotalFiles() != 0:
			m.secondaryPane.SetContent(m.dirTreePreview.View())
		case m.sourcecode.GetContent() == "" && m.markdown.GetContent() == "" && m.colorimage.GetImage() == nil && m.dirTreePreview.GetTotalFiles() == 0:
			m.secondaryPane.SetContent(lipgloss.NewStyle().
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize()).
				Render(m.help.View(m.keys)),
			)
		}

		if !m.ready {
			m.ready = true
			m.secondaryPane.SetContent(lipgloss.NewStyle().
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.Style.GetHorizontalFrameSize()).
				Render(m.help.View(m.keys)),
			)
		}

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			// The command bar is not open and the primary pane is active
			// so scroll the dirtree up and update the primary panes content.
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				m.previousKey = tea.KeyMsg{}

				return m, m.getDirectoryItemSizeCmd(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineUp(3)

			return m, nil

		case tea.MouseWheelDown:
			// Command bar is not shown and the primary pane is active
			// so scroll the dirtree down and update the primary panes content.
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				m.previousKey = tea.KeyMsg{}

				return m, m.getDirectoryItemSizeCmd(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineDown(3)

			return m, nil
		}

	// tea.KeyMsg is received anytime a key event is received.
	case tea.KeyMsg:
		// when gg is pressed.
		if msg.String() == "g" && m.previousKey.String() == "g" && !m.showCommandInput {
			if m.primaryPane.GetIsActive() {
				m.previousKey = tea.KeyMsg{}
				m.dirTree.GotoTop()
				m.primaryPane.GotoTop()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")

				return m, m.getDirectoryItemSizeCmd(m.dirTree.GetSelectedFile().Name())
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
			if !m.showCommandInput {
				return m, tea.Quit
			}

		// Go back a directory.
		case key.Matches(msg, m.keys.Left):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.statusBar.SetItemSize("")
				m.previousKey = tea.KeyMsg{}
				previousDirectory, err := dirfs.GetWorkingDirectory()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				m.previousDirectory = previousDirectory

				return m, m.updateDirectoryListingCmd(
					fmt.Sprintf("%s/%s", m.previousDirectory, dirfs.PreviousDirectory),
				)
			}

		// Scroll down in directory tree or pane.
		case key.Matches(msg, m.keys.Down):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				m.previousKey = tea.KeyMsg{}

				return m, m.getDirectoryItemSizeCmd(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineDown(1)

		// Scroll up in directory tree or pane.
		case key.Matches(msg, m.keys.Up):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				m.updateStatusBarContent()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				m.previousKey = tea.KeyMsg{}

				return m, m.getDirectoryItemSizeCmd(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.LineUp(1)

		// Open directory or read file content and display it.
		case key.Matches(msg, m.keys.Right):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				m.previousKey = tea.KeyMsg{}

				switch {
				case m.dirTree.GetSelectedFile().IsDir() && !m.statusBar.CommandBarFocused():
					m.statusBar.SetItemSize("")
					currentDir, err := dirfs.GetWorkingDirectory()
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					return m, m.updateDirectoryListingCmd(
						fmt.Sprintf("%s/%s", currentDir, m.dirTree.GetSelectedFile().Name()),
					)
				case m.dirTree.GetSelectedFile().Mode()&os.ModeSymlink == os.ModeSymlink:
					m.statusBar.SetItemSize("")
					symlinkFile, err := os.Readlink(m.dirTree.GetSelectedFile().Name())
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					fileInfo, err := os.Stat(symlinkFile)
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					if fileInfo.IsDir() {
						return m, m.updateDirectoryListingCmd(symlinkFile)
					}

					return m, m.readFileContentCmd(
						fileInfo,
						m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(),
						m.secondaryPane.GetHeight(),
					)
				default:
					return m, m.readFileContentCmd(
						m.dirTree.GetSelectedFile(),
						m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(),
						m.secondaryPane.GetHeight(),
					)
				}
			}

		// Jump to bottom of directory tree or pane.
		case key.Matches(msg, m.keys.GotoBottom):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.GotoBottom()
				m.primaryPane.GotoBottom()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				m.previousKey = tea.KeyMsg{}

				return m, m.getDirectoryItemSizeCmd(m.dirTree.GetSelectedFile().Name())
			}

			m.secondaryPane.GotoBottom()

		// Handle confirmation of actions (move, create and rename).
		case key.Matches(msg, m.keys.Enter):
			switch {
			case m.moveMode:
				return m, m.moveDirectoryItemCmd(m.itemToMove.Name())
			case m.createFileMode:
				return m, tea.Sequentially(
					m.createFileCmd(m.statusBar.CommandBarValue()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			case m.createDirectoryMode:
				return m, tea.Sequentially(
					m.createDirectoryCmd(m.statusBar.CommandBarValue()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			case m.renameMode:
				return m, tea.Sequentially(
					m.renameDirectoryItemCmd(m.dirTree.GetSelectedFile().Name(), m.statusBar.CommandBarValue()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			default:
				return m, nil
			}

		// Delete the currently selected directory item.
		case key.Matches(msg, m.keys.Delete):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, tea.Sequentially(
						m.deleteDirectoryCmd(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.deleteFileCmd(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Create file mode.
		case key.Matches(msg, m.keys.CreateFile):
			if !m.moveMode && !m.createDirectoryMode && !m.showCommandInput {
				m.createFileMode = true
				m.showCommandInput = true
				m.statusBar.FocusCommandBar()
				m.updateStatusBarContent()

				return m, nil
			}

		// Create directory mode.
		case key.Matches(msg, m.keys.CreateDirectory):
			if !m.moveMode && !m.createFileMode && !m.showCommandInput {
				m.createDirectoryMode = true
				m.showCommandInput = true
				m.statusBar.FocusCommandBar()
				m.updateStatusBarContent()

				return m, nil
			}

		// Rename directory mode.
		case key.Matches(msg, m.keys.Rename):
			if !m.moveMode && !m.createFileMode && !m.createDirectoryMode && !m.showCommandInput {
				m.renameMode = true
				m.showCommandInput = true
				m.statusBar.FocusCommandBar()
				m.updateStatusBarContent()

				return m, nil
			}

		// Go to home directory.
		case key.Matches(msg, m.keys.OpenHomeDirectory):
			if !m.showCommandInput {
				homeDir, err := dirfs.GetHomeDirectory()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.updateDirectoryListingCmd(homeDir)
			}

		// Go to previous directory.
		case key.Matches(msg, m.keys.OpenPreviousDirectory):
			if !m.showCommandInput && m.previousDirectory != "" {
				return m, m.updateDirectoryListingCmd(m.previousDirectory)
			}

		// Go to root directory.
		case key.Matches(msg, m.keys.OpenRootDirectory):
			if !m.showCommandInput {
				return m, m.updateDirectoryListingCmd(dirfs.RootDirectory)
			}

		// Toggle hidden files and folders.
		case key.Matches(msg, m.keys.ToggleHidden):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.ToggleHidden()
				return m, m.updateDirectoryListingCmd(dirfs.CurrentDirectory)
			}

		// Toggle between the two panes if the command bar is not currently active.
		case key.Matches(msg, m.keys.Tab):
			if !m.showCommandInput {
				m.primaryPane.SetActive(!m.primaryPane.GetIsActive())
				m.secondaryPane.SetActive(!m.secondaryPane.GetIsActive())
			}

		// Move mode.
		case key.Matches(msg, m.keys.EnterMoveMode):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				m.moveMode = true
				m.primaryPane.ShowAlternateBorder(true)
				initialMoveDirectory, err := dirfs.GetWorkingDirectory()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				m.initialMoveDirectory = initialMoveDirectory
				m.itemToMove = m.dirTree.GetSelectedFile()
				m.updateStatusBarContent()
			}

		// Create zip of currently selected directory item.
		case key.Matches(msg, m.keys.Zip):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				return m, tea.Sequentially(
					m.zipDirectoryCmd(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Unzip the currently selected zip file.
		case key.Matches(msg, m.keys.Unzip):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				return m, tea.Sequentially(
					m.unzipDirectoryCmd(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Copy the currently selected directory item.
		case key.Matches(msg, m.keys.Copy):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				if m.dirTree.GetSelectedFile().IsDir() {
					return m, tea.Sequentially(
						m.copyDirectoryCmd(m.dirTree.GetSelectedFile().Name()),
						m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.copyFileCmd(m.dirTree.GetSelectedFile().Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Edit the currently selected file.
		case key.Matches(msg, m.keys.EditFile):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && !m.dirTree.GetSelectedFile().IsDir() {
				editorPath := os.Getenv("EDITOR")
				if editorPath == "" {
					return m, m.handleErrorCmd(errors.New("$EDITOR not set"))
				}

				editorCmd := exec.Command(editorPath, m.dirTree.GetSelectedFile().Name())
				editorCmd.Stdin = os.Stdin
				editorCmd.Stdout = os.Stdout
				editorCmd.Stderr = os.Stderr
				err := editorCmd.Start()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				err = editorCmd.Wait()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.updateDirectoryListingCmd(dirfs.CurrentDirectory)
			}

		// Preview the currently selected directory.
		case key.Matches(msg, m.keys.PreviewDirectory):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetSelectedFile().IsDir() {
				return m, m.previewDirectoryListingCmd(m.dirTree.GetSelectedFile().Name())
			}

		// Reset to initial state.
		case key.Matches(msg, m.keys.Escape):
			m.showCommandInput = false
			m.moveMode = false
			m.itemToMove = nil
			m.initialMoveDirectory = ""
			m.help.ShowAll = true
			m.createFileMode = false
			m.createDirectoryMode = false
			m.renameMode = false
			m.previousKey = tea.KeyMsg{}
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
			m.dirTreePreview.SetContent(nil)
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
