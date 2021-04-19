package components

import (
	"fmt"

	"github.com/knipferrc/fm/src/config"
	"github.com/knipferrc/fm/src/icons"

	"github.com/charmbracelet/lipgloss"
)

func DirItem(label string, selected, isDir bool, ext string) string {
	config := config.GetConfig()

	if !config.ShowIcons && !selected {
		return label
	} else if !config.ShowIcons && selected {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(config.SelectedItemColor)).Render(label)
	} else if selected && isDir {
		listing := fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), label)

		return lipgloss.NewStyle().Foreground(lipgloss.Color(config.SelectedItemColor)).Render(listing)
	} else if !selected && isDir {
		listing := fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), label)

		return lipgloss.NewStyle().Render(listing)
	} else if selected && !isDir {
		listing := fmt.Sprintf("%s %s", icons.Icon_Def["file"].GetGlyph(), label)

		return lipgloss.NewStyle().Foreground(lipgloss.Color(config.SelectedItemColor)).Render(listing)
	} else {
		listing := fmt.Sprintf("%s %s", icons.Icon_Def["file"].GetGlyph(), label)

		return listing
	}
}
