package main

import (
	"io/fs"

	"github.com/charmbracelet/bubbles/viewport"
)

type model struct {
	Files                []fs.FileInfo
	Cursor               int
	Viewport             viewport.Model
	ViewportReady        bool
	CurrentlyHighlighted string
	ScreenWidth          int
	Move                 bool
}

func createInitialModel() model {
	return model{
		make([]fs.FileInfo, 0),
		0,
		viewport.Model{},
		false,
		"",
		0,
		false,
	}
}
