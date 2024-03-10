package filetree

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height - 2
		m.width = msg.Width
		m.max = m.height - 1
	case getDirectoryListingMsg:
		if msg != nil {
			m.files = msg
			m.max = max(m.max, m.height-1)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Down):
			m.Cursor++
			if m.Cursor >= len(m.files) {
				m.Cursor = len(m.files) - 1
			}

			if m.Cursor > m.max {
				m.min++
				m.max++
			}
		case key.Matches(msg, m.keyMap.Up):
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = 0
			}

			if m.Cursor < m.min {
				m.min--
				m.max--
			}
		}
	}

	return m, tea.Batch(cmds...)
}
