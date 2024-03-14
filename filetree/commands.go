package filetree

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/fm/filesystem"
)

type getDirectoryListingMsg []DirectoryItem
type errorMsg string
type copyToClipboardMsg string
type statusMessageTimeoutMsg struct{}
type editorFinishedMsg struct{ err error }

// getDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func getDirectoryListingCmd(directoryName string, showHidden, directoriesOnly, filesOnly bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		var directoryItems []DirectoryItem
		var files []fs.DirEntry

		if directoryName == filesystem.HomeDirectory {
			directoryName, err = filesystem.GetHomeDirectory()
			if err != nil {
				return errorMsg(err.Error())
			}
		}

		directoryInfo, err := os.Stat(directoryName)
		if err != nil {
			return errorMsg(err.Error())
		}

		if !directoryInfo.IsDir() {
			return nil
		}

		err = os.Chdir(directoryName)
		if err != nil {
			return errorMsg(err.Error())
		}

		if !directoriesOnly && !filesOnly {
			files, err = filesystem.GetDirectoryListing(directoryName, showHidden)
			if err != nil {
				return errorMsg(err.Error())
			}
		} else {
			listingType := filesystem.DirectoriesListingType

			if filesOnly {
				listingType = filesystem.FilesListingType
			}

			files, err = filesystem.GetDirectoryListingByType(directoryName, listingType, showHidden)
			if err != nil {
				return errorMsg(err.Error())
			}
		}

		workingDirectory, err := filesystem.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			status := fmt.Sprintf("%s %s",
				ConvertBytesToSizeString(fileInfo.Size()),
				fileInfo.Mode().String())

			directoryItems = append(directoryItems, DirectoryItem{
				Name:             file.Name(),
				Details:          status,
				Path:             filepath.Join(workingDirectory, file.Name()),
				Extension:        filepath.Ext(fileInfo.Name()),
				IsDirectory:      fileInfo.IsDir(),
				CurrentDirectory: workingDirectory,
			})
		}

		return getDirectoryListingMsg(directoryItems)
	}
}

// deleteDirectoryItemCmd deletes a directory based on the name provided.
func deleteDirectoryItemCmd(name string, isDirectory bool) tea.Cmd {
	return func() tea.Msg {
		if isDirectory {
			if err := filesystem.DeleteDirectory(name); err != nil {
				return errorMsg(err.Error())
			}
		} else {
			if err := filesystem.DeleteFile(name); err != nil {
				return errorMsg(err.Error())
			}
		}

		return nil
	}
}

// createDirectoryCmd creates a directory based on the name provided.
func createDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := filesystem.CreateDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// createFileCmd creates a file based on the name provided.
func createFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := filesystem.CreateFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// zipDirectoryCmd zips a directory based on the name provided.
func zipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := filesystem.Zip(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// unzipDirectoryCmd unzips a directory based on the name provided.
func unzipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := filesystem.Unzip(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyDirectoryItemCmd copies a directory based on the name provided.
func copyDirectoryItemCmd(name string, isDirectory bool) tea.Cmd {
	return func() tea.Msg {
		if isDirectory {
			if err := filesystem.CopyDirectory(name); err != nil {
				return errorMsg(err.Error())
			}
		} else {
			if err := filesystem.CopyFile(name); err != nil {
				return errorMsg(err.Error())
			}
		}

		return nil
	}
}

// copyToClipboardCmd copies the provided string to the clipboard.
func copyToClipboardCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := filesystem.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		filePath := filepath.Join(workingDir, name)
		err = clipboard.WriteAll(filePath)
		if err != nil {
			return errorMsg(err.Error())
		}

		return copyToClipboardMsg(fmt.Sprintf("%s %s %s", "Successfully copied", filePath, "to clipboard"))
	}
}

// writeSelectionPathCmd writes content to the file specified.
func writeSelectionPathCmd(selectionPath, filePath string) tea.Cmd {
	return func() tea.Msg {
		if err := filesystem.WriteToFile(selectionPath, filePath); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time. Note that this also returns a command.
func (m *Model) NewStatusMessage(s string) tea.Cmd {
	m.StatusMessage = s

	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}

	m.statusMessageTimer = time.NewTimer(m.StatusMessageLifetime)

	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}

func openEditorCmd(file string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	c := exec.Command(editor, file)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}
