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
	ShowHidden          bool
	SelectedItemColor   string
	UnselectedItemColor string
}

// Create a new instance of a dirtree
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

// Go to the top of the tree
func (m *Model) GotoTop() {
	m.Cursor = 0
}

// Go to the bottom of the tree which is the length of all the files
// minus one
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

// dirItem is each individual item within the tree
func (m Model) dirItem(selected bool, file fs.FileInfo) string {
	// Get the icon and color based on the current file
	icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
	fileIcon := fmt.Sprintf("%s%s", color, icon)

	if m.ShowIcons && selected {
		// Reset the color of the text after getting the color of the icon
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name()))
	} else if m.ShowIcons && !selected {
		// Reset the color of the text after getting the color of the icon
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name()))

	} else if !m.ShowIcons && selected {
		// If icons are not enabled but the item is selected
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name())
	} else {
		// If icons are not enabled and the item is not currently selected
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name())
	}
}

// Display the dirtree and all of its dir items
func (m Model) View() string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range m.Files {
		curFiles += fmt.Sprintf("%s\n", m.dirItem(m.Cursor == i, file))
	}

	doc.WriteString(curFiles)

	return doc.String()
}
