package filetree

import (
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/mistakenelf/fm/filesystem"
)

type DirectoryItem struct {
	Name        string
	Details     string
	Path        string
	Extension   string
	FileSize    string
	IsDirectory bool
	FileInfo    os.FileInfo
}

type Model struct {
	files                 []DirectoryItem
	Cursor                int
	min                   int
	max                   int
	height                int
	width                 int
	Disabled              bool
	showHidden            bool
	showDirectoriesOnly   bool
	showFilesOnly         bool
	showIcons             bool
	CreatingNewFile       bool
	CreatingNewDirectory  bool
	keyMap                KeyMap
	startDir              string
	StatusMessage         string
	selectionPath         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
	selectedItemColor     lipgloss.AdaptiveColor
	unselectedItemColor   lipgloss.AdaptiveColor
	inactiveItemColor     lipgloss.AdaptiveColor
	err                   error
	CurrentDirectory      string
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
		CreatingNewFile:       false,
		CreatingNewDirectory:  false,
	}
}
