package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init intializes the UI.
func (m model) Init() tea.Cmd {
	return m.filetree.Init()
}
