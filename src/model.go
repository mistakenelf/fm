package main

import (
	"io/fs"

	"github.com/charmbracelet/bubbles/viewport"
)

type model struct {
	Files         []fs.FileInfo
	Selected      map[int]struct{}
	Cursor        int
	Quitting      bool
	FileContent   string
	Viewport      viewport.Model
	ViewportReady bool
}

func createInitialModel() model {
	return model{
		make([]fs.FileInfo, 0),
		make(map[int]struct{}),
		0,
		false,
		"",
		viewport.Model{},
		false,
	}
}
