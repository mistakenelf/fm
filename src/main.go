package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(createInitialModel())

	p.EnterAltScreen()
	defer p.ExitAltScreen()

	p.EnableMouseCellMotion()
	defer p.DisableMouseCellMotion()

	if err := p.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
