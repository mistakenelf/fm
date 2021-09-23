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

	"github.com/knipferrc/fm/directory"
	"github.com/knipferrc/fm/internal/colorimage"
	"github.com/knipferrc/fm/internal/markdown"
	"github.com/knipferrc/fm/internal/sourcecode"
	"github.com/knipferrc/fm/strfmt"

	tea "github.com/charmbracelet/bubbletea"
)

type updateDirectoryListingMsg []fs.DirEntry
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

// updateDirectoryListing updates the directory listing based on the name of the direcoctory provided.
func (m Model) updateDirectoryListing(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := directory.GetDirectoryListing(name, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(files)
	}
}

// renameFileOrDir renames a file or directory based on the name and value provided.
func (m Model) renameFileOrDir(name, value string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.RenameDirOrFile(name, value); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// moveDirectoryItem moves a file to the current working directory.
func (m Model) moveDirectoryItem(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := directory.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same file name.
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had.
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		if err = directory.MoveDirectoryItem(src, dst); err != nil {
			return errorMsg(err.Error())
		}

		files, err := directory.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveDirItemMsg(files)
	}
}

// deleteDir deletes a directory based on the name provided.
func (m Model) deleteDir(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.DeleteDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// deleteFile deletes a file based on the name provided.
func (m Model) deleteFile(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.DeleteFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// readFileContent reads the content of a file and returns it.
func (m Model) readFileContent(file fs.DirEntry, width, height int) tea.Cmd {
	return func() tea.Msg {
		content, err := directory.ReadFileContent(file.Name())
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

			imageString, err := colorimage.ImageToString(uint(width), uint(height), img)
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
			code, err := sourcecode.Highlight(content, filepath.Ext(file.Name()))
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

// redrawImage redraws the image based on the width and height provided.
func (m Model) redrawImage(width, height int) tea.Cmd {
	return func() tea.Msg {
		imageString, err := colorimage.ImageToString(uint(width), uint(height), m.colorimage.Image)
		if err != nil {
			return errorMsg(err.Error())
		}

		return convertImageToStringMsg(imageString)
	}
}

// createDir creates a directory based on the name provided.
func (m Model) createDir(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.CreateDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// createFile creates a file based on the name provided.
func (m Model) createFile(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.CreateFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// zipDirectory zips a directory based on the name provided.
func (m Model) zipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.Zip(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// unzipDirectory unzips a directory based on the name provided.
func (m Model) unzipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.Unzip(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyFile copies a file based on the name provided.
func (m Model) copyFile(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.CopyFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyDirectory copies a directory based on the name provided.
func (m Model) copyDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		if err := directory.CopyDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// getDirectoryItemSize calculates the size of a directory or file.
func (m Model) getDirectoryItemSize(name string) tea.Cmd {
	if m.directoryItemSizeCtx != nil && m.directoryItemSizeCtx.cancel != nil {
		m.directoryItemSizeCtx.cancel()
	}

	ctx, cancel := context.WithTimeout(m.directoryItemSizeCtx.ctx, 300*time.Millisecond)
	m.directoryItemSizeCtx.cancel = cancel

	return func() tea.Msg {
		defer cancel()
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			size, err := directory.GetDirectoryItemSize(name)
			if err != nil {
				return directoryItemSizeMsg("N/A")
			}

			sizeString := strfmt.ConvertBytesToSizeString(size)

			return directoryItemSizeMsg(sizeString)
		}

		return nil
	}
}
