package ui

import (
	"context"
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	// "image/jpeg" is needed for the image.Decode function.
	_ "image/jpeg"
	// "image/png" is needed for the image.Decode function.
	_ "image/png"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/internal/colorimage"
	"github.com/knipferrc/fm/internal/markdown"
	"github.com/knipferrc/fm/internal/sourcecode"
	"github.com/knipferrc/fm/strfmt"

	tea "github.com/charmbracelet/bubbletea"
)

type updateDirectoryListingMsg []fs.DirEntry
type previewDirectoryListingMsg []fs.DirEntry
type leftKeyMsg struct {
	files             []fs.DirEntry
	previousDirectory string
}
type handleDownKeyMsg struct{}
type handleUpKeyMsg struct{}
type handleRightKeyMsg struct{}
type moveDirectoryItemMsg []fs.DirEntry
type errorMsg string
type convertImageToStringMsg string
type directoryItemSizeMsg string
type readFileContentMsg struct {
	rawContent  string
	markdown    string
	code        string
	imageString string
	image       image.Image
}

// getDirectoryListing returns a list of files and directories in the directory provided.
func getDirectoryListing(name string, showHidden bool) ([]fs.DirEntry, error) {
	files, err := dirfs.GetDirectoryListing(name, showHidden)
	if err != nil {
		return nil, err
	}

	err = os.Chdir(name)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// readFileContent reads the content of the file provided.
func readFileContentCmd(file fs.FileInfo, width, height int, prettyMarkdown bool, syntaxTheme string) (readFileContentMsg, error) {
	content, err := dirfs.ReadFileContent(file.Name())
	if err != nil {
		return readFileContentMsg{}, err
	}

	switch {
	case filepath.Ext(file.Name()) == ".md" && prettyMarkdown:
		markdownContent, err := markdown.RenderMarkdown(width, content)
		if err != nil {
			return readFileContentMsg{}, err
		}

		return readFileContentMsg{
			rawContent:  content,
			markdown:    markdownContent,
			code:        "",
			imageString: "",
			image:       nil,
		}, nil
	case filepath.Ext(file.Name()) == ".png" || filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == ".jpeg":
		imageContent, err := os.Open(file.Name())
		if err != nil {
			return readFileContentMsg{}, err
		}

		img, _, err := image.Decode(imageContent)
		if err != nil {
			return readFileContentMsg{}, err
		}

		imageString, err := colorimage.ImageToString(width, height, img)
		if err != nil {
			return readFileContentMsg{}, err
		}

		return readFileContentMsg{
			rawContent:  content,
			code:        "",
			markdown:    "",
			imageString: imageString,
			image:       img,
		}, nil
	default:
		code, err := sourcecode.Highlight(content, filepath.Ext(file.Name()), syntaxTheme)
		if err != nil {
			return readFileContentMsg{}, err
		}

		return readFileContentMsg{
			rawContent:  content,
			code:        code,
			markdown:    "",
			imageString: "",
			image:       nil,
		}, nil
	}
}

// updateDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (m Model) updateDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := getDirectoryListing(name, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(files)
	}
}

// handleLeftKeyCmd handles when the left key is pressed.
func (m Model) handleLeftKeyCmd() tea.Cmd {
	return func() tea.Msg {
		workingDirectory, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		directoryName := fmt.Sprintf("%s/%s", workingDirectory, dirfs.PreviousDirectory)

		files, err := getDirectoryListing(directoryName, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return leftKeyMsg{
			files:             files,
			previousDirectory: workingDirectory,
		}
	}
}

// handleDownKeyCmd handles when the down key is pressed.
func (m Model) handleDownKeyCmd() tea.Cmd {
	return func() tea.Msg {
		return handleDownKeyMsg{}
	}
}

// handleUpKeyCmd handles when the up key is pressed.
func (m Model) handleUpKeyCmd() tea.Cmd {
	return func() tea.Msg {
		return handleUpKeyMsg{}
	}
}

// handleRightKeyCmd handles when the right key is pressed.
func (m Model) handleRightKeyCmd() tea.Cmd {
	return func() tea.Msg {
		switch {
		case m.dirTree.GetSelectedFile().IsDir() && !m.statusBar.CommandBarFocused():
			currentDir, err := dirfs.GetWorkingDirectory()
			if err != nil {
				return errorMsg(err.Error())
			}

			directoryName := fmt.Sprintf("%s/%s", currentDir, m.dirTree.GetSelectedFile().Name())
			files, err := getDirectoryListing(directoryName, m.dirTree.ShowHidden)
			if err != nil {
				return errorMsg(err.Error())
			}

			return updateDirectoryListingMsg(files)
		case m.dirTree.GetSelectedFile().Mode()&os.ModeSymlink == os.ModeSymlink:
			symlinkFile, err := os.Readlink(m.dirTree.GetSelectedFile().Name())
			if err != nil {
				return errorMsg(err.Error())
			}

			fileInfo, err := os.Stat(symlinkFile)
			if err != nil {
				return errorMsg(err.Error())
			}

			if fileInfo.IsDir() {
				files, err := getDirectoryListing(symlinkFile, m.dirTree.ShowHidden)
				if err != nil {
					return errorMsg(err.Error())
				}

				return updateDirectoryListingMsg(files)
			}

			fileContent, err := readFileContentCmd(fileInfo, m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(), m.secondaryPane.GetHeight(), m.appConfig.Settings.PrettyMarkdown, m.appConfig.Settings.SyntaxTheme)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  fileContent.rawContent,
				code:        fileContent.code,
				markdown:    fileContent.markdown,
				imageString: fileContent.imageString,
				image:       fileContent.image,
			}
		default:
			fileContent, err := readFileContentCmd(m.dirTree.GetSelectedFile(), m.secondaryPane.GetWidth()-m.secondaryPane.Style.GetHorizontalFrameSize(), m.secondaryPane.GetHeight(), m.appConfig.Settings.PrettyMarkdown, m.appConfig.Settings.SyntaxTheme)
			if err != nil {
				return errorMsg(err.Error())
			}

			return readFileContentMsg{
				rawContent:  fileContent.rawContent,
				code:        fileContent.code,
				markdown:    fileContent.markdown,
				imageString: fileContent.imageString,
				image:       fileContent.image,
			}
		}
	}
}

// previewDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (m Model) previewDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		fileName := fmt.Sprintf("%s/%s", currentDir, name)
		files, err := dirfs.GetDirectoryListing(fileName, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		return previewDirectoryListingMsg(files)
	}
}

// renameDirectoryItemCmd renames a file or directory based on the name and value provided.
func (m Model) renameDirectoryItemCmd(name, value string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.RenameDirectoryItem(name, value); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// moveDirectoryItemCmd moves a file to the current working directory.
func (m Model) moveDirectoryItemCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		// Get the directory from which the move was intiated from
		// and give it the same file name.
		src := fmt.Sprintf("%s/%s", m.initialMoveDirectory, name)

		// Destination is the current working directory with
		// the same file name that it had.
		dst := fmt.Sprintf("%s/%s", workingDir, name)

		if err = dirfs.MoveDirectoryItem(src, dst); err != nil {
			return errorMsg(err.Error())
		}

		files, err := dirfs.GetDirectoryListing(m.initialMoveDirectory, m.dirTree.ShowHidden)
		if err != nil {
			return errorMsg(err.Error())
		}

		err = os.Chdir(m.initialMoveDirectory)
		if err != nil {
			return errorMsg(err.Error())
		}

		return moveDirectoryItemMsg(files)
	}
}

// deleteDirectoryCmd deletes a directory based on the name provided.
func (m Model) deleteDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// deleteFileCmd deletes a file based on the name provided.
func (m Model) deleteFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.DeleteFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// redrawImageCmd redraws the image based on the width and height provided.
func (m Model) redrawImageCmd(width, height int) tea.Cmd {
	return func() tea.Msg {
		imageString, err := colorimage.ImageToString(width, height, m.colorimage.Image)
		if err != nil {
			return errorMsg(err.Error())
		}

		return convertImageToStringMsg(imageString)
	}
}

// createDirectoryCmd creates a directory based on the name provided.
func (m Model) createDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// createFileCmd creates a file based on the name provided.
func (m Model) createFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CreateFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// zipDirectoryCmd zips a directory based on the name provided.
func (m Model) zipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		dirToZip := fmt.Sprintf("%s/%s", currentDir, name)
		if err := dirfs.Zip(dirToZip); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// unzipDirectoryCmd unzips a directory based on the name provided.
func (m Model) unzipDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dirfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		dirToUnzip := fmt.Sprintf("%s/%s", currentDir, name)
		if err := dirfs.Unzip(dirToUnzip); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyFileCmd copies a file based on the name provided.
func (m Model) copyFileCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyFile(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// copyDirectoryCmd copies a directory based on the name provided.
func (m Model) copyDirectoryCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if err := dirfs.CopyDirectory(name); err != nil {
			return errorMsg(err.Error())
		}

		return nil
	}
}

// getDirectoryItemSizeCmd calculates the size of a directory or file.
func (m Model) getDirectoryItemSizeCmd(name string) tea.Cmd {
	if m.directoryItemSizeCtx != nil && m.directoryItemSizeCtx.cancel != nil {
		m.directoryItemSizeCtx.cancel()
	}

	ctx, cancel := context.WithTimeout(m.directoryItemSizeCtx.ctx, 300*time.Millisecond)
	m.directoryItemSizeCtx.cancel = cancel

	return func() tea.Msg {
		defer cancel()
		<-ctx.Done()
		if ctx.Err() == context.DeadlineExceeded {
			size, err := dirfs.GetDirectoryItemSize(name)
			if err != nil {
				return directoryItemSizeMsg("N/A")
			}

			sizeString := strfmt.ConvertBytesToSizeString(size)

			return directoryItemSizeMsg(sizeString)
		}

		return nil
	}
}

// handleErrorCmd returns an error message to the UI.
func (m Model) handleErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err.Error())
	}
}
