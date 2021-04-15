package main

import (
	"io/fs"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

type model struct {
	Files         []fs.FileInfo
	Viewport      viewport.Model
	TextInput     textinput.Model
	Cursor        int
	ScreenWidth   int
	ScreenHeight  int
	ViewportReady bool
	Move          bool
	Rename        bool
	Delete        bool
	ShowHelp      bool
}

func createInitialModel() model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Width = 50

	return model{
		Files:         make([]fs.FileInfo, 0),
		Viewport:      viewport.Model{},
		TextInput:     input,
		Cursor:        0,
		ScreenWidth:   0,
		ScreenHeight:  0,
		ViewportReady: false,
		Move:          false,
		Rename:        false,
		Delete:        false,
		ShowHelp:      false,
	}
}
