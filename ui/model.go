package ui

import (
	"fmt"
	"io/fs"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/help"
	"github.com/knipferrc/fm/pane"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	PrimaryPane    pane.Model
	SecondaryPane  pane.Model
	Files          []fs.FileInfo
	Textinput      textinput.Model
	Spinner        spinner.Model
	Help           help.Model
	Cursor         int
	ScreenWidth    int
	ScreenHeight   int
	ShowCommandBar bool
	Ready          bool
}

func NewModel() Model {
	cfg := config.GetConfig()

	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Width = 50

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Colors.Components.Spinner))

	h := help.NewModel()
	h.HeaderText = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(constants.White)).
		MarginBottom(1).
		Render("FM (File Manager)")

	h.BodyText = lipgloss.NewStyle().
		Foreground(lipgloss.Color(constants.White)).
		Render(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
			"h or left arrow - go back a directory",
			"j or down arrow - move cursor down",
			"k or up arrow - move cursor up",
			"l or right arrow - open selected folder / view file",
			": - open command bar",
			"mkdir /new/dir - create directory in current directory",
			"touch filename.txt - create file in current directory",
			"mv newname.txt - rename currently selected file or directory",
			"cp /dir/to/move/to - move file or directory",
			"rm - remove file or directory",
			"tab - toggle between panes"),
		)

	return Model{
		PrimaryPane:    pane.Model{},
		SecondaryPane:  pane.Model{},
		Files:          make([]fs.FileInfo, 0),
		Textinput:      input,
		Spinner:        s,
		Help:           h,
		Cursor:         0,
		ScreenWidth:    0,
		ScreenHeight:   0,
		ShowCommandBar: false,
		Ready:          false,
	}
}
