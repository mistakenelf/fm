package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/fm/statusbar"
)

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.filetree, cmd = m.filetree.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		halfSize := msg.Width / 2
		bubbleHeight := msg.Height - statusbar.Height

		cmds = append(cmds, m.image.SetSize(halfSize, bubbleHeight))
		cmds = append(cmds, m.markdown.SetSize(halfSize, bubbleHeight))

		m.filetree.SetSize(halfSize, bubbleHeight)
		m.help.SetSize(halfSize, bubbleHeight)
		m.code.SetSize(halfSize, bubbleHeight)
		m.pdf.SetSize(halfSize, bubbleHeight)
		m.statusbar.SetSize(msg.Width)

		m.filetree, cmd = m.filetree.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.OpenFile):
			cmds = append(cmds, tea.Batch(m.openFile()...))
		case key.Matches(msg, m.keys.TogglePane):
			m.togglePane()
		}
	}

	m.updateStatusbarContent()

	m.code, cmd = m.code.Update(msg)
	cmds = append(cmds, cmd)

	m.markdown, cmd = m.markdown.Update(msg)
	cmds = append(cmds, cmd)

	m.image, cmd = m.image.Update(msg)
	cmds = append(cmds, cmd)

	m.pdf, cmd = m.pdf.Update(msg)
	cmds = append(cmds, cmd)

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
