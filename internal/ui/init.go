package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI and sets up initial data.
func (m Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.fileTree.Init())
	cmds = append(cmds, m.statusBar.Init())

	return tea.Batch(cmds...)
}
