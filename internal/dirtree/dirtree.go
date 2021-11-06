package dirtree

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/icons"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Model is a struct to represent the properties of a dirtree.
type Model struct {
	Files               []fs.DirEntry
	FilePaths           []string
	Width               int
	Cursor              int
	ShowIcons           bool
	ShowHidden          bool
	SelectedItemColor   lipgloss.AdaptiveColor
	UnselectedItemColor lipgloss.AdaptiveColor
}

// NewModel creates a new instance of a dirtree.
func NewModel(showIcons bool, selectedItemColor, unselectedItemColor lipgloss.AdaptiveColor) Model {
	return Model{
		Cursor:              0,
		ShowIcons:           showIcons,
		ShowHidden:          true,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
	}
}

// SetContent sets the files currently displayed in the tree.
func (m *Model) SetContent(files []fs.DirEntry) {
	m.Files = files
}

// SetFilePaths sets an array of file paths.
func (m *Model) SetFilePaths(filePaths []string) {
	m.FilePaths = filePaths
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
func (m Model) GetSelectedFile() (os.FileInfo, error) {
	if len(m.Files) > 0 {
		fileInfo, err := m.Files[m.Cursor].Info()
		if err != nil {
			return nil, err
		}

		return fileInfo, nil
	}

	return nil, nil
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
	icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
	fileIcon := fmt.Sprintf("%s%s", color, icon)

	switch {
	case m.ShowIcons && selected:
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Bold(true).
			Foreground(m.SelectedItemColor).
			Render(fileInfo.Name()))
	case m.ShowIcons && !selected:
		return fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Bold(true).
			Foreground(m.UnselectedItemColor).
			Render(fileInfo.Name()))
	case !m.ShowIcons && selected:
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(m.SelectedItemColor).
			Render(fileInfo.Name())
	default:
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(m.UnselectedItemColor).
			Render(fileInfo.Name())
	}
}

// View returns a string representation of the current tree.
func (m Model) View() string {
	curFiles := ""

	if len(m.Files) == 0 {
		return "Directory is empty"
	}

	for i, file := range m.Files {
		var modTimeColor lipgloss.AdaptiveColor

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
			Foreground(modTimeColor).
			Render(fileInfo.ModTime().
				Format("2006-01-02 15:04:05"),
			)

		dirItem := lipgloss.NewStyle().Width(m.Width - lipgloss.Width(modTime) - 2).Render(
			truncate.StringWithTail(
				m.dirItem(m.Cursor == i, fileInfo), uint(m.Width-lipgloss.Width(modTime)), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, modTime)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	return curFiles
}
