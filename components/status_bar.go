package components

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/helpers"
	"github.com/knipferrc/fm/icons"

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

	totalFilesStyle = statusItemStyle.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	logoStyle = statusItemStyle.Copy().Background(lipgloss.Color("#6124DF"))
)

func StatusBar(screenWidth, cursor, totalFiles int, currentFile fs.FileInfo) string {
	doc := strings.Builder{}
	width := lipgloss.Width
	cfg := config.GetConfig()
	fileTotals := totalFilesStyle.Render(fmt.Sprintf("%d/%d", cursor+1, totalFiles))
	logo := ""

	if cfg.Settings.ShowIcons {
		logo = logoStyle.Render(fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM"))
	} else {
		logo = logoStyle.Render("FM")
	}

	status := statusText.Copy().
		Width(screenWidth - width(selectedFileStyle.Render(currentFile.Name())) - width(fileTotals) - width(logo)).
		Render(fmt.Sprintf("%s %s", helpers.ConvertBytes(currentFile.Size()), currentFile.Mode().String()))

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileStyle.Render(currentFile.Name()),
		status,
		fileTotals,
		logo,
	)

	doc.WriteString(statusBarStyle.Render(bar))

	return doc.String()
}
