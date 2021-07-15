package dirtree

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/pkg/icons"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type Model struct {
	Files               []fs.FileInfo
	Width               int
	Cursor              int
	ShowIcons           bool
	ShowHidden          bool
	SelectedItemColor   string
	UnselectedItemColor string
}

func NewModel(showIcons bool, selectedItemColor, unselectedItemColor string) Model {
	return Model{
		Cursor:              0,
		ShowIcons:           showIcons,
		ShowHidden:          true,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
	}
}

// Update the set of files the tree is currently displaying
func (m *Model) SetContent(files []fs.FileInfo) {
	m.Files = files
}

func (m *Model) SetSize(width int) {
	m.Width = width
}

// Go to the top of the tree
func (m *Model) GotoTop() {
	m.Cursor = 0
}

// Go to the bottom of the tree
func (m *Model) GotoBottom() {
	m.Cursor = len(m.Files) - 1
}

// Get the currently selected file
func (m Model) GetSelectedFile() fs.FileInfo {
	return m.Files[m.Cursor]
}

// Get the current position of the cursor in the tree
func (m Model) GetCursor() int {
	return m.Cursor
}

// Move down the tree by 1
func (m *Model) GoDown() {
	m.Cursor++
}

// Move up the tree by one
func (m *Model) GoUp() {
	m.Cursor--
}

// Get the total number of files currently being displayed in the tree
func (m Model) GetTotalFiles() int {
	return len(m.Files)
}

// Toggle whether or not to show hidden files and folders
func (m *Model) ToggleHidden() {
	m.ShowHidden = !m.ShowHidden
}

// Individual tree items
func (m Model) dirItem(selected bool, file fs.FileInfo) string {
	// Get the icon and color based on the current file
	icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
	fileIcon := fmt.Sprintf("%s%s", color, icon)

	if m.ShowIcons && selected {
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name()))
	} else if m.ShowIcons && !selected {
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name()))

	} else if !m.ShowIcons && selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name())
	} else {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name())
	}
}

// Display the directory tree
func (m Model) View() string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range m.Files {
		curFiles += fmt.Sprintf("%s\n", truncate.StringWithTail(m.dirItem(m.Cursor == i, file), uint(m.Width-8), "..."))
	}

	doc.WriteString(curFiles)

	return doc.String()
}
