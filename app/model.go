package app

import (
	"io/fs"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Files             []fs.FileInfo
	Viewport          viewport.Model
	SecondaryViewport viewport.Model
	Textinput         textinput.Model
	Spinner           spinner.Model
	Cursor            int
	ScreenWidth       int
	ScreenHeight      int
	Move              bool
	Rename            bool
	Delete            bool
	Ready             bool
}

func CreateModel() Model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Width = 50

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		Files:             make([]fs.FileInfo, 0),
		Viewport:          viewport.Model{},
		SecondaryViewport: viewport.Model{},
		Textinput:         input,
		Spinner:           s,
		Cursor:            0,
		ScreenWidth:       0,
		ScreenHeight:      0,
		Move:              false,
		Rename:            false,
		Delete:            false,
		Ready:             false,
	}
}
