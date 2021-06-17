package app

import (
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
	// Setup default and load the users
	// app configuration
	config.SetDefaults()
	config.LoadConfig()

	cfg := config.GetConfig()

	var files []fs.FileInfo
	var startDir string

	// If logging is enabled in the config, log to the
	// debug.log file
	if cfg.Settings.EnableLogging {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		defer f.Close()
	}

	// FM can be started with a starting directory as an argument,
	// check if an argument exists and use it as a starting directory
	if len(os.Args) > 1 {
		startDir = os.Args[1]
	}

	// Get the initial directory listing to be displayed
	if _, err := os.Stat(startDir); err == nil {
		files, err = utils.GetDirectoryListing(startDir, true)
		if err != nil {
			log.Fatal(err)
		}
	} else if cfg.Settings.StartDir == constants.HomeDirectory {
		homeDir, err := utils.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		files, err = utils.GetDirectoryListing(homeDir, true)
		if err != nil {
			log.Fatal(err)
		}
	} else if _, err := os.Stat(cfg.Settings.StartDir); err == nil {
		files, err = utils.GetDirectoryListing(cfg.Settings.StartDir, true)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		files, err = utils.GetDirectoryListing(constants.CurrentDirectory, true)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create new model of UI passing it the inital direcotry listing
	m := ui.NewModel(files)

	var opts []tea.ProgramOption

	// Always append alt screen program option
	opts = append(opts, tea.WithAltScreen())

	// If mousewheel is enabled, append it to the program options
	if cfg.Settings.EnableMouseWheel {
		opts = append(opts, tea.WithMouseAllMotion())
	}

	// Initialize new app
	p := tea.NewProgram(m, opts...)

	// If the program fails to start, exit the program
	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}
