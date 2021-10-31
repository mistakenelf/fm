package dirfs

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Constants to represent different directories.
const (
	CurrentDirectory  = "."
	PreviousDirectory = ".."
	HomeDirectory     = "~"
	RootDirectory     = "/"
)

// RenameDirectoryItem renames a directory or files given a source and destination.
func RenameDirectoryItem(src, dst string) error {
	err := os.Rename(src, dst)

	return err
}

// CreateDirectory creates a new directory given a name.
func CreateDirectory(name string) error {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(name, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetDirectoryListing returns a list of files and directories within a given directory.
func GetDirectoryListing(dir string, showHidden bool) ([]fs.DirEntry, error) {
	n := 0

	// Read files from the directory.
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if !showHidden {
		for _, file := range files {
			// If the file or directory starts with a dot,
			// we know its hidden so dont add it to the array
			// of files to return.
			if !strings.HasPrefix(file.Name(), ".") {
				files[n] = file
				n++
			}
		}

		// Set files to the list that does not include hidden files.
		files = files[:n]
	}

	return files, nil
}

// GetDirectoryListingByType returns a directory listing based on type (directories | files).
func GetDirectoryListingByType(dir, listType string, showHidden bool) ([]fs.DirEntry, error) {
	n := 0

	// Read files from the directory.
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		switch {
		case file.IsDir() && listType == "directories" && !showHidden:
			if !strings.HasPrefix(file.Name(), ".") {
				files[n] = file
				n++
			}
		case file.IsDir() && listType == "directories" && showHidden:
			files[n] = file
			n++
		case !file.IsDir() && listType == "files" && !showHidden:
			if !strings.HasPrefix(file.Name(), ".") {
				files[n] = file
				n++
			}
		case !file.IsDir() && listType == "files" && showHidden:
			files[n] = file
			n++
		}
	}

	return files[:n], nil
}

// DeleteDirectory deletes a directory given a name.
func DeleteDirectory(name string) error {
	err := os.RemoveAll(name)

	return err
}

// GetHomeDirectory returns the users home directory.
func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}

// GetWorkingDirectory returns the current working directory.
func GetWorkingDirectory() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return workingDir, nil
}

// DeleteFile deletes a file given a name.
func DeleteFile(name string) error {
	err := os.Remove(name)

	return err
}

// MoveDirectoryItem moves a file from one place to another.
func MoveDirectoryItem(src, dst string) error {
	err := os.Rename(src, dst)

	return err
}

// ReadFileContent returns the contents of a file given a name.
func ReadFileContent(name string) (string, error) {
	fileContent, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}

// CreateFile creates a file given a name.
func CreateFile(name string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return err
}

// Zip zips a directory given a name.
func Zip(name string) error {
	// Generate output name based on the directories current name along
	// with a timestamp to make names unique.
	output := fmt.Sprintf("%s_%d.zip", strings.Split(name, ".")[0], time.Now().Unix())
	newfile, err := os.Create(output)
	if err != nil {
		return err
	}

	defer func() {
		err = newfile.Close()
	}()

	zipWriter := zip.NewWriter(newfile)
	defer func() {
		err = zipWriter.Close()
	}()

	info, err := os.Stat(name)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err = filepath.Walk(name, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if err != nil {
				return err
			}

			relPath := strings.TrimPrefix(filePath, name)
			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			fsFile, err := os.Open(filePath)
			if err != nil {
				return err
			}

			_, err = io.Copy(zipFile, fsFile)
			if err != nil {
				return err
			}

			return nil
		})
	} else {
		var files []string

		// Walk the directory to get a list of files within it and append it to
		// the array of files.
		err := filepath.Walk(name, func(path string, f fs.FileInfo, err error) error {
			if f.Name() != "." && !f.IsDir() {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			return err
		}

		for _, file := range files {
			zipfile, err := os.Open(file)
			if err != nil {
				return err
			}

			defer func() {
				err = zipfile.Close()
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
	}

	if err != nil {
		return err
	}

	err = zipWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

// Unzip unzips a directory given a name.
func Unzip(name string) error {
	r, err := zip.OpenReader(name)
	if err != nil {
		return err
	}

	defer func() {
		err = r.Close()
	}()

	// Generate the name to unzip to based on its current name
	// minus the extension.
	output := strings.Split(name, ".")[0]

	for _, f := range r.File {
		archiveFile := f.Name
		fpath := filepath.Join(output, archiveFile)

		if !strings.HasPrefix(fpath, filepath.Clean(output)+string(os.PathSeparator)) {
			return err
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

// CopyFile copies a file given a name.
func CopyFile(name string) error {
	var splitName []string
	var output string

	srcFile, err := os.Open(name)
	if err != nil {
		return err
	}

	defer func() {
		err = srcFile.Close()
	}()

	fileExtension := filepath.Ext(name)
	switch {
	case strings.HasPrefix(name, ".") && fileExtension != "" && fileExtension == name:
		output = fmt.Sprintf("%s_%d", name, time.Now().Unix())
	case strings.HasPrefix(name, ".") && fileExtension != "" && fileExtension != name:
		splitName = strings.Split(name, ".")
		output = fmt.Sprintf(".%s_%d.%s", splitName[1], time.Now().Unix(), splitName[2])
	case fileExtension != "":
		splitName = strings.Split(name, ".")
		output = fmt.Sprintf("%s_%d.%s", splitName[0], time.Now().Unix(), splitName[1])
	default:
		output = fmt.Sprintf("%s_%d", name, time.Now().Unix())
	}

	destFile, err := os.Create(output)
	if err != nil {
		return err
	}

	defer func() {
		err = destFile.Close()
	}()

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

// CopyDirectory copies a directory given a name.
func CopyDirectory(name string) error {
	// Generate a unique name for the output folder.
	output := fmt.Sprintf("%s_%d", name, time.Now().Unix())

	err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		relPath := strings.Replace(path, name, "", 1)

		if info.IsDir() {
			return os.Mkdir(filepath.Join(output, relPath), os.ModePerm)
		}

		var data, err1 = os.ReadFile(filepath.Join(name, relPath))
		if err1 != nil {
			return err1
		}

		return os.WriteFile(filepath.Join(output, relPath), data, os.ModePerm)
	})

	return err
}

// GetDirectoryItemSize calculates the size of a directory or file.
func GetDirectoryItemSize(path string) (int64, error) {
	curFile, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !curFile.IsDir() {
		return curFile.Size(), nil
	}

	var size int64
	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if !d.IsDir() {
			size += fileInfo.Size()
		}

		return err
	})

	return size, err
}

// WriteToFile writes content to a file.
func WriteToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s\n", filepath.Join(workingDir, content)))
	if err != nil {
		f.Close()
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
