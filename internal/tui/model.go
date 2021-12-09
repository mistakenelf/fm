package tui

import (
	"image"
	"io/fs"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/theme"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Bubble represents the state of the UI.
type Bubble struct {
	appConfig              config.Config
	theme                  theme.Theme
	currentImage           image.Image
	spinner                spinner.Model
	textinput              textinput.Model
	primaryViewport        viewport.Model
	secondaryViewport      viewport.Model
	treeFiles              []fs.DirEntry
	treePreviewFiles       []fs.DirEntry
	treeItemToMove         fs.DirEntry
	previousKey            tea.KeyMsg
	width                  int
	height                 int
	activeBox              int
	treeCursor             int
	showHiddenFiles        bool
	ready                  bool
	showCommandInput       bool
	showFilesOnly          bool
	showDirectoriesOnly    bool
	showFileTreePreview    bool
	createFileMode         bool
	createDirectoryMode    bool
	renameMode             bool
	moveMode               bool
	findMode               bool
	deleteMode             bool
	showBoxSpinner         bool
	showHelp               bool
	foundFilesPaths        []string
	moveInitiatedDirectory string
	primaryContent         string
	secondaryContent       string
	errorMsg               string
}

// NewBubble create an instance of the entire application.
func NewBubble() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.Settings.Theme)

	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.SpinnerColor)

	t := textinput.NewModel()
	t.Prompt = "❯ "
	t.CharLimit = 250
	t.PlaceholderStyle = lipgloss.NewStyle().
		Background(theme.StatusBarBarBackgroundColor).
		Foreground(theme.StatusBarBarForegroundColor)

	return Bubble{
		appConfig:       cfg,
		theme:           theme,
		showHiddenFiles: true,
		spinner:         s,
		textinput:       t,
		showHelp:        true,
	}
}
