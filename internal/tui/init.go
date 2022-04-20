package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI and sets up initial data.
// Init intializes the UI.
func (b Bubble) Init() tea.Cmd {
	initTreeCmd := b.filetree.Init()
	toggleTreeIconsCmd := b.filetree.ToggleShowIcons(b.config.Settings.ShowIcons)

	return tea.Sequentially(initTreeCmd, toggleTreeIconsCmd)
}
