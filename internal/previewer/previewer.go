package previewer

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/commands"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/strfmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Bubble is a struct to represent the properties of a filetree.
type Bubble struct {
	Viewport            viewport.Model
	AppConfig           config.Config
	Style               lipgloss.Style
	UnselectedItemColor lipgloss.AdaptiveColor
	SelectedItemColor   lipgloss.AdaptiveColor
	ActiveBorderColor   lipgloss.AdaptiveColor
	InactiveBorderColor lipgloss.AdaptiveColor
	Files               []fs.DirEntry
	ShowIcons           bool
	ShowHidden          bool
	Borderless          bool
	IsActive            bool
}

// NewBubble creates a new instance of a filetree.
func NewBubble(
	showIcons, borderless, isActive, showHidden bool,
	selectedItemColor, unselectedItemColor, activeBorderColor, inactiveBorderColor lipgloss.AdaptiveColor,
	appConfig config.Config,
) Bubble {
	border := lipgloss.NormalBorder()
	padding := 1

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	style := lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border)

	return Bubble{
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

// SetContent sets the files currently displayed in the tree.
func (m *Bubble) SetContent(files []fs.DirEntry) {
	var directoryItem string
	curFiles := ""

	m.Files = files

	for _, file := range files {
		var fileSizeColor lipgloss.AdaptiveColor

		fileInfo, _ := file.Info()

		fileSize := lipgloss.NewStyle().
			Foreground(fileSizeColor).
			Render(strfmt.ConvertBytesToSizeString(fileInfo.Size()))

		icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)

		switch {
		case m.ShowIcons:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(m.UnselectedItemColor).
				Render(fileInfo.Name()))
		case !m.ShowIcons:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(m.UnselectedItemColor).
				Render(fileInfo.Name())
		default:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(m.UnselectedItemColor).
				Render(fileInfo.Name())
		}

		dirItem := lipgloss.NewStyle().Width(m.Viewport.Width - lipgloss.Width(fileSize) - m.Style.GetHorizontalPadding()).Render(
			truncate.StringWithTail(
				directoryItem, uint(m.Viewport.Width-lipgloss.Width(fileSize)), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	m.Viewport.SetContent(curFiles)
}

// SetSize updates the size of the filetree, useful when resizing the terminal.
func (m *Bubble) SetSize(width, height int) {
	m.Viewport.Width = (width / 2) - m.Style.GetHorizontalBorderSize()
	m.Viewport.Height = height - m.Style.GetVerticalBorderSize() - statusbar.StatusbarHeight
}

// GetIsActive returns the active state of the filetree.
func (m Bubble) GetIsActive() bool {
	return m.IsActive
}

// SetIsActive sets the active state of the filetree.
func (m *Bubble) SetIsActive(isActive bool) {
	m.IsActive = isActive
}

// Update updates the statusbar.
func (m Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.PreviewDirectoryListingMsg:
		m.SetContent(msg)
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.SetContent(m.Files)
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if m.IsActive {
				m.Viewport.LineUp(1)
				m.SetContent(m.Files)
			}
		case tea.MouseWheelDown:
			if m.IsActive {
				m.Viewport.LineDown(1)
				m.SetContent(m.Files)
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.IsActive {
				m.Viewport.LineUp(1)
				m.SetContent(m.Files)
			}
		case "down", "j":
			if m.IsActive {
				m.Viewport.LineDown(1)
				m.SetContent(m.Files)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the current tree.
func (m Bubble) View() string {
	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()
	content := m.Viewport.View()

	if len(m.Files) == 0 {
		return m.Style.Copy().
			BorderForeground(borderColor).
			Border(border).
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render("Directory is empty")
	}

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return m.Style.Copy().
		BorderForeground(borderColor).
		Border(border).
		Width(m.Viewport.Width).
		Height(m.Viewport.Height).
		Render(content)
}
