package ui

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/utils"

	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type directoryMsg []fs.FileInfo
type moveMsg []fs.FileInfo
type markdownMsg string
type errorMsg string
type fileContentMsg struct {
	markdownContent string
	fileContent     string
}

// Get an updated directory listing based on the name
// of a directory passed in
func (m model) updateDirectoryListing(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := utils.GetDirectoryListing(name, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryMsg(files)
	}
}

// Rename a file or directory based on its current filename
// and its new value, returning an updated directory listing
func (m model) renameFileOrDir(name, value string) tea.Cmd {
	return func() tea.Msg {
		err := utils.RenameDirOrFile(name, value)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Move a directory to the current working directory
// returning an updated directory listing
func (m model) moveDir(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := utils.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same folder name
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same folder name that it had
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		err = utils.MoveDirectory(src, dst)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveMsg(files)
	}
}

// Move a file to the current working directory
// returning an updated directory listing
func (m model) moveFile(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := utils.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same file name
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		err = utils.MoveFile(src, dst)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveMsg(files)
	}
}

// Delete a directory based on name and
// return an updated directory listing
func (m model) deleteDir(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.DeleteDirectory(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Delete a file based on name and return an
// updated directory listing
func (m model) deleteFile(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.DeleteFile(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Read a files content based on its name and return its content as a string.
// If the file is markdown and pretty markdown is enabled, run the content
// through glamour else run the content through chroma to get syntax highlighting
func (m model) readFileContent(file fs.FileInfo) tea.Cmd {
	cfg := config.GetConfig()
	width := m.secondaryPane.Width

	return func() tea.Msg {
		content, err := utils.ReadFileContent(file.Name())
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return both the pretty markdown as well as the plain content without glamour
		// to use later when resizing the window
		if filepath.Ext(file.Name()) == ".md" && cfg.Settings.PrettyMarkdown {
			return fileContentMsg{
				fileContent:     renderMarkdown(width, content),
				markdownContent: content,
			}
		} else {
			buf := new(bytes.Buffer)
			err := quick.Highlight(buf, content, filepath.Ext(file.Name()), "terminal256", "dracula")
			if err != nil {
				return errorMsg(err.Error())
			}

			// Return the syntax highlighted content and markdown content is empty
			// since were not dealing with markdown
			return fileContentMsg{
				fileContent:     buf.String(),
				markdownContent: "",
			}
		}
	}
}

// Render some markdown content. Need to specify a width for
// when the terminal is resized
func renderMarkdownContent(width int, content string) tea.Cmd {
	return func() tea.Msg {
		return markdownMsg(renderMarkdown(width, content))
	}
}

// Render some markdown passing it a width and content
func renderMarkdown(width int, content string) string {
	bg := "light"

	// if the terminal has a dark background, use a dark background
	// for glamour
	if lipgloss.HasDarkBackground() {
		bg = "dark"
	}

	// Create a new glamour instance with word wrapping
	// and custom background based on terminal color
	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithStandardStyle(bg),
	)

	out, err := r.Render(content)

	// TODO need to handle this error
	if err != nil {
		log.Fatal(err)
	}

	return out
}

// Create a new directory given a name and return an
// updated directory listing
func (m model) createDir(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CreateDirectory(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Create a new file given a name and return
// an updated directory listing
func (m model) createFile(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CreateFile(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Create a zipped directory given a name and return
// an updated directory listing
func (m model) zipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.ZipDirectory(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Unzip a zipped directory given a name and
// return an updated directory listing
func (m model) unzipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.UnzipDirectory(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Copy a file given a name and return an updated directory listing
func (m model) copyFile(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CopyFile(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Copy a directory given a name and return an updated directory listing
func (m model) copyDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CopyDirectory(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}
