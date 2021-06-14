package utils

import (
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func RenameDirOrFile(src, dst string) error {
	err := os.Rename(src, dst)

	if err != nil {
		return err
	}

	return nil
}

func CreateDirectory(name string) error {
	_, err := os.Stat(name)

	if os.IsNotExist(err) {
		err := os.MkdirAll(name, 0755)

		if err != nil {
			return err
		}

	}

	return nil
}

func GetDirectoryListing(dir string, showHidden bool) ([]fs.FileInfo, error) {
	n := 0

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	err = os.Chdir(dir)
	if err != nil {
		return nil, err
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

	return files, nil
}

func DeleteDirectory(dirname string) error {
	err := os.RemoveAll(dirname)

	if err != nil {
		return err
	}

	return nil
}

func MoveDirectory(src, dst string) error {
	err := os.Rename(src, dst)

	if err != nil {
		return err
	}

	return nil
}

func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

func GetWorkingDirectory() (string, error) {
	directory, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return directory, nil
}

func DeleteFile(filename string) error {
	err := os.Remove(filename)

	if err != nil {
		return err
	}

	return nil
}

func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)

	if err != nil {
		return err
	}

	return nil
}

func ReadFileContent(name string) (string, error) {
	dat, err := os.ReadFile(name)

	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func CreateFile(name string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	f.Close()

	return nil
}
