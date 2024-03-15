package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Init intializes the UI.
func (m model) Init() tea.Cmd {
	return tea.Batch(m.filetree.Init(), textinput.Blink)
}
