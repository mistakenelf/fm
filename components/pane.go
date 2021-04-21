package components

import "github.com/charmbracelet/lipgloss"

func Pane(width int, content string) string {
	return lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Width(width).
		Render(content)
}
