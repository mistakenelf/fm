package utils

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func RenameDirOrFile(src, dst string) {
	os.Rename(src, dst)
}

func CreateDirectory(name string) {
	_, err := os.Stat(name)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(name, 0755)

		if errDir != nil {
			log.Fatal(err)
		}

	}
}

func GetDirectoryListing(dir string, showHidden bool) []fs.FileInfo {
	n := 0

	files, err := ioutil.ReadDir(dir)
	os.Chdir(dir)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if !showHidden {
		for _, file := range files {
			if !strings.HasPrefix(file.Name(), ".") {
				files[n] = file
				n++
			}
		}

		files = files[:n]
	}

	return files
}

func DeleteDirectory(dirname string) {
	removeError := os.RemoveAll(dirname)

	if removeError != nil {
		log.Fatal("Error deleting directory", removeError)
	}
}

func MoveDirectory(src, dst string) {
	os.Rename(src, dst)
}

func GetHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return home
}

func GetWorkingDirectory() string {
	directory, err := os.Getwd()

	if err != nil {
		log.Fatal("error getting working directory")
	}

	return directory
}

func DeleteFile(filename string) {
	removeError := os.Remove(filename)

	if removeError != nil {
		log.Fatal("Error deleting file", removeError)
	}
}

func MoveFile(src, dst string) {
	os.Rename(src, dst)
}

func ReadFileContent(name string) string {
	dat, err := os.ReadFile(name)

	if err != nil {
		log.Fatal("Error occured reading file")
	}

	return string(dat)
}

func CreateFile(name string) {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	f.Close()
}
