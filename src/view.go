package main

import (
	"fmt"

	"github.com/knipferrc/fm/src/components"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if !m.ready || len(m.files) <= 0 {
		return fmt.Sprintf("%s%s", m.spinner.View(), "loading...")
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		components.StatusBar(m.screenwidth, m.files[m.cursor], m.move, m.rename, m.delete, m.textinput),
	)
}
