package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func ZipDirectory(dir string) error {
	var files []string

	filepath.Walk(dir, func(path string, f fs.FileInfo, err error) error {
		if f.Name() != "." && !f.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	output := fmt.Sprintf("%s_%d.zip", dir, time.Now().Unix())
	newfile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer newfile.Close()

	zipWriter := zip.NewWriter(newfile)
	defer zipWriter.Close()

	for _, file := range files {
		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, zipfile)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnzipDirectory(dir string) error {
	r, err := zip.OpenReader(dir)
	if err != nil {
		return err
	}
	defer r.Close()

	output := strings.Split(dir, ".")[0]

	for _, f := range r.File {
		fpath := filepath.Join(output, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(output)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
