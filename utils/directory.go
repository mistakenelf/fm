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

// Rename a file or directory given a source and destination,
// returning an error if it exists
func RenameDirOrFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}

	return nil
}

// Create a new directory given a name,
// returning an error if it exists
func CreateDirectory(name string) error {
	_, err := os.Stat(name)

	// If the directory does not already exist, create it
	if os.IsNotExist(err) {
		err := os.MkdirAll(name, 0755)
		if err != nil {
			return err
		}

	}

	return nil
}

// Get directory listing based on the name and weather or not to show hidden files
// and folders, returning the new file listing and an error if it exists
func GetDirectoryListing(dir string, showHidden bool) ([]fs.FileInfo, error) {
	n := 0

	// Read files from the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Change the apps directory to the one passed in
	err = os.Chdir(dir)
	if err != nil {
		return nil, err
	}

	// Dont want to show hidden files and directories
	if !showHidden {
		for _, file := range files {
			// If the file or directory starts with a dot,
			// we know its hidden so dont add it to the array
			// of files to return
			if !strings.HasPrefix(file.Name(), ".") {
				files[n] = file
				n++
			}
		}

		// Set files to the list that does not include hidden files
		files = files[:n]
	}

	// return the files and nil since no error occured
	return files, nil
}

// Delete a directory given a name,
// returning an error if it exists
func DeleteDirectory(name string) error {
	err := os.RemoveAll(name)

	// Something went wrong removing the directory
	if err != nil {
		return err
	}

	return nil
}

// Move a directory given a source and destination,
// returning an error if it exists
func MoveDirectory(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}

	return nil
}

// Get the users home directoring returning an
// error if it exists
func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Return the home directory
	return home, nil
}

// Get the users current working directory
// returning an error if it exists
func GetWorkingDirectory() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return directory, nil
}

// Delete a file given the name
// returning an error if it exists
func DeleteFile(name string) error {
	err := os.Remove(name)
	if err != nil {
		return err
	}

	return nil
}

// Move file from one place to another given a source
// and destination, returning an error if it exists
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}

	return nil
}

// Read a files content given a name returning its content and
// an error if it exists
func ReadFileContent(name string) (string, error) {
	dat, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	// Return file data as a string and no error
	return string(dat), nil
}

// Create a new file given a name and return an
// error if it exists
func CreateFile(name string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Close the file that was created
	f.Close()

	return nil
}

// Zip a directory given a name and return an error
// if it exists
func ZipDirectory(name string) error {
	var files []string

	// Walk the directory to get a list of files within it and append it to
	// the array of files to return
	filepath.Walk(name, func(path string, f fs.FileInfo, err error) error {
		if f.Name() != "." && !f.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	// Generate output name based on the directorys current name along
	// with a timestamp to make names unique
	output := fmt.Sprintf("%s_%d.zip", name, time.Now().Unix())
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

// Unzip a directory given a name returning an error
// if it exists
func UnzipDirectory(name string) error {
	r, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer r.Close()

	// Generate the name to unzip to based on its current name
	// minus the extension
	output := strings.Split(name, ".")[0]

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
