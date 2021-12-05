package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI and sets up initial data.
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.fileTree.Init())

	return tea.Batch(cmds...)
}
