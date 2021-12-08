package filetree

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/commands"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/strfmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
	"github.com/spf13/viper"
)

// Bubble is a struct to represent the properties of a filetree.
type Bubble struct {
	Viewport            viewport.Model
	AppConfig           config.Config
	Style               lipgloss.Style
	UnselectedItemColor lipgloss.AdaptiveColor
	SelectedItemColor   lipgloss.AdaptiveColor
	ActiveBorderColor   lipgloss.AdaptiveColor
	InactiveBorderColor lipgloss.AdaptiveColor
	Files               []fs.DirEntry
	FilePaths           []string
	Cursor              int
	ShowIcons           bool
	ShowHidden          bool
	Borderless          bool
	IsActive            bool
	AlternateBorder     bool
	ShowLoading         bool
	ShowDirectoriesOnly bool
	ShowFilesOnly       bool
	CreateFileMode      bool
}

// Init intializes the filetree.
func (m Bubble) Init() tea.Cmd {
	startDir := viper.GetString("start-dir")

	switch {
	case startDir != "":
		_, err := os.Stat(startDir)
		if err != nil {
			return nil
		}

		if strings.HasPrefix(startDir, "/") {
			return commands.UpdateDirectoryListingCmd(startDir, m.ShowHidden)
		} else {
			path, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			filePath := filepath.Join(path, startDir)

			return commands.UpdateDirectoryListingCmd(filePath, m.ShowHidden)
		}
	case m.AppConfig.Settings.StartDir == dirfs.HomeDirectory:
		homeDir, err := dirfs.GetHomeDirectory()
		if err != nil {
			log.Fatal(err)
		}

		return commands.UpdateDirectoryListingCmd(homeDir, m.ShowHidden)
	default:
		return commands.UpdateDirectoryListingCmd(m.AppConfig.Settings.StartDir, m.ShowHidden)
	}
}

// NewBubble creates a new instance of a filetree.
func NewBubble(
	showIcons, borderless, isActive, showHidden bool,
	selectedItemColor, unselectedItemColor, activeBorderColor, inactiveBorderColor lipgloss.AdaptiveColor,
	appConfig config.Config,
) Bubble {
	border := lipgloss.NormalBorder()
	padding := 1

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	style := lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border)

	return Bubble{
		Cursor:              0,
		ShowIcons:           showIcons,
		Borderless:          borderless,
		IsActive:            isActive,
		ShowHidden:          showHidden,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		AppConfig:           appConfig,
		Style:               style,
	}
}

// scrollFiletree moves handles wrapping of the filetree and
// scrolling of the viewport.
func (m *Bubble) scrollFiletree() {
	top := m.Viewport.YOffset
	bottom := m.Viewport.Height + m.Viewport.YOffset - 1

	if m.Cursor < top {
		m.Viewport.LineUp(1)
	} else if m.Cursor > bottom {
		m.Viewport.LineDown(1)
	}

	if m.Cursor > m.GetTotalFiles()-1 {
		m.GotoTop()
		m.Viewport.GotoTop()
	} else if m.Cursor < top {
		m.GotoBottom()
		m.Viewport.GotoBottom()
	}
}

// SetContent sets the files currently displayed in the tree.
func (m *Bubble) SetContent(files []fs.DirEntry) {
	var directoryItem string
	curFiles := ""

	m.Files = files

	for i, file := range files {
		var fileSizeColor lipgloss.AdaptiveColor

		if m.Cursor == i {
			fileSizeColor = m.SelectedItemColor
		} else {
			fileSizeColor = m.UnselectedItemColor
		}

		fileInfo, _ := file.Info()

		fileSize := lipgloss.NewStyle().
			Foreground(fileSizeColor).
			Render(strfmt.ConvertBytesToSizeString(fileInfo.Size()))

		icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)

		switch {
		case m.ShowIcons && m.Cursor == i:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(m.SelectedItemColor).
				Render(fileInfo.Name()))
		case m.ShowIcons && m.Cursor != i:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(m.UnselectedItemColor).
				Render(fileInfo.Name()))
		case !m.ShowIcons && m.Cursor == i:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(m.SelectedItemColor).
				Render(fileInfo.Name())
		default:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(m.UnselectedItemColor).
				Render(fileInfo.Name())
		}

		dirItem := lipgloss.NewStyle().Width(m.Viewport.Width - lipgloss.Width(fileSize) - m.Style.GetHorizontalPadding()).Render(
			truncate.StringWithTail(
				directoryItem, uint(m.Viewport.Width-lipgloss.Width(fileSize)), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	m.Viewport.SetContent(curFiles)
}

// SetFilePaths sets an array of file paths.
func (m *Bubble) SetFilePaths(filePaths []string) {
	m.FilePaths = filePaths
}

// GetFilePaths returns an array of file paths.
func (m Bubble) GetFilePaths() []string {
	return m.FilePaths
}

// SetSize updates the size of the filetree, useful when resizing the terminal.
func (m *Bubble) SetSize(width, height int) {
	m.Viewport.Width = (width / 2) - m.Style.GetHorizontalBorderSize()
	m.Viewport.Height = height - m.Style.GetVerticalBorderSize() - statusbar.StatusbarHeight
}

// GotoTop goes to the top of the tree.
func (m *Bubble) GotoTop() {
	m.Cursor = 0
}

// GotoBottom goes to the bottom of the tree.
func (m *Bubble) GotoBottom() {
	m.Cursor = len(m.Files) - 1
}

// GetSelectedFile returns the currently selected file in the tree.
func (m Bubble) GetSelectedFile() (os.FileInfo, error) {
	if len(m.Files) > 0 {
		fileInfo, err := m.Files[m.Cursor].Info()
		if err != nil {
			return nil, err
		}

		return fileInfo, nil
	}

	return nil, nil
}

// GetCursor gets the position of the cursor in the tree.
func (m Bubble) GetCursor() int {
	return m.Cursor
}

// GoDown goes down the tree by one.
func (m *Bubble) GoDown() {
	m.Cursor++
}

// GoUp goes up the tree by one.
func (m *Bubble) GoUp() {
	m.Cursor--
}

// GetTotalFiles returns the total number of files in the tree.
func (m Bubble) GetTotalFiles() int {
	return len(m.Files)
}

// GetIsActive returns the active state of the filetree.
func (m Bubble) GetIsActive() bool {
	return m.IsActive
}

// SetIsActive sets the active state of the filetree.
func (m *Bubble) SetIsActive(isActive bool) {
	m.IsActive = isActive
}

// Update updates the statusbar.
func (m Bubble) Update(msg tea.Msg) (Bubble, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.UpdateDirectoryListingMsg:
		m.GotoTop()
		m.SetFilePaths(nil)
		m.Viewport.GotoTop()
		m.SetContent(msg)
		m.CreateFileMode = false

		return m, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths)
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.SetContent(m.Files)
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if m.IsActive && !m.CreateFileMode {
				m.GoUp()
				m.scrollFiletree()
				m.SetContent(m.Files)
			}
		case tea.MouseWheelDown:
			if m.IsActive && !m.CreateFileMode {
				m.GoDown()
				m.scrollFiletree()
				m.SetContent(m.Files)
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.IsActive && !m.CreateFileMode {
				m.GoUp()
				m.scrollFiletree()
				m.SetContent(m.Files)
				cmds = append(cmds, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths))
			}
		case "down", "j":
			if m.IsActive && !m.CreateFileMode {
				m.GoDown()
				m.scrollFiletree()
				m.SetContent(m.Files)
				cmds = append(cmds, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths))
			}
		case "right", "l":
			if m.IsActive && !m.CreateFileMode && len(m.Files) > 0 {
				selectedFile, err := m.GetSelectedFile()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				switch {
				case selectedFile.IsDir():
					currentDir, err := dirfs.GetWorkingDirectory()
					if err != nil {
						cmds = append(cmds, commands.HandleErrorCmd(err))
					}

					directoryToOpen := filepath.Join(currentDir, selectedFile.Name())

					if len(m.GetFilePaths()) > 0 {
						directoryToOpen = m.GetFilePaths()[m.GetCursor()]
					}

					cmds = append(cmds, commands.UpdateDirectoryListingCmd(directoryToOpen, m.ShowHidden))
				case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
					symlinkFile, err := os.Readlink(selectedFile.Name())
					if err != nil {
						cmds = append(cmds, commands.HandleErrorCmd(err))
					}

					fileInfo, err := os.Stat(symlinkFile)
					if err != nil {
						cmds = append(cmds, commands.HandleErrorCmd(err))
					}

					if fileInfo.IsDir() {
						currentDir, err := dirfs.GetWorkingDirectory()
						if err != nil {
							cmds = append(cmds, commands.HandleErrorCmd(err))
						}

						cmds = append(
							cmds,
							commands.UpdateDirectoryListingCmd(
								filepath.Join(currentDir, fileInfo.Name()),
								m.ShowHidden),
						)
					}

					cmds = append(cmds, commands.ReadFileContentCmd(
						fileInfo.Name(),
						m.AppConfig.Settings.SyntaxTheme,
						m.Viewport.Width,
						m.AppConfig.Settings.PrettyMarkdown,
					))
				default:
					fileToRead := selectedFile.Name()

					if len(m.GetFilePaths()) > 0 {
						fileToRead = m.GetFilePaths()[m.GetCursor()]
					}

					cmds = append(cmds, commands.ReadFileContentCmd(
						fileToRead,
						m.AppConfig.Settings.SyntaxTheme,
						m.Viewport.Width,
						m.AppConfig.Settings.PrettyMarkdown,
					))
				}

				cmds = append(cmds, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths))
			}
		case "left", "h":
			if m.IsActive && !m.CreateFileMode {
				m.ShowFilesOnly = false
				m.ShowDirectoriesOnly = false
				workingDirectory, err := dirfs.GetWorkingDirectory()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				cmds = append(
					cmds,
					commands.UpdateDirectoryListingCmd(
						filepath.Join(workingDirectory, dirfs.PreviousDirectory),
						m.ShowHidden),
				)

				cmds = append(cmds, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths))
			}
		case "ctrl+g":
			if m.IsActive && !m.CreateFileMode {
				m.GotoTop()
				m.Viewport.GotoTop()
				m.SetContent(m.Files)
				cmds = append(cmds, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths))
			}
		case "G":
			if m.IsActive && !m.CreateFileMode {
				m.GotoBottom()
				m.Viewport.GotoBottom()
				m.SetContent(m.Files)
				cmds = append(cmds, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths))
			}
		case "~":
			if m.IsActive && !m.CreateFileMode {
				homeDir, err := dirfs.GetHomeDirectory()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				cmds = append(cmds, commands.UpdateDirectoryListingCmd(homeDir, m.ShowHidden))
			}
		case "/":
			if m.IsActive && !m.CreateFileMode {
				cmds = append(cmds, commands.UpdateDirectoryListingCmd(dirfs.RootDirectory, m.ShowHidden))
			}
		case ".":
			if m.IsActive && !m.CreateFileMode {
				m.ShowHidden = !m.ShowHidden

				switch {
				case m.ShowDirectoriesOnly:
					cmds = append(cmds, commands.GetDirectoryListingByTypeCmd(dirfs.DirectoriesListingType, m.ShowHidden))
				case m.ShowFilesOnly:
					cmds = append(cmds, commands.GetDirectoryListingByTypeCmd(dirfs.FilesListingType, m.ShowHidden))
				default:
					cmds = append(cmds, commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden))
				}
			}
		case "s":
			if m.IsActive && !m.CreateFileMode {
				m.ShowFilesOnly = !m.ShowFilesOnly
				m.ShowDirectoriesOnly = false

				if m.ShowFilesOnly {
					cmds = append(cmds, commands.GetDirectoryListingByTypeCmd(dirfs.FilesListingType, m.ShowHidden))
				}

				cmds = append(cmds, commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden))
			}
		case "S":
			if m.IsActive && !m.CreateFileMode {
				m.ShowDirectoriesOnly = !m.ShowDirectoriesOnly
				m.ShowFilesOnly = false

				if m.ShowDirectoriesOnly {
					cmds = append(cmds, commands.GetDirectoryListingByTypeCmd(dirfs.DirectoriesListingType, m.ShowHidden))
				}

				cmds = append(cmds, commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden))
			}
		case "C":
			if m.IsActive && !m.CreateFileMode {
				selectedFile, err := m.GetSelectedFile()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				if selectedFile.IsDir() {
					cmds = append(cmds, tea.Sequentially(
						commands.CopyDirectoryCmd(selectedFile.Name()),
						commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden),
					))
				} else {
					cmds = append(cmds, tea.Sequentially(
						commands.CopyFileCmd(selectedFile.Name()),
						commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden),
					))
				}
			}
		case "Z":
			if m.IsActive && !m.CreateFileMode {
				selectedFile, err := m.GetSelectedFile()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				cmds = append(cmds, tea.Sequentially(
					commands.ZipDirectoryCmd(selectedFile.Name()),
					commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden),
				))
			}
		case "U":
			if m.IsActive && !m.CreateFileMode {
				selectedFile, err := m.GetSelectedFile()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				cmds = append(cmds, tea.Sequentially(
					commands.UnzipDirectoryCmd(selectedFile.Name()),
					commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden),
				))
			}
		case "Y":
			if m.IsActive && !m.CreateFileMode {
				selectedFile, err := m.GetSelectedFile()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				cmds = append(cmds, commands.CopyToClipboardCmd(selectedFile.Name()))

			}
		case "E":
			if m.IsActive && !m.CreateFileMode {
				selectedFile, err := m.GetSelectedFile()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				selectionPath := viper.GetString("selection-path")

				if selectionPath == "" && !selectedFile.IsDir() {
					editorPath := os.Getenv("EDITOR")
					if editorPath == "" {
						cmds = append(cmds, commands.HandleErrorCmd(errors.New("$EDITOR not set")))
					}

					editorCmd := exec.Command(editorPath, selectedFile.Name())
					editorCmd.Stdin = os.Stdin
					editorCmd.Stdout = os.Stdout
					editorCmd.Stderr = os.Stderr

					err := editorCmd.Start()
					if err != nil {
						cmds = append(cmds, commands.HandleErrorCmd(err))
					}

					err = editorCmd.Wait()
					if err != nil {
						cmds = append(cmds, commands.HandleErrorCmd(err))
					}

					cmds = append(cmds, commands.UpdateDirectoryListingCmd(dirfs.CurrentDirectory, m.ShowHidden))
				} else {
					cmds = append(cmds, tea.Sequentially(commands.WriteSelectionPathCmd(selectionPath, selectedFile.Name()), tea.Quit))
				}
			}
		case "p":
			if m.IsActive && !m.CreateFileMode {
				selectedFile, err := m.GetSelectedFile()
				if err != nil {
					cmds = append(cmds, commands.HandleErrorCmd(err))
				}

				switch {
				case selectedFile.IsDir():
					cmds = append(cmds, commands.PreviewDirectoryListingCmd(selectedFile.Name(), m.ShowHidden))
				case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
					symlinkFile, err := os.Readlink(selectedFile.Name())
					if err != nil {
						cmds = append(cmds, commands.HandleErrorCmd(err))
					}

					fileInfo, err := os.Stat(symlinkFile)
					if err != nil {
						cmds = append(cmds, commands.HandleErrorCmd(err))
					}

					if fileInfo.IsDir() {
						cmds = append(cmds, commands.PreviewDirectoryListingCmd(fileInfo.Name(), m.ShowHidden))
					}
				default:
					return m, nil
				}
			}
		case "n":
			if m.IsActive && !m.CreateFileMode {
				m.CreateFileMode = true
				return m, commands.UpdateStatusbarCmd(m.Files, m.Cursor, m.FilePaths)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the current tree.
func (m Bubble) View() string {
	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()
	content := m.Viewport.View()
	alternateBorder := lipgloss.Border{
		Top:         "-",
		Bottom:      "-",
		Left:        "|",
		Right:       "|",
		TopLeft:     "*",
		TopRight:    "*",
		BottomLeft:  "*",
		BottomRight: "*",
	}

	if len(m.Files) == 0 {
		return m.Style.Copy().
			BorderForeground(borderColor).
			Border(border).
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render("Directory is empty")
	}

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	if m.AlternateBorder {
		border = alternateBorder
	}

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return m.Style.Copy().
		BorderForeground(borderColor).
		Border(border).
		Width(m.Viewport.Width).
		Height(m.Viewport.Height).
		Render(content)
}
