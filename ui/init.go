package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// App initialization, enable the text input as well as spinner
func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, spinner.Tick)
}
