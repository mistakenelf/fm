package statusbar

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/utils"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width          int
	Cursor         int
	TotalFiles     int
	CurrentFile    fs.FileInfo
	ShowCommandBar bool
	TextInput      string
}

func NewModel(width, cursor, totalFiles int, currentFile fs.FileInfo, showCommandBar bool, textInput string) Model {
	return Model{
		Width:          width,
		Cursor:         cursor,
		TotalFiles:     totalFiles,
		CurrentFile:    currentFile,
		ShowCommandBar: showCommandBar,
		TextInput:      textInput,
	}
}

func (m Model) View() string {
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
		Render(m.CurrentFile.Name())

	fileTotals := lipgloss.NewStyle().
		Foreground(lipgloss.Color(cfg.Colors.StatusBar.TotalFiles.Foreground)).
		Background(lipgloss.Color(cfg.Colors.StatusBar.TotalFiles.Background)).
		Align(lipgloss.Right).
		Padding(0, 1).
		Render(fmt.Sprintf("%d/%d", m.Cursor+1, m.TotalFiles))

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
		Width(m.Width - width(selectedFile) - width(fileTotals) - width(logo)).
		Render(fmt.Sprintf("%s %s %s",
			utils.ConvertBytesToSizeString(m.CurrentFile.Size()),
			m.CurrentFile.Mode().String(),
			currentPath),
		)

	if m.ShowCommandBar {
		status = lipgloss.NewStyle().
			Padding(0, 1).
			Width(m.Width - width(selectedFile) - width(fileTotals) - width(logo)).
			Render(m.TextInput)
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFile,
		status,
		fileTotals,
		logo,
	)

	doc.WriteString(bar)

	return doc.String()
}

func (m *Model) SetContent(width, cursor, totalFiles int, currentFile fs.FileInfo, showCommandBar bool, textInput string) {
	m.Width = width
	m.Cursor = cursor
	m.TotalFiles = totalFiles
	m.CurrentFile = currentFile
	m.ShowCommandBar = showCommandBar
	m.TextInput = textInput
}
