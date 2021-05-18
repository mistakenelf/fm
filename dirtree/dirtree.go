package dirtree

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/icons"

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

func View(files []fs.FileInfo, cursor, width int) string {
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
