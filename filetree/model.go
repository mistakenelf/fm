package filetree

import "github.com/mistakenelf/fm/filesystem"

type DirectoryItem struct {
	Name             string
	Details          string
	Path             string
	Extension        string
	IsDirectory      bool
	CurrentDirectory string
}

type Model struct {
	Cursor   int
	files    []DirectoryItem
	active   bool
	keyMap   KeyMap
	min      int
	max      int
	height   int
	width    int
	startDir string
}

func New(active bool, startDir string) Model {
	startingDirectory := filesystem.CurrentDirectory

	if startDir != "" {
		startingDirectory = startDir
	}

	return Model{
		Cursor:   0,
		active:   active,
		keyMap:   DefaultKeyMap(),
		min:      0,
		max:      0,
		startDir: startingDirectory,
	}
}
