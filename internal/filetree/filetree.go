package filetree

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/renderer"
	"github.com/spf13/viper"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

type updateDirectoryListingMsg []fs.DirEntry
type errorMsg string
type directoryItemSizeMsg struct {
	size  string
	index int
}

// Model is a struct to represent the properties of a filetree.
type Model struct {
	Viewport            viewport.Model
	AppConfig           config.Config
	Files               []fs.DirEntry
	Style               lipgloss.Style
	FilePaths           []string
	Cursor              int
	ShowIcons           bool
	ShowHidden          bool
	Borderless          bool
	SelectedItemColor   lipgloss.AdaptiveColor
	UnselectedItemColor lipgloss.AdaptiveColor
	IsActive            bool
	AlternateBorder     bool
	ShowLoading         bool
	ActiveBorderColor   lipgloss.AdaptiveColor
	InactiveBorderColor lipgloss.AdaptiveColor
}

func (m Model) Init() tea.Cmd {
	startDir := viper.GetString("start-dir")

	switch {
	case startDir != "":
		_, err := os.Stat(startDir)
		if err != nil {
			return nil
		}

		if strings.HasPrefix(startDir, "/") {
			return m.updateDirectoryListingCmd(startDir)
		} else {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			filePath := filepath.Join(path, startDir)

			return m.updateDirectoryListingCmd(filePath)
		}
	case m.AppConfig.Settings.StartDir == dirfs.HomeDirectory:
		homeDir, err := dirfs.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		return m.updateDirectoryListingCmd(homeDir)
	default:
		return m.updateDirectoryListingCmd(m.AppConfig.Settings.StartDir)
	}
}

// NewModel creates a new instance of a filetree.
func NewModel(showIcons, borderless bool, selectedItemColor, unselectedItemColor lipgloss.AdaptiveColor, appConfig config.Config) Model {
	border := lipgloss.NormalBorder()
	padding := 1

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	style := lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border)

	return Model{
		Cursor:              0,
		ShowIcons:           showIcons,
		ShowHidden:          true,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
		AppConfig:           appConfig,
		Style:               style,
	}
}

// updateDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (m Model) updateDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := dirfs.GetDirectoryListing(name, m.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		err = os.Chdir(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(files)
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

func (m *Model) scrollPrimaryPane() {
	top := m.Viewport.YOffset
	bottom := m.Viewport.Height + m.Viewport.YOffset - 1

	// If the cursor is above the top of the viewport scroll up on the viewport
	// else were at the bottom and need to scroll the viewport down.
	if m.Cursor < top {
		m.Viewport.LineUp(1)
	} else if m.Cursor > bottom {
		m.Viewport.LineDown(1)
	}

	// If the cursor of the dirtree is at the bottom of the files
	// set the cursor to 0 to go to the top of the dirtree and
	// scroll the pane to the top else, were at the top of the dirtree and pane so
	// scroll the pane to the bottom and set the cursor to the bottom.
	if m.Cursor > m.GetTotalFiles()-1 {
		m.GotoTop()
		m.Viewport.GotoTop()
	} else if m.Cursor < top {
		m.GotoBottom()
		m.Viewport.GotoBottom()
	}
}

// SetContent sets the files currently displayed in the tree.
func (m *Model) SetContent(files []fs.DirEntry) {
	m.Files = files
	var directoryItem string
	curFiles := ""

	for i, file := range files {
		var fileSizeColor lipgloss.AdaptiveColor

		if m.Cursor == i {
			fileSizeColor = m.SelectedItemColor
		} else {
			fileSizeColor = m.UnselectedItemColor
		}

		fileInfo, _ := file.Info()

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

		dirItem := lipgloss.NewStyle().Width(m.Viewport.Width - lipgloss.Width(fileSize) - 2).Render(
			truncate.StringWithTail(
				directoryItem, uint(m.Viewport.Width-lipgloss.Width(fileSize)), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	m.Viewport.SetContent(curFiles)
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
func (m *Model) SetSize(width, height int) {
	m.Viewport.Width = width - m.Style.GetHorizontalBorderSize()
	m.Viewport.Height = height - m.Style.GetVerticalBorderSize() - 1
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

// Update updates the statusbar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case updateDirectoryListingMsg:
		m.SetContent(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.SetContent(m.Files)
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.GoUp()
			m.scrollPrimaryPane()
			m.SetContent(m.Files)
		case "down", "j":
			m.GoDown()
			m.scrollPrimaryPane()
			m.SetContent(m.Files)
		}
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the current tree.
func (m Model) View() string {
	if len(m.Files) == 0 {
		return "Directory is empty"
	}

	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()
	padding := 1
	content := m.Viewport.View()
	alternateBorder := lipgloss.Border{
		Top:         "-",
		Bottom:      "-",
		Left:        "|",
		Right:       "|",
		TopLeft:     "*",
		TopRight:    "*",
		BottomLeft:  "*",
		BottomRight: "*",
	}

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	if m.AlternateBorder {
		border = alternateBorder
	}

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return m.Style.Copy().
		BorderForeground(borderColor).
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		Width(m.Viewport.Width).
		Height(m.Viewport.Height).
		Render(content)
}
