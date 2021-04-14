package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/fm/src/components"
)

func (m model) View() string {
	doc := strings.Builder{}
	files := ""

	for i, file := range m.Files {
		files += fmt.Sprintf("%s\n", components.FileListing(file.Name(), m.Cursor == i, file.IsDir(), filepath.Ext(file.Name())))
	}

	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Width(50).Align(lipgloss.Left).Render(files),
	))

	m.Viewport.SetContent(doc.String())

	view := fmt.Sprintf("%s%s", m.Viewport.View(), components.StatusBar(m.ScreenWidth, m.CurrentlyHighlighted))

	return view
}
