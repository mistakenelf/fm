package helpers

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

	return string(dat), nil
}

// Create a new file given a name and return an
// error if it exists
func CreateFile(name string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// Zip a directory given a name and return an error
// if it exists
func ZipDirectory(name string) error {
	var files []string

	// Walk the directory to get a list of files within it and append it to
	// the array of files to return
	err := filepath.Walk(name, func(path string, f fs.FileInfo, err error) error {
		if f.Name() != "." && !f.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Generate output name based on the directorys current name along
	// with a timestamp to make names unique
	output := fmt.Sprintf("%s_%d.zip", name, time.Now().Unix())
	newfile, err := os.Create(output)
	if err != nil {
		return err
	}

	defer func() {
		newfile.Close()
	}()

	zipWriter := zip.NewWriter(newfile)

	defer func() {
		zipWriter.Close()
	}()

	for _, file := range files {
		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}

		defer func() {
			zipfile.Close()
		}()

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

	defer func() error {
		err = r.Close()
		if err != nil {
			return err
		}

		return nil
	}()

	// Generate the name to unzip to based on its current name
	// minus the extension
	output := strings.Split(name, ".")[0]

	for _, f := range r.File {
		fpath := filepath.Join(output, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(output)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}

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
		if err != nil {
			return err
		}

		err = outFile.Close()
		if err != nil {
			return err
		}

		err = rc.Close()
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Copy a file given a name
func CopyFile(name string) error {
	srcFile, err := os.Open(name)
	if err != nil {
		return err
	}

	defer func() error {
		err = srcFile.Close()
		if err != nil {
			return err
		}

		return nil
	}()

	splitName := strings.Split(name, ".")
	output := fmt.Sprintf("%s_%d.%s", splitName[0], time.Now().Unix(), splitName[1])
	destFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// Copy a directory given a name
func CopyDirectory(name string) error {
	// Generate a unique name for the output folder
	output := fmt.Sprintf("%s_%d", name, time.Now().Unix())

	f, err := os.Open(name)
	if err != nil {
		return err
	}

	file, err := f.Stat()
	if err != nil {
		return err
	}

	if !file.IsDir() {
		return fmt.Errorf("Source " + file.Name() + " is not a directory!")
	}

	// Create the output folder
	err = os.Mkdir(output, 0755)
	if err != nil {
		return err
	}

	// Read all files in the directory
	files, err := ioutil.ReadDir(name)
	if err != nil {
		return err
	}

	// Loop through the directory getting a list of all its files
	for _, f := range files {
		// If its a directory, copy it
		if f.IsDir() {
			err = CopyDirectory(name + "/" + f.Name())
			if err != nil {
				return err
			}
		}

		// If its not a directory, read the file and write it to the new folder
		if !f.IsDir() {
			content, err := ioutil.ReadFile(name + "/" + f.Name())
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(output+"/"+f.Name(), content, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
