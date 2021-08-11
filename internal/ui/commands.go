package ui

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/helpers"

	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
)

type directoryUpdateMsg []fs.FileInfo
type moveDirItemMsg []fs.FileInfo
type markdownMsg string
type imageMsg string
type errorMsg string
type readFileContentMsg struct {
	markdownContent string
	fileContent     string
	imgFile         image.Image
}

// updateDirectoryListing updates the directory listing based on the name of the direcoctory provided.
func (m Model) updateDirectoryListing(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := helpers.GetDirectoryListing(name, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryUpdateMsg(files)
	}
}

// renameFileOrDir renames a file or directory based on the name and value provided.
func (m Model) renameFileOrDir(name, value string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.RenameDirOrFile(name, value); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// moveDir moves a directory to the current working directory.
func (m Model) moveDir(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := helpers.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same folder name.
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same folder name that it had.
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		if err = helpers.MoveDirectory(src, dst); err != nil {
			return errorMsg(err.Error())
		}

		files, err := helpers.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveDirItemMsg(files)
	}
}

// moveFile moves a file to the current working directory.
func (m Model) moveFile(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := helpers.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same file name.
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had.
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		if err = helpers.MoveFile(src, dst); err != nil {
			return errorMsg(err.Error())
		}

		files, err := helpers.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveDirItemMsg(files)
	}
}

// deleteDir deletes a directory based on the name provided.
func (m Model) deleteDir(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.DeleteDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// deleteFile deletes a file based on the name provided.
func (m Model) deleteFile(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.DeleteFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// readFileContent reads the content of a file and returns it.
func (m Model) readFileContent(file fs.FileInfo, width int) tea.Cmd {
	cfg := config.GetConfig()

	return func() tea.Msg {
		content, err := helpers.ReadFileContent(file.Name())
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return both the pretty markdown as well as the plain content without glamour
		// to use later when resizing the window.
		if filepath.Ext(file.Name()) == ".md" && cfg.Settings.PrettyMarkdown {
			markdownContent, err := helpers.RenderMarkdown(width, content)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				fileContent:     markdownContent,
				markdownContent: content,
			}
		}

		if filepath.Ext(file.Name()) == ".png" {
			imageContent, err := os.Open(file.Name())
			if err != nil {
				return errorMsg(err.Error())
			}

			img, err := png.Decode(imageContent)
			if err != nil {
				return errorMsg(err.Error())
			}

			imageString := helpers.ConvertToAscii(helpers.ScaleImage(img, width-2))

			return readFileContentMsg{
				fileContent:     imageString,
				markdownContent: "",
				imgFile:         img,
			}
		}

		buf := new(bytes.Buffer)
		if err = quick.Highlight(buf, content, filepath.Ext(file.Name()), "terminal256", "dracula"); err != nil {
			return errorMsg(err.Error())
		}

		// Return the syntax highlighted content and markdown content as empty
		// since were not dealing with markdown.
		return readFileContentMsg{
			fileContent:     buf.String(),
			markdownContent: "",
		}
	}
}

// renderMarkdownContent renders the markdown content and returns it.
func renderMarkdownContent(width int, content string) tea.Cmd {
	return func() tea.Msg {
		markdownContent, err := helpers.RenderMarkdown(width, content)
		if err != nil {
			return errorMsg(err.Error())
		}

		return markdownMsg(markdownContent)
	}
}

// renderMarkdownContent renders the markdown content and returns it.
func renderImageContent(img image.Image, width int) tea.Cmd {
	return func() tea.Msg {
		imageString := helpers.ConvertToAscii(helpers.ScaleImage(img, width-2))

		return imageMsg(imageString)
	}
}

// createDir creates a directory based on the name provided.
func (m Model) createDir(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.CreateDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// createFile creates a file based on the name provided.
func (m Model) createFile(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.CreateFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// zipDirectory zips a directory based on the name provided.
func (m Model) zipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.ZipDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// unzipDirectory unzips a directory based on the name provided.
func (m Model) unzipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.UnzipDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyFile copies a file based on the name provided.
func (m Model) copyFile(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.CopyFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyDirectory copies a directory based on the name provided.
func (m Model) copyDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		if err := helpers.CopyDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}
