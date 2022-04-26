package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init intializes the UI.
func (b Bubble) Init() tea.Cmd {
	return b.filetree.Init()
}
