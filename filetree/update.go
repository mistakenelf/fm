package filetree

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mistakenelf/fm/filesystem"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

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
	case copyToClipboardMsg:
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Bold(true).
				Render(string(msg))))
	case createFileMsg:
		m.CreatingNewFile = false
		m.CreatingNewDirectory = false

		return m, getDirectoryListingCmd(
			filesystem.CurrentDirectory,
			m.showHidden,
			m.showDirectoriesOnly,
			m.showFilesOnly,
		)
	case createDirectoryMsg:
		m.CreatingNewDirectory = false
		m.CreatingNewFile = false

		return m, getDirectoryListingCmd(
			filesystem.CurrentDirectory,
			m.showHidden,
			m.showDirectoriesOnly,
			m.showFilesOnly,
		)
	case getDirectoryListingMsg:
		if msg != nil {
			m.files = msg
			m.Cursor = 0
			m.min = 0
			m.max = m.height
			m.max = max(m.max, m.height)
		}
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

			return m, getDirectoryListingCmd(
				filesystem.HomeDirectory,
				m.showHidden,
				m.showDirectoriesOnly,
				m.showFilesOnly,
			)
		case key.Matches(msg, m.keyMap.GoToRootDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, getDirectoryListingCmd(
				filesystem.RootDirectory,
				m.showHidden,
				m.showDirectoriesOnly,
				m.showFilesOnly,
			)
		case key.Matches(msg, m.keyMap.ToggleHidden):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.showHidden = !m.showHidden

			return m, getDirectoryListingCmd(
				filesystem.CurrentDirectory,
				m.showHidden,
				m.showDirectoriesOnly,
				m.showFilesOnly,
			)
		case key.Matches(msg, m.keyMap.OpenDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			if m.files[m.Cursor].IsDirectory {
				return m, getDirectoryListingCmd(
					m.files[m.Cursor].Path,
					m.showHidden,
					m.showDirectoriesOnly,
					m.showFilesOnly,
				)
			}
		case key.Matches(msg, m.keyMap.PreviousDirectory):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, getDirectoryListingCmd(
				filepath.Dir(m.files[m.Cursor].CurrentDirectory),
				m.showHidden,
				m.showDirectoriesOnly,
				m.showFilesOnly,
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
				getDirectoryListingCmd(
					filesystem.CurrentDirectory,
					m.showHidden,
					m.showDirectoriesOnly,
					m.showFilesOnly,
				),
			)
		case key.Matches(msg, m.keyMap.DeleteDirectoryItem):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, tea.Sequence(
				deleteDirectoryItemCmd(m.files[m.Cursor].Name, m.files[m.Cursor].IsDirectory),
				getDirectoryListingCmd(
					filesystem.CurrentDirectory,
					m.showHidden,
					m.showDirectoriesOnly,
					m.showFilesOnly,
				),
			)
		case key.Matches(msg, m.keyMap.ZipDirectoryItem):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, tea.Sequence(
				zipDirectoryCmd(m.files[m.Cursor].Name),
				getDirectoryListingCmd(
					filesystem.CurrentDirectory,
					m.showHidden,
					m.showDirectoriesOnly,
					m.showFilesOnly,
				),
			)
		case key.Matches(msg, m.keyMap.UnzipDirectoryItem):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			return m, tea.Sequence(
				unzipDirectoryCmd(m.files[m.Cursor].Name),
				getDirectoryListingCmd(
					filesystem.CurrentDirectory,
					m.showHidden,
					m.showDirectoriesOnly,
					m.showFilesOnly,
				),
			)
		case key.Matches(msg, m.keyMap.ShowDirectoriesOnly):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.showDirectoriesOnly = !m.showDirectoriesOnly
			m.showFilesOnly = false

			return m, getDirectoryListingCmd(
				filesystem.CurrentDirectory,
				m.showHidden,
				m.showDirectoriesOnly,
				m.showFilesOnly,
			)
		case key.Matches(msg, m.keyMap.ShowFilesOnly):
			if m.CreatingNewFile || m.CreatingNewDirectory {
				return m, nil
			}

			m.showFilesOnly = !m.showFilesOnly
			m.showDirectoriesOnly = false

			return m, getDirectoryListingCmd(
				filesystem.CurrentDirectory,
				m.showHidden,
				m.showDirectoriesOnly,
				m.showFilesOnly,
			)
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
