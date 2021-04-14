package components

import (
	"fmt"
	"strings"

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

func StatusBar(screenWidth int, currentFile string) string {
	doc := strings.Builder{}
	w := lipgloss.Width

	statusNugget := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1)

	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle := lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF5F87")).
		Padding(0, 1).
		MarginRight(1)

	encodingStyle := statusNugget.Copy().
		Background(lipgloss.Color("#A550DF")).
		Align(lipgloss.Right)

	statusText := lipgloss.NewStyle().Inherit(statusBarStyle)

	fmStyle := statusNugget.Copy().Background(lipgloss.Color("#6124DF"))

	statusKey := statusStyle.Render(currentFile)
	encoding := encodingStyle.Render("UTF-8")
	fm := fmStyle.Render(fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM"))
	statusVal := statusText.Copy().
		Width(screenWidth - w(statusKey) - w(encoding) - w(fm)).
		Render("M - Move, D - Delete, R - Rename")

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fm,
	)

	doc.WriteString(statusBarStyle.Render(bar))

	return doc.String()
}
