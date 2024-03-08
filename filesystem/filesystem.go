// Package filesystem is a collection of various different filesystem
// helper functions.
package filesystem

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

// Directory shortcuts.
const (
	CurrentDirectory  = "."
	PreviousDirectory = ".."
	HomeDirectory     = "~"
	RootDirectory     = "/"
)

// Different types of listings.
const (
	DirectoriesListingType = "directories"
	FilesListingType       = "files"
)

// RenameDirectoryItem renames a directory or files given a source and destination.
func RenameDirectoryItem(src, dst string) error {
	err := os.Rename(src, dst)

	return errors.Unwrap(err)
}

// CreateDirectory creates a new directory given a name.
func CreateDirectory(name string) error {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(name, os.ModePerm)
		if err != nil {
			return errors.Unwrap(err)
		}
	}

	return nil
}

// GetDirectoryListing returns a list of files and directories within a given directory.
func GetDirectoryListing(dir string, showHidden bool) ([]fs.DirEntry, error) {
	index := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	if !showHidden {
		for _, file := range files {
			// If the file or directory starts with a dot,
			// we know its hidden so dont add it to the array
			// of files to return.
			if !strings.HasPrefix(file.Name(), ".") {
				files[index] = file
				index++
			}
		}

		// Set files to the list that does not include hidden files.
		files = files[:index]
	}

	return files, nil
}

// GetDirectoryListingByType returns a directory listing based on type (directories | files).
func GetDirectoryListingByType(dir, listingType string, showHidden bool) ([]fs.DirEntry, error) {
	index := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Unwrap(err)
	}

	for _, file := range files {
		switch {
		case file.IsDir() && listingType == DirectoriesListingType && !showHidden:
			if !strings.HasPrefix(file.Name(), ".") {
				files[index] = file
				index++
			}
		case file.IsDir() && listingType == DirectoriesListingType && showHidden:
			files[index] = file
			index++
		case !file.IsDir() && listingType == FilesListingType && !showHidden:
			if !strings.HasPrefix(file.Name(), ".") {
				files[index] = file
				index++
			}
		case !file.IsDir() && listingType == FilesListingType && showHidden:
			files[index] = file
			index++
		}
	}

	return files[:index], nil
}

// DeleteDirectory deletes a directory given a name.
func DeleteDirectory(name string) error {
	err := os.RemoveAll(name)

	return errors.Unwrap(err)
}

// GetHomeDirectory returns the users home directory.
func GetHomeDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return home, nil
}

// GetWorkingDirectory returns the current working directory.
func GetWorkingDirectory() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return workingDir, nil
}

// DeleteFile deletes a file given a name.
func DeleteFile(name string) error {
	err := os.Remove(name)

	return errors.Unwrap(err)
}

// MoveDirectoryItem moves a file from one place to another.
func MoveDirectoryItem(src, dst string) error {
	err := os.Rename(src, dst)

	return errors.Unwrap(err)
}

// ReadFileContent returns the contents of a file given a name.
func ReadFileContent(name string) (string, error) {
	fileContent, err := os.ReadFile(filepath.Clean(name))
	if err != nil {
		return "", errors.Unwrap(err)
	}

	return string(fileContent), nil
}

// CreateFile creates a file given a name.
func CreateFile(name string) error {
	f, err := os.Create(filepath.Clean(name))
	if err != nil {
		return errors.Unwrap(err)
	}

	if err = f.Close(); err != nil {
		return errors.Unwrap(err)
	}

	return errors.Unwrap(err)
}

// Zip zips a directory given a name.
func Zip(name string) error {
	var splitName []string
	var output string

	srcFile, err := os.Open(filepath.Clean(name))
	if err != nil {
		return errors.Unwrap(err)
	}

	defer func() {
		err = srcFile.Close()
	}()

	fileExtension := filepath.Ext(name)
	splitFileName := strings.Split(name, "/")
	fileName := splitFileName[len(splitFileName)-1]
	switch {
	case strings.HasPrefix(fileName, ".") && fileExtension != "" && fileExtension == fileName:
		output = fmt.Sprintf("%s_%d.zip", fileName, time.Now().Unix())
	case strings.HasPrefix(fileName, ".") && fileExtension != "" && fileExtension != fileName:
		splitName = strings.Split(fileName, ".")
		output = fmt.Sprintf(".%s_%d.zip", splitName[1], time.Now().Unix())
	case fileExtension != "":
		splitName = strings.Split(fileName, ".")
		output = fmt.Sprintf("%s_%d.zip", splitName[0], time.Now().Unix())
	default:
		output = fmt.Sprintf("%s_%d.zip", fileName, time.Now().Unix())
	}

	newfile, err := os.Create(filepath.Clean(output))
	if err != nil {
		return errors.Unwrap(err)
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
		return errors.Unwrap(err)
	}

	if info.IsDir() {
		err = filepath.Walk(name, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if err != nil {
				return errors.Unwrap(err)
			}

			relPath := strings.TrimPrefix(filePath, name)
			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return errors.Unwrap(err)
			}

			fsFile, err := os.Open(filepath.Clean(filePath))
			if err != nil {
				return errors.Unwrap(err)
			}

			_, err = io.Copy(zipFile, fsFile)
			if err != nil {
				return errors.Unwrap(err)
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
			return errors.Unwrap(err)
		}

		for _, file := range files {
			zipfile, err := os.Open(filepath.Clean(file))
			if err != nil {
				return errors.Unwrap(err)
			}

			defer func() {
				err = zipfile.Close()
			}()

			info, err := zipfile.Stat()
			if err != nil {
				return errors.Unwrap(err)
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return errors.Unwrap(err)
			}

			header.Method = zip.Deflate
			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return errors.Unwrap(err)
			}

			_, err = io.Copy(writer, zipfile)
			if err != nil {
				return errors.Unwrap(err)
			}
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return errors.Unwrap(err)
	}

	return errors.Unwrap(err)
}

// Unzip unzips a directory given a name.
func Unzip(name string) error {
	var output string

	reader, err := zip.OpenReader(name)
	if err != nil {
		return errors.Unwrap(err)
	}

	defer func() {
		err = reader.Close()
	}()

	if strings.HasPrefix(name, ".") {
		output = strings.Split(name, ".")[1]
	} else {
		output = strings.Split(name, ".")[0]
	}

	for _, file := range reader.File {
		archiveFile := file.Name
		fpath := filepath.Join(output, archiveFile)

		if !strings.HasPrefix(fpath, filepath.Clean(output)+string(os.PathSeparator)) {
			return errors.Unwrap(err)
		}

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return errors.Unwrap(err)
			}

			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return errors.Unwrap(err)
		}

		outFile, err := os.OpenFile(filepath.Clean(fpath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return errors.Unwrap(err)
		}

		outputFile, err := file.Open()
		if err != nil {
			return errors.Unwrap(err)
		}

		_, err = io.Copy(outFile, outputFile)
		if err != nil {
			return errors.Unwrap(err)
		}

		err = outFile.Close()
		if err != nil {
			return errors.Unwrap(err)
		}

		err = outputFile.Close()
		if err != nil {
			return errors.Unwrap(err)
		}
	}

	return errors.Unwrap(err)
}

// CopyFile copies a file given a name.
func CopyFile(name string) error {
	var splitName []string
	var output string

	srcFile, err := os.Open(filepath.Clean(name))
	if err != nil {
		return errors.Unwrap(err)
	}

	defer func() {
		err = srcFile.Close()
	}()

	fileExtension := filepath.Ext(name)
	splitFileName := strings.Split(name, "/")
	fileName := splitFileName[len(splitFileName)-1]
	switch {
	case strings.HasPrefix(fileName, ".") && fileExtension != "" && fileExtension == fileName:
		output = fmt.Sprintf("%s_%d", fileName, time.Now().Unix())
	case strings.HasPrefix(fileName, ".") && fileExtension != "" && fileExtension != fileName:
		splitName = strings.Split(fileName, ".")
		output = fmt.Sprintf(".%s_%d.%s", splitName[1], time.Now().Unix(), splitName[2])
	case fileExtension != "":
		splitName = strings.Split(fileName, ".")
		output = fmt.Sprintf("%s_%d.%s", splitName[0], time.Now().Unix(), splitName[1])
	default:
		output = fmt.Sprintf("%s_%d", fileName, time.Now().Unix())
	}

	destFile, err := os.Create(filepath.Clean(output))
	if err != nil {
		return errors.Unwrap(err)
	}

	defer func() {
		err = destFile.Close()
	}()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return errors.Unwrap(err)
	}

	err = destFile.Sync()
	if err != nil {
		return errors.Unwrap(err)
	}

	return errors.Unwrap(err)
}

// CopyDirectory copies a directory given a name.
func CopyDirectory(name string) error {
	output := fmt.Sprintf("%s_%d", name, time.Now().Unix())

	err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		relPath := strings.Replace(path, name, "", 1)

		if info.IsDir() {
			return fmt.Errorf("%w", os.Mkdir(filepath.Join(output, relPath), os.ModePerm))
		}

		var data, err1 = os.ReadFile(filepath.Join(filepath.Clean(name), filepath.Clean(relPath)))
		if err1 != nil {
			return errors.Unwrap(err)
		}

		return fmt.Errorf("%w", os.WriteFile(filepath.Join(output, relPath), data, os.ModePerm))
	})

	return errors.Unwrap(err)
}

// GetDirectoryItemSize calculates the size of a directory or file.
func GetDirectoryItemSize(path string) (int64, error) {
	curFile, err := os.Stat(path)
	if err != nil {
		return 0, errors.Unwrap(err)
	}

	if !curFile.IsDir() {
		return curFile.Size(), nil
	}

	var size int64
	err = filepath.WalkDir(path, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return errors.Unwrap(err)
		}

		fileInfo, err := entry.Info()
		if err != nil {
			return errors.Unwrap(err)
		}

		if !entry.IsDir() {
			size += fileInfo.Size()
		}

		return errors.Unwrap(err)
	})

	return size, errors.Unwrap(err)
}

// FindFilesByName returns files found based on a name.
func FindFilesByName(name, dir string) ([]string, []fs.DirEntry, error) {
	var paths []string
	var entries []fs.DirEntry

	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return filepath.SkipDir
		}

		if strings.Contains(entry.Name(), name) {
			paths = append(paths, path)
			entries = append(entries, entry)
		}

		return errors.Unwrap(err)
	})

	return paths, entries, errors.Unwrap(err)
}

// WriteToFile writes content to a file, overwriting content if it exists.
func WriteToFile(path, content string) error {
	file, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.Unwrap(err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return errors.Unwrap(err)
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", filepath.Join(workingDir, content)))
	if err != nil {
		err = file.Close()
		if err != nil {
			return errors.Unwrap(err)
		}

		return errors.Unwrap(err)
	}

	err = file.Close()
	if err != nil {
		return errors.Unwrap(err)
	}

	return errors.Unwrap(err)
}
