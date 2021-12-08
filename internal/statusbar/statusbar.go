package statusbar

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

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
type Bubble struct {
	Width               int
	Files               []fs.DirEntry
	Cursor              int
	ShowIcons           bool
	ShowCommandInput    bool
	InMoveMode          bool
	SimpleMode          bool
	CreateFileMode      bool
	CreateDirectoryMode bool
	DeleteMode          bool
	RenameMode          bool
	FilePaths           []string
	SelectedFile        os.FileInfo
	ItemToMove          os.FileInfo
	FirstColumnColors   Color
	SecondColumnColors  Color
	ThirdColumnColors   Color
	FourthColumnColors  Color
	Textinput           textinput.Model
}

// NewBubble creates an instance of a statusbar.
func NewBubble(
	firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color, showIcons, simpleMode bool,
) Bubble {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Placeholder = "enter a name"

	if !simpleMode {
		input.PlaceholderStyle.Background(secondColumnColors.Background)
	}

	return Bubble{
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
		CreateFileMode:     true,
	}
}

// Init initializes the statusbar.
func (m Bubble) Init() tea.Cmd {
	return textinput.Blink
}

// SetCommandInputPlaceholderText sets the placeholder text of the command input.
func (m *Bubble) SetCommandInputPlaceholderText(text string) {
	m.Textinput.Placeholder = text
}

// SetContent sets the content of the statusbar.
func (m *Bubble) SetContent(
	files []fs.DirEntry,
	cursor int,
	selectedFile, itemToMove os.FileInfo, filePaths []string,
) {
	m.Files = files
	m.Cursor = cursor
	m.SelectedFile = selectedFile
	m.ItemToMove = itemToMove
	m.FilePaths = filePaths
}

// SetSize sets the size of the statusbar, useful when the terminal is resized.
func (m *Bubble) SetSize(width int) {
	m.Width = width
}

// Update updates the statusbar.
func (m Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.UpdateStatusbarMsg:
		m.SetContent(
			msg.Files,
			msg.Cursor,
			msg.SelectedFile,
			msg.ItemToMove,
			msg.FilePaths,
		)
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selectedFile, err := m.Files[m.Cursor].Info()
			if err != nil {
				cmds = append(cmds, commands.HandleErrorCmd(err))
			}

			switch {
			case m.CreateFileMode:
				cmds = append(cmds, tea.Sequentially(
					commands.CreateFileCmd(m.Textinput.Value()),
					commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, true),
				))
			case m.CreateDirectoryMode:
				cmds = append(cmds, tea.Sequentially(
					commands.CreateDirectoryCmd(m.Textinput.Value()),
					commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, true),
				))
			case m.RenameMode:
				cmds = append(cmds, tea.Sequentially(
					commands.RenameDirectoryItemCmd(selectedFile.Name(), m.Textinput.Value()),
					commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, true),
				))
			case m.DeleteMode:
				if strings.ToLower(m.Textinput.Value()) == "y" || strings.ToLower(m.Textinput.Value()) == "yes" {
					if selectedFile.IsDir() {
						cmds = append(cmds, tea.Sequentially(
							commands.DeleteDirectoryCmd(selectedFile.Name()),
							commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, true),
						))
					} else {
						cmds = append(cmds, tea.Sequentially(
							commands.DeleteFileCmd(selectedFile.Name()),
							commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, true),
						))
					}
				}

				m.Textinput.Blur()
				m.Textinput.Reset()
			default:
				return m, nil
			}
		}
	}

	m.Textinput, cmd = m.Textinput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the statusbar.
func (m Bubble) View() string {
	var logo string
	var status string

	width := lipgloss.Width
	selectedFile := "N/A"
	fileCount := "0/0"

	if len(m.Files) > 0 && m.SelectedFile != nil {
		selectedFile = m.SelectedFile.Name()
		fileCount = fmt.Sprintf("%d/%d", m.Cursor+1, len(m.Files))

		currentPath, err := dirfs.GetWorkingDirectory()
		if err != nil {
			currentPath = dirfs.CurrentDirectory
		}

		if len(m.FilePaths) > 0 {
			currentPath = m.FilePaths[m.Cursor]
		}

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
		Height(StatusbarHeight).
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
		Height(StatusbarHeight).
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
		Height(StatusbarHeight).
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
		Height(StatusbarHeight).
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
