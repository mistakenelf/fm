package app

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/constants"
	"github.com/knipferrc/fm/dirtree"
	"github.com/knipferrc/fm/ui"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	config.SetDefaults()
	config.LoadConfig()

	cfg := config.GetConfig()
	m := ui.NewModel()

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
		files = utils.GetDirectoryListing(startDir, true)
	} else if cfg.Settings.StartDir == constants.HomeDirectory {
		files = utils.GetDirectoryListing(utils.GetHomeDirectory(), true)
	} else if _, err := os.Stat(cfg.Settings.StartDir); err == nil {
		files = utils.GetDirectoryListing(cfg.Settings.StartDir, true)
	} else {
		files = utils.GetDirectoryListing(".", true)
	}

	m.DirTree = dirtree.NewModel(
		files,
		cfg.Settings.ShowIcons,
		cfg.Colors.DirTree.SelectedItem,
		cfg.Colors.DirTree.UnselectedItem,
	)

	var cmds []tea.ProgramOption

	cmds = append(cmds, tea.WithAltScreen())

	if cfg.Settings.EnableMouseWheel {
		cmds = append(cmds, tea.WithMouseAllMotion())
	}

	p := tea.NewProgram(m, cmds...)

	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}
