package filetree

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return m.GetDirectoryListingCmd(m.startDir)
}
