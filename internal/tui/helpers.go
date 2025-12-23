package tui

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/fm/polish"
)

func contains(s []string, str string) bool {
	return slices.Contains(s, str)
}

func (m *model) disableAllViewports() {
	m.code.SetViewportDisabled(true)
	m.pdf.SetViewportDisabled(true)
	m.markdown.SetViewportDisabled(true)
	m.help.SetViewportDisabled(true)
	m.image.SetViewportDisabled(true)
	m.csv.SetViewportDisabled(true)
}

func (m *model) resetViewports() {
	m.code.GotoTop()
	m.pdf.GotoTop()
	m.markdown.GotoTop()
	m.help.GotoTop()
	m.image.GotoTop()
	m.csv.GotoTop()
}

func (m *model) updateStatusBar() {
	if m.filetree.GetSelectedItem().Name != "" {
		statusMessage :=
			m.filetree.CurrentDirectory +
				lipgloss.NewStyle().
					Padding(0, 1).
					Foreground(polish.Colors.Yellow500).
					Render(m.filetree.GetSelectedItem().Details) +
				lipgloss.NewStyle().
					Padding(0, 1).
					Foreground(polish.Colors.Pink).
					Render(m.filetree.GetSelectedItem().
						FileInfo.
						ModTime().Format("2006-01-02"))

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
			fmt.Sprintf("%s", m.filetree.GetSelectedItem().FileSize),
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
