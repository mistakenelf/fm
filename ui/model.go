package ui

import (
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

type Model struct {
	PrimaryPane       pane.Model
	SecondaryPane     pane.Model
	Textinput         textinput.Model
	Spinner           spinner.Model
	DirTree           dirtree.Model
	StatusBar         statusbar.Model
	PreviousKey       tea.KeyMsg
	PreviousDirectory string
	ScreenWidth       int
	ScreenHeight      int
	ShowCommandBar    bool
	Ready             bool

	activeMarkdownSource string
}

func NewModel() Model {
	cfg := config.GetConfig()

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Colors.Components.Spinner))

	return Model{
		PrimaryPane:       pane.Model{},
		SecondaryPane:     pane.Model{},
		Textinput:         input,
		Spinner:           s,
		DirTree:           dirtree.Model{},
		StatusBar:         statusbar.Model{},
		PreviousKey:       tea.KeyMsg{},
		PreviousDirectory: "",
		ScreenWidth:       0,
		ScreenHeight:      0,
		ShowCommandBar:    false,
		Ready:             false,

		activeMarkdownSource: constants.HelpText,
	}
}
