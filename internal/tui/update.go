package tui

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/dirfs"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

// scrollFiletree handles wrapping of the filetree and
// scrolling of the viewport.
func (b *Bubble) scrollFileTree() {
	top := b.primaryViewport.YOffset
	bottom := b.primaryViewport.Height + b.primaryViewport.YOffset - 1

	if b.treeCursor < top {
		b.primaryViewport.LineUp(1)
	} else if b.treeCursor > bottom {
		b.primaryViewport.LineDown(1)
	}

	if b.treeCursor > len(b.treeFiles)-1 {
		b.treeCursor = 0
		b.primaryViewport.GotoTop()
	} else if b.treeCursor < top {
		b.treeCursor = len(b.treeFiles) - 1
		b.primaryViewport.GotoBottom()
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
		b.fileSizes = make([]string, len(msg))

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
		b.treePreviewFiles = msg
		b.secondaryViewport.SetContent(b.fileTreePreviewView(msg))

		return b, nil
	case convertImageToStringMsg:
		b.showHelp = false
		b.secondaryViewport.SetContent(b.textContentView(string(msg)))

		return b, nil
	case moveDirItemMsg:
		b.moveMode = false
		b.treeItemToMove = nil
		b.primaryViewport.SetContent(b.fileTreeView(msg))

		return b, nil
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
		b.primaryViewport.SetContent(b.fileTreeView(msg.entries))

		return b, nil
	case errorMsg:
		b.showHelp = false
		b.errorMsg = string(msg)
		b.secondaryViewport.SetContent(b.errorView(string(msg)))

		return b, nil
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		b.primaryViewport.Width = (msg.Width / 2) - box.GetHorizontalBorderSize()
		b.primaryViewport.Height = msg.Height - StatusBarHeight - box.GetVerticalBorderSize()
		b.secondaryViewport.Width = (msg.Width / 2) - box.GetHorizontalBorderSize()
		b.secondaryViewport.Height = msg.Height - StatusBarHeight - box.GetVerticalBorderSize()

		b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))

		switch {
		case b.showFileTreePreview:
			b.secondaryViewport.SetContent(b.fileTreePreviewView(b.treePreviewFiles))
		case b.currentImage != nil:
			return b, b.convertImageToStringCmd(b.secondaryViewport.Width - box.GetHorizontalFrameSize())
		case b.errorMsg != "":
			b.secondaryViewport.SetContent(b.errorView(b.errorMsg))
		case b.showHelp:
			b.secondaryViewport.SetContent(b.helpView())
		default:
			b.secondaryViewport.SetContent(b.textContentView(b.secondaryBoxContent))
		}

		if !b.ready {
			b.ready = true
		}

		return b, nil
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if b.activeBox == 0 {
				b.treeCursor--
				b.scrollFileTree()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			if b.activeBox == 1 {
				b.secondaryViewport.LineUp(1)
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}
		case tea.MouseWheelDown:
			if b.activeBox == 0 {
				b.treeCursor++
				b.scrollFileTree()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			if b.activeBox == 1 {
				b.secondaryViewport.LineDown(1)
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}
		}
	case tea.KeyMsg:
		if msg.String() == "g" && b.previousKey.String() == "g" {
			if !b.showCommandInput && b.activeBox == 0 && !b.showBoxSpinner {
				b.treeCursor = 0
				b.primaryViewport.GotoTop()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			if !b.showCommandInput && b.activeBox == 1 {
				b.secondaryViewport.GotoTop()
			}

			return b, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		case "j", "up":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor++
				b.scrollFileTree()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}
		case "k", "down":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor--
				b.scrollFileTree()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}
		case "h", "left":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = 0
				b.showFilesOnly = false
				b.showDirectoriesOnly = false
				workingDirectory, err := dirfs.GetWorkingDirectory()
				if err != nil {
					return b, b.handleErrorCmd(err)
				}

				return b, b.updateDirectoryListingCmd(
					filepath.Join(workingDirectory, dirfs.PreviousDirectory),
				)
			}
		case "l", "right":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				selectedFile, err := b.treeFiles[b.treeCursor].Info()
				if err != nil {
					return b, b.handleErrorCmd(err)
				}

				switch {
				case selectedFile.IsDir():
					currentDir, err := dirfs.GetWorkingDirectory()
					if err != nil {
						return b, b.handleErrorCmd(err)
					}

					directoryToOpen := filepath.Join(currentDir, selectedFile.Name())

					if len(b.foundFilesPaths) > 0 {
						directoryToOpen = b.foundFilesPaths[b.treeCursor]
					}

					return b, b.updateDirectoryListingCmd(directoryToOpen)
				case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
					symlinkFile, err := os.Readlink(selectedFile.Name())
					if err != nil {
						return b, b.handleErrorCmd(err)
					}

					fileInfo, err := os.Stat(symlinkFile)
					if err != nil {
						return b, b.handleErrorCmd(err)
					}

					if fileInfo.IsDir() {
						currentDir, err := dirfs.GetWorkingDirectory()
						if err != nil {
							return b, b.handleErrorCmd(err)
						}

						return b, b.updateDirectoryListingCmd(filepath.Join(currentDir, fileInfo.Name()))
					}

					return b, b.readFileContentCmd(
						fileInfo.Name(),
						b.secondaryViewport.Width-box.GetHorizontalFrameSize(),
					)

				default:
					fileToRead := selectedFile.Name()

					if len(b.foundFilesPaths) > 0 {
						fileToRead = b.foundFilesPaths[b.treeCursor]
					}

					return b, b.readFileContentCmd(
						fileToRead,
						b.secondaryViewport.Width-box.GetHorizontalFrameSize(),
					)
				}
			}
		case "p":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				selectedFile, err := b.treeFiles[b.treeCursor].Info()
				if err != nil {
					return b, b.handleErrorCmd(err)
				}

				switch {
				case selectedFile.IsDir():
					return b, b.previewDirectoryListingCmd(selectedFile.Name())
				case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
					symlinkFile, err := os.Readlink(selectedFile.Name())
					if err != nil {
						return b, b.handleErrorCmd(err)
					}

					fileInfo, err := os.Stat(symlinkFile)
					if err != nil {
						return b, b.handleErrorCmd(err)
					}

					if fileInfo.IsDir() {
						return b, b.previewDirectoryListingCmd(fileInfo.Name())
					}
				default:
					return b, nil
				}
			}
		case "G":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = len(b.treeFiles) - 1
				b.primaryViewport.GotoBottom()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			if b.activeBox == 1 && !b.showCommandInput && !b.showBoxSpinner {
				b.secondaryViewport.GotoBottom()
			}
		case "~":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = 0
				b.fileSizes = nil
				homeDir, err := dirfs.GetHomeDirectory()
				if err != nil {
					return b, b.handleErrorCmd(err)
				}

				return b, b.updateDirectoryListingCmd(homeDir)
			}
		case "/":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = 0
				b.fileSizes = nil
				return b, b.updateDirectoryListingCmd(dirfs.RootDirectory)
			}
		case ".":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.showHiddenFiles = !b.showHiddenFiles

				switch {
				case b.showDirectoriesOnly:
					return b, b.getDirectoryListingByTypeCmd(dirfs.DirectoriesListingType)
				case b.showFilesOnly:
					return b, b.getDirectoryListingByTypeCmd(dirfs.FilesListingType)
				default:
					return b, b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
				}
			}
		case "S":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.showDirectoriesOnly = !b.showDirectoriesOnly
				b.showFilesOnly = false

				if b.showDirectoriesOnly {
					return b, b.getDirectoryListingByTypeCmd(dirfs.DirectoriesListingType)
				}

				return b, b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
			}
		case "s":
			if b.activeBox == 0 && !b.showCommandInput && !b.showBoxSpinner {
				b.showFilesOnly = !b.showFilesOnly
				b.showDirectoriesOnly = false

				if b.showFilesOnly {
					return b, b.getDirectoryListingByTypeCmd(dirfs.FilesListingType)
				}

				return b, b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
			}
		case "y":
			if b.activeBox == 0 && len(b.treeFiles) > 0 && !b.showCommandInput && !b.showBoxSpinner {
				selectedFile := b.treeFiles[b.treeCursor]

				return b, b.copyToClipboardCmd(selectedFile.Name())
			}
		case "Z":
			if b.activeBox == 0 && len(b.treeFiles) > 0 && !b.showCommandInput && !b.showBoxSpinner {
				selectedFile := b.treeFiles[b.treeCursor]

				return b, tea.Sequentially(
					b.zipDirectoryCmd(selectedFile.Name()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}
		case "U":
			if b.activeBox == 0 && len(b.treeFiles) > 0 && !b.showCommandInput && !b.showBoxSpinner {
				selectedFile := b.treeFiles[b.treeCursor]

				return b, tea.Sequentially(
					b.unzipDirectoryCmd(selectedFile.Name()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}
		case "n":
			if !b.showCommandInput && !b.showBoxSpinner {
				b.createFileMode = true
				b.showCommandInput = true
				b.textinput.Placeholder = "Enter file name"
				b.textinput.Focus()

				return b, textinput.Blink
			}
		case "N":
			if !b.showCommandInput && !b.showBoxSpinner {
				b.createDirectoryMode = true
				b.showCommandInput = true
				b.textinput.Placeholder = "Enter directory name"
				b.textinput.Focus()

				return b, textinput.Blink
			}
		case "ctrl+d":
			if !b.showCommandInput && !b.showBoxSpinner {
				b.deleteMode = true
				b.showCommandInput = true
				b.textinput.Placeholder = "Are you sure you want to delete this? (y/n)"
				b.textinput.Focus()

				return b, nil
			}
		case "M":
			if !b.showCommandInput && !b.showBoxSpinner {
				b.moveMode = true
				b.treeItemToMove = b.treeFiles[b.treeCursor]
				workingDir, err := dirfs.GetWorkingDirectory()
				if err != nil {
					b.moveInitiatedDirectory = dirfs.CurrentDirectory
				}

				b.moveInitiatedDirectory = workingDir

				return b, nil
			}
		case "enter":
			switch {
			case b.moveMode:
				return b, b.moveDirectoryItemCmd(b.treeItemToMove.Name())
			case b.createFileMode:
				return b, tea.Sequentially(
					b.createFileCmd(b.textinput.Value()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			case b.createDirectoryMode:
				return b, tea.Sequentially(
					b.createDirectoryCmd(b.textinput.Value()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			case b.renameMode:
				selectedFile := b.treeFiles[b.treeCursor]

				return b, tea.Sequentially(
					b.renameDirectoryItemCmd(selectedFile.Name(), b.textinput.Value()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			case b.findMode:
				b.showCommandInput = false
				b.showBoxSpinner = true

				return b, b.findFilesByNameCmd(b.textinput.Value())
			case b.deleteMode:
				selectedFile := b.treeFiles[b.treeCursor]

				if strings.ToLower(b.textinput.Value()) == "y" || strings.ToLower(b.textinput.Value()) == "yes" {
					if selectedFile.IsDir() {
						return b, tea.Sequentially(
							b.deleteDirectoryCmd(selectedFile.Name()),
							b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
						)
					} else {
						return b, tea.Sequentially(
							b.deleteFileCmd(selectedFile.Name()),
							b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
						)
					}
				}
			default:
				return b, nil
			}
		case "E":
			selectedFile := b.treeFiles[b.treeCursor]

			if !b.showCommandInput && b.activeBox == 0 && !b.showBoxSpinner {
				selectionPath := viper.GetString("selection-path")

				if selectionPath == "" && !selectedFile.IsDir() {
					editorPath := os.Getenv("EDITOR")
					if editorPath == "" {
						return b, b.handleErrorCmd(errors.New("$EDITOR not set"))
					}

					editorCmd := exec.Command(editorPath, selectedFile.Name())
					editorCmd.Stdin = os.Stdin
					editorCmd.Stdout = os.Stdout
					editorCmd.Stderr = os.Stderr

					err := editorCmd.Start()
					if err != nil {
						return b, b.handleErrorCmd(err)
					}

					err = editorCmd.Wait()
					if err != nil {
						return b, b.handleErrorCmd(err)
					}

					return b, b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
				} else {
					return b, tea.Sequentially(
						b.writeSelectionPathCmd(selectionPath, selectedFile.Name()),
						tea.Quit,
					)
				}
			}
		case "C":
			if !b.showCommandInput && b.activeBox == 0 && len(b.treeFiles) > 0 && !b.showBoxSpinner {
				selectedFile := b.treeFiles[b.treeCursor]

				if selectedFile.IsDir() {
					return b, tea.Sequentially(
						b.copyDirectoryCmd(selectedFile.Name()),
						b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
					)
				}

				return b, tea.Sequentially(
					b.copyFileCmd(selectedFile.Name()),
					b.updateDirectoryListingCmd(dirfs.CurrentDirectory),
				)
			}
		case "ctrl+f":
			if !b.showCommandInput && !b.showBoxSpinner {
				b.findMode = true
				b.showCommandInput = true
				b.textinput.Placeholder = "Enter a search term"
				b.textinput.Focus()

				return b, textinput.Blink
			}
		case "R":
			if b.activeBox == 0 && !b.showBoxSpinner && len(b.treeFiles) > 0 {
				selectedFile := b.treeFiles

				if selectedFile != nil {
					b.renameMode = true
					b.showCommandInput = true
					b.textinput.Placeholder = "Enter new name"
					b.textinput.Focus()

					return b, textinput.Blink
				}
			}
		case "esc":
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
			b.foundFilesPaths = nil
			b.showBoxSpinner = false
			b.secondaryViewport.SetContent(b.helpView())
			b.textinput.Blur()
			b.textinput.Reset()

			return b, b.updateDirectoryListingCmd(dirfs.CurrentDirectory)
		case "tab":
			b.activeBox = (b.activeBox + 1) % 2
		}

		b.previousKey = msg
	}

	if b.activeBox != 0 {
		b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	b.textinput, cmd = b.textinput.Update(msg)
	cmds = append(cmds, cmd)

	b.spinner, cmd = b.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}
