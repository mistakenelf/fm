package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type fileStatus []fs.FileInfo

func (m model) getInitialDirectoryListing() tea.Msg {
	files, err := ioutil.ReadDir("./")
	os.Chdir("./")

	if err != nil {
		log.Fatal(err)
	}

	m.Files = append(m.Files, files...)

	return fileStatus(m.Files)
}

func getUpdatedDirectoryListing(dir string) []fs.FileInfo {
	files, err := ioutil.ReadDir(dir)
	curFiles := make([]fs.FileInfo, 0)
	os.Chdir(dir)

	if err != nil {
		log.Fatal(err)
	}

	curFiles = append(curFiles, files...)

	return curFiles
}
