package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all UI interactions and events for updating the screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.fileTree.GetIsActive() {
				m.fileTree.SetIsActive(false)
				m.renderer.SetIsActive(true)
			} else {
				m.fileTree.SetIsActive(true)
				m.renderer.SetIsActive(false)
			}
		}
	}

	m.fileTree, cmd = m.fileTree.Update(msg)
	cmds = append(cmds, cmd)

	m.renderer, cmd = m.renderer.Update(msg)
	cmds = append(cmds, cmd)

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
