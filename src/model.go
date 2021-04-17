package main

import (
	"io/fs"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

type model struct {
	files        []fs.FileInfo
	viewport     viewport.Model
	textinput    textinput.Model
	cursor       int
	screenwidth  int
	screenheight int
	move         bool
	rename       bool
	delete       bool
	showhelp     bool
	ready        bool
}

func createModel() model {
	input := textinput.NewModel()
	input.Prompt = "❯ "
	input.CharLimit = 250
	input.Width = 50

	return model{
		files:        make([]fs.FileInfo, 0),
		viewport:     viewport.Model{},
		textinput:    input,
		cursor:       0,
		screenwidth:  0,
		screenheight: 0,
		move:         false,
		rename:       false,
		delete:       false,
		showhelp:     false,
		ready:        false,
	}
}
