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
	Files               []fs.DirEntry
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
func (m *Model) SetContent(files []fs.DirEntry) {
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
func (m Model) GetSelectedFile() fs.DirEntry {
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
func (m Model) dirItem(selected bool, fileInfo fs.FileInfo) string {
	// Get the icon and color based on the current file.
	icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
	fileIcon := fmt.Sprintf("%s%s", color, icon)

	if m.ShowIcons && selected {
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(fileInfo.Name()))
	} else if m.ShowIcons && !selected {
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(fileInfo.Name()))
	} else if !m.ShowIcons && selected {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(fileInfo.Name())
	} else {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(fileInfo.Name())
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

		fileInfo, err := file.Info()
		if err != nil {
			return err.Error()
		}

		modTime := lipgloss.NewStyle().
			Align(lipgloss.Right).
			Foreground(lipgloss.Color(modTimeColor)).
			Render(fileInfo.ModTime().
				Format("2006-01-02 15:04:05"),
			)

		dirItem := lipgloss.NewStyle().Width(m.Width - lipgloss.Width(modTime) - 2).Render(
			truncate.StringWithTail(
				m.dirItem(m.Cursor == i, fileInfo), uint(m.Width-lipgloss.Width(modTime)-constants.Dimensions.PanePadding), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, modTime)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	return curFiles
}
