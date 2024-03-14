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
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.OpenFile):
			cmds = append(cmds, m.openFile())
		case key.Matches(msg, m.keyMap.ResetState):
			m.state = idleState
			m.disableAllViewports()
			m.resetViewports()
		case key.Matches(msg, m.keyMap.TogglePane):
			m.activePane = (m.activePane + 1) % 2

			if m.activePane == 0 {
				m.filetree.SetDisabled(false)
				m.disableAllViewports()
			} else {
				m.filetree.SetDisabled(true)

				switch m.state {
				case idleState:
					m.disableAllViewports()
					m.help.SetViewportDisabled(false)
				case showCodeState:
					m.disableAllViewports()
					m.code.SetViewportDisabled(false)
				case showImageState:
					m.disableAllViewports()
					m.image.SetViewportDisabled(false)
				case showPdfState:
					m.disableAllViewports()
					m.pdf.SetViewportDisabled(false)
				case showMarkdownState:
					m.disableAllViewports()
					m.markdown.SetViewportDisabled(false)
				}
			}
		}
	}

	m.filetree, cmd = m.filetree.Update(msg)
	cmds = append(cmds, cmd)

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

	m.updateStatusbarContent()

	return m, tea.Batch(cmds...)
}
