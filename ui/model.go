package ui

import (
	"io/fs"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/dirtree"
	"github.com/knipferrc/fm/pane"
	"github.com/knipferrc/fm/statusbar"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	primaryPane          pane.Model
	secondaryPane        pane.Model
	textInput            textinput.Model
	loader               spinner.Model
	dirTree              dirtree.Model
	statusBar            statusbar.Model
	previousKey          tea.KeyMsg
	previousDirectory    string
	showCommandBar       bool
	activeMarkdownSource string
	inMoveMode           bool
	initialMoveDirectory string
	itemToMove           fs.FileInfo
	ready                bool
}

func NewModel(files []fs.FileInfo) model {
	cfg := config.GetConfig()

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250

	l := spinner.NewModel()
	l.Spinner = spinner.Dot
	l.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Colors.Spinner))

	dirTree := dirtree.NewModel(
		files,
		cfg.Settings.ShowIcons,
		cfg.Colors.DirTree.SelectedItem,
		cfg.Colors.DirTree.UnselectedItem,
	)

	return model{
		primaryPane:          pane.Model{},
		secondaryPane:        pane.Model{},
		textInput:            input,
		loader:               l,
		dirTree:              dirTree,
		statusBar:            statusbar.Model{},
		previousKey:          tea.KeyMsg{},
		previousDirectory:    "",
		showCommandBar:       false,
		activeMarkdownSource: constants.HelpText,
		inMoveMode:           false,
		initialMoveDirectory: "",
		itemToMove:           nil,
		ready:                false,
	}
}
