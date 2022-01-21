package tui

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/strfmt"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type updateDirectoryListingMsg []fs.DirEntry
type previewDirectoryListingMsg []fs.DirEntry
type moveDirItemMsg []fs.DirEntry
type errorMsg string
type copyToClipboardMsg string
type convertImageToStringMsg string
type directoryItemSizeMsg struct {
	index int
	size  string
}
type findFilesByNameMsg struct {
	paths   []string
	entries []fs.DirEntry
}
type readFileContentMsg struct {
	rawContent  string
	markdown    string
	code        string
	imageString string
	pdfContent  string
	image       image.Image
}

// updateDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (b Bubble) updateDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := dirfs.GetDirectoryListing(name, b.showHiddenFiles)
		if err != nil {
			return errorMsg(err.Error())
		}

		err = os.Chdir(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(files)
	}
}

// previewDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (b Bubble) previewDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(filepath.Join(currentDir, name), b.showHiddenFiles)
		if err != nil {
			return errorMsg(err.Error())
		}

		return previewDirectoryListingMsg(files)
	}
}

// renameDirectoryItemCmd renames a file or directory based on the name and value provided.
func (b Bubble) renameDirectoryItemCmd(name, value string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.RenameDirectoryItem(name, value); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// moveDirectoryItemCmd moves a file to the current working directory.
func (b Bubble) moveDirectoryItemCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory frob which the move was intiated from
		// and give it the same file name.
		src := filepath.Join(b.moveInitiatedDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had.
		dst := filepath.Join(workingDir, name)

		if err = dirfs.MoveDirectoryItem(src, dst); err != nil {
			return errorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(b.moveInitiatedDirectory, b.showHiddenFiles)
		if err != nil {
			return errorMsg(err.Error())
		}

		err = os.Chdir(b.moveInitiatedDirectory)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveDirItemMsg(files)
	}
}

// convertImageToStringCmd redraws the image based on the width provided.
func (b Bubble) convertImageToStringCmd(width int) tea.Cmd {
	return func() tea.Msg {
		imageString := strfmt.ImageToString(width, b.currentImage)

		return convertImageToStringMsg(imageString)
	}
}

// deleteDirectoryCmd deletes a directory based on the name provided.
func (b Bubble) deleteDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// deleteFileCmd deletes a file based on the name provided.
func (b Bubble) deleteFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// readFileContentCmd reads the content of a file and returns it.
func (b Bubble) readFileContentCmd(fileName string, width int) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(fileName)
		if err != nil {
			return errorMsg(err.Error())
		}

		switch {
		case filepath.Ext(fileName) == ".md" && b.appConfig.Settings.PrettyMarkdown:
			markdownContent, err := strfmt.RenderMarkdown(width, content)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  content,
				markdown:    markdownContent,
				code:        "",
				imageString: "",
				pdfContent:  "",
				image:       nil,
			}
		case filepath.Ext(fileName) == ".png" || filepath.Ext(fileName) == ".jpg" || filepath.Ext(fileName) == ".jpeg":
			imageContent, err := os.Open(fileName)
			if err != nil {
				return errorMsg(err.Error())
			}

			img, _, err := image.Decode(imageContent)
			if err != nil {
				return errorMsg(err.Error())
			}

			imageString := strfmt.ImageToString(width, img)

			return readFileContentMsg{
				rawContent:  content,
				code:        "",
				markdown:    "",
				imageString: imageString,
				pdfContent:  "",
				image:       img,
			}
		case filepath.Ext(fileName) == ".pdf":
			pdfContent, err := strfmt.ReadPdf(fileName)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  content,
				code:        "",
				markdown:    "",
				imageString: "",
				pdfContent:  pdfContent,
				image:       nil,
			}
		default:
			syntaxTheme := b.appConfig.Settings.SyntaxTheme.Light
			if lipgloss.HasDarkBackground() {
				syntaxTheme = b.appConfig.Settings.SyntaxTheme.Dark
			}

			code, err := strfmt.Highlight(content, filepath.Ext(fileName), syntaxTheme)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  content,
				code:        code,
				markdown:    "",
				imageString: "",
				pdfContent:  "",
				image:       nil,
			}
		}
	}
}

// createDirectoryCmd creates a directory based on the name provided.
func (b Bubble) createDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// createFileCmd creates a file based on the name provided.
func (b Bubble) createFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// zipDirectoryCmd zips a directory based on the name provided.
func (b Bubble) zipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.Zip(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// unzipDirectoryCmd unzips a directory based on the name provided.
func (b Bubble) unzipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.Unzip(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyFileCmd copies a file based on the name provided.
func (b Bubble) copyFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyDirectoryCmd copies a directory based on the name provided.
func (b Bubble) copyDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// getDirectoryItemSizeCmd calculates the size of a directory or file.
func (b Bubble) getDirectoryItemSizeCmd(name string, i int) tea.Cmd {
	return func() tea.Msg {
		size, err := dirfs.GetDirectoryItemSize(name)
		if err != nil {
			return directoryItemSizeMsg{size: "N/A", index: i}
		}

		sizeString := strfmt.ConvertBytesToSizeString(size)

		return directoryItemSizeMsg{
			size:  sizeString,
			index: i,
		}
	}
}

// handleErrorCmd returns an error message to the UI.
func (b Bubble) handleErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err.Error())
	}
}

// copyToClipboardCmd copies the provided string to the clipboard.
func (b Bubble) copyToClipboardCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		filePath := filepath.Join(workingDir, name)
		err = clipboard.WriteAll(filePath)
		if err != nil {
			return errorMsg(err.Error())
		}

		return copyToClipboardMsg(fmt.Sprintf("%s %s %s", "Successfully copied", filePath, "to clipboard"))
	}
}

// getDirectoryListingByTypeCmd returns only directories in the current directory.
func (b Bubble) getDirectoryListingByTypeCmd(listType string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		directories, err := dirfs.GetDirectoryListingByType(workingDir, listType, b.showHiddenFiles)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(directories)
	}
}

// findFilesByNameCmd finds files based on name.
func (b Bubble) findFilesByNameCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		paths, entries, err := dirfs.FindFilesByName(name, workingDir)
		if err != nil {
			return errorMsg(err.Error())
		}

		return findFilesByNameMsg{
			paths:   paths,
			entries: entries,
		}
	}
}

// writeSelectionPathCmd writes content to the file specified.
func (b Bubble) writeSelectionPathCmd(selectionPath, filePath string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.WriteToFile(selectionPath, filePath); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// redrawCmd redraws the UI.
func (b Bubble) redrawCmd() tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  b.width,
			Height: b.height,
		}
	}
}
