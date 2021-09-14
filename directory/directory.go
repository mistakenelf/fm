package directory

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

// RenameDirOrFile renames a directory or files given a source and destination.
func RenameDirOrFile(src, dst string) error {
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

	// Update the apps directory to the directory currently being read.
	err = os.Chdir(dir)
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

// DeleteDirectory deletes a directory given a name.
func DeleteDirectory(name string) error {
	err := os.RemoveAll(name)

	return err
}

// MoveDirectory moves a directory from one place to another.
func MoveDirectory(src, dst string) error {
	err := os.Rename(src, dst)

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
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return directory, nil
}

// DeleteFile deletes a file given a name.
func DeleteFile(name string) error {
	err := os.Remove(name)

	return err
}

// MoveFile moves a file from one place to another.
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)

	return err
}

// ReadFileContent returns the contents of a file given a name.
func ReadFileContent(name string) (string, error) {
	dat, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(dat), nil
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

// ZipDirectory zips a directory given a name.
func ZipDirectory(name string) error {
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

	return nil
}

// UnzipDirectory unzips a directory given a name.
func UnzipDirectory(name string) error {
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
	srcFile, err := os.Open(name)
	if err != nil {
		return err
	}

	defer func() {
		err = srcFile.Close()
	}()

	splitName := strings.Split(name, ".")
	output := fmt.Sprintf("%s_%d.%s", splitName[0], time.Now().Unix(), splitName[1])
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

	var err error = filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, name, "", 1)

		if info.IsDir() {
			return os.Mkdir(filepath.Join(output, relPath), os.ModePerm)
		} else {
			var data, err1 = os.ReadFile(filepath.Join(name, relPath))
			if err1 != nil {
				return err1
			}
			return os.WriteFile(filepath.Join(output, relPath), data, os.ModePerm)
		}
	})

	return err
}

// GetTreeSize calculates the size of a directory recursively.
func GetTreeSize(path string) (int64, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	var total int64
	for _, entry := range entries {
		if entry.IsDir() {
			size, err := GetTreeSize(filepath.Join(path, entry.Name()))
			if err != nil {
				return 0, err
			}

			total += size
		} else {
			info, err := entry.Info()

			if err != nil {
				return 0, err
			}

			total += info.Size()
		}
	}

	return total, nil
}
