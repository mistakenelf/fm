package tui

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/internal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/viper"
)

// checkPrimaryViewportBounds handles wrapping of the filetree and
// scrolling of the viewport.
func (b *Bubble) checkPrimaryViewportBounds() {
	top := b.primaryViewport.YOffset
	bottom := b.primaryViewport.Height + b.primaryViewport.YOffset - 1

	if b.treeCursor < top {
		b.primaryViewport.LineUp(1)
	} else if b.treeCursor > bottom {
		b.primaryViewport.LineDown(1)
	}

	if b.treeCursor > len(b.treeFiles)-1 {
		b.primaryViewport.GotoTop()
		b.treeCursor = 0
	} else if b.treeCursor < top {
		b.primaryViewport.GotoBottom()
		b.treeCursor = len(b.treeFiles) - 1
	}
}

// writeLog writes a message to the log.
func (b *Bubble) writeLog(msg string) {
	if b.appConfig.Settings.EnableLogging {
		b.logs = append(b.logs, msg)

		if b.showLogs {
			bottom := b.secondaryViewport.Height + b.secondaryViewport.YOffset - 1
			if lipgloss.Height(b.logView()) > bottom {
				b.secondaryViewport.GotoBottom()
			}
			b.secondaryViewport.SetContent(b.logView())
		}
	}
}

// handleKeys handles all keypresses.
func (b *Bubble) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Jump to top of box.
	if msg.String() == "g" && b.previousKey.String() == "g" {
		if !b.showCommandInput && b.activeBox == PrimaryBoxActive && !b.showBoxSpinner {
			b.treeCursor = 0
			b.primaryViewport.GotoTop()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}

		if !b.showCommandInput && b.activeBox == SecondaryBoxActive {
			b.secondaryViewport.GotoTop()
		}

		return nil
	}

	// Reload config file.
	if msg.String() == "c" && b.previousKey.String() == "r" {
		if !b.showCommandInput && b.activeBox == PrimaryBoxActive && !b.showBoxSpinner {
			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					log.Fatal(err)
				}
			}

			b.appConfig = config.GetConfig()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}

		return nil
	}

	switch {
	case key.Matches(msg, b.keyMap.Quit):
		return tea.Quit
	case key.Matches(msg, b.keyMap.Down):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.treeCursor++
			b.checkPrimaryViewportBounds()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}
	case key.Matches(msg, b.keyMap.Up):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.treeCursor--
			b.checkPrimaryViewportBounds()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}
	case key.Matches(msg, b.keyMap.Left):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.treeCursor = 0
			b.showFilesOnly = false
			b.showDirectoriesOnly = false
			b.foundFilesPaths = nil
			workingDirectory, err := dirfs.GetWorkingDirectory()
			if err != nil {
				return b.handleErrorCmd(err)
			}

			return b.updateDirectoryListingCmd(
				filepath.Join(workingDirectory, dirfs.PreviousDirectory),
			)
		}
	case key.Matches(msg, b.keyMap.Right):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			selectedFile, err := b.treeFiles[b.treeCursor].Info()
			if err != nil {
				return b.handleErrorCmd(err)
			}

			switch {
			case selectedFile.IsDir():
				currentDir, err := dirfs.GetWorkingDirectory()
				if err != nil {
					return b.handleErrorCmd(err)
				}

				directoryToOpen := filepath.Join(currentDir, selectedFile.Name())

				if len(b.foundFilesPaths) > 0 {
					directoryToOpen = b.foundFilesPaths[b.treeCursor]
				}

				return b.updateDirectoryListingCmd(directoryToOpen)
			case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
				symlinkFile, err := os.Readlink(selectedFile.Name())
				if err != nil {
					return b.handleErrorCmd(err)
				}

				fileInfo, err := os.Stat(symlinkFile)
				if err != nil {
					return b.handleErrorCmd(err)
				}

				if fileInfo.IsDir() {
					currentDir, err := dirfs.GetWorkingDirectory()
					if err != nil {
						return b.handleErrorCmd(err)
					}

					return b.updateDirectoryListingCmd(filepath.Join(currentDir, fileInfo.Name()))
				}

				return b.readFileContentCmd(
					fileInfo.Name(),
					b.secondaryViewport.Width,
				)

			default:
				fileToRead := selectedFile.Name()

				if len(b.foundFilesPaths) > 0 {
					fileToRead = b.foundFilesPaths[b.treeCursor]
				}

				return b.readFileContentCmd(
					fileToRead,
					b.secondaryViewport.Width,
				)
			}
		}
	case key.Matches(msg, b.keyMap.Preview):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			selectedFile, err := b.treeFiles[b.treeCursor].Info()
			if err != nil {
				return b.handleErrorCmd(err)
			}

			switch {
			case selectedFile.IsDir():
				return b.previewDirectoryListingCmd(selectedFile.Name())
			case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
				symlinkFile, err := os.Readlink(selectedFile.Name())
				if err != nil {
					return b.handleErrorCmd(err)
				}

				fileInfo, err := os.Stat(symlinkFile)
				if err != nil {
					return b.handleErrorCmd(err)
				}

				if fileInfo.IsDir() {
					return b.previewDirectoryListingCmd(fileInfo.Name())
				}
			default:
				return nil
			}
		}
	case key.Matches(msg, b.keyMap.GotoBottom):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.treeCursor = len(b.treeFiles) - 1
			b.primaryViewport.GotoBottom()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}

		if b.activeBox == SecondaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.secondaryViewport.GotoBottom()
		}
	case key.Matches(msg, b.keyMap.HomeShortcut):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.treeCursor = 0
			b.fileSizes = nil
			homeDir, err := dirfs.GetHomeDirectory()
			if err != nil {
				return b.handleErrorCmd(err)
			}

			return b.updateDirectoryListingCmd(homeDir)
		}
	case key.Matches(msg, b.keyMap.RootShortcut):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.treeCursor = 0
			b.fileSizes = nil

			return b.updateDirectoryListingCmd(dirfs.RootDirectory)
		}
	case key.Matches(msg, b.keyMap.ToggleHidden):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.showHiddenFiles = !b.showHiddenFiles

			switch {
			case b.showDirectoriesOnly:
				return b.getDirectoryListingByTypeCmd(dirfs.DirectoriesListingType)
			case b.showFilesOnly:
				return b.getDirectoryListingByTypeCmd(dirfs.FilesListingType)
			default:
				return b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
			}
		}
	case key.Matches(msg, b.keyMap.ShowDirectoriesOnly):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.showDirectoriesOnly = !b.showDirectoriesOnly
			b.showFilesOnly = false

			if b.showDirectoriesOnly {
				return b.getDirectoryListingByTypeCmd(dirfs.DirectoriesListingType)
			}

			return b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
		}
	case key.Matches(msg, b.keyMap.ShowFilesOnly):
		if b.activeBox == PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
			b.showFilesOnly = !b.showFilesOnly
			b.showDirectoriesOnly = false

			if b.showFilesOnly {
				return b.getDirectoryListingByTypeCmd(dirfs.FilesListingType)
			}

			b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
		}
	case key.Matches(msg, b.keyMap.CopyPathToClipboard):
		if b.activeBox == PrimaryBoxActive && len(b.treeFiles) > 0 && !b.showCommandInput && !b.showBoxSpinner {
			selectedFile := b.treeFiles[b.treeCursor]

			return b.copyToClipboardCmd(selectedFile.Name())
		}
	case key.Matches(msg, b.keyMap.Zip):
		if b.activeBox == PrimaryBoxActive && len(b.treeFiles) > 0 && !b.showCommandInput && !b.showBoxSpinner {
			selectedFile := b.treeFiles[b.treeCursor]

			return tea.Sequentially(
				b.zipDirectoryCmd(selectedFile.Name()),
				b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
			)
		}
	case key.Matches(msg, b.keyMap.Unzip):
		if b.activeBox == PrimaryBoxActive && len(b.treeFiles) > 0 && !b.showCommandInput && !b.showBoxSpinner {
			selectedFile := b.treeFiles[b.treeCursor]

			return tea.Sequentially(
				b.unzipDirectoryCmd(selectedFile.Name()),
				b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
			)
		}
	case key.Matches(msg, b.keyMap.NewFile):
		if !b.showCommandInput && !b.showBoxSpinner {
			b.createFileMode = true
			b.showCommandInput = true
			b.textinput.Placeholder = "Enter file name"
			b.textinput.Focus()

			return textinput.Blink
		}
	case key.Matches(msg, b.keyMap.NewDirectory):
		if !b.showCommandInput && !b.showBoxSpinner {
			b.createDirectoryMode = true
			b.showCommandInput = true
			b.textinput.Placeholder = "Enter directory name"
			b.textinput.Focus()

			return textinput.Blink
		}
	case key.Matches(msg, b.keyMap.Delete):
		if !b.showCommandInput && !b.showBoxSpinner {
			b.deleteMode = true
			b.showCommandInput = true
			b.textinput.Placeholder = "Are you sure you want to delete this? (y/n)"
			b.textinput.Focus()
		}
	case key.Matches(msg, b.keyMap.Move):
		if !b.showCommandInput && !b.showBoxSpinner {
			b.moveMode = true
			b.treeItemToMove = b.treeFiles[b.treeCursor]
			workingDir, err := dirfs.GetWorkingDirectory()
			if err != nil {
				b.moveInitiatedDirectory = dirfs.CurrentDirectory
			}

			b.moveInitiatedDirectory = workingDir

			return nil
		}
	case key.Matches(msg, b.keyMap.Enter):
		switch {
		case b.moveMode:
			return b.moveDirectoryItemCmd(b.treeItemToMove.Name())
		case b.createFileMode:
			return tea.Sequentially(
				b.createFileCmd(b.textinput.Value()),
				b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
			)
		case b.createDirectoryMode:
			return tea.Sequentially(
				b.createDirectoryCmd(b.textinput.Value()),
				b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
			)
		case b.renameMode:
			selectedFile := b.treeFiles[b.treeCursor]

			return tea.Sequentially(
				b.renameDirectoryItemCmd(selectedFile.Name(), b.textinput.Value()),
				b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
			)
		case b.findMode:
			b.showCommandInput = false
			b.showBoxSpinner = true

			return b.findFilesByNameCmd(b.textinput.Value())
		case b.deleteMode:
			selectedFile := b.treeFiles[b.treeCursor]

			if strings.ToLower(b.textinput.Value()) == "y" || strings.ToLower(b.textinput.Value()) == "yes" {
				if selectedFile.IsDir() {
					return tea.Sequentially(
						b.deleteDirectoryCmd(selectedFile.Name()),
						b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
					)
				} else {
					return tea.Sequentially(
						b.deleteFileCmd(selectedFile.Name()),
						b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
					)
				}
			}
		default:
			return nil
		}
	case key.Matches(msg, b.keyMap.Edit):
		selectedFile := b.treeFiles[b.treeCursor]

		if !b.showCommandInput && b.activeBox == PrimaryBoxActive && !b.showBoxSpinner {
			selectionPath := viper.GetString("selection-path")

			if selectionPath == "" && !selectedFile.IsDir() {
				editorPath := os.Getenv("EDITOR")
				if editorPath == "" {
					return b.handleErrorCmd(errors.New("$EDITOR not set"))
				}

				editorCmd := exec.Command(editorPath, selectedFile.Name())
				editorCmd.Stdin = os.Stdin
				editorCmd.Stdout = os.Stdout
				editorCmd.Stderr = os.Stderr

				err := editorCmd.Run()
				termenv.AltScreen()

				if err != nil {
					return b.handleErrorCmd(err)
				}

				return tea.Batch(b.redrawCmd(), tea.HideCursor)
			} else {
				return tea.Sequentially(
					b.writeSelectionPathCmd(selectionPath, selectedFile.Name()),
					tea.Quit,
				)
			}
		}
	case key.Matches(msg, b.keyMap.Copy):
		if !b.showCommandInput && b.activeBox == PrimaryBoxActive && len(b.treeFiles) > 0 && !b.showBoxSpinner {
			selectedFile := b.treeFiles[b.treeCursor]

			if selectedFile.IsDir() {
				return tea.Sequentially(
					b.copyDirectoryCmd(selectedFile.Name()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			} else {
				return tea.Sequentially(
					b.copyFileCmd(selectedFile.Name()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}
		}
	case key.Matches(msg, b.keyMap.Find):
		if !b.showCommandInput && !b.showBoxSpinner {
			b.findMode = true
			b.showCommandInput = true
			b.textinput.Placeholder = "Enter a search term"
			b.textinput.Focus()

			return textinput.Blink
		}
	case key.Matches(msg, b.keyMap.Rename):
		if b.activeBox == PrimaryBoxActive && !b.showBoxSpinner && len(b.treeFiles) > 0 {
			selectedFile := b.treeFiles

			if selectedFile != nil {
				b.renameMode = true
				b.showCommandInput = true
				b.textinput.Placeholder = "Enter new name"
				b.textinput.Focus()

				return textinput.Blink
			}
		}
	case key.Matches(msg, b.keyMap.Escape):
		b.showCommandInput = false
		b.moveMode = false
		b.createFileMode = false
		b.createDirectoryMode = false
		b.renameMode = false
		b.showFilesOnly = false
		b.showHiddenFiles = true
		b.showDirectoriesOnly = false
		b.findMode = false
		b.deleteMode = false
		b.errorMsg = ""
		b.showHelp = true
		b.showLogs = false
		b.foundFilesPaths = nil
		b.showBoxSpinner = false
		b.currentImage = nil
		b.secondaryViewport.GotoTop()
		b.secondaryViewport.SetContent(b.helpView())
		b.textinput.Blur()
		b.textinput.Reset()

		return b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
	case key.Matches(msg, b.keyMap.ShowLogs):
		if !b.showCommandInput && b.appConfig.Settings.EnableLogging {
			b.showLogs = true
			b.currentImage = nil
			b.secondaryViewport.SetContent(b.logView())
		}
	case key.Matches(msg, b.keyMap.ToggleBox):
		b.activeBox = (b.activeBox + 1) % 2
	}

	b.previousKey = msg

	if b.activeBox != PrimaryBoxActive {
		b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	b.textinput, cmd = b.textinput.Update(msg)
	cmds = append(cmds, cmd)

	b.spinner, cmd = b.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// handleMouse handles all mouse interaction.
func (b *Bubble) handleMouse(msg tea.MouseMsg) {
	switch msg.Type {
	case tea.MouseWheelUp:
		if b.activeBox == PrimaryBoxActive {
			b.treeCursor--
			b.checkPrimaryViewportBounds()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}

		if b.activeBox == SecondaryBoxActive {
			b.secondaryViewport.LineUp(1)
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}
	case tea.MouseWheelDown:
		if b.activeBox == PrimaryBoxActive {
			b.treeCursor++
			b.checkPrimaryViewportBounds()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}

		if b.activeBox == SecondaryBoxActive {
			b.secondaryViewport.LineDown(1)
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}
	}
}

// Update handles all UI interactions and events for updating the screen.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case updateDirectoryListingMsg:
		b.showCommandInput = false
		b.createFileMode = false
		b.createDirectoryMode = false
		b.deleteMode = false
		b.renameMode = false
		b.treeCursor = 0
		b.treeFiles = msg
		b.showFileTreePreview = false
		b.fileSizes = make([]string, len(msg))
		b.writeLog("Directory listing updated.")

		for i, file := range msg {
			cmds = append(cmds, b.getDirectoryItemSizeCmd(file.Name(), i))
		}

		b.primaryViewport.SetContent(b.fileTreeView(msg))
		b.textinput.Blur()
		b.textinput.Reset()

		return b, tea.Batch(cmds...)
	case directoryItemSizeMsg:
		if len(b.fileSizes) > 0 && msg.index < len(b.fileSizes) {
			b.fileSizes[msg.index] = msg.size
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}

		return b, nil
	case readFileContentMsg:
		b.showFileTreePreview = false
		b.showHelp = false
		b.showLogs = false
		b.currentImage = nil
		b.secondaryViewport.GotoTop()

		switch {
		case msg.code != "":
			b.secondaryBoxContent = msg.code
		case msg.pdfContent != "":
			b.secondaryBoxContent = msg.pdfContent
		case msg.markdown != "":
			b.secondaryBoxContent = msg.markdown
		case msg.image != nil:
			b.currentImage = msg.image
			b.secondaryBoxContent = msg.imageString
		default:
			b.secondaryBoxContent = msg.rawContent
		}

		b.secondaryViewport.SetContent(b.textContentView(b.secondaryBoxContent))

		return b, nil
	case previewDirectoryListingMsg:
		b.showFileTreePreview = true
		b.showHelp = false
		b.showLogs = false
		b.treePreviewFiles = msg
		b.secondaryViewport.GotoTop()
		b.secondaryViewport.SetContent(b.fileTreePreviewView(msg))

		return b, nil
	case convertImageToStringMsg:
		b.showHelp = false
		b.showLogs = false
		b.secondaryViewport.GotoTop()
		b.secondaryViewport.SetContent(b.textContentView(string(msg)))

		return b, nil
	case copyToClipboardMsg:
		b.secondaryViewport.SetContent(b.textContentView(string(msg)))

		return b, nil
	case moveDirItemMsg:
		b.moveMode = false
		b.treeItemToMove = nil
		b.moveInitiatedDirectory = ""
		b.treeFiles = msg
		b.treeCursor = 0
		b.fileSizes = make([]string, len(msg))

		for i, file := range msg {
			cmds = append(cmds, b.getDirectoryItemSizeCmd(file.Name(), i))
		}

		b.primaryViewport.SetContent(b.fileTreeView(msg))

		return b, tea.Batch(cmds...)
	case findFilesByNameMsg:
		b.showCommandInput = false
		b.createFileMode = false
		b.createDirectoryMode = false
		b.renameMode = false
		b.findMode = false
		b.treeCursor = 0
		b.treeFiles = msg.entries
		b.foundFilesPaths = msg.paths
		b.showBoxSpinner = false
		b.textinput.Blur()
		b.textinput.Reset()
		b.fileSizes = make([]string, len(msg.entries))

		for i, file := range msg.entries {
			cmds = append(cmds, b.getDirectoryItemSizeCmd(file.Name(), i))
		}

		b.primaryViewport.SetContent(b.fileTreeView(msg.entries))

		return b, tea.Batch(cmds...)
	case errorMsg:
		b.showHelp = false
		b.showLogs = false
		b.errorMsg = string(msg)
		b.secondaryViewport.SetContent(b.errorView(string(msg)))

		return b, nil
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height

		b.primaryViewport.Width = (msg.Width / 2) - b.primaryViewport.Style.GetHorizontalFrameSize()
		b.primaryViewport.Height = msg.Height - StatusBarHeight - b.primaryViewport.Style.GetVerticalFrameSize()
		b.secondaryViewport.Width = (msg.Width / 2) - b.secondaryViewport.Style.GetHorizontalFrameSize()
		b.secondaryViewport.Height = msg.Height - StatusBarHeight - b.secondaryViewport.Style.GetVerticalFrameSize()

		b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))

		switch {
		case b.showFileTreePreview && !b.showLogs:
			b.secondaryViewport.SetContent(b.fileTreePreviewView(b.treePreviewFiles))
		case b.currentImage != nil && !b.showLogs:
			return b, b.convertImageToStringCmd(b.secondaryViewport.Width)
		case b.errorMsg != "":
			b.secondaryViewport.SetContent(b.errorView(b.errorMsg))
		case b.showHelp && !b.showLogs:
			b.secondaryViewport.SetContent(b.helpView())
		case b.showLogs && b.currentImage == nil:
			b.secondaryViewport.SetContent(b.logView())
		default:
			b.secondaryViewport.SetContent(b.textContentView(b.secondaryBoxContent))
		}

		if !b.ready {
			b.ready = true
		}

		return b, nil
	case tea.MouseMsg:
		b.handleMouse(msg)
	case tea.KeyMsg:
		cmd = b.handleKeys(msg)
		cmds = append(cmds, cmd)

		return b, tea.Batch(cmds...)
	}

	if b.activeBox != PrimaryBoxActive {
		b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	b.textinput, cmd = b.textinput.Update(msg)
	cmds = append(cmds, cmd)

	b.spinner, cmd = b.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}
