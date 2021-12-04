package commands

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/strfmt"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type UpdateDirectoryListingMsg []fs.DirEntry
type PreviewDirectoryListingMsg []fs.DirEntry
type MoveDirItemMsg []fs.DirEntry
type ErrorMsg string
type CopyToClipboardMsg string
type ConvertImageToStringMsg string
type FindFilesByNameMsg struct {
	paths   []string
	entries []fs.DirEntry
}
type DirectoryItemSizeMsg struct {
	size  string
	index int
}
type ReadFileContentMsg struct {
	RawContent  string
	Markdown    string
	Code        string
	ImageString string
	PDFContent  string
	Image       image.Image
}

// UpdateDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func UpdateDirectoryListingCmd(name string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		files, err := dirfs.GetDirectoryListing(name, showHidden)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		err = os.Chdir(name)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		return UpdateDirectoryListingMsg(files)
	}
}

// PreviewDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func PreviewDirectoryListingCmd(name string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return ErrorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(filepath.Join(currentDir, name), showHidden)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		return PreviewDirectoryListingMsg(files)
	}
}

// RenameDirectoryItemCmd renames a file or directory based on the name and value provided.
func RenameDirectoryItemCmd(name, value string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.RenameDirectoryItem(name, value); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// MoveDirectoryItemCmd moves a file to the current working directory.
func MoveDirectoryItemCmd(name, initialMoveDirectory string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return ErrorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same file name.
		src := filepath.Join(initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had.
		dst := filepath.Join(workingDir, name)

		if err = dirfs.MoveDirectoryItem(src, dst); err != nil {
			return ErrorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(initialMoveDirectory, showHidden)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		err = os.Chdir(initialMoveDirectory)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		return MoveDirItemMsg(files)
	}
}

// RedrawImageCmd redraws the image based on the width provided.
func RedrawImageCmd(width int, img image.Image) tea.Cmd {
	return func() tea.Msg {
		imageString := strfmt.ImageToString(width, img)

		return ConvertImageToStringMsg(imageString)
	}
}

// DeleteDirectoryCmd deletes a directory based on the name provided.
func DeleteDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteDirectory(name); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// DeleteFileCmd deletes a file based on the name provided.
func DeleteFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteFile(name); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// ReadFileContentCmd reads the content of a file and returns it.
func ReadFileContentCmd(fileName, syntaxTheme string, width int, prettyMarkdown bool) tea.Cmd {
	return func() tea.Msg {
		content, err := dirfs.ReadFileContent(fileName)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		switch {
		case filepath.Ext(fileName) == ".md" && prettyMarkdown:
			markdownContent, err := strfmt.RenderMarkdown(width, content)
			if err != nil {
				return ErrorMsg(err.Error())
			}

			return ReadFileContentMsg{
				RawContent:  content,
				Markdown:    markdownContent,
				Code:        "",
				ImageString: "",
				PDFContent:  "",
				Image:       nil,
			}
		case filepath.Ext(fileName) == ".png" || filepath.Ext(fileName) == ".jpg" || filepath.Ext(fileName) == ".jpeg":
			imageContent, err := os.Open(fileName)
			if err != nil {
				return ErrorMsg(err.Error())
			}

			img, _, err := image.Decode(imageContent)
			if err != nil {
				return ErrorMsg(err.Error())
			}

			imageString := strfmt.ImageToString(width, img)

			return ReadFileContentMsg{
				RawContent:  content,
				Code:        "",
				Markdown:    "",
				ImageString: imageString,
				PDFContent:  "",
				Image:       img,
			}
		case filepath.Ext(fileName) == ".pdf":
			pdfContent, err := strfmt.ReadPdf(fileName)
			if err != nil {
				return ErrorMsg(err.Error())
			}

			return ReadFileContentMsg{
				RawContent:  content,
				Code:        "",
				Markdown:    "",
				ImageString: "",
				PDFContent:  pdfContent,
				Image:       nil,
			}
		default:
			code, err := strfmt.Highlight(content, filepath.Ext(fileName), syntaxTheme)
			if err != nil {
				return ErrorMsg(err.Error())
			}

			return ReadFileContentMsg{
				RawContent:  content,
				Code:        code,
				Markdown:    "",
				ImageString: "",
				PDFContent:  "",
				Image:       nil,
			}
		}
	}
}

// CreateDirectoryCmd creates a directory based on the name provided.
func CreateDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateDirectory(name); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// CreateFileCmd creates a file based on the name provided.
func CreateFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateFile(name); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// ZipDirectoryCmd zips a directory based on the name provided.
func ZipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return ErrorMsg(err.Error())
		}

		dirToZip := filepath.Join(currentDir, name)
		if err := dirfs.Zip(dirToZip); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// UnzipDirectoryCmd unzips a directory based on the name provided.
func UnzipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return ErrorMsg(err.Error())
		}

		dirToUnzip := filepath.Join(currentDir, name)
		if err := dirfs.Unzip(dirToUnzip); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// CopyFileCmd copies a file based on the name provided.
func CopyFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyFile(name); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// CopyDirectoryCmd copies a directory based on the name provided.
func CopyDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyDirectory(name); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}

// GetDirectoryItemSizeCmd calculates the size of a directory or file.
func GetDirectoryItemSizeCmd(name string, i int) tea.Cmd {
	return func() tea.Msg {
		size, err := dirfs.GetDirectoryItemSize(name)
		if err != nil {
			return DirectoryItemSizeMsg{size: "N/A", index: i}
		}

		sizeString := strfmt.ConvertBytesToSizeString(size)

		return DirectoryItemSizeMsg{
			size:  sizeString,
			index: i,
		}
	}
}

// HandleErrorCmd returns an error message to the UI.
func HandleErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err.Error())
	}
}

// CopyToClipboardCmd copies the provided string to the clipboard.
func CopyToClipboardCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return ErrorMsg(err.Error())
		}

		filePath := filepath.Join(workingDir, name)
		clipboard.Write(clipboard.FmtText, []byte(filePath))

		return CopyToClipboardMsg(fmt.Sprintf("%s %s %s", "Successfully copied", filePath, "to clipboard"))
	}
}

// GetDirectoryListingByTypeCmd returns only directories in the current directory.
func GetDirectoryListingByTypeCmd(listType string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return ErrorMsg(err.Error())
		}

		directories, err := dirfs.GetDirectoryListingByType(workingDir, listType, showHidden)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		return UpdateDirectoryListingMsg(directories)
	}
}

// FindFilesByNameCmd finds files based on name.
func FindFilesByNameCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return ErrorMsg(err.Error())
		}

		paths, entries, err := dirfs.FindFilesByName(name, workingDir)
		if err != nil {
			return ErrorMsg(err.Error())
		}

		return FindFilesByNameMsg{
			paths:   paths,
			entries: entries,
		}
	}
}

// WriteSelectionPathCmd writes content to the file specified.
func WriteSelectionPathCmd(selectionPath, filePath string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.WriteToFile(selectionPath, filePath); err != nil {
			return ErrorMsg(err.Error())
		}

		return nil
	}
}
