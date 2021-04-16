package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type fileStatus []fs.FileInfo

func getDirectoryListing() tea.Msg {
	files, err := ioutil.ReadDir("./")
	os.Chdir("./")

	if err != nil {
		log.Fatal(err)
	}

	return fileStatus(files)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getDirectoryListing, textinput.Blink)
}
