package app

import (
	"fmt"

	"github.com/knipferrc/fm/components"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	border := lipgloss.NormalBorder()
	halfScreenWidth := m.ScreenWidth / 2

	if !m.Ready || len(m.Files) <= 0 {
		return fmt.Sprintf("%s%s", m.Spinner.View(), "loading...")
	}

	leftPane := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(halfScreenWidth).
		Render(m.Viewport.View())

	rightPane := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(halfScreenWidth - lipgloss.Width(border.Left+border.Right+border.Top+border.Bottom)).
		Render(m.SecondaryViewport.View())

	panes := lipgloss.JoinHorizontal(0, leftPane, rightPane)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		panes,
		components.StatusBar(m.ScreenWidth, m.Cursor, len(m.Files), m.Files[m.Cursor], m.Move, m.Rename, m.Delete, m.Textinput),
	)
}
