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

type Model struct {
	Files  []fs.FileInfo
	Cursor int
}

func NewModel(files []fs.FileInfo, cursor int) Model {
	return Model{
		Files:  files,
		Cursor: cursor,
	}
}

func (m Model) View() string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range m.Files {
		curFiles += fmt.Sprintf("%s\n",
			dirItem(
				m.Cursor == i,
				file,
			))
	}

	doc.WriteString(curFiles)

	return doc.String()
}

func (m *Model) SetContent(files []fs.FileInfo, cursor int) {
	m.Files = files
	m.Cursor = cursor
}

func dirItem(selected bool, file fs.FileInfo) string {
	cfg := config.GetConfig()

	if !cfg.Settings.ShowIcons && !selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.UnselectedItem)).
			Render(file.Name())
	} else if !cfg.Settings.ShowIcons && selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(file.Name())
	} else if selected && file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(file.Name()))

		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(listing)
	} else if !selected && file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.UnselectedItem)).
			Render(file.Name()))

		return listing
	} else if selected && !file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.SelectedItem)).
			Render(file.Name()))

		return listing
	} else {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(cfg.Colors.DirTree.UnselectedItem)).
			Render(file.Name()))

		return listing
	}
}
