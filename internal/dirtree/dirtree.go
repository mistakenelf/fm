package dirtree

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/renderer"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type directoryItemSizeMsg struct {
	size  string
	index int
}

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

// getDirectoryItemSizeCmd calculates the size of a directory or file.
func (m Model) getDirectoryItemSizeCmd(name string, i int) tea.Cmd {
	return func() tea.Msg {
		size, err := dirfs.GetDirectoryItemSize(name)
		if err != nil {
			return directoryItemSizeMsg{size: "N/A", index: i}
		}

		sizeString := renderer.ConvertBytesToSizeString(size)

		return directoryItemSizeMsg{
			size:  sizeString,
			index: i,
		}
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

// GetFilePaths returns an array of file paths.
func (m Model) GetFilePaths() []string {
	return m.FilePaths
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

// View returns a string representation of the current tree.
func (m Model) View() string {
	curFiles := ""
	var directoryItem string

	if len(m.Files) == 0 {
		return "Directory is empty"
	}

	for i, file := range m.Files {
		var fileSizeColor lipgloss.AdaptiveColor

		if m.Cursor == i {
			fileSizeColor = m.SelectedItemColor
		} else {
			fileSizeColor = m.UnselectedItemColor
		}

		fileInfo, err := file.Info()
		if err != nil {
			return err.Error()
		}

		fileSize := lipgloss.NewStyle().
			Foreground(fileSizeColor).
			Render(renderer.ConvertBytesToSizeString(fileInfo.Size()))

		icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)

		switch {
		case m.ShowIcons && m.Cursor == i:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(m.SelectedItemColor).
				Render(fileInfo.Name()))
		case m.ShowIcons && m.Cursor != i:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(m.UnselectedItemColor).
				Render(fileInfo.Name()))
		case !m.ShowIcons && m.Cursor == i:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(m.SelectedItemColor).
				Render(fileInfo.Name())
		default:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(m.UnselectedItemColor).
				Render(fileInfo.Name())
		}

		dirItem := lipgloss.NewStyle().Width(m.Width - lipgloss.Width(fileSize) - 2).Render(
			truncate.StringWithTail(
				directoryItem, uint(m.Width-lipgloss.Width(fileSize)), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	return curFiles
}
