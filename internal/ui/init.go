package ui

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/dirfs"

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
			cmds = append(cmds, m.updateDirectoryListingCmd(startDir))
		} else {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			filePath := filepath.Join(path, startDir)

			cmds = append(cmds, m.updateDirectoryListingCmd(filePath))
		}
	case m.appConfig.Settings.StartDir == dirfs.HomeDirectory:
		homeDir, err := dirfs.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		cmds = append(cmds, m.updateDirectoryListingCmd(homeDir))
	default:
		cmds = append(cmds, m.updateDirectoryListingCmd(m.appConfig.Settings.StartDir))
	}

	cmds = append(cmds, spinner.Tick)

	return tea.Batch(cmds...)
}
