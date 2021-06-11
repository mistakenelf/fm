package ui

import (
	"bytes"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/utils"

	"github.com/alecthomas/chroma/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type directoryMsg []fs.FileInfo
type fileContentMsg struct {
	markdownContent string
	fileContent     string
}
type markdownMsg string

func updateDirectoryListing(dir string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		files := utils.GetDirectoryListing(dir, showHidden)

		return directoryMsg(files)
	}
}

func renameFileOrDir(filename, value string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		utils.RenameDirOrFile(filename, value)
		files := utils.GetDirectoryListing(constants.CurrentDirectory, showHidden)

		return directoryMsg(files)
	}
}

func moveDir(dir, value string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		utils.CopyDir(dir, value, true)
		files := utils.GetDirectoryListing(constants.CurrentDirectory, showHidden)

		return directoryMsg(files)
	}
}

func moveFile(file, value string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		utils.CopyFile(file, value, true)
		files := utils.GetDirectoryListing(constants.CurrentDirectory, showHidden)

		return directoryMsg(files)
	}
}

func deleteDir(dir string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		utils.DeleteDirectory(dir)
		files := utils.GetDirectoryListing(constants.CurrentDirectory, showHidden)

		return directoryMsg(files)
	}
}

func deleteFile(file string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		utils.DeleteFile(file)
		files := utils.GetDirectoryListing(constants.CurrentDirectory, showHidden)

		return directoryMsg(files)
	}
}

func (m model) readFileContent(file fs.FileInfo) tea.Cmd {
	width := m.secondaryPane.Width

	return func() tea.Msg {
		content := utils.ReadFileContent(file.Name())

		if filepath.Ext(file.Name()) == ".md" {
			return fileContentMsg{
				fileContent:     renderMarkdown(width, content),
				markdownContent: content,
			}
		} else {
			buf := new(bytes.Buffer)
			err := quick.Highlight(buf, content, filepath.Ext(file.Name()), "terminal256", "dracula")

			if err != nil {
				log.Fatal("error")
			}

			return fileContentMsg{
				fileContent:     buf.String(),
				markdownContent: "",
			}
		}
	}
}

func renderMarkdownContent(width int, content string) tea.Cmd {
	return func() tea.Msg {
		return markdownMsg(renderMarkdown(width, content))
	}
}

func renderMarkdown(width int, content string) string {
	bg := "light"
	if lipgloss.HasDarkBackground() {
		bg = "dark"
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithStandardStyle(bg),
	)

	out, err := r.Render(content)
	if err != nil {
		log.Fatal(err)
	}

	return out
}

func createDir(dir string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		utils.CreateDirectory(dir)
		files := utils.GetDirectoryListing(constants.CurrentDirectory, showHidden)

		return directoryMsg(files)
	}
}

func createFile(name string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		utils.CreateFile(name)
		files := utils.GetDirectoryListing(constants.CurrentDirectory, showHidden)

		return directoryMsg(files)
	}
}
