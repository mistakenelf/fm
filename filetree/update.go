package filetree

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mistakenelf/fm/filesystem"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.Disabled {
		return m, nil
	}

	switch msg := msg.(type) {
	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
	case errorMsg:
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cc241d")).
				Bold(true).
				Render(string(msg))))
	case statusMessageTimeoutMsg:
		m.StatusMessage = ""
	case moveDirectoryItemMsg:
		m.CreatingNewFile = false
		m.CreatingNewDirectory = false

		return m, m.GetDirectoryListingCmd(filesystem.CurrentDirectory)
	case copyToClipboardMsg:
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Bold(true).
				Render(string(msg))))
	case createFileMsg:
		m.CreatingNewFile = false
		m.CreatingNewDirectory = false

		return m, m.GetDirectoryListingCmd(filesystem.CurrentDirectory)
	case createDirectoryMsg:
		m.CreatingNewDirectory = false
		m.CreatingNewFile = false

		return m, m.GetDirectoryListingCmd(filesystem.CurrentDirectory)
	case getDirectoryListingMsg:
		if msg.files != nil {
			m.files = msg.files
		} else {
			m.files = make([]DirectoryItem, 0)
		}

		m.CurrentDirectory = msg.workingDirectory
		m.Cursor = 0
		m.min = 0
		m.max = max(m.max, m.height-1)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Down):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.Cursor++

			if m.Cursor >= len(m.files) {
				m.Cursor = len(m.files) - 1
			}

			if m.Cursor > m.max {
				m.min++
				m.max++
			}
		case key.Matches(msg, m.keyMap.Up):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.Cursor--

			if m.Cursor < 0 {
				m.Cursor = 0
			}

			if m.Cursor < m.min {
				m.min--
				m.max--
			}
		case key.Matches(msg, m.keyMap.GoToTop):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.Cursor = 0
			m.min = 0
			m.max = m.height
		case key.Matches(msg, m.keyMap.GoToBottom):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.Cursor = len(m.files) - 1
			m.min = len(m.files) - m.height
			m.max = len(m.files) - 1
		case key.Matches(msg, m.keyMap.PageDown):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.Cursor += m.height

			if m.Cursor >= len(m.files) {
				m.Cursor = len(m.files) - 1
			}

			m.min += m.height
			m.max += m.height

			if m.max >= len(m.files) {
				m.max = len(m.files) - 1
				m.min = m.max - m.height
			}
		case key.Matches(msg, m.keyMap.PageUp):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.Cursor -= m.height

			if m.Cursor < 0 {
				m.Cursor = 0
			}

			m.min -= m.height
			m.max -= m.height

			if m.min < 0 {
				m.min = 0
				m.max = m.min + m.height
			}
		case key.Matches(msg, m.keyMap.GoToHomeDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, m.GetDirectoryListingCmd(filesystem.HomeDirectory)
		case key.Matches(msg, m.keyMap.GoToRootDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, m.GetDirectoryListingCmd(filesystem.RootDirectory)
		case key.Matches(msg, m.keyMap.ToggleHidden):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.showHidden = !m.showHidden

			return m, m.GetDirectoryListingCmd(filesystem.CurrentDirectory)
		case key.Matches(msg, m.keyMap.OpenDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			if m.files[m.Cursor].IsDirectory {
				return m, m.GetDirectoryListingCmd(m.files[m.Cursor].Path)
			}
		case key.Matches(msg, m.keyMap.PreviousDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			if len(m.files) == 0 {
				return m, m.GetDirectoryListingCmd(filepath.Dir(m.CurrentDirectory))
			}

			return m, m.GetDirectoryListingCmd(
				filepath.Dir(m.files[m.Cursor].CurrentDirectory),
			)
		case key.Matches(msg, m.keyMap.CopyPathToClipboard):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, copyToClipboardCmd(m.files[m.Cursor].Name)
		case key.Matches(msg, m.keyMap.CopyDirectoryItem):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, tea.Sequence(
				copyDirectoryItemCmd(m.files[m.Cursor].Name, m.files[m.Cursor].IsDirectory),
				m.GetDirectoryListingCmd(filesystem.CurrentDirectory),
			)
		case key.Matches(msg, m.keyMap.DeleteDirectoryItem):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, tea.Sequence(
				deleteDirectoryItemCmd(m.files[m.Cursor].Name, m.files[m.Cursor].IsDirectory),
				m.GetDirectoryListingCmd(filesystem.CurrentDirectory),
			)
		case key.Matches(msg, m.keyMap.ZipDirectoryItem):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, tea.Sequence(
				zipDirectoryCmd(m.files[m.Cursor].Name),
				m.GetDirectoryListingCmd(filesystem.CurrentDirectory),
			)
		case key.Matches(msg, m.keyMap.UnzipDirectoryItem):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, tea.Sequence(
				unzipDirectoryCmd(m.files[m.Cursor].Name),
				m.GetDirectoryListingCmd(filesystem.CurrentDirectory),
			)
		case key.Matches(msg, m.keyMap.ShowDirectoriesOnly):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.showDirectoriesOnly = !m.showDirectoriesOnly
			m.showFilesOnly = false

			return m, m.GetDirectoryListingCmd(filesystem.CurrentDirectory)
		case key.Matches(msg, m.keyMap.ShowFilesOnly):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.showFilesOnly = !m.showFilesOnly
			m.showDirectoriesOnly = false

			return m, m.GetDirectoryListingCmd(filesystem.CurrentDirectory)
		case key.Matches(msg, m.keyMap.WriteSelectionPath):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			if m.selectionPath != "" {
				return m, tea.Sequence(
					writeSelectionPathCmd(m.selectionPath, m.files[m.Cursor].Name),
					tea.Quit,
				)
			}
		case key.Matches(msg, m.keyMap.OpenInEditor):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, openEditorCmd(m.files[m.Cursor].Name)
		case key.Matches(msg, m.keyMap.CreateFile):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.CreatingNewFile = true
			m.CreatingNewDirectory = false

			return m, nil
		case key.Matches(msg, m.keyMap.CreateDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.CreatingNewDirectory = true
			m.CreatingNewFile = false

			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}
