package app

import (
	"fmt"

	"github.com/knipferrc/fm/components"
	"github.com/knipferrc/fm/constants"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	border := lipgloss.NormalBorder()
	borderRightWidth := lipgloss.Width(border.Right + border.Top)
	borderLeftWidth := lipgloss.Width(border.Left + border.Top)
	halfScreenWidth := m.ScreenWidth / 2
	leftPaneActive := false
	rightPaneActive := false

	if m.ActivePane == constants.PrimaryPane {
		leftPaneActive = true
	} else {
		rightPaneActive = true
	}

	if !m.Ready || len(m.Files) <= 0 {
		return fmt.Sprintf("%s%s", m.Spinner.View(), "loading...")
	}

	leftPane := components.Pane(halfScreenWidth-borderLeftWidth, leftPaneActive, m.PrimaryViewport.View())
	rightPane := components.Pane(halfScreenWidth-borderRightWidth, rightPaneActive, m.SecondaryViewport.View())
	panes := lipgloss.JoinHorizontal(0, leftPane, rightPane)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		panes,
		components.StatusBar(m.ScreenWidth, m.Cursor, len(m.Files), m.Files[m.Cursor]),
	)
}
