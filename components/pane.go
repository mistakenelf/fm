package components

import (
	"github.com/knipferrc/fm/config"

	"github.com/charmbracelet/lipgloss"
)

func Pane(width, height int, isActive bool, content string) string {
	cfg := config.GetConfig()
	borderColor := cfg.Colors.Pane.InactivePane

	if isActive {
		borderColor = cfg.Colors.Pane.ActivePane
	}

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		Border(lipgloss.NormalBorder()).
		Width(width).
		Height(height).
		Render(content)
}
