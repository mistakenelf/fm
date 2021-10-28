package ui

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/knipferrc/fm/dirfs"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// updateStatusBarContent updates the content of the statusbar.
func (m *Model) updateStatusBarContent() error {
	selectedFile, err := m.dirTree.GetSelectedFile()
	if err != nil {
		return err
	}

	m.statusBar.SetContent(
		m.dirTree.GetTotalFiles(),
		m.dirTree.GetCursor(),
		m.showCommandInput,
		m.moveMode,
		selectedFile,
		m.itemToMove,
	)

	return nil
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
	// updateDirectoryListingMsg is received when a directory is read from.
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
		err := m.updateStatusBarContent()
		if err != nil {
			return m, m.handleErrorCmd(err)
		}

		if len(msg) > 0 {
			selectedFile, err := m.dirTree.GetSelectedFile()
			if err != nil {
				return m, m.handleErrorCmd(err)
			}

			return m, m.getDirectoryItemSizeCmd(selectedFile.Name())
		}

		return m, nil

	// previewDirectoryListingMsg shows a preview of a directory in the secondary pane.
	case previewDirectoryListingMsg:
		m.showCommandInput = false
		m.createFileMode = false
		m.createDirectoryMode = false
		m.renameMode = false

		m.renderer.SetContent("")
		m.renderer.SetImage(nil)
		m.dirTreePreview.GotoTop()
		m.dirTreePreview.SetContent(msg)
		m.secondaryPane.SetContent(m.dirTreePreview.View())
		m.secondaryPane.GotoTop()
		m.statusBar.BlurCommandBar()
		m.statusBar.ResetCommandBar()

		return m, nil

	// moveDirItemMsg is received any time a file or directory has been moved.
	case moveDirItemMsg:
		m.moveMode = false
		m.initialMoveDirectory = ""
		m.itemToMove = nil

		m.primaryPane.ShowAlternateBorder(false)
		m.dirTree.SetContent(msg)
		m.primaryPane.SetContent(m.dirTree.View())
		err := m.updateStatusBarContent()
		if err != nil {
			return m, m.handleErrorCmd(err)
		}

		return m, nil

	// readFileContentMsg is received when a file is read from.
	case readFileContentMsg:
		switch {
		case msg.code != "":
			m.secondaryPane.GotoTop()
			m.dirTreePreview.SetContent(nil)
			m.renderer.SetContent(msg.code)
			m.renderer.SetImage(nil)
			m.secondaryPane.SetContent(m.renderer.View())
		case msg.pdfContent != "":
			m.secondaryPane.GotoTop()
			m.dirTreePreview.SetContent(nil)
			m.renderer.SetContent(msg.pdfContent)
			m.renderer.SetImage(nil)
			m.secondaryPane.SetContent(m.renderer.View())
		case msg.markdown != "":
			m.secondaryPane.GotoTop()
			m.renderer.SetImage(nil)
			m.dirTreePreview.SetContent(nil)
			m.renderer.SetContent(msg.markdown)
			m.secondaryPane.SetContent(m.renderer.View())
		case msg.image != nil:
			m.secondaryPane.GotoTop()
			m.dirTreePreview.SetContent(nil)
			m.renderer.SetImage(msg.image)
			m.renderer.SetContent(msg.imageString)
			m.secondaryPane.SetContent(m.renderer.View())
		default:
			m.secondaryPane.GotoTop()
			m.secondaryPane.SetContent(msg.rawContent)
		}

		return m, nil

	// convertImageToStringMsg is received when an image is to be converted to a string.
	case convertImageToStringMsg:
		m.renderer.SetContent(string(msg))
		m.secondaryPane.SetContent(m.renderer.View())

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

	// directoryItemSizeMsg is received whenever the directory size needs calculated.
	case directoryItemSizeMsg:
		m.statusBar.SetItemSize(string(msg))

		return m, nil

	// copyToClipboardMsg when the selected directory item is copyied to the clipboard.
	case copyToClipboardMsg:
		m.renderer.SetContent(string(msg))
		m.secondaryPane.SetContent(m.renderer.View())

		return m, nil

	// tea.WindowSizeMsg is received whenever the window size changes.
	case tea.WindowSizeMsg:
		m.primaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
		m.secondaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
		m.dirTree.SetSize(m.primaryPane.GetWidth())
		m.dirTreePreview.SetSize(m.secondaryPane.GetWidth())
		m.renderer.SetSize(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize())
		m.primaryPane.SetContent(m.dirTree.View())
		m.help.Width = msg.Width

		switch {
		case m.renderer.GetImage() != nil:
			cmds = append(cmds, m.redrawImageCmd(m.renderer.GetWidth()))
		case m.renderer.GetContent() != "":
			m.secondaryPane.SetContent(m.renderer.View())
		case m.dirTreePreview.GetTotalFiles() != 0:
			m.secondaryPane.SetContent(m.dirTreePreview.View())
		case m.renderer.GetContent() == "" && m.dirTreePreview.GetTotalFiles() == 0:
			m.secondaryPane.SetContent(lipgloss.NewStyle().
				Width(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize()).
				Render(m.help.View(m.keys)),
			)
		}

		if !m.ready {
			m.ready = true
		}

	// tea.MouseMsg is received whenever a mouse event is triggered.
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				err := m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.getDirectoryItemSizeCmd(selectedFile.Name())
			}

			m.secondaryPane.LineUp(3)

			return m, nil

		case tea.MouseWheelDown:
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				err := m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.getDirectoryItemSizeCmd(selectedFile.Name())
			}

			m.secondaryPane.LineDown(3)

			return m, nil
		}

	case tea.KeyMsg:
		switch {
		// Exit FM.
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit

		// Exit FM if the command bar is not open.
		case key.Matches(msg, m.keys.Quit):
			if !m.showCommandInput {
				return m, tea.Quit
			}

		// If the command bar is not shown and the primary pane is active
		// set the previous directory to the current directory,
		// and update the directory listing to go back one directory.
		case key.Matches(msg, m.keys.Left):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.statusBar.SetItemSize("")
				workingDirectory, err := dirfs.GetWorkingDirectory()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				m.previousDirectory = workingDirectory

				return m, m.updateDirectoryListingCmd(filepath.Join(workingDirectory, dirfs.PreviousDirectory))
			}

		// Scroll pane down.
		case key.Matches(msg, m.keys.Down):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
				m.dirTree.GoDown()
				m.scrollPrimaryPane()
				err := m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.getDirectoryItemSizeCmd(selectedFile.Name())
			}

			m.secondaryPane.LineDown(1)

		// Scroll pane up.
		case key.Matches(msg, m.keys.Up):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
				m.dirTree.GoUp()
				m.scrollPrimaryPane()
				err := m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.getDirectoryItemSizeCmd(selectedFile.Name())
			}

			m.secondaryPane.LineUp(1)

		// Open directory or read file content and display in secondary pane.
		case key.Matches(msg, m.keys.Right):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				switch {
				case selectedFile.IsDir() && !m.statusBar.CommandBarFocused():
					m.statusBar.SetItemSize("")
					currentDir, err := dirfs.GetWorkingDirectory()
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					return m, m.updateDirectoryListingCmd(filepath.Join(currentDir, selectedFile.Name()))
				case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
					m.statusBar.SetItemSize("")
					symlinkFile, err := os.Readlink(selectedFile.Name())
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					fileInfo, err := os.Stat(symlinkFile)
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					if fileInfo.IsDir() {
						currentDir, err := dirfs.GetWorkingDirectory()
						if err != nil {
							return m, m.handleErrorCmd(err)
						}

						return m, m.updateDirectoryListingCmd(filepath.Join(currentDir, fileInfo.Name()))
					}

					return m, m.readFileContentCmd(
						fileInfo,
						m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(),
					)
				default:
					return m, m.readFileContentCmd(
						selectedFile,
						m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(),
					)
				}
			}

		case key.Matches(msg, m.keys.JumpToTop):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.GotoTop()
				m.primaryPane.GotoTop()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.getDirectoryItemSizeCmd(selectedFile.Name())
			}

			m.secondaryPane.GotoTop()

		// Jump to the bottom of a pane.
		case key.Matches(msg, m.keys.JumpToBottom):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.GotoBottom()
				m.primaryPane.GotoBottom()
				m.primaryPane.SetContent(m.dirTree.View())
				m.statusBar.SetItemSize("")
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.getDirectoryItemSizeCmd(selectedFile.Name())
			}

			m.secondaryPane.GotoBottom()

		// process command bar input.
		case key.Matches(msg, m.keys.Enter):
			selectedFile, err := m.dirTree.GetSelectedFile()
			if err != nil {
				return m, m.handleErrorCmd(err)
			}

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
					m.renameDirectoryItemCmd(selectedFile.Name(), m.statusBar.CommandBarValue()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			default:
				return m, nil
			}

		// Delete the currently selected item.
		case key.Matches(msg, m.keys.Delete):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				if selectedFile.IsDir() {
					return m, tea.Sequentially(
						m.deleteDirectoryCmd(selectedFile.Name()),
						m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.deleteFileCmd(selectedFile.Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Enter create file mode.
		case key.Matches(msg, m.keys.CreateFile):
			if !m.moveMode && !m.createDirectoryMode && !m.showCommandInput {
				m.createFileMode = true
				m.showCommandInput = true
				m.statusBar.FocusCommandBar()
				err := m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, nil
			}

		// Enter create directory mode.
		case key.Matches(msg, m.keys.CreateDirectory):
			if !m.moveMode && !m.createFileMode && !m.showCommandInput {
				m.createDirectoryMode = true
				m.showCommandInput = true
				m.statusBar.FocusCommandBar()
				err := m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, nil
			}

		// Enter create directory mode.
		case key.Matches(msg, m.keys.Rename):
			if !m.moveMode && !m.createFileMode && !m.createDirectoryMode && !m.showCommandInput {
				m.renameMode = true
				m.showCommandInput = true
				m.statusBar.FocusCommandBar()
				err := m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, nil
			}

		// Shortcut to get back to the home directory if the
		// command bar is not curently open.
		case key.Matches(msg, m.keys.OpenHomeDirectory):
			if !m.showCommandInput {
				homeDir, err := dirfs.GetHomeDirectory()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, m.updateDirectoryListingCmd(homeDir)
			}

		// Shortcut to go back to the previous directory.
		case key.Matches(msg, m.keys.OpenPreviousDirectory):
			if !m.showCommandInput && m.previousDirectory != "" {
				return m, m.updateDirectoryListingCmd(m.previousDirectory)
			}

		// Shortcut to go back to the root directory.
		case key.Matches(msg, m.keys.OpenRootDirectory):
			if !m.showCommandInput {
				return m, m.updateDirectoryListingCmd(dirfs.RootDirectory)
			}

		// Toggle hidden files and folders.
		case key.Matches(msg, m.keys.ToggleHidden):
			if !m.showCommandInput && m.primaryPane.GetIsActive() {
				m.dirTree.ToggleHidden()
				m.showHidden = !m.showHidden

				switch {
				case m.showDirectoriesOnly:
					return m, m.getDirectoryListingByType("directories", m.showHidden)
				case m.showFilesOnly:
					return m, m.getDirectoryListingByType("files", m.showHidden)
				default:
					return m, m.updateDirectoryListingCmd(dirfs.CurrentDirectory)
				}
			}

		// Toggle between the two panes if the command bar is not currently active.
		case key.Matches(msg, m.keys.Tab):
			if !m.showCommandInput {
				m.primaryPane.SetActive(!m.primaryPane.GetIsActive())
				m.secondaryPane.SetActive(!m.secondaryPane.GetIsActive())
			}

		// Enter move mode.
		case key.Matches(msg, m.keys.EnterMoveMode):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				m.moveMode = true
				m.primaryPane.ShowAlternateBorder(true)
				initialMoveDirectory, err := dirfs.GetWorkingDirectory()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				m.initialMoveDirectory = initialMoveDirectory
				m.itemToMove = selectedFile
				err = m.updateStatusBarContent()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}
			}

		// Zip up the currently selected item.
		case key.Matches(msg, m.keys.Zip):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, tea.Sequentially(
					m.zipDirectoryCmd(selectedFile.Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Unzip the currently selected zip file.
		case key.Matches(msg, m.keys.Unzip):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				return m, tea.Sequentially(
					m.unzipDirectoryCmd(selectedFile.Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Copy the currently selected item.
		case key.Matches(msg, m.keys.Copy):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				if selectedFile.IsDir() {
					return m, tea.Sequentially(
						m.copyDirectoryCmd(selectedFile.Name()),
						m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
					)
				}

				return m, tea.Sequentially(
					m.copyFileCmd(selectedFile.Name()),
					m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}

		// Edit the currently selected file.
		case key.Matches(msg, m.keys.EditFile):
			selectedFile, err := m.dirTree.GetSelectedFile()
			if err != nil {
				return m, m.handleErrorCmd(err)
			}

			if !m.showCommandInput && m.primaryPane.GetIsActive() && !selectedFile.IsDir() {
				editorPath := os.Getenv("EDITOR")
				if editorPath == "" {
					return m, m.handleErrorCmd(errors.New("$EDITOR not set"))
				}

				editorCmd := exec.Command(editorPath, selectedFile.Name())
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

		case key.Matches(msg, m.keys.PreviewDirectory):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				selectedFile, err := m.dirTree.GetSelectedFile()
				if err != nil {
					return m, m.handleErrorCmd(err)
				}

				switch {
				case selectedFile.IsDir() && !m.statusBar.CommandBarFocused():
					return m, m.previewDirectoryListingCmd(selectedFile.Name())
				case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
					symlinkFile, err := os.Readlink(selectedFile.Name())
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					fileInfo, err := os.Stat(symlinkFile)
					if err != nil {
						return m, m.handleErrorCmd(err)
					}

					if fileInfo.IsDir() {
						return m, m.previewDirectoryListingCmd(fileInfo.Name())
					}
				default:
					return m, m.previewDirectoryListingCmd(selectedFile.Name())
				}
			}

		case key.Matches(msg, m.keys.CopyToClipboard):
			selectedFile, err := m.dirTree.GetSelectedFile()
			if err != nil {
				return m, m.handleErrorCmd(err)
			}

			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				return m, m.copyToClipboardCmd(selectedFile.Name())
			}

		case key.Matches(msg, m.keys.ShowOnlyDirectories):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				m.showDirectoriesOnly = !m.showDirectoriesOnly
				m.showFilesOnly = false

				if m.showDirectoriesOnly {
					return m, m.getDirectoryListingByType("directories", m.showHidden)
				}

				return m, m.updateDirectoryListingCmd(dirfs.CurrentDirectory)
			}

		case key.Matches(msg, m.keys.ShowOnlyFiles):
			if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
				m.showFilesOnly = !m.showFilesOnly
				m.showDirectoriesOnly = false

				if m.showFilesOnly {
					return m, m.getDirectoryListingByType("files", m.showHidden)
				}

				return m, m.updateDirectoryListingCmd(dirfs.CurrentDirectory)
			}

		// Reset FM to its initial state.
		case key.Matches(msg, m.keys.Escape):
			m.showCommandInput = false
			m.moveMode = false
			m.itemToMove = nil
			m.initialMoveDirectory = ""
			m.help.ShowAll = true
			m.createFileMode = false
			m.createDirectoryMode = false
			m.renameMode = false
			m.showFilesOnly = false
			m.showHidden = false
			m.showDirectoriesOnly = false
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
			m.renderer.SetImage(nil)
			m.renderer.SetContent("")
			m.dirTreePreview.SetContent(nil)
			err := m.updateStatusBarContent()
			if err != nil {
				return m, m.handleErrorCmd(err)
			}
			cmds = append(cmds, m.updateDirectoryListingCmd(dirfs.CurrentDirectory))
		}
	}

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
