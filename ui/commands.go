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
type moveFileMsg []fs.FileInfo
type markdownMsg string
type fileContentMsg struct {
	markdownContent string
	fileContent     string
}

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

func (m model) moveDir(dir string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, dir)
		dst := fmt.Sprintf("%s/%s", utils.GetWorkingDirectory(), dir)

		utils.MoveDirectory(src, dst)
		files := utils.GetDirectoryListing(m.initialMoveDirectory, showHidden)

		return moveFileMsg(files)
	}
}

func (m model) moveFile(file string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, file)
		dst := fmt.Sprintf("%s/%s", utils.GetWorkingDirectory(), file)

		utils.MoveFile(src, dst)
		files := utils.GetDirectoryListing(m.initialMoveDirectory, showHidden)

		return moveFileMsg(files)
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
	cfg := config.GetConfig()
	width := m.secondaryPane.Width

	return func() tea.Msg {
		content := utils.ReadFileContent(file.Name())

		if filepath.Ext(file.Name()) == ".md" && cfg.Settings.PrettyMarkdown {
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
