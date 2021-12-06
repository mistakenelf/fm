package statusbar

import (
	"fmt"
	"os"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/commands"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// StatusbarHeight is the height of the statusbar.
const StatusbarHeight = 1

// Color is a struct that contains the foreground and background colors of the statusbar.
type Color struct {
	Background lipgloss.AdaptiveColor
	Foreground lipgloss.AdaptiveColor
}

// Model is a struct that contains all the properties of the statusbar.
type Model struct {
	Width              int
	Height             int
	TotalFiles         int
	Cursor             int
	ShowIcons          bool
	ShowCommandInput   bool
	InMoveMode         bool
	SimpleMode         bool
	FilePaths          []string
	SelectedFile       os.FileInfo
	ItemToMove         os.FileInfo
	FirstColumnColors  Color
	SecondColumnColors Color
	ThirdColumnColors  Color
	FourthColumnColors Color
	Textinput          textinput.Model
}

// NewModel creates an instance of a statusbar.
func NewModel(
	firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color, showIcons, simpleMode bool,
) Model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Placeholder = "enter a name"

	if !simpleMode {
		input.PlaceholderStyle.Background(secondColumnColors.Background)
	}

	return Model{
		Height:             StatusbarHeight,
		TotalFiles:         0,
		Cursor:             0,
		ShowIcons:          showIcons,
		ShowCommandInput:   false,
		InMoveMode:         false,
		SimpleMode:         simpleMode,
		SelectedFile:       nil,
		ItemToMove:         nil,
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
		Textinput:          input,
	}
}

// GetHeight returns the height of the statusbar.
func (m Model) GetHeight() int {
	return m.Height
}

// BlurCommandInput blurs the textinput used for the command input.
func (m *Model) BlurCommandInput() {
	m.Textinput.Blur()
}

// ResetCommandInput resets the textinput used for the command input.
func (m *Model) ResetCommandInput() {
	m.Textinput.Reset()
}

// CommandInputValue returns the value of the command input.
func (m Model) CommandInputValue() string {
	return m.Textinput.Value()
}

// CommandInputFocused returns true if the command input is focused.
func (m Model) CommandInputFocused() bool {
	return m.Textinput.Focused()
}

// FocusCommandInput focuses the command input.
func (m *Model) FocusCommandInput() {
	m.Textinput.Focus()
}

// SetCommandInputPlaceholderText sets the placeholder text of the command input.
func (m *Model) SetCommandInputPlaceholderText(text string) {
	m.Textinput.Placeholder = text
}

// SetContent sets the content of the statusbar.
func (m *Model) SetContent(
	totalFiles, cursor int,
	showCommandInput, inMoveMode bool,
	selectedFile, itemToMove os.FileInfo, filePaths []string,
) {
	m.TotalFiles = totalFiles
	m.Cursor = cursor
	m.ShowCommandInput = showCommandInput
	m.InMoveMode = inMoveMode
	m.SelectedFile = selectedFile
	m.ItemToMove = itemToMove
	m.FilePaths = filePaths
}

// SetSize sets the size of the statusbar, useful when the terminal is resized.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// Update updates the statusbar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.UpdateStatusbarMsg:
		m.SetContent(msg.TotalFiles, msg.Cursor, msg.ShowCommandInput, msg.InMoveMode, msg.SelectedFile, msg.ItemToMove, msg.FilePaths)
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width)
	}

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the statusbar.
func (m Model) View() string {
	var logo string
	var status string

	width := lipgloss.Width
	selectedFile := "N/A"
	fileCount := "0/0"

	if m.TotalFiles > 0 && m.SelectedFile != nil {
		selectedFile = m.SelectedFile.Name()
		fileCount = fmt.Sprintf("%d/%d", m.Cursor+1, m.TotalFiles)

		currentPath, err := dirfs.GetWorkingDirectory()
		if err != nil {
			currentPath = dirfs.CurrentDirectory
		}

		if len(m.FilePaths) > 0 {
			currentPath = m.FilePaths[m.Cursor]
		}

		// Display some information about the currently seleted file including
		// its size, the mode and the current path.
		status = fmt.Sprintf("%s %s %s",
			m.SelectedFile.ModTime().Format("2006-01-02 15:04:05"),
			m.SelectedFile.Mode().String(),
			currentPath,
		)
	}

	if m.ShowCommandInput {
		status = m.Textinput.View()
	}

	if m.InMoveMode {
		status = fmt.Sprintf("%s %s", "Currently moving:", m.ItemToMove.Name())
	}

	if m.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	// Selected file styles
	selectedFileStyle := lipgloss.NewStyle().
		Foreground(m.FirstColumnColors.Foreground).
		Background(m.FirstColumnColors.Background)

	if m.SimpleMode {
		selectedFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	selectedFileColumn := selectedFileStyle.
		Padding(0, 1).
		Height(m.Height).
		Render(truncate.StringWithTail(selectedFile, 30, "..."))

	// File count styles
	fileCountStyle := lipgloss.NewStyle().
		Foreground(m.ThirdColumnColors.Foreground).
		Background(m.ThirdColumnColors.Background)

	if m.SimpleMode {
		fileCountStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	fileCountColumn := fileCountStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(m.Height).
		Render(fileCount)

	// Logo styles
	logoStyle := lipgloss.NewStyle().
		Foreground(m.FourthColumnColors.Foreground).
		Background(m.FourthColumnColors.Background)

	if m.SimpleMode {
		logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	logoColumn := logoStyle.
		Padding(0, 1).
		Height(m.Height).
		Render(logo)

	// Status styles
	statusStyle := lipgloss.NewStyle().
		Foreground(m.SecondColumnColors.Foreground).
		Background(m.SecondColumnColors.Background)

	if m.SimpleMode {
		statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	statusColumn := statusStyle.
		Padding(0, 1).
		Height(m.Height).
		Width(m.Width - width(selectedFileColumn) - width(fileCountColumn) - width(logoColumn)).
		Render(truncate.StringWithTail(
			status,
			uint(m.Width-width(selectedFileColumn)-width(fileCountColumn)-width(logoColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileColumn,
		statusColumn,
		fileCountColumn,
		logoColumn,
	)
}
