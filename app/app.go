package app

import (
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
	m := ui.NewModel()

	if cfg.Settings.StartDir == constants.HomeDirectory {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		m.Files = utils.GetDirectoryListing(home, cfg.Settings.ShowHidden)
	} else {
		m.Files = utils.GetDirectoryListing(cfg.Settings.StartDir, cfg.Settings.ShowHidden)
	}

	m.DirTree.SetContent(m.Files, m.Cursor)

	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}
