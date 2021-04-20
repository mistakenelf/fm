package components

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/icons"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	selectedFileStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#FF5F87")).
				Padding(0, 1).
				MarginRight(1)

	fileEncodingStyle = statusItemStyle.Copy().
				Background(lipgloss.Color("#A550DF")).
				Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	logoStyle = statusItemStyle.Copy().Background(lipgloss.Color("#6124DF"))
)

func StatusBar(screenWidth int, currentFile fs.FileInfo, isMoving, isRenaming, isDeleting bool, textInput textinput.Model) string {
	doc := strings.Builder{}
	width := lipgloss.Width
	cfg := config.GetConfig()
	fileEncoding := fileEncodingStyle.Render("UTF-8")
	status := ""
	logo := ""

	if cfg.ShowIcons {
		logo = logoStyle.Render(fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM"))
	} else {
		logo = logoStyle.Render("FM")
	}

	if isMoving {
		status = fmt.Sprintf("%s %s", "Where would you like to move this to?", textInput.View())
	} else if isRenaming {
		status = fmt.Sprintf("%s %s", "What would you like to name this file?", textInput.View())
	} else if isDeleting {
		status = fmt.Sprintf("%s %s? [y/n] %s", "Are you sure you want to delete", currentFile.Name(), textInput.View())
	} else {
		status = "m - move, d - delete, r - rename, i - help"
	}

	status = statusText.Copy().
		Width(screenWidth - width(selectedFileStyle.Render(currentFile.Name())) - width(fileEncoding) - width(logo)).
		Render(status)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileStyle.Render(currentFile.Name()),
		status,
		fileEncoding,
		logo,
	)

	doc.WriteString(statusBarStyle.Render(bar))

	return doc.String()
}
