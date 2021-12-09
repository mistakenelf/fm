package tui

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/knipferrc/fm/dirfs"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/strfmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

func (b Bubble) statusBarView() string {
	var logo string
	var status string

	width := lipgloss.Width
	selectedFileName := "N/A"
	fileCount := "0/0"

	if len(b.treeFiles) > 0 && b.treeFiles[b.treeCursor] != nil {
		selectedFile, err := b.treeFiles[b.treeCursor].Info()
		if err != nil {
			return "error"
		}
		fileCount = fmt.Sprintf("%d/%d", b.treeCursor+1, len(b.treeFiles))
		selectedFileName = selectedFile.Name()

		currentPath, err := dirfs.GetWorkingDirectory()
		if err != nil {
			currentPath = dirfs.CurrentDirectory
		}

		status = fmt.Sprintf("%s %s %s",
			selectedFile.ModTime().Format("2006-01-02 15:04:05"),
			selectedFile.Mode().String(),
			currentPath,
		)
	}

	if b.showCommandInput {
		status = b.textinput.View()
	}

	if b.moveMode {
		status = fmt.Sprintf("%s %s", "Currently moving:", b.treeFiles[b.treeCursor].Name())
	}

	if b.appConfig.Settings.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	// Selected file styles
	selectedFileStyle := lipgloss.NewStyle().
		Foreground(b.theme.StatusBarSelectedFileForegroundColor).
		Background(b.theme.StatusBarSelectedFileBackgroundColor)

	if b.simpleMode {
		selectedFileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	selectedFileColumn := selectedFileStyle.
		Padding(0, 1).
		Height(StatusBarHeight).
		Render(truncate.StringWithTail(selectedFileName, 30, "..."))

	// File count styles
	fileCountStyle := lipgloss.NewStyle().
		Foreground(b.theme.StatusBarTotalFilesForegroundColor).
		Background(b.theme.StatusBarTotalFilesBackgroundColor)

	if b.simpleMode {
		fileCountStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	fileCountColumn := fileCountStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(StatusBarHeight).
		Render(fileCount)

	// Logo styles
	logoStyle := lipgloss.NewStyle().
		Foreground(b.theme.StatusBarLogoForegroundColor).
		Background(b.theme.StatusBarLogoBackgroundColor)

	if b.simpleMode {
		logoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	logoColumn := logoStyle.
		Padding(0, 1).
		Height(StatusBarHeight).
		Render(logo)

	// Status styles
	statusStyle := lipgloss.NewStyle().
		Foreground(b.theme.StatusBarBarForegroundColor).
		Background(b.theme.StatusBarBarBackgroundColor)

	if b.simpleMode {
		statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
	}

	statusColumn := statusStyle.
		Padding(0, 1).
		Height(StatusBarHeight).
		Width(b.width - width(selectedFileColumn) - width(fileCountColumn) - width(logoColumn)).
		Render(truncate.StringWithTail(
			status,
			uint(b.width-width(selectedFileColumn)-width(fileCountColumn)-width(logoColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileColumn,
		statusColumn,
		fileCountColumn,
		logoColumn,
	)
}

func (b Bubble) fileTreeView(files []fs.DirEntry) string {
	var directoryItem string
	curFiles := ""

	for i, file := range files {
		var fileSizeColor lipgloss.AdaptiveColor

		if b.treeCursor == i {
			fileSizeColor = b.theme.SelectedTreeItemColor
		} else {
			fileSizeColor = b.theme.UnselectedTreeItemColor
		}

		fileInfo, _ := file.Info()

		fileSize := lipgloss.NewStyle().
			Foreground(fileSizeColor).
			Render(strfmt.ConvertBytesToSizeString(fileInfo.Size()))

		icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)

		switch {
		case b.appConfig.Settings.ShowIcons && b.treeCursor == i:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(b.theme.SelectedTreeItemColor).
				Render(fileInfo.Name()))
		case b.appConfig.Settings.ShowIcons && b.treeCursor != i:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(b.theme.UnselectedTreeItemColor).
				Render(fileInfo.Name()))
		case !b.appConfig.Settings.ShowIcons && b.treeCursor == i:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(b.theme.SelectedTreeItemColor).
				Render(fileInfo.Name())
		default:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(b.theme.UnselectedTreeItemColor).
				Render(fileInfo.Name())
		}

		dirItem := lipgloss.NewStyle().Width(b.primaryViewport.Width - lipgloss.Width(fileSize) - box.GetHorizontalPadding()).Render(
			truncate.StringWithTail(
				directoryItem, uint(b.primaryViewport.Width-lipgloss.Width(fileSize)), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	return curFiles
}

func (b Bubble) fileTreePreviewView(files []fs.DirEntry) string {
	var directoryItem string
	curFiles := ""

	for _, file := range files {
		fileColor := b.theme.UnselectedTreeItemColor

		fileInfo, _ := file.Info()

		fileSize := lipgloss.NewStyle().
			Foreground(fileColor).
			Render(strfmt.ConvertBytesToSizeString(fileInfo.Size()))

		icon, color := icons.GetIcon(fileInfo.Name(), filepath.Ext(fileInfo.Name()), icons.GetIndicator(fileInfo.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)

		switch {
		case b.appConfig.Settings.ShowIcons:
			directoryItem = fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
				Bold(true).
				Foreground(fileColor).
				Render(fileInfo.Name()))
		case !b.appConfig.Settings.ShowIcons:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(fileColor).
				Render(fileInfo.Name())
		default:
			directoryItem = lipgloss.NewStyle().
				Bold(true).
				Foreground(fileColor).
				Render(fileInfo.Name())
		}

		dirItem := lipgloss.NewStyle().Width(b.secondaryViewport.Width - lipgloss.Width(fileSize) - box.GetHorizontalPadding()).Render(
			truncate.StringWithTail(
				directoryItem, uint(b.secondaryViewport.Width-lipgloss.Width(fileSize)), "...",
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	return curFiles
}

func (b Bubble) textContentView(content string) string {
	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width - box.GetHorizontalPadding()).
		Height(b.secondaryViewport.Height - box.GetVerticalPadding()).
		Render(content)
}

func (b Bubble) errorView(msg string) string {
	return lipgloss.NewStyle().
		Foreground(b.theme.ErrorColor).
		Width(b.secondaryViewport.Width - box.GetHorizontalPadding()).
		Height(b.secondaryViewport.Height - box.GetVerticalPadding()).
		Render(msg)
}

// View returns a string representation of the entire application UI.
func (b Bubble) View() string {
	if !b.ready {
		return fmt.Sprintf("%s %s", b.spinner.View(), "loading...")
	}

	primaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorder := lipgloss.NormalBorder()
	primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
	secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor

	if b.activeBox == 0 {
		primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
	}

	if b.activeBox == 1 {
		secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor
	}

	if b.appConfig.Settings.Borderless {
		primaryBoxBorder = lipgloss.HiddenBorder()
		secondaryBoxBorder = lipgloss.HiddenBorder()
	}

	if b.moveMode {
		primaryBoxBorder = lipgloss.Border{
			Top:         "-",
			Bottom:      "-",
			Left:        "|",
			Right:       "|",
			TopLeft:     "*",
			TopRight:    "*",
			BottomLeft:  "*",
			BottomRight: "*",
		}
	}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			box.Copy().
				Border(primaryBoxBorder).
				BorderForeground(primaryBoxBorderColor).
				Width(b.primaryViewport.Width).
				Height(b.primaryViewport.Height).
				Render(b.primaryViewport.View()),
			box.Copy().
				Border(secondaryBoxBorder).
				BorderForeground(secondaryBoxBorderColor).
				Width(b.secondaryViewport.Width).
				Height(b.secondaryViewport.Height).
				Render(b.secondaryViewport.View()),
		),
		b.statusBarView(),
	)

	return view
}
