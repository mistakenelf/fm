package tui

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/dirfs"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

// Init initializes the UI and sets up initial data.
func (b Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd
	startDir := viper.GetString("start-dir")

	switch {
	case startDir != "":
		_, err := os.Stat(startDir)
		if err != nil {
			return nil
		}

		if strings.HasPrefix(startDir, "/") {
			cmds = append(cmds, b.updateDirectoryListingCmd(startDir))
		} else {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			filePath := filepath.Join(path, startDir)

			cmds = append(cmds, b.updateDirectoryListingCmd(filePath))
		}
	case b.appConfig.Settings.StartDir == dirfs.HomeDirectory:
		homeDir, err := dirfs.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		cmds = append(cmds, b.updateDirectoryListingCmd(homeDir))
	default:
		cmds = append(cmds, b.updateDirectoryListingCmd(b.appConfig.Settings.StartDir))
	}

	cmds = append(cmds, spinner.Tick)
	cmds = append(cmds, textinput.Blink)

	return tea.Batch(cmds...)
}
