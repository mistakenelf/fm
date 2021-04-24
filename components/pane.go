package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/fm/config"
)

func Pane(width int, isActive bool, content string) string {
	cfg := config.GetConfig()
	borderColor := cfg.Colors.InactivePane

	if isActive {
		borderColor = cfg.Colors.ActivePane
	}

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		Border(lipgloss.NormalBorder()).
		Width(width).
		Render(content)
}
