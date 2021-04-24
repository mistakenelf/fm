package filesystem

import (
	"io"
	"log"
	"os"
)

func DeleteFile(filename string) {
	removeError := os.Remove(filename)

	if removeError != nil {
		log.Fatal("Error deleting file", removeError)
	}
}

func CopyFile(src, dst string, shouldRemove bool) (err error) {
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

	if shouldRemove {
		removeError := os.Remove(src)

		if removeError != nil {
			log.Fatal("error removing file", removeError)
		}
	}

	return
}

func ReadFileContent(name string) string {
	dat, err := os.ReadFile(name)

	if err != nil {
		log.Fatal("Error occured reading file")
	}

	return string(dat)
}
