package ui

import (
	"context"
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	// "image/jpeg" is needed for the image.Decode function.
	_ "image/jpeg"
	// "image/png" is needed for the image.Decode function.
	_ "image/png"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/internal/colorimage"
	"github.com/knipferrc/fm/internal/markdown"
	"github.com/knipferrc/fm/internal/sourcecode"
	"github.com/knipferrc/fm/strfmt"

	tea "github.com/charmbracelet/bubbletea"
)

type updateDirectoryListingMsg []fs.DirEntry
type previewDirectoryListingMsg []fs.DirEntry
type moveDirItemMsg []fs.DirEntry
type errorMsg string
type convertImageToStringMsg string
type directoryItemSizeMsg string
type readFileContentMsg struct {
	rawContent  string
	markdown    string
	code        string
	imageString string
	image       image.Image
}

// updateDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (m Model) updateDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := dirfs.GetDirectoryListing(name, m.dirTree.ShowHidden, true)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(files)
	}
}

// previewDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (m Model) previewDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := dirfs.GetDirectoryListing(name, m.dirTree.ShowHidden, false)
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
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had.
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		if err = dirfs.MoveDirectoryItem(src, dst); err != nil {
			return errorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden, true)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveDirItemMsg(files)
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
func (m Model) readFileContentCmd(file os.FileInfo, width, height int) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(file.Name())
		if err != nil {
			return errorMsg(err.Error())
		}

		switch {
		case filepath.Ext(file.Name()) == ".md" && m.appConfig.Settings.PrettyMarkdown:
			markdownContent, err := markdown.RenderMarkdown(width, content)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  content,
				markdown:    markdownContent,
				code:        "",
				imageString: "",
				image:       nil,
			}
		case filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == ".jpeg":
			imageContent, err := os.Open(file.Name())
			if err != nil {
				return errorMsg(err.Error())
			}

			img, _, err := image.Decode(imageContent)
			if err != nil {
				return errorMsg(err.Error())
			}

			imageString, err := colorimage.ImageToString(width, height, img)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  content,
				code:        "",
				markdown:    "",
				imageString: imageString,
				image:       img,
			}
		default:
			code, err := sourcecode.Highlight(content, filepath.Ext(file.Name()), m.appConfig.Settings.SyntaxTheme)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  content,
				code:        code,
				markdown:    "",
				imageString: "",
				image:       nil,
			}
		}
	}
}

// redrawImageCmd redraws the image based on the width and height provided.
func (m Model) redrawImageCmd(width, height int) tea.Cmd {
	return func() tea.Msg {
		imageString, err := colorimage.ImageToString(width, height, m.colorimage.Image)
		if err != nil {
			return errorMsg(err.Error())
		}

		return convertImageToStringMsg(imageString)
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

		dirToZip := fmt.Sprintf("%s/%s", currentDir, name)
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

		dirToUnzip := fmt.Sprintf("%s/%s", currentDir, name)
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
