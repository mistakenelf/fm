package components

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/knipferrc/fm/src/icons"

	"github.com/charmbracelet/bubbles/textinput"
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

func StatusBar(screenWidth int, currentFile fs.FileInfo, isMoving, isRenaming, isDeleting bool, textInput *textinput.Model) string {
	doc := strings.Builder{}
	w := lipgloss.Width

	statusKey := ""

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

	if currentFile != nil {
		statusKey = statusStyle.Render(currentFile.Name())
	}

	encoding := encodingStyle.Render("UTF-8")
	fm := fmStyle.Render(fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM"))

	statusVal := statusText.Copy().
		Width(screenWidth - w(statusKey) - w(encoding) - w(fm)).
		Render("M - Move, D - Delete, R - Rename")

	if isMoving {
		movePrompt := fmt.Sprintf("%s %s", "Where would you like to move this to?", textInput.View())
		statusVal = statusText.Copy().
			Width(screenWidth - w(statusKey) - w(encoding) - w(fm)).
			Render(movePrompt)
	}

	if isRenaming {
		movePrompt := fmt.Sprintf("%s %s", "What would you like to name this file?", textInput.View())
		statusVal = statusText.Copy().
			Width(screenWidth - w(statusKey) - w(encoding) - w(fm)).
			Render(movePrompt)
	}

	if isDeleting {
		movePrompt := fmt.Sprintf("%s %s %s", "Are you sure you want to delete this? [y/n]", currentFile.Name(), textInput.View())
		statusVal = statusText.Copy().
			Width(screenWidth - w(statusKey) - w(encoding) - w(fm)).
			Render(movePrompt)
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fm,
	)

	doc.WriteString(statusBarStyle.Render(bar))

	return doc.String()
}
