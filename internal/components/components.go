package components

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/utils"

	"github.com/charmbracelet/lipgloss"
)

func DirItem(selected, isDir bool, name, ext, indicator string) string {
	cfg := config.GetConfig()

	if !cfg.Settings.ShowIcons && !selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.UnselectedDirItem)).
			Render(name)
	} else if !cfg.Settings.ShowIcons && selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(name)
	} else if selected && isDir {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(name))

		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(listing)
	} else if !selected && isDir {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.UnselectedDirItem)).
			Render(name))

		return listing
	} else if selected && !isDir {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(name))

		return listing
	} else {
		icon, color := icons.GetIcon(name, ext, indicator)
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.UnselectedDirItem)).
			Render(name))

		return listing
	}
}

func DirTree(files []fs.FileInfo, cursor, width int) string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range files {
		curFiles += fmt.Sprintf("%s\n",
			DirItem(
				cursor == i, file.IsDir(),
				file.Name(),
				filepath.Ext(file.Name()),
				icons.GetIndicator(file.Mode()),
			))
	}

	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Left).
			Render(curFiles),
	))

	return doc.String()
}

func Pane(width, height int, isActive bool, content string) string {
	cfg := config.GetConfig()
	borderColor := cfg.Colors.Pane.InactivePane
	borderType := lipgloss.NormalBorder()

	if cfg.Settings.RoundedPanes {
		borderType = lipgloss.RoundedBorder()
	} else {
		borderType = lipgloss.NormalBorder()
	}

	if isActive {
		borderColor = cfg.Colors.Pane.ActivePane
	}

	return lipgloss.NewStyle().
		BorderForeground(lipgloss.Color(borderColor)).
		Border(borderType).
		Width(width).
		Height(height).
		Render(content)
}

func StatusBar(screenWidth, cursor, totalFiles int, currentFile fs.FileInfo, showCommandBar bool, textInput string) string {
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
		Render(fmt.Sprintf("%s %s %s",
			utils.ConvertBytesToSizeString(currentFile.Size()),
			currentFile.Mode().String(),
			currentPath),
		)

	if showCommandBar {
		status = lipgloss.NewStyle().
			Padding(0, 1).
			Width(screenWidth - width(selectedFile) - width(fileTotals) - width(logo)).
			Render(textInput)
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
