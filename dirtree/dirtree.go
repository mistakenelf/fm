package dirtree

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/icons"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Files               []fs.FileInfo
	Cursor              int
	ShowIcons           bool
	SelectedItemColor   string
	UnselectedItemColor string
}

func NewModel(files []fs.FileInfo, showIcons bool, selectedItemColor, unselectedItemColor string) Model {
	return Model{
		Files:               files,
		ShowIcons:           showIcons,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
	}
}

func (m *Model) SetContent(files []fs.FileInfo) {
	m.Files = files
}

func (m *Model) GotoTop() {
	m.Cursor = 0
}

func (m *Model) GotoBottom() {
	m.Cursor = len(m.Files) - 1
}

func (m Model) GetSelectedFile() fs.FileInfo {
	return m.Files[m.Cursor]
}

func (m Model) GetCursor() int {
	return m.Cursor
}

func (m *Model) GoDown() {
	m.Cursor++
}

func (m *Model) GoUp() {
	m.Cursor--
}

func (m Model) dirItem(selected bool, file fs.FileInfo) string {
	if !m.ShowIcons && !selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name())
	} else if !m.ShowIcons && selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name())
	} else if selected && file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name()))

		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(listing)
	} else if !selected && file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name()))

		return listing
	} else if selected && !file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name()))

		return listing
	} else {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name()))

		return listing
	}
}

func (m Model) View() string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range m.Files {
		curFiles += fmt.Sprintf("%s\n", m.dirItem(m.Cursor == i, file))
	}

	doc.WriteString(curFiles)

	return doc.String()
}
