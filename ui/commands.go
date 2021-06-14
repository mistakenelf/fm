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
type errorMsg string
type fileContentMsg struct {
	markdownContent string
	fileContent     string
}

func (m model) updateDirectoryListing(dir string) tea.Cmd {
	return func() tea.Msg {
		files, err := utils.GetDirectoryListing(dir, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryMsg(files)
	}
}

func (m model) renameFileOrDir(filename, value string) tea.Cmd {
	return func() tea.Msg {
		err := utils.RenameDirOrFile(filename, value)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryMsg(files)
	}
}

func (m model) moveDir(dir string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := utils.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, dir)
		dst := fmt.Sprintf("%s/%s", workingDir, dir)

		err = utils.MoveDirectory(src, dst)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveFileMsg(files)
	}
}

func (m model) moveFile(file string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := utils.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, file)
		dst := fmt.Sprintf("%s/%s", workingDir, file)

		err = utils.MoveFile(src, dst)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)

		if err != nil {
			return errorMsg(err.Error())
		}

		return moveFileMsg(files)
	}
}

func (m model) deleteDir(dir string) tea.Cmd {
	return func() tea.Msg {
		err := utils.DeleteDirectory(dir)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryMsg(files)
	}
}

func (m model) deleteFile(file string) tea.Cmd {
	return func() tea.Msg {
		err := utils.DeleteFile(file)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryMsg(files)
	}
}

func (m model) readFileContent(file fs.FileInfo) tea.Cmd {
	cfg := config.GetConfig()
	width := m.secondaryPane.Width

	return func() tea.Msg {
		content, err := utils.ReadFileContent(file.Name())
		if err != nil {
			return errorMsg(err.Error())
		}

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

func (m model) createDir(dir string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CreateDirectory(dir)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryMsg(files)
	}
}

func (m model) createFile(name string) tea.Cmd {
	return func() tea.Msg {
		err := utils.CreateFile(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := utils.GetDirectoryListing(constants.CurrentDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return directoryMsg(files)
	}
}
