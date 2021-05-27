package ui

import (
	"fmt"
	"io/fs"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/dirtree"
	"github.com/knipferrc/fm/pane"
	"github.com/knipferrc/fm/statusbar"
	"github.com/knipferrc/fm/text"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	PrimaryPane       pane.Model
	SecondaryPane     pane.Model
	Files             []fs.FileInfo
	Textinput         textinput.Model
	Spinner           spinner.Model
	HelpText          text.Model
	DirTree           dirtree.Model
	StatusBar         statusbar.Model
	LastKey           tea.KeyMsg
	PreviousDirectory string
	ScreenWidth       int
	ScreenHeight      int
	ShowCommandBar    bool
	Ready             bool
}

func NewModel() Model {
	cfg := config.GetConfig()

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Colors.Components.Spinner))

	t := text.NewModel()
	t.HeaderText = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(constants.White)).
		MarginBottom(1).
		Render("FM (File Manager)")

	t.BodyText = lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.White)).
		Render(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
			"h or left arrow    | go back a directory",
			"j or down arrow    | move cursor down",
			"k or up arrow      | move cursor up",
			"l or right arrow   | open selected folder / view file",
			"gg                 | go to top of pane",
			"G                  | go to botom of pane",
			"~                  | switch to home directory",
			"-                  | Go To previous directory",
			":                  | open command bar",
			"mkdir /new/dir     | create directory in current directory",
			"touch filename.txt | create file in current directory",
			"mv newname.txt     | rename currently selected file or directory",
			"cp /dir/to/move/to | move file or directory",
			"rm                 | remove file or directory",
			"tab                | toggle between panes"),
		)

	return Model{
		PrimaryPane:       pane.Model{},
		SecondaryPane:     pane.Model{},
		Files:             make([]fs.FileInfo, 0),
		Textinput:         input,
		Spinner:           s,
		HelpText:          t,
		DirTree:           dirtree.Model{},
		StatusBar:         statusbar.Model{},
		LastKey:           tea.KeyMsg{},
		PreviousDirectory: "",
		ScreenWidth:       0,
		ScreenHeight:      0,
		ShowCommandBar:    false,
		Ready:             false,
	}
}
