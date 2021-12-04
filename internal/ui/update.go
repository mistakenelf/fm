package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// updateStatusBarContent updates the content of the statusbar.
func (m *Model) updateStatusBarContent() {
	selectedFile, _ := m.fileTree.GetSelectedFile()

	m.statusBar.SetContent(
		m.fileTree.GetTotalFiles(),
		m.fileTree.GetCursor(),
		false,
		false,
		selectedFile,
		nil,
		m.fileTree.GetFilePaths(),
	)
}

// handleWindowSizeMsg is received whenever the window size changes.
func (m *Model) handleWindowSizeMsg() {
	if !m.ready {
		m.ready = true
	}
}

// Update handles all UI interactions and events for updating the screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowSizeMsg()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	m.updateStatusBarContent()

	m.fileTree, cmd = m.fileTree.Update(msg)
	cmds = append(cmds, cmd)

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)

	m.loader, cmd = m.loader.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
