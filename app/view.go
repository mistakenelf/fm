package app

import (
	"fmt"

	"github.com/knipferrc/fm/components"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.Ready || len(m.Files) <= 0 {
		return fmt.Sprintf("%s%s", m.Spinner.View(), "loading...")
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.Viewport.View(),
		components.StatusBar(m.ScreenWidth, m.Cursor, len(m.Files), m.Files[m.Cursor], m.Move, m.Rename, m.Delete, m.Textinput),
	)
}
