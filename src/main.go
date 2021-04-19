package main

import (
	"log"
	"os"

	"github.com/knipferrc/fm/src/config"
	"github.com/knipferrc/fm/src/filesystem"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/knipferrc/fm/src/components"
)

func main() {
	config.SetDefaults()
	config.LoadConfig()

	config := config.GetConfig()
	m := createModel()

	if config.StartDir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		m.files = filesystem.GetDirectoryListing(home)
	} else {
		m.files = filesystem.GetDirectoryListing(config.StartDir)
	}
	m.viewport.SetContent(components.DirTree(m.files, m.cursor, m.screenwidth))

	p := tea.NewProgram(m)

	p.EnterAltScreen()
	defer p.ExitAltScreen()

	p.EnableMouseCellMotion()
	defer p.DisableMouseCellMotion()

	if err := p.Start(); err != nil {
		log.Fatal("Failed to start fm", err)
		os.Exit(1)
	}
}
