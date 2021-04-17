package main

import (
	"io/fs"

	"github.com/knipferrc/fm/src/components"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var file fs.FileInfo = nil
	view := ""

	if len(m.files) > 0 {
		file = m.files[m.cursor]
	}

	if !m.ready {
		return "Loading..."
	} else if m.showhelp {
		m.viewport.SetContent(components.Help())

		view = lipgloss.JoinVertical(lipgloss.Top, m.viewport.View(), components.StatusBar(m.screenwidth, file, m.move, m.rename, m.delete, &m.textinput))
	} else {
		m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

		view = lipgloss.JoinVertical(lipgloss.Top, m.viewport.View(), components.StatusBar(m.screenwidth, file, m.move, m.rename, m.delete, &m.textinput))
	}

	return view
}
