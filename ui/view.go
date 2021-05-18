package ui

import (
	"fmt"

	"github.com/knipferrc/fm/statusbar"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.Ready || len(m.Files) <= 0 {
		return fmt.Sprintf("%s%s", m.Spinner.View(), "loading...")
	}

	leftPane := m.PrimaryPane.View()
	rightPane := m.SecondaryPane.View()
	panes := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		panes,
		statusbar.View(
			m.ScreenWidth,
			m.Cursor,
			len(m.Files),
			m.Files[m.Cursor],
			m.ShowCommandBar,
			m.Textinput.View(),
		),
	)
}
