package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Init intializes the UI.
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.filetree[m.activeWorkspace].Init(),
		m.secondaryFiletree.Init(),
		textinput.Blink,
		tea.SetWindowTitle("FM"),
	)
}
