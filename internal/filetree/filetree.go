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
	"github.com/knipferrc/fm/internal/commands"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/renderer"
	"github.com/knipferrc/fm/internal/statusbar"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/spf13/viper"
)

// Model is a struct to represent the properties of a filetree.
type Model struct {
	Viewport            viewport.Model
	AppConfig           config.Config
	Style               lipgloss.Style
	UnselectedItemColor lipgloss.AdaptiveColor
	SelectedItemColor   lipgloss.AdaptiveColor
	ActiveBorderColor   lipgloss.AdaptiveColor
	InactiveBorderColor lipgloss.AdaptiveColor
	Files               []fs.DirEntry
	FilePaths           []string
	Cursor              int
	ShowIcons           bool
	ShowHidden          bool
	Borderless          bool
	IsActive            bool
	AlternateBorder     bool
	ShowLoading         bool
}

// Init intializes the filetree.
func (m Model) Init() tea.Cmd {
	startDir := viper.GetString("start-dir")

	switch {
	case startDir != "":
		_, err := os.Stat(startDir)
		if err != nil {
			return nil
		}

		if strings.HasPrefix(startDir, "/") {
			return commands.UpdateDirectoryListingCmd(startDir, m.ShowHidden)
		} else {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			filePath := filepath.Join(path, startDir)

			return commands.UpdateDirectoryListingCmd(filePath, m.ShowHidden)
		}
	case m.AppConfig.Settings.StartDir == dirfs.HomeDirectory:
		homeDir, err := dirfs.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		return commands.UpdateDirectoryListingCmd(homeDir, m.ShowHidden)
	default:
		return commands.UpdateDirectoryListingCmd(m.AppConfig.Settings.StartDir, m.ShowHidden)
	}
}

// NewModel creates a new instance of a filetree.
func NewModel(
	showIcons, borderless, isActive, showHidden bool,
	selectedItemColor, unselectedItemColor, activeBorderColor, inactiveBorderColor lipgloss.AdaptiveColor,
	appConfig config.Config,
) Model {
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
		Borderless:          borderless,
		IsActive:            isActive,
		ShowHidden:          showHidden,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		AppConfig:           appConfig,
		Style:               style,
	}
}

// scrollFiletree moves handles wrapping of the filetree and
// scrolling of the viewport.
func (m *Model) scrollFiletree() {
	top := m.Viewport.YOffset
	bottom := m.Viewport.Height + m.Viewport.YOffset - 1

	if m.Cursor < top {
		m.Viewport.LineUp(1)
	} else if m.Cursor > bottom {
		m.Viewport.LineDown(1)
	}

	if m.Cursor > m.GetTotalFiles()-1 {
		m.GotoTop()
		m.Viewport.GotoTop()
	} else if m.Cursor < top {
		m.GotoBottom()
		m.Viewport.GotoBottom()
	}
}

// handleRightKeyPress opens directory if it is one or reads a files content.
func (m *Model) handleRightKeyPress(cmds *[]tea.Cmd) {
	if m.IsActive && m.GetTotalFiles() > 0 {
		selectedFile, err := m.GetSelectedFile()
		if err != nil {
			*cmds = append(*cmds, commands.HandleErrorCmd(err))
		}

		if selectedFile.IsDir() {
			currentDir, err := dirfs.GetWorkingDirectory()
			if err != nil {
				*cmds = append(*cmds, commands.HandleErrorCmd(err))
			}

			directoryToOpen := filepath.Join(currentDir, selectedFile.Name())

			if len(m.GetFilePaths()) > 0 {
				directoryToOpen = m.GetFilePaths()[m.GetCursor()]
			}

			*cmds = append(*cmds, commands.UpdateDirectoryListingCmd(directoryToOpen, m.ShowHidden))
		}
	}
}

// handleLeftKeyPress goes to the previous directory.
func (m *Model) handleLeftKeyPress(cmds *[]tea.Cmd) {
	if m.IsActive && m.GetTotalFiles() > 0 {
		workingDirectory, err := dirfs.GetWorkingDirectory()
		if err != nil {
			*cmds = append(*cmds, commands.HandleErrorCmd(err))
		}

		*cmds = append(*cmds, commands.UpdateDirectoryListingCmd(filepath.Join(workingDirectory, dirfs.PreviousDirectory), m.ShowHidden))
	}
}

// SetContent sets the files currently displayed in the tree.
func (m *Model) SetContent(files []fs.DirEntry) {
	var directoryItem string
	curFiles := ""

	m.Files = files

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
	m.Viewport.Width = (width / 2) - m.Style.GetHorizontalBorderSize()
	m.Viewport.Height = height - m.Style.GetVerticalBorderSize() - statusbar.StatusbarHeight
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

// GetIsActive returns the active state of the filetree.
func (m Model) GetIsActive() bool {
	return m.IsActive
}

// SetIsActive sets the active state of the filetree.
func (m *Model) SetIsActive(isActive bool) {
	m.IsActive = isActive
}

// Update updates the statusbar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.UpdateDirectoryListingMsg:
		m.GotoTop()
		m.SetFilePaths(nil)
		m.Viewport.GotoTop()
		m.SetContent(msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.SetContent(m.Files)
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.IsActive {
				m.GoUp()
				m.scrollFiletree()
				m.SetContent(m.Files)
			}
		case "down", "j":
			if m.IsActive {
				m.GoDown()
				m.scrollFiletree()
				m.SetContent(m.Files)
			}
		case "right", "l":
			if m.IsActive {
				m.handleRightKeyPress(&cmds)
			}
		case "left", "h":
			if m.IsActive {
				m.handleLeftKeyPress(&cmds)
			}
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
