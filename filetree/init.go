package filetree

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return getDirectoryListingCmd(m.startDir, true, false, false)
}
