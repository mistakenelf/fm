package filetree

import (
	"time"

	"github.com/charmbracelet/lipgloss"
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
	Disabled              bool
	keyMap                KeyMap
	min                   int
	max                   int
	height                int
	width                 int
	startDir              string
	showHidden            bool
	showDirectoriesOnly   bool
	showFilesOnly         bool
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
	selectedItemColor     lipgloss.AdaptiveColor
	unselectedItemColor   lipgloss.AdaptiveColor
	inactiveItemColor     lipgloss.AdaptiveColor
	selectionPath         string
	showIcons             bool
}

func New(startDir string) Model {
	startingDirectory := filesystem.CurrentDirectory

	if startDir != "" {
		startingDirectory = startDir
	}

	return Model{
		Cursor:                0,
		Disabled:              false,
		keyMap:                DefaultKeyMap(),
		min:                   0,
		max:                   0,
		startDir:              startingDirectory,
		showHidden:            true,
		StatusMessage:         "",
		StatusMessageLifetime: time.Second,
		showFilesOnly:         false,
		showDirectoriesOnly:   false,
		selectedItemColor:     lipgloss.AdaptiveColor{Light: "212", Dark: "212"},
		unselectedItemColor:   lipgloss.AdaptiveColor{Light: "ffffff", Dark: "#000000"},
		inactiveItemColor:     lipgloss.AdaptiveColor{Light: "243", Dark: "243"},
		showIcons:             true,
	}
}
