package main

import (
	"io/fs"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initialModel := model{make([]fs.FileInfo, 0), make(map[int]struct{}), 0, false, ""}
	p := tea.NewProgram(initialModel)

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
