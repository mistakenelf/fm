package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type fileStatus []fs.FileInfo

func (m model) getDirectoryListing() tea.Msg {
	files, err := ioutil.ReadDir("./")
	os.Chdir("./")

	if err != nil {
		log.Fatal(err)
	}

	m.Files = append(m.Files, files...)

	return fileStatus(m.Files)
}

func (m model) Init() tea.Cmd {
	return m.getDirectoryListing
}
