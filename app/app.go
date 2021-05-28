package app

import (
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

	if cfg.Settings.StartDir == constants.HomeDirectory {
		files = utils.GetDirectoryListing(utils.GetHomeDirectory(), cfg.Settings.ShowHidden)
	} else {
		files = utils.GetDirectoryListing(cfg.Settings.StartDir, cfg.Settings.ShowHidden)
	}

	m.DirTree = dirtree.NewModel(files, cfg.Settings.ShowIcons, cfg.Colors.DirTree.SelectedItem, cfg.Colors.DirTree.UnselectedItem)

	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}
