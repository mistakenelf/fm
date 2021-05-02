package components

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/knipferrc/fm/config"
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

	selectedFile := lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.Colors.StatusBar.SelectedFile.Foreground)).
		Background(lipgloss.Color(cfg.Colors.StatusBar.SelectedFile.Background)).
		Padding(0, 1).
		Render(currentFile.Name())

	fileTotals := lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.Colors.StatusBar.TotalFiles.Foreground)).
		Background(lipgloss.Color(cfg.Colors.StatusBar.TotalFiles.Background)).
		Align(lipgloss.Right).
		Padding(0, 1).
		Render(fmt.Sprintf("%d/%d", cursor+1, totalFiles))

	logoStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(lipgloss.Color(cfg.Colors.StatusBar.Logo.Foreground)).
		Background(lipgloss.Color(cfg.Colors.StatusBar.Logo.Background))

	logo := ""
	if cfg.Settings.ShowIcons {
		logo = logoStyle.Render(fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM"))
	} else {
		logo = logoStyle.Render("FM")
	}

	status := lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.Colors.StatusBar.Bar.Foreground)).
		Background(lipgloss.Color(cfg.Colors.StatusBar.Bar.Background)).
		Padding(0, 1).
		Width(screenWidth - width(selectedFile) - width(fileTotals) - width(logo)).
		Render(fmt.Sprintf("%s %s %s", helpers.ConvertBytesToSizeString(currentFile.Size()), currentFile.Mode().String(), currentPath))

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFile,
		status,
		fileTotals,
		logo,
	)

	doc.WriteString(bar)

	return doc.String()
}
