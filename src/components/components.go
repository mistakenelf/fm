package components

import (
	"fmt"

	"github.com/knipferrc/fm/src/icons"

	"github.com/charmbracelet/lipgloss"
)

func FileListing(label string, selected bool, isDir bool, ext string) string {
	if selected && isDir {
		listing := fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), label)

		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(listing)
	} else if isDir && !selected {
		listing := fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), label)

		return lipgloss.NewStyle().Render(listing)
	} else if !isDir && selected {
		listing := fmt.Sprintf("%s %s", icons.Icon_Def["file"].GetGlyph(), label)

		return lipgloss.NewStyle().Foreground(lipgloss.Color("#F25D94")).Render(listing)
	} else {
		return fmt.Sprintf(lipgloss.NewStyle().Render("\uf723 %s"), label)
	}
}
