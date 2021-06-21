package ui

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/utils"

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
		// Get directory listing and determine to show hidden files/folders or not
		files, err := utils.GetDirectoryListing(name, m.dirTree.ShowHidden)

		// Something went wrong getting the directory listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return the new directory listing back to
		// the UI to be displayed in the dirtree
		return directoryMsg(files)
	}
}

// Rename a file or directory based on its current filename
// and its new value
func (m model) renameFileOrDir(name, value string) tea.Cmd {
	return func() tea.Msg {
		err := utils.RenameDirOrFile(name, value)

		// Something went wrong renaming the file or directory
		// return the error back to the UI to be displayed
		if err != nil {
			return errorMsg(err.Error())
		}

		// Need to get the updated directory listing after renaming
		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the listing, return
		// the error back to the UI and display it
		if err != nil {
			return errorMsg(err.Error())
		}

		// Now that the file or folder has been renamed, return
		// the updated directory listing to reflect those changes
		return directoryMsg(files)
	}
}

// Move a directory from one directory to another
func (m model) moveDir(name string) tea.Cmd {
	return func() tea.Msg {
		// Get the working directory so we know where to
		// move to
		workingDir, err := utils.GetWorkingDirectory()

		// Failed to get the current working directory
		// send the error back to the UI to display
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

		// Something went wrong moving the directory
		// send the error back to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Need to get an updated directory listing now that the directory has been moved
		files, err := utils.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the directory listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Now that the directory has been moved, return
		// the updated directory listing to the UI
		return moveMsg(files)
	}
}

// Move a file from one directory to another
func (m model) moveFile(name string) tea.Cmd {
	return func() tea.Msg {
		// Get the current working directory
		// so we know where to move the file to
		workingDir, err := utils.GetWorkingDirectory()

		// Something went wrong getting the working directory
		// send the error back to the UI to display
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

		// Something went wrong moving the file
		// send the error to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Need to get an updated directory listing now that a file has been moved
		files, err := utils.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Now that the file has been moved, return
		// the updated directory listing to the UI
		return moveMsg(files)
	}
}

// Delete a directory based on name
func (m model) deleteDir(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.DeleteDirectory(name)

		// Something went wrong deleting the directory
		// send the error to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Need to get an updated directory listing now that the file has been deleted
		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return the updated directory listing to the UI
		return directoryMsg(files)
	}
}

// Delete a file based on name
func (m model) deleteFile(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.DeleteFile(name)

		// Something went wrong deleting a file
		// send the error to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Now that the file has been deleted, we need an updated directory listing
		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Send updated directory listing to UI
		return directoryMsg(files)
	}
}

// Read a files content based on its name
func (m model) readFileContent(file fs.FileInfo) tea.Cmd {
	cfg := config.GetConfig()
	width := m.secondaryPane.Width

	return func() tea.Msg {
		content, err := utils.ReadFileContent(file.Name())

		// Something went wrong reading the files content
		// send the error to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// If the file is markdown and we want pretty_markdown displayed
		// return the pretty markdown as well as returning the unprettified
		// version so that we can use it later
		if filepath.Ext(file.Name()) == ".md" && cfg.Settings.PrettyMarkdown {
			return fileContentMsg{
				fileContent:     renderMarkdown(width, content),
				markdownContent: content,
			}
		} else {
			// Its not a markdown file so we want to try to run it through syntax highlighting
			buf := new(bytes.Buffer)
			err := quick.Highlight(buf, content, filepath.Ext(file.Name()), "terminal256", "dracula")

			// Somethign went wrong highlighting the file so
			// lets send it to the UI to display
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
// returning the pretty markdown
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

	// Return pretty markdown
	return out
}

// Create a new directory given a name
func (m model) createDir(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CreateDirectory(name)

		// Something went wrong creating the directory
		// send the error the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Now that a new directory has been created, we need an updated directory listing
		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)

		// Somethign went wrong getting the updated directory listing
		// return it to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return updated directory listing to the UI
		return directoryMsg(files)
	}
}

// Create a new file given a name
func (m model) createFile(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CreateFile(name)

		// Something went wrong creating a new file
		// lets send the error to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Need an updated directory listing now that a new file has been created
		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the directory listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return updated directory listing to the UI
		return directoryMsg(files)
	}
}

// Create a zipped directory given a name
func (m model) zipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.ZipDirectory(name)

		// Something went wrong zipping up the directory
		// send the error to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Since a new zipped directory has been created, we need an updated directory listing
		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the directory listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return the updated directory listing to the UI
		return directoryMsg(files)
	}
}

// Unzip a zipped directory given a name
func (m model) unzipDirectory(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.UnzipDirectory(name)

		// Something went wrong unzipping the directory
		// lets send the error to the UI to display
		if err != nil {
			return errorMsg(err.Error())
		}

		// Since the folder has been unzipped, we need an updated directory listing
		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)

		// Something went wrong getting the directory listing, return
		// the error back to the UI and display it in a pane
		if err != nil {
			return errorMsg(err.Error())
		}

		// Return the updated directory listing to the UI
		return directoryMsg(files)
	}
}
