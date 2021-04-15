package main

import (
	"io/fs"

	"github.com/knipferrc/fm/src/components"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	var file fs.FileInfo = nil

	if len(m.Files) > 0 {
		file = m.Files[m.Cursor]
	}

	if m.ShowHelp {
		view := lipgloss.JoinVertical(lipgloss.Top, m.Viewport.View(), components.StatusBar(m.ScreenWidth, file, m.Move, m.Rename, m.Delete, &m.TextInput))

		return view
	} else {
		m.Viewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))

		view := lipgloss.JoinVertical(lipgloss.Top, m.Viewport.View(), components.StatusBar(m.ScreenWidth, file, m.Move, m.Rename, m.Delete, &m.TextInput))

		return view
	}
}
