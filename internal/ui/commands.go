package ui

import (
	"io/fs"

	"github.com/knipferrc/fm/internal/constants"
	"github.com/knipferrc/fm/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type updateDirMsg []fs.FileInfo
type actionMsg []fs.FileInfo
type fileContentMsg string

func updateDirectoryListing(dir string) tea.Cmd {
	return func() tea.Msg {
		files := utils.GetDirectoryListing(dir)

		return updateDirMsg(files)
	}
}

func renameFileOrDir(filename, value string) tea.Cmd {
	return func() tea.Msg {
		utils.RenameDirOrFile(filename, value)
		files := utils.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func moveDir(dir, value string) tea.Cmd {
	return func() tea.Msg {
		utils.CopyDir(dir, value, true)
		files := utils.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func moveFile(file, value string) tea.Cmd {
	return func() tea.Msg {
		utils.CopyFile(file, value, true)
		files := utils.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func deleteDir(dir string) tea.Cmd {
	return func() tea.Msg {
		utils.DeleteDirectory(dir)
		files := utils.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func deleteFile(file string) tea.Cmd {
	return func() tea.Msg {
		utils.DeleteFile(file)
		files := utils.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func readFileContent(file string) tea.Cmd {
	return func() tea.Msg {
		content := utils.ReadFileContent(file)

		return fileContentMsg(content)
	}
}

func createDir(dir string) tea.Cmd {
	return func() tea.Msg {
		utils.CreateDirectory(dir)
		files := utils.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}
