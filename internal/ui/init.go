package ui

import (
	"log"
	"os"

	"github.com/knipferrc/fm/directory"
	"github.com/knipferrc/fm/internal/constants"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the UI and sets up initial data.
func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	startDir := ""

	// If a starting directory was specified, use it.
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	}

	// Get the initial directory listing to be displayed
	if _, err := os.Stat(startDir); err == nil {
		cmds = append(cmds, m.updateDirectoryListing(startDir))
	} else if m.appConfig.Settings.StartDir == constants.Directories.HomeDirectory {
		homeDir, err := directory.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		cmds = append(cmds, m.updateDirectoryListing(homeDir))
	} else {
		cmds = append(cmds, m.updateDirectoryListing(m.appConfig.Settings.StartDir))
	}

	cmds = append(cmds, spinner.Tick)

	return tea.Batch(cmds...)
}
