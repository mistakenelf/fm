package ui

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/helpers"

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
		files, err := helpers.GetDirectoryListing(name, m.dirTree.ShowHidden)
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
		err := helpers.RenameDirOrFile(name, value)
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
		workingDir, err := helpers.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same folder name
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same folder name that it had
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		err = helpers.MoveDirectory(src, dst)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := helpers.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
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
		workingDir, err := helpers.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same file name
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		err = helpers.MoveFile(src, dst)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := helpers.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
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
		err := helpers.DeleteDirectory(name)
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
		err := helpers.DeleteFile(name)
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
		content, err := helpers.ReadFileContent(file.Name())
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return both the pretty markdown as well as the plain content without glamour
		// to use later when resizing the window
		if filepath.Ext(file.Name()) == ".md" && cfg.Settings.PrettyMarkdown {
			markdownContent, err := renderMarkdown(width, content)
			if err != nil {
				return errorMsg(err.Error())
			}

			return fileContentMsg{
				fileContent:     markdownContent,
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
		markdownContent, err := renderMarkdown(width, content)
		if err != nil {
			return errorMsg(err.Error())
		}

		return markdownMsg(markdownContent)
	}
}

// Render some markdown passing it a width and content
func renderMarkdown(width int, content string) (string, error) {
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
	if err != nil {
		return "", err
	}

	return out, nil
}

// Create a new directory given a name and return an
// updated directory listing
func (m model) createDir(name string) tea.Cmd {
	return func() tea.Msg {
		err := helpers.CreateDirectory(name)
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
		err := helpers.CreateFile(name)
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
		err := helpers.ZipDirectory(name)
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
		err := helpers.UnzipDirectory(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Copy a file given a name and return an updated directory listing
func (m model) copyFile(name string) tea.Cmd {
	return func() tea.Msg {
		err := helpers.CopyFile(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// Copy a directory given a name and return an updated directory listing
func (m model) copyDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		err := helpers.CopyDirectory(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}