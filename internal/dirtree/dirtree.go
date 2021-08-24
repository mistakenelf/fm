package dirtree

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/constants"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Model is a struct to represent the properties on a dirtree.
type Model struct {
	Files               []fs.FileInfo
	Width               int
	Cursor              int
	ShowIcons           bool
	ShowHidden          bool
	SelectedItemColor   string
	UnselectedItemColor string
}

// NewModel creates a new instance of a dirtree.
func NewModel(showIcons bool, selectedItemColor, unselectedItemColor string) Model {
	return Model{
		Cursor:              0,
		ShowIcons:           showIcons,
		ShowHidden:          true,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
	}
}

// SetContent update the files currently displayed in the tree.
func (m *Model) SetContent(files []fs.FileInfo) {
	m.Files = files
}

// SetSize updates the size of the dirtree, useful when resizing the terminal.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// GotoTop goes to the top of the tree.
func (m *Model) GotoTop() {
	m.Cursor = 0
}

// GotoBottom goes to the bottom of the tree.
func (m *Model) GotoBottom() {
	m.Cursor = len(m.Files) - 1
}

// GetSelectedFile returns the currently selected file in the tree.
func (m Model) GetSelectedFile() fs.FileInfo {
	return m.Files[m.Cursor]
}

// GetCursor gets the position of the cursor in the tree.
func (m Model) GetCursor() int {
	return m.Cursor
}

// GoDown goes down the tree by one.
func (m *Model) GoDown() {
	m.Cursor++
}

// GoUp goes up the tree by one.
func (m *Model) GoUp() {
	m.Cursor--
}

// GetTotalFiles returns the total number of files in the tree.
func (m Model) GetTotalFiles() int {
	return len(m.Files)
}

// ToggleHidden toggles the visibility of hidden files.
func (m *Model) ToggleHidden() {
	m.ShowHidden = !m.ShowHidden
}

// dirItem returns a string representation of a directory item.
func (m Model) dirItem(selected bool, file fs.FileInfo) string {
	// Get the icon and color based on the current file.
	icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
	fileIcon := fmt.Sprintf("%s%s", color, icon)

	if m.ShowIcons && selected {
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name()))
	} else if m.ShowIcons && !selected {
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name()))
	} else if !m.ShowIcons && selected {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name())
	} else {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name())
	}
}

// View returns a string representation of the current tree.
func (m Model) View() string {
	curFiles := ""

	for i, file := range m.Files {
		modTimeColor := ""

		if m.Cursor == i {
			modTimeColor = m.SelectedItemColor
		} else {
			modTimeColor = m.UnselectedItemColor
		}

		modTime := lipgloss.NewStyle().
			Align(lipgloss.Right).
			Foreground(lipgloss.Color(modTimeColor)).
			Render(file.ModTime().
				Format("2006-01-02 15:04:05"),
			)

		dirItem := lipgloss.NewStyle().Width(m.Width - lipgloss.Width(modTime) - 2).Render(
			truncate.StringWithTail(
				m.dirItem(m.Cursor == i, file), uint(m.Width-lipgloss.Width(modTime)-constants.Dimensions.PanePadding), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, modTime)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	return curFiles
}
