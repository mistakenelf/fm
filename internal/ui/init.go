package ui

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/directory"

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

	switch {
	case startDir != "":
		_, err := os.Stat(startDir)
		if err != nil {
			return nil
		}

		if strings.HasPrefix(startDir, "/") {
			cmds = append(cmds, m.updateDirectoryListing(startDir))
		} else {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			filePath := filepath.Join(path, startDir)

			cmds = append(cmds, m.updateDirectoryListing(filePath))
		}
	case m.appConfig.Settings.StartDir == directory.HomeDirectory:
		homeDir, err := directory.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		cmds = append(cmds, m.updateDirectoryListing(homeDir))
	default:
		cmds = append(cmds, m.updateDirectoryListing(m.appConfig.Settings.StartDir))
	}

	cmds = append(cmds, spinner.Tick)

	return tea.Batch(cmds...)
}
