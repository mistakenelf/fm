package ui

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/internal/renderer"
	"github.com/knipferrc/fm/strfmt"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type updateDirectoryListingMsg []fs.DirEntry
type previewDirectoryListingMsg []fs.DirEntry
type moveDirItemMsg []fs.DirEntry
type errorMsg string
type directoryItemSizeMsg string
type copyToClipboardMsg string
type convertImageToStringMsg string
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
func (m Model) updateDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := dirfs.GetDirectoryListing(name, m.dirTree.ShowHidden)
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
func (m Model) previewDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(filepath.Join(currentDir, name), m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return previewDirectoryListingMsg(files)
	}
}

// renameDirectoryItemCmd renames a file or directory based on the name and value provided.
func (m Model) renameDirectoryItemCmd(name, value string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.RenameDirectoryItem(name, value); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// moveDirectoryItemCmd moves a file to the current working directory.
func (m Model) moveDirectoryItemCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same file name.
		src := filepath.Join(m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had.
		dst := filepath.Join(workingDir, name)

		if err = dirfs.MoveDirectoryItem(src, dst); err != nil {
			return errorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		err = os.Chdir(m.initialMoveDirectory)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveDirItemMsg(files)
	}
}

// redrawImageCmd redraws the image based on the width provided.
func (m Model) redrawImageCmd(width int) tea.Cmd {
	return func() tea.Msg {
		imageString := renderer.ImageToString(width, m.renderer.Image)

		return convertImageToStringMsg(imageString)
	}
}

// deleteDirectoryCmd deletes a directory based on the name provided.
func (m Model) deleteDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// deleteFileCmd deletes a file based on the name provided.
func (m Model) deleteFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// readFileContentCmd reads the content of a file and returns it.
func (m Model) readFileContentCmd(fileName string, width int) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(fileName)
		if err != nil {
			return errorMsg(err.Error())
		}

		switch {
		case filepath.Ext(fileName) == ".md" && m.appConfig.Settings.PrettyMarkdown:
			markdownContent, err := renderer.RenderMarkdown(width, content)
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

			imageString := renderer.ImageToString(width, img)

			return readFileContentMsg{
				rawContent:  content,
				code:        "",
				markdown:    "",
				imageString: imageString,
				pdfContent:  "",
				image:       img,
			}
		case filepath.Ext(fileName) == ".pdf":
			pdfContent, err := renderer.ReadPdf(fileName)
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
			code, err := renderer.Highlight(content, filepath.Ext(fileName), m.appConfig.Settings.SyntaxTheme)
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
func (m Model) createDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// createFileCmd creates a file based on the name provided.
func (m Model) createFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// zipDirectoryCmd zips a directory based on the name provided.
func (m Model) zipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		dirToZip := filepath.Join(currentDir, name)
		if err := dirfs.Zip(dirToZip); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// unzipDirectoryCmd unzips a directory based on the name provided.
func (m Model) unzipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		dirToUnzip := filepath.Join(currentDir, name)
		if err := dirfs.Unzip(dirToUnzip); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyFileCmd copies a file based on the name provided.
func (m Model) copyFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyDirectoryCmd copies a directory based on the name provided.
func (m Model) copyDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// getDirectoryItemSizeCmd calculates the size of a directory or file.
func (m Model) getDirectoryItemSizeCmd(name string) tea.Cmd {
	if m.directoryItemSizeCtx != nil && m.directoryItemSizeCtx.cancel != nil {
		m.directoryItemSizeCtx.cancel()
	}

	ctx, cancel := context.WithTimeout(m.directoryItemSizeCtx.ctx, 300*time.Millisecond)
	m.directoryItemSizeCtx.cancel = cancel

	return func() tea.Msg {
		defer cancel()
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			size, err := dirfs.GetDirectoryItemSize(name)
			if err != nil {
				return directoryItemSizeMsg("N/A")
			}

			sizeString := strfmt.ConvertBytesToSizeString(size)

			return directoryItemSizeMsg(sizeString)
		}

		return nil
	}
}

// handleErrorCmd returns an error message to the UI.
func (m Model) handleErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err.Error())
	}
}

// copyToClipboardCmd copies the provided string to the clipboard.
func (m Model) copyToClipboardCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		filePath := filepath.Join(workingDir, name)
		clipboard.Write(clipboard.FmtText, []byte(filePath))

		return copyToClipboardMsg(fmt.Sprintf("%s %s %s", "Successfully copied", filePath, "to clipboard"))
	}
}

// getDirectoryListingByTypeCmd returns only directories in the current directory.
func (m Model) getDirectoryListingByTypeCmd(listType string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		directories, err := dirfs.GetDirectoryListingByType(workingDir, listType, showHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(directories)
	}
}

// findFilesByNameCmd finds files based on name.
func (m Model) findFilesByNameCmd(name string) tea.Cmd {
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
func (m Model) writeSelectionPathCmd(selectionPath, filePath string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.WriteToFile(selectionPath, filePath); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}
