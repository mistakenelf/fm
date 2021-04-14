package directory

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

func GetDirectoryListing(dir string) []fs.FileInfo {
	files, err := ioutil.ReadDir(dir)
	curFiles := make([]fs.FileInfo, 0)
	os.Chdir(dir)

	if err != nil {
		log.Fatal(err)
	}

	curFiles = append(curFiles, files...)

	return curFiles
}
