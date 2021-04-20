package components

import (
	"fmt"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/icons"

	"github.com/charmbracelet/lipgloss"
)

func DirItem(selected, isDir bool, name, ext, indicator string) string {
	cfg := config.GetConfig()

	if !cfg.ShowIcons && !selected {
		return name
	} else if !cfg.ShowIcons && selected {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.SelectedItemColor)).Render(name)
	} else if selected && isDir {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.SelectedItemColor)).Render(name))

		return lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.SelectedItemColor)).Render(listing)
	} else if !selected && isDir {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, name)

		return lipgloss.NewStyle().Render(listing)
	} else if selected && !isDir {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.SelectedItemColor)).Render(name))

		return listing
	} else {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, name)

		return listing
	}
}
