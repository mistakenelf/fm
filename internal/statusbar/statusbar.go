package statusbar

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/knipferrc/fm/directory"
	"github.com/knipferrc/fm/formatter"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/constants"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Color is a struct that contains the foreground and background colors of the statusbar.
type Color struct {
	Background string
	Foreground string
}

// Model is a struct that contains all the properties of the statusbar.
type Model struct {
	Width              int
	Height             int
	TotalFiles         int
	Cursor             int
	TextInput          string
	ShowIcons          bool
	ShowCommandBar     bool
	InMoveMode         bool
	SelectedFile       fs.DirEntry
	ItemToMove         fs.DirEntry
	FirstColumnColors  Color
	SecondColumnColors Color
	ThirdColumnColors  Color
	FourthColumnColors Color
}

// NewModel creates an instance of a statusbar.
func NewModel(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color) Model {
	return Model{
		Height:             1,
		TotalFiles:         0,
		Cursor:             0,
		TextInput:          "",
		ShowIcons:          true,
		ShowCommandBar:     false,
		InMoveMode:         false,
		SelectedFile:       nil,
		ItemToMove:         nil,
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
	}
}

// ParseCommand parses the command and returns the command name and the arguments.
func ParseCommand(command string) (string, string) {
	// Split the command string into an array.
	cmdString := strings.Split(command, " ")

	// If theres only one item in the array, its a singular
	// command such as rm.
	if len(cmdString) == 1 {
		cmdName := cmdString[0]

		return cmdName, ""
	}

	// This command has two values, first one is the name
	// of the command, other is the value to pass back
	// to the UI to update.
	if len(cmdString) == 2 {
		cmdName := cmdString[0]
		cmdValue := cmdString[1]

		return cmdName, cmdValue
	}

	return "", ""
}

func (m Model) getStatusbarContent() (string, string, string, string) {
	currentPath, err := directory.GetWorkingDirectory()
	if err != nil {
		currentPath = constants.Directories.CurrentDirectory
	}

	if m.TotalFiles == 0 {
		return "", "", "", ""
	}

	logo := ""

	// If icons are enabled, show the directory icon next to the logo text
	// else just show the text of the logo.
	if m.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	// Display some information about the currently seleted file including
	// its size, the mode and the current path.
	fileInfo, err := m.SelectedFile.Info()
	if err != nil {
		return "", "", "", ""
	}

	status := fmt.Sprintf("%s %s %s",
		formatter.ConvertBytesToSizeString(fileInfo.Size()),
		fileInfo.Mode().String(),
		currentPath,
	)

	// If the command bar is shown, show the text input.
	if m.ShowCommandBar {
		status = m.TextInput
	}

	// If in move mode, update the status text to indicate move mode is enabled
	// and the name of the file or directory being moved.
	if m.InMoveMode {
		status = fmt.Sprintf("Currently moving %s", m.ItemToMove.Name())
	}

	return m.SelectedFile.Name(),
		status,
		fmt.Sprintf("%d/%d", m.Cursor+1, m.TotalFiles),
		logo
}

// GetHeight returns the height of the statusbar.
func (m Model) GetHeight() int {
	return m.Height
}

// SetContent sets the content of the statusbar.
func (m *Model) SetContent(totalFiles, cursor int, textInput string, showIcons, showCommandBar, inMoveMode bool, selectedFile, itemToMove fs.DirEntry) {
	m.TotalFiles = totalFiles
	m.Cursor = cursor
	m.TextInput = textInput
	m.ShowIcons = showIcons
	m.ShowCommandBar = showCommandBar
	m.InMoveMode = inMoveMode
	m.SelectedFile = selectedFile
	m.ItemToMove = itemToMove
}

// SetSize sets the size of the statusbar, useful when the terminal is resized.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// View returns a string representation of the statusbar.
func (m Model) View() string {
	width := lipgloss.Width

	firstColumnContent, secondColumnContent, thirdColumnContent, fourthColumnContent := m.getStatusbarContent()

	firstColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FirstColumnColors.Foreground)).
		Background(lipgloss.Color(m.FirstColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Render(truncate.StringWithTail(firstColumnContent, 30, "..."))

	thirdColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.ThirdColumnColors.Foreground)).
		Background(lipgloss.Color(m.ThirdColumnColors.Background)).
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(m.Height).
		Render(thirdColumnContent)

	fourthColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FourthColumnColors.Foreground)).
		Background(lipgloss.Color(m.FourthColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Render(fourthColumnContent)

	// Second column of the status bar displayed in the center with configurable
	// foreground and background colors and some padding. Also calculate the
	// width of the other three columns so that this one can take up the rest of the space.
	secondColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.SecondColumnColors.Foreground)).
		Background(lipgloss.Color(m.SecondColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Width(m.Width - width(firstColumn) - width(thirdColumn) - width(fourthColumn)).
		Render(truncate.StringWithTail(secondColumnContent, uint(m.Width-width(firstColumn)-width(thirdColumn)-width(fourthColumn)-3), "..."))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		firstColumn,
		secondColumn,
		thirdColumn,
		fourthColumn,
	)
}
