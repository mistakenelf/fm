package filetree

import (
	"time"

	"github.com/mistakenelf/fm/filesystem"
)

type DirectoryItem struct {
	Name             string
	Details          string
	Path             string
	Extension        string
	IsDirectory      bool
	CurrentDirectory string
}

type Model struct {
	Cursor                int
	files                 []DirectoryItem
	Active                bool
	keyMap                KeyMap
	min                   int
	max                   int
	height                int
	width                 int
	startDir              string
	showHidden            bool
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
}

func New(active bool, startDir string) Model {
	startingDirectory := filesystem.CurrentDirectory

	if startDir != "" {
		startingDirectory = startDir
	}

	return Model{
		Cursor:                0,
		Active:                active,
		keyMap:                DefaultKeyMap(),
		min:                   0,
		max:                   0,
		startDir:              startingDirectory,
		showHidden:            true,
		StatusMessage:         "",
		StatusMessageLifetime: time.Second,
	}
}
