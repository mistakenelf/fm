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

	if !m.Active {
		return m, nil
	}

	switch msg := msg.(type) {
	case errorMsg:
		cmds = append(cmds, m.NewStatusMessage(lipgloss.NewStyle().Foreground(lipgloss.Color("#cc241d")).Bold(true).Render(string(msg))))
	case statusMessageTimeoutMsg:
		m.StatusMessage = ""
	case copyToClipboardMsg:
		m.StatusMessage = string(msg)
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
			m.Cursor++
			if m.Cursor >= len(m.files) {
				m.Cursor = len(m.files) - 1
			}

			if m.Cursor > m.max {
				m.min++
				m.max++
			}
		case key.Matches(msg, m.keyMap.Up):
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = 0
			}

			if m.Cursor < m.min {
				m.min--
				m.max--
			}
		case key.Matches(msg, m.keyMap.GoToTop):
			m.Cursor = 0
			m.min = 0
			m.max = m.height
		case key.Matches(msg, m.keyMap.GoToBottom):
			m.Cursor = len(m.files) - 1
			m.min = len(m.files) - m.height
			m.max = len(m.files) - 1
		case key.Matches(msg, m.keyMap.PageDown):
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
			return m, getDirectoryListingCmd(filesystem.HomeDirectory, m.showHidden)
		case key.Matches(msg, m.keyMap.GoToRootDirectory):
			return m, getDirectoryListingCmd(filesystem.RootDirectory, m.showHidden)
		case key.Matches(msg, m.keyMap.ToggleHidden):
			m.showHidden = !m.showHidden

			return m, getDirectoryListingCmd(filesystem.CurrentDirectory, m.showHidden)
		case key.Matches(msg, m.keyMap.OpenDirectory):
			if m.files[m.Cursor].IsDirectory {
				return m, getDirectoryListingCmd(m.files[m.Cursor].Path, m.showHidden)
			}
		case key.Matches(msg, m.keyMap.PreviousDirectory):
			return m, getDirectoryListingCmd(filepath.Dir(m.files[m.Cursor].CurrentDirectory), m.showHidden)
		case key.Matches(msg, m.keyMap.CopyPathToClipboard):
			return m, copyToClipboardCmd(m.files[m.Cursor].Name)
		case key.Matches(msg, m.keyMap.CopyDirectoryItem):
			return m, tea.Sequence(
				copyDirectoryItemCmd(m.files[m.Cursor].Name, m.files[m.Cursor].IsDirectory),
				getDirectoryListingCmd(filesystem.CurrentDirectory, m.showHidden),
			)
		case key.Matches(msg, m.keyMap.DeleteDirectoryItem):
			return m, tea.Sequence(
				deleteDirectoryItemCmd(m.files[m.Cursor].Name, m.files[m.Cursor].IsDirectory),
				getDirectoryListingCmd(filesystem.CurrentDirectory, m.showHidden),
			)
		case key.Matches(msg, m.keyMap.ZipDirectoryItem):
			return m, tea.Sequence(
				zipDirectoryCmd(m.files[m.Cursor].Name),
				getDirectoryListingCmd(filesystem.CurrentDirectory, m.showHidden),
			)
		case key.Matches(msg, m.keyMap.UnzipDirectoryItem):
			return m, tea.Sequence(
				unzipDirectoryCmd(m.files[m.Cursor].Name),
				getDirectoryListingCmd(filesystem.CurrentDirectory, m.showHidden),
			)
		}
	}

	return m, tea.Batch(cmds...)
}
