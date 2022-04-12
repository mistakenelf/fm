package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI and sets up initial data.
// Init intializes the UI.
func (b Bubble) Init() tea.Cmd {
	return b.filetree.Init()
}
