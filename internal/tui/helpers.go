package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (m *model) disableAllViewports() {
	m.code.SetViewportDisabled(true)
	m.pdf.SetViewportDisabled(true)
	m.markdown.SetViewportDisabled(true)
	m.help.SetViewportDisabled(true)
	m.image.SetViewportDisabled(true)
}

func (m *model) resetViewports() {
	m.code.GotoTop()
	m.pdf.GotoTop()
	m.markdown.GotoTop()
	m.help.GotoTop()
	m.image.GotoTop()
}

func (m *model) updateStatusBar() {
	if m.filetree.GetSelectedItem().Name != "" {
		statusMessage :=
			m.filetree.CurrentDirectory +
				lipgloss.NewStyle().
					Padding(0, 1).
					Foreground(lipgloss.Color("#eab308")).
					Render(m.filetree.GetSelectedItem().Details)

		if m.filetree.StatusMessage != "" {
			statusMessage = m.filetree.StatusMessage
		}

		if m.code.StatusMessage != "" {
			statusMessage = m.code.StatusMessage
		}

		if m.markdown.StatusMessage != "" {
			statusMessage = m.markdown.StatusMessage
		}

		if m.pdf.StatusMessage != "" {
			statusMessage = m.pdf.StatusMessage
		}

		if m.image.StatusMessage != "" {
			statusMessage = m.image.StatusMessage
		}

		if m.statusMessage != "" {
			statusMessage = m.statusMessage
		}

		if m.showTextInput {
			statusMessage = m.textinput.View()
		}

		m.statusbar.SetContent(
			m.filetree.GetSelectedItem().Name,
			statusMessage,
			fmt.Sprintf("%d/%d", m.filetree.Cursor+1, m.filetree.GetTotalItems()),
			fmt.Sprintf(m.filetree.GetSelectedItem().FileSize),
		)
	} else {
		statusMessage := "Directory is empty"

		if m.showTextInput {
			statusMessage = m.textinput.View()
		}

		m.statusbar.SetContent(
			"N/A",
			statusMessage,
			fmt.Sprintf("%d/%d", 0, 0),
			"FM",
		)
	}
}
