package app

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/ui"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	config.SetDefaults()
	config.LoadConfig()

	cfg := config.GetConfig()

	var files []fs.FileInfo
	var startDir string

	if cfg.Settings.EnableLogging {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}

		defer f.Close()
	}

	if len(os.Args) > 1 {
		startDir = os.Args[1]
	}

	if startDir != "" {
		files, _ = utils.GetDirectoryListing(startDir, true)
	} else if cfg.Settings.StartDir == constants.HomeDirectory {
		homeDir, _ := utils.GetHomeDirectory()
		files, _ = utils.GetDirectoryListing(homeDir, true)
	} else if _, err := os.Stat(cfg.Settings.StartDir); err == nil {
		files, _ = utils.GetDirectoryListing(cfg.Settings.StartDir, true)
	} else {
		files, _ = utils.GetDirectoryListing(constants.CurrentDirectory, true)
	}

	m := ui.NewModel(files)

	var opts []tea.ProgramOption
	opts = append(opts, tea.WithAltScreen())

	if cfg.Settings.EnableMouseWheel {
		opts = append(opts, tea.WithMouseAllMotion())
	}

	p := tea.NewProgram(m, opts...)

	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}
