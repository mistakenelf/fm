package fm

import (
	"log"
	"os"

	"github.com/knipferrc/fm/internal/components"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/constants"
	"github.com/knipferrc/fm/internal/ui"
	"github.com/knipferrc/fm/internal/utils"

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

		m.Files = utils.GetDirectoryListing(home)
	} else {
		m.Files = utils.GetDirectoryListing(cfg.Settings.StartDir)
	}

	m.PrimaryViewport.SetContent(components.DirTree(m.Files, m.Cursor, m.ScreenWidth))
	m.SecondaryViewport.SetContent(components.Instructions())

	p := tea.NewProgram(m)

	p.EnableMouseCellMotion()
	defer p.DisableMouseCellMotion()

	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}
