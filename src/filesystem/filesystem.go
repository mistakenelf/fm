package filesystem

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

func RenameDirOrFile(currentName string, newName string) {
	os.Rename(currentName, newName)
}

func MoveFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}

	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	removeError := os.Remove(src)

	if removeError != nil {
		log.Fatal("error removing file", removeError)
	}

	return
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}

	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func MoveDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = MoveDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	removeError := os.RemoveAll(src)

	if removeError != nil {
		log.Fatal("error removing directory", removeError)
	}

	return
}

func DeleteDirectory(dirname string) {
	removeError := os.RemoveAll(dirname)

	if removeError != nil {
		log.Fatal("Error deleting directory", removeError)
	}
}

func DeleteFile(filename string) {
	removeError := os.Remove(filename)

	if removeError != nil {
		log.Fatal("Error deleting file", removeError)
	}
}
