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

// handleUpdateDirectoryListingMsg is received when a directory is read from .
func (m *Model) handleUpdateDirectoryListingMsg(msg updateDirectoryListingMsg) (tea.Model, tea.Cmd) {
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
}

// handlePreviewDirectoryListingMsg is received when a preview of a directory is requested.
func (m *Model) handlePreviewDirectoryListingMsg(msg previewDirectoryListingMsg) (tea.Model, tea.Cmd) {
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
}

// handleMoveDirItemMsg is received any time a file or directory is moved.
func (m *Model) handleMoveDirItemMsg(msg moveDirItemMsg) (tea.Model, tea.Cmd) {
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
}

// handleReadFileContentMsg is received when a file is read from.
func (m *Model) handleReadFileContentMsg(msg readFileContentMsg) (tea.Model, tea.Cmd) {
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
}

// handleConvertImageToStringMsg is received when an image is converted to a string.
func (m *Model) handleConvertImageToStringMsg(msg convertImageToStringMsg) (tea.Model, tea.Cmd) {
	m.renderer.SetContent(string(msg))
	m.secondaryPane.SetContent(m.renderer.View())

	return m, nil
}

// handleErrorMsg is received any time something goes wrong.
func (m *Model) handleErrorMsg(msg errorMsg) (tea.Model, tea.Cmd) {
	m.secondaryPane.SetContent(
		lipgloss.NewStyle().
			Bold(true).
			Foreground(m.theme.ErrorColor).
			Width(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize()).
			Render(string(msg)),
	)

	return m, nil
}

// handleDirectoryItemSizeMsg is received whenever the directory size has been calculated.
func (m *Model) handleDirectoryItemSizeMsg(msg directoryItemSizeMsg) (tea.Model, tea.Cmd) {
	m.statusBar.SetItemSize(string(msg))

	return m, nil
}

// handleCopyToClipboardMsg is received when the selected directory item is copied to the clipboard.
func (m *Model) handleCopyToClipboardMsg(msg copyToClipboardMsg) (tea.Model, tea.Cmd) {
	m.renderer.SetContent(string(msg))
	m.secondaryPane.SetContent(m.renderer.View())

	return m, nil
}

// handleWindowSizeMsg is received whenever the window size changes.
func (m *Model) handleWindowSizeMsg(msg tea.WindowSizeMsg, cmds *[]tea.Cmd) {
	m.primaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
	m.secondaryPane.SetSize(msg.Width/2, msg.Height-m.statusBar.GetHeight())
	m.dirTree.SetSize(m.primaryPane.GetWidth())
	m.dirTreePreview.SetSize(m.secondaryPane.GetWidth())
	m.renderer.SetSize(m.secondaryPane.GetWidth() - m.secondaryPane.GetHorizontalFrameSize())
	m.primaryPane.SetContent(m.dirTree.View())
	m.help.Width = msg.Width

	switch {
	case m.renderer.GetImage() != nil:
		*cmds = append(*cmds, m.redrawImageCmd(m.renderer.GetWidth()))
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
}

// handleMouseMsg is received whenever a mouse event is triggered.
func (m *Model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
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

	return m, nil
}

// handleLeftKeyPress goes back to the previous directory when pressed.
func (m *Model) handleLeftKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() {
		m.statusBar.SetItemSize("")
		workingDirectory, err := dirfs.GetWorkingDirectory()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		m.previousDirectory = workingDirectory

		*cmds = append(*cmds, m.updateDirectoryListingCmd(filepath.Join(workingDirectory, dirfs.PreviousDirectory)))
	}
}

// handleDownKeyPress goes down one in the directory tree.
func (m *Model) handleDownKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
		m.dirTree.GoDown()
		m.scrollPrimaryPane()

		err := m.updateStatusBarContent()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		m.primaryPane.SetContent(m.dirTree.View())
		m.statusBar.SetItemSize("")

		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, m.getDirectoryItemSizeCmd(selectedFile.Name()))
	}

	if !m.showCommandInput && m.secondaryPane.GetIsActive() {
		m.secondaryPane.LineDown(1)
	}
}

// handleUpKeyPress goes up in the directory tree or scrolls the secondary pane up.
func (m *Model) handleUpKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
		m.dirTree.GoUp()
		m.scrollPrimaryPane()
		err := m.updateStatusBarContent()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}
		m.primaryPane.SetContent(m.dirTree.View())
		m.statusBar.SetItemSize("")
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, m.getDirectoryItemSizeCmd(selectedFile.Name()))
	}

	if !m.showCommandInput && m.secondaryPane.GetIsActive() {
		m.secondaryPane.LineUp(1)
	}
}

// handleRightKeyPress opens directory if it is one or reads a files content.
func (m *Model) handleRightKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		switch {
		case selectedFile.IsDir() && !m.statusBar.CommandBarFocused():
			m.statusBar.SetItemSize("")
			currentDir, err := dirfs.GetWorkingDirectory()
			if err != nil {
				*cmds = append(*cmds, m.handleErrorCmd(err))
			}

			*cmds = append(*cmds, m.updateDirectoryListingCmd(filepath.Join(currentDir, selectedFile.Name())))
		case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
			m.statusBar.SetItemSize("")
			symlinkFile, err := os.Readlink(selectedFile.Name())
			if err != nil {
				*cmds = append(*cmds, m.handleErrorCmd(err))
			}

			fileInfo, err := os.Stat(symlinkFile)
			if err != nil {
				*cmds = append(*cmds, m.handleErrorCmd(err))
			}

			if fileInfo.IsDir() {
				currentDir, err := dirfs.GetWorkingDirectory()
				if err != nil {
					*cmds = append(*cmds, m.handleErrorCmd(err))
				}

				*cmds = append(*cmds, m.updateDirectoryListingCmd(filepath.Join(currentDir, fileInfo.Name())))
			}

			*cmds = append(*cmds, m.readFileContentCmd(
				fileInfo,
				m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(),
			))
		default:
			*cmds = append(*cmds, m.readFileContentCmd(
				selectedFile,
				m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(),
			))
		}
	}
}

// handleJumpToTopKeyPress jumps to the top of a pane.
func (m *Model) handleJumpToTopKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
		m.dirTree.GotoTop()
		m.primaryPane.GotoTop()
		m.primaryPane.SetContent(m.dirTree.View())
		m.statusBar.SetItemSize("")
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, m.getDirectoryItemSizeCmd(selectedFile.Name()))
	}

	if !m.showCommandInput && m.secondaryPane.GetIsActive() {
		m.secondaryPane.GotoTop()
	}
}

// handleJumpToBottomKeyPress jumps to the bottom of a pane.
func (m *Model) handleJumpToBottomKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 1 {
		m.dirTree.GotoBottom()
		m.primaryPane.GotoBottom()
		m.primaryPane.SetContent(m.dirTree.View())
		m.statusBar.SetItemSize("")
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, m.getDirectoryItemSizeCmd(selectedFile.Name()))
	}

	if !m.showCommandInput && m.secondaryPane.GetIsActive() {
		m.secondaryPane.GotoBottom()
	}
}

// handleEnterKeyPress processes command input.
func (m *Model) handleEnterKeyPress(cmds *[]tea.Cmd) {
	selectedFile, err := m.dirTree.GetSelectedFile()
	if err != nil {
		*cmds = append(*cmds, m.handleErrorCmd(err))
	}

	switch {
	case m.moveMode:
		*cmds = append(*cmds, m.moveDirectoryItemCmd(m.itemToMove.Name()))
	case m.createFileMode:
		*cmds = append(*cmds, tea.Sequentially(
			m.createFileCmd(m.statusBar.CommandBarValue()),
			m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
		))
	case m.createDirectoryMode:
		*cmds = append(*cmds, tea.Sequentially(
			m.createDirectoryCmd(m.statusBar.CommandBarValue()),
			m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
		))
	case m.renameMode:
		*cmds = append(*cmds, tea.Sequentially(
			m.renameDirectoryItemCmd(selectedFile.Name(), m.statusBar.CommandBarValue()),
			m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
		))
	default:
		return
	}
}

// handleDeleteKeyPress deletes the selected directory item.
func (m *Model) handleDeleteKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() {
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		if selectedFile.IsDir() {
			*cmds = append(*cmds, tea.Sequentially(
				m.deleteDirectoryCmd(selectedFile.Name()),
				m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
			))
		}

		*cmds = append(*cmds, tea.Sequentially(
			m.deleteFileCmd(selectedFile.Name()),
			m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
		))
	}
}

// handleCreateFileKeyPress creates a new file.
func (m *Model) handleCreateFileKeyPress(cmds *[]tea.Cmd) {
	m.createFileMode = true
	m.showCommandInput = true
	m.statusBar.FocusCommandBar()

	err := m.updateStatusBarContent()
	if err != nil {
		*cmds = append(*cmds, m.handleErrorCmd(err))
	}
}

// handleCreateDirectoryKeyPress creates a new directory.
func (m *Model) handleCreateDirectoryKeyPress(cmds *[]tea.Cmd) {
	m.createDirectoryMode = true
	m.showCommandInput = true

	m.statusBar.FocusCommandBar()

	err := m.updateStatusBarContent()
	if err != nil {
		*cmds = append(*cmds, m.handleErrorCmd(err))
	}
}

// handleRenameKeyPress renames the selected directory item.
func (m *Model) handleRenameKeyPress(cmds *[]tea.Cmd) {
	m.renameMode = true
	m.showCommandInput = true

	m.statusBar.FocusCommandBar()

	err := m.updateStatusBarContent()
	if err != nil {
		*cmds = append(*cmds, m.handleErrorCmd(err))
	}
}

// handleOpenHomeDirectoryKeyPress opens the home directory.
func (m *Model) handleOpenHomeDirectoryKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput {
		homeDir, err := dirfs.GetHomeDirectory()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, m.updateDirectoryListingCmd(homeDir))
	}
}

// handleOpenPreviousDirectoryKeyPress opens the previous directory.
func (m *Model) handleOpenPreviousDirectoryKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.previousDirectory != "" {
		*cmds = append(*cmds, m.updateDirectoryListingCmd(m.previousDirectory))
	}
}

// handleOpenRootDirectoryKeyPress opens the root directory.
func (m *Model) handleOpenRootDirectoryKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput {
		*cmds = append(*cmds, m.updateDirectoryListingCmd(dirfs.RootDirectory))
	}
}

// handleToggleHiddenKeyPress toggles between hidden files and directories.
func (m *Model) handleToggleHiddenKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() {
		m.dirTree.ToggleHidden()
		m.showHidden = !m.showHidden

		switch {
		case m.showDirectoriesOnly:
			*cmds = append(*cmds, m.getDirectoryListingByType("directories", m.showHidden))
		case m.showFilesOnly:
			*cmds = append(*cmds, m.getDirectoryListingByType("files", m.showHidden))
		default:
			*cmds = append(*cmds, m.updateDirectoryListingCmd(dirfs.CurrentDirectory))
		}
	}
}

// handleTabKeyPress switches between panes.
func (m *Model) handleTabKeyPress() {
	if !m.showCommandInput {
		m.primaryPane.SetActive(!m.primaryPane.GetIsActive())
		m.secondaryPane.SetActive(!m.secondaryPane.GetIsActive())
	}
}

// handleEnterMoveModeKeyPress enters move mode.
func (m *Model) handleEnterMoveModeKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		m.moveMode = true
		m.primaryPane.ShowAlternateBorder(true)
		initialMoveDirectory, err := dirfs.GetWorkingDirectory()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		m.initialMoveDirectory = initialMoveDirectory
		m.itemToMove = selectedFile
		err = m.updateStatusBarContent()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}
	}
}

// handleZipKeyPress zips the selected directory item.
func (m *Model) handleZipKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, tea.Sequentially(
			m.zipDirectoryCmd(selectedFile.Name()),
			m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
		))
	}
}

// handleUnzipKeyPress unzips the selected directory item.
func (m *Model) handleUnzipKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, tea.Sequentially(
			m.unzipDirectoryCmd(selectedFile.Name()),
			m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
		))
	}
}

// handleCopyKeyPress copies the selected directory item.
func (m *Model) handleCopyKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		if selectedFile.IsDir() {
			*cmds = append(*cmds, tea.Sequentially(
				m.copyDirectoryCmd(selectedFile.Name()),
				m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
			))
		}

		*cmds = append(*cmds, tea.Sequentially(
			m.copyFileCmd(selectedFile.Name()),
			m.updateDirectoryListingCmd(dirfs.CurrentDirectory),
		))
	}
}

// handleEditFileKeyPress edits the selected directory item.
func (m *Model) handleEditFileKeyPress(cmds *[]tea.Cmd) {
	selectedFile, err := m.dirTree.GetSelectedFile()
	if err != nil {
		*cmds = append(*cmds, m.handleErrorCmd(err))
	}

	if !m.showCommandInput && m.primaryPane.GetIsActive() && !selectedFile.IsDir() {
		editorPath := os.Getenv("EDITOR")
		if editorPath == "" {
			*cmds = append(*cmds, m.handleErrorCmd(errors.New("$EDITOR not set")))
		}

		editorCmd := exec.Command(editorPath, selectedFile.Name())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		err := editorCmd.Start()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		err = editorCmd.Wait()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		*cmds = append(*cmds, m.updateDirectoryListingCmd(dirfs.CurrentDirectory))
	}
}

// handlePreviewDirectoryKeyPress previews the selected directory item.
func (m *Model) handlePreviewDirectoryKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		selectedFile, err := m.dirTree.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, m.handleErrorCmd(err))
		}

		switch {
		case selectedFile.IsDir() && !m.statusBar.CommandBarFocused():
			*cmds = append(*cmds, m.previewDirectoryListingCmd(selectedFile.Name()))
		case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
			symlinkFile, err := os.Readlink(selectedFile.Name())
			if err != nil {
				*cmds = append(*cmds, m.handleErrorCmd(err))
			}

			fileInfo, err := os.Stat(symlinkFile)
			if err != nil {
				*cmds = append(*cmds, m.handleErrorCmd(err))
			}

			if fileInfo.IsDir() {
				*cmds = append(*cmds, m.previewDirectoryListingCmd(fileInfo.Name()))
			}
		default:
			*cmds = append(*cmds, m.previewDirectoryListingCmd(selectedFile.Name()))
		}
	}
}

// handleCopyToClipboardKeyPres copies the selected directory item to the clipboard.
func (m *Model) handleCopyToClipboardKeyPress(cmds *[]tea.Cmd) {
	selectedFile, err := m.dirTree.GetSelectedFile()
	if err != nil {
		*cmds = append(*cmds, m.handleErrorCmd(err))
	}

	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		*cmds = append(*cmds, m.copyToClipboardCmd(selectedFile.Name()))
	}
}

// handleShowOnlyDirectoriesKeyPress shows only directories in the directory tree.
func (m *Model) handleShowOnlyDirectoriesKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		m.showDirectoriesOnly = !m.showDirectoriesOnly
		m.showFilesOnly = false

		if m.showDirectoriesOnly {
			*cmds = append(*cmds, m.getDirectoryListingByType("directories", m.showHidden))
		}

		*cmds = append(*cmds, m.updateDirectoryListingCmd(dirfs.CurrentDirectory))
	}
}

// handleShowOnlyFilesKeyPress shows only files in the directory tree.
func (m *Model) handleShowOnlyFilesKeyPress(cmds *[]tea.Cmd) {
	if !m.showCommandInput && m.primaryPane.GetIsActive() && m.dirTree.GetTotalFiles() > 0 {
		m.showFilesOnly = !m.showFilesOnly
		m.showDirectoriesOnly = false

		if m.showFilesOnly {
			*cmds = append(*cmds, m.getDirectoryListingByType("files", m.showHidden))
		}

		*cmds = append(*cmds, m.updateDirectoryListingCmd(dirfs.CurrentDirectory))
	}
}

// handleEscapeKeyPress resets FM to its initial state.
func (m *Model) handleEscapeKeyPress(cmds *[]tea.Cmd) {
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
		*cmds = append(*cmds, m.handleErrorCmd(err))
	}

	*cmds = append(*cmds, m.updateDirectoryListingCmd(dirfs.CurrentDirectory))
}

// Update handles all UI interactions and events for updating the screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case updateDirectoryListingMsg:
		return m.handleUpdateDirectoryListingMsg(msg)
	case previewDirectoryListingMsg:
		return m.handlePreviewDirectoryListingMsg(msg)
	case moveDirItemMsg:
		return m.handleMoveDirItemMsg(msg)
	case readFileContentMsg:
		return m.handleReadFileContentMsg(msg)
	case convertImageToStringMsg:
		return m.handleConvertImageToStringMsg(msg)
	case errorMsg:
		return m.handleErrorMsg(msg)
	case directoryItemSizeMsg:
		return m.handleDirectoryItemSizeMsg(msg)
	case copyToClipboardMsg:
		return m.handleCopyToClipboardMsg(msg)
	case tea.WindowSizeMsg:
		m.handleWindowSizeMsg(msg, &cmds)
	case tea.MouseMsg:
		return m.handleMouseMsg(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Quit):
			if !m.showCommandInput {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keys.Left):
			m.handleLeftKeyPress(&cmds)
		case key.Matches(msg, m.keys.Down):
			m.handleDownKeyPress(&cmds)
		case key.Matches(msg, m.keys.Up):
			m.handleUpKeyPress(&cmds)
		case key.Matches(msg, m.keys.Right):
			m.handleRightKeyPress(&cmds)
		case key.Matches(msg, m.keys.JumpToTop):
			m.handleJumpToTopKeyPress(&cmds)
		case key.Matches(msg, m.keys.JumpToBottom):
			m.handleJumpToBottomKeyPress(&cmds)
		case key.Matches(msg, m.keys.Enter):
			m.handleEnterKeyPress(&cmds)
		case key.Matches(msg, m.keys.Delete):
			m.handleDeleteKeyPress(&cmds)
		case key.Matches(msg, m.keys.CreateFile):
			if !m.moveMode && !m.createDirectoryMode && !m.showCommandInput {
				m.handleCreateFileKeyPress(&cmds)
				return m, nil
			}
		case key.Matches(msg, m.keys.CreateDirectory):
			if !m.moveMode && !m.createFileMode && !m.showCommandInput {
				m.handleCreateDirectoryKeyPress(&cmds)
				return m, nil
			}
		case key.Matches(msg, m.keys.Rename):
			if !m.moveMode && !m.createFileMode && !m.createDirectoryMode && !m.showCommandInput {
				m.handleRenameKeyPress(&cmds)
				return m, nil
			}
		case key.Matches(msg, m.keys.OpenHomeDirectory):
			m.handleOpenHomeDirectoryKeyPress(&cmds)
		case key.Matches(msg, m.keys.OpenPreviousDirectory):
			m.handleOpenPreviousDirectoryKeyPress(&cmds)
		case key.Matches(msg, m.keys.OpenRootDirectory):
			m.handleOpenRootDirectoryKeyPress(&cmds)
		case key.Matches(msg, m.keys.ToggleHidden):
			m.handleToggleHiddenKeyPress(&cmds)
		case key.Matches(msg, m.keys.Tab):
			m.handleTabKeyPress()
		case key.Matches(msg, m.keys.EnterMoveMode):
			m.handleEnterMoveModeKeyPress(&cmds)
		case key.Matches(msg, m.keys.Zip):
			m.handleZipKeyPress(&cmds)
		case key.Matches(msg, m.keys.Unzip):
			m.handleUnzipKeyPress(&cmds)
		case key.Matches(msg, m.keys.Copy):
			m.handleCopyKeyPress(&cmds)
		case key.Matches(msg, m.keys.EditFile):
			m.handleEditFileKeyPress(&cmds)
		case key.Matches(msg, m.keys.PreviewDirectory):
			m.handlePreviewDirectoryKeyPress(&cmds)
		case key.Matches(msg, m.keys.CopyToClipboard):
			m.handleCopyToClipboardKeyPress(&cmds)
		case key.Matches(msg, m.keys.ShowOnlyDirectories):
			m.handleShowOnlyDirectoriesKeyPress(&cmds)
		case key.Matches(msg, m.keys.ShowOnlyFiles):
			m.handleShowOnlyFilesKeyPress(&cmds)
		case key.Matches(msg, m.keys.Escape):
			m.handleEscapeKeyPress(&cmds)
		}
	}

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
