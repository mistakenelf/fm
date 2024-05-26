package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mistakenelf/fm/polish"
	"github.com/mistakenelf/fm/statusbar"
)

func (m *Model) disableAllViewports() {
	m.code.SetViewportDisabled(true)
	m.pdf.SetViewportDisabled(true)
	m.markdown.SetViewportDisabled(true)
	m.help.SetViewportDisabled(true)
	m.image.SetViewportDisabled(true)
	m.csv.SetViewportDisabled(true)
}

func (m *Model) resetViewports() {
	m.code.GotoTop()
	m.pdf.GotoTop()
	m.markdown.GotoTop()
	m.help.GotoTop()
	m.image.GotoTop()
	m.csv.GotoTop()
}

func (m *Model) updateStatusBar() {
	selectedItem := m.filetree[m.activeWorkspace].GetSelectedItem()

	if selectedItem.Name != "" {
		statusMessage := m.filetree[m.activeWorkspace].CurrentDirectory +
			lipgloss.NewStyle().
				Padding(0, 1).
				Foreground(polish.Colors.Yellow500).
				Render(selectedItem.Details)

		switch {
		case m.filetree[m.activeWorkspace].StatusMessage != "":
			statusMessage = m.filetree[m.activeWorkspace].StatusMessage
		case m.code.StatusMessage != "":
			statusMessage = m.code.StatusMessage
		case m.markdown.StatusMessage != "":
			statusMessage = m.markdown.StatusMessage
		case m.pdf.StatusMessage != "":
			statusMessage = m.pdf.StatusMessage
		case m.image.StatusMessage != "":
			statusMessage = m.image.StatusMessage
		case m.statusMessage != "":
			statusMessage = m.statusMessage
		case m.showTextInput:
			statusMessage = m.textinput.View()
		}

		m.statusbar.SetContent(
			selectedItem.Name,
			statusMessage,
			strings.Trim(strings.Join(strings.Fields(fmt.Sprint(m.workspaces)), ","), "[]"),
			fmt.Sprintf("%d/%d", m.filetree[m.activeWorkspace].Cursor+1, m.filetree[m.activeWorkspace].GetTotalItems()),
			fmt.Sprintf(selectedItem.FileSize),
		)
	} else {
		statusMessage := "Directory is empty"

		if m.showTextInput {
			statusMessage = m.textinput.View()
		}

		m.statusbar.SetContent(
			"N/A",
			statusMessage,
			fmt.Sprintf("%d", m.activeWorkspace+1),
			fmt.Sprintf("%d/%d", 0, 0),
			"FM",
		)
	}
}

func (m Model) handleWindowResizie(msg tea.WindowSizeMsg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.width = msg.Width
	m.height = msg.Height
	halfSize := msg.Width / 2
	height := msg.Height - statusbar.Height

	cmds = append(cmds, m.image.SetSizeCmd(halfSize, height))
	cmds = append(cmds, m.markdown.SetSizeCmd(halfSize, height))
	cmds = append(cmds, m.csv.SetSizeCmd(halfSize, height))

	m.filetree[m.activeWorkspace].SetSize(halfSize, height-3)
	m.secondaryFiletree.SetSize(halfSize, height-3)
	m.code.SetSize(halfSize, height)
	m.pdf.SetSize(halfSize, height)
	m.statusbar.SetSize(msg.Width)
	m.help.SetSize(halfSize, height)

	return m, tea.Batch(cmds...)
}
