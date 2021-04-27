package app

import (
	"io/fs"

	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/filesystem"

	tea "github.com/charmbracelet/bubbletea"
)

type updateDirMsg []fs.FileInfo
type actionMsg []fs.FileInfo
type fileContentMsg string

func updateDirectoryListing(dir string) tea.Cmd {
	return func() tea.Msg {
		files := filesystem.GetDirectoryListing(dir)

		return updateDirMsg(files)
	}
}

func renameFileOrDir(filename, value string) tea.Cmd {
	return func() tea.Msg {
		filesystem.RenameDirOrFile(filename, value)
		files := filesystem.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func moveDir(dir, value string) tea.Cmd {
	return func() tea.Msg {
		filesystem.CopyDir(dir, value, true)
		files := filesystem.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func moveFile(file, value string) tea.Cmd {
	return func() tea.Msg {
		filesystem.CopyFile(file, value, true)
		files := filesystem.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func deleteDir(dir string) tea.Cmd {
	return func() tea.Msg {
		filesystem.DeleteDirectory(dir)
		files := filesystem.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func deleteFile(file string) tea.Cmd {
	return func() tea.Msg {
		filesystem.DeleteFile(file)
		files := filesystem.GetDirectoryListing(constants.CurrentDirectory)

		return actionMsg(files)
	}
}

func readFileContent(file string) tea.Cmd {
	return func() tea.Msg {
		content := filesystem.ReadFileContent(file)

		return fileContentMsg(content)
	}
}
