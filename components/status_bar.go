package components

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/helpers"
	"github.com/knipferrc/fm/icons"

	"github.com/charmbracelet/lipgloss"
)

func StatusBar(screenWidth, cursor, totalFiles int, currentFile fs.FileInfo) string {
	cfg := config.GetConfig()
	doc := strings.Builder{}
	width := lipgloss.Width
	currentPath, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	statusItemStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.White)).
		Padding(0, 1)

	statusBarStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.White)).
		Background(lipgloss.Color(constants.DarkGray))

	selectedFileStyle := lipgloss.NewStyle().
		Inherit(statusBarStyle).
		Foreground(lipgloss.Color(constants.White)).
		Background(lipgloss.Color(constants.Pink)).
		Padding(0, 1).
		MarginRight(1)

	totalFilesStyle := statusItemStyle.Copy().
		Background(lipgloss.Color(constants.LightPurple)).
		Align(lipgloss.Right)

	statusText := lipgloss.NewStyle().Inherit(statusBarStyle)
	logoStyle := statusItemStyle.Copy().Background(lipgloss.Color(constants.DarkPurple))
	fileTotals := totalFilesStyle.Render(fmt.Sprintf("%d/%d", cursor+1, totalFiles))
	logo := ""

	if cfg.Settings.ShowIcons {
		logo = logoStyle.Render(fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM"))
	} else {
		logo = logoStyle.Render("FM")
	}

	status := statusText.Copy().
		Width(screenWidth - width(selectedFileStyle.Render(currentFile.Name())) - width(fileTotals) - width(logo)).
		Render(fmt.Sprintf("%s %s %s", helpers.ConvertBytes(currentFile.Size()), currentFile.Mode().String(), currentPath))

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileStyle.Render(currentFile.Name()),
		status,
		fileTotals,
		logo,
	)

	doc.WriteString(statusBarStyle.Render(bar))

	return doc.String()
}
