package tui

import (
	"image"
	"io/fs"

	"github.com/knipferrc/fm/help"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/constants"
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
	help                   help.Bubble
	treeFiles              []fs.DirEntry
	treePreviewFiles       []fs.DirEntry
	treeItemToMove         fs.DirEntry
	previousKey            tea.KeyMsg
	keyMap                 KeyMap
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
	showLogs               bool
	foundFilesPaths        []string
	fileSizes              []string
	logs                   []string
	moveInitiatedDirectory string
	secondaryBoxContent    string
	errorMsg               string
}

// New creates an instance of the entire application.
func New() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.Theme.AppTheme)

	primaryBoxBorder := lipgloss.NormalBorder()
	secondaryBoxBorder := lipgloss.NormalBorder()
	primaryBoxBorderColor := theme.ActiveBoxBorderColor
	secondaryBoxBorderColor := theme.InactiveBoxBorderColor

	if cfg.Settings.Borderless {
		primaryBoxBorder = lipgloss.HiddenBorder()
		secondaryBoxBorder = lipgloss.HiddenBorder()
	}

	if cfg.Settings.SimpleMode {
		primaryBoxBorder = lipgloss.HiddenBorder()
		secondaryBoxBorder = lipgloss.HiddenBorder()
	}

	pvp := viewport.New(0, 0)
	pvp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder).
		BorderForeground(primaryBoxBorderColor)

	svp := viewport.New(0, 0)
	svp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(secondaryBoxBorder).
		BorderForeground(secondaryBoxBorderColor)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.SpinnerColor)

	t := textinput.New()
	t.Prompt = "❯ "
	t.CharLimit = 250
	t.PlaceholderStyle = lipgloss.NewStyle().
		Background(theme.StatusBarBarBackgroundColor).
		Foreground(theme.StatusBarBarForegroundColor)

	h := help.New(
		theme.DefaultTextColor,
		"Welcome to FM!",
		[]help.HelpEntry{
			{Key: "ctrl+c", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
			{Key: "h/left", Description: "Go back a directory"},
			{Key: "l/right", Description: "Read file or enter directory"},
			{Key: "p", Description: "Preview directory"},
			{Key: "gg", Description: "Go to top of filetree or box"},
			{Key: "G", Description: "Go to bottom of filetree or box"},
			{Key: "~", Description: "Go to home directory"},
			{Key: "/", Description: "Go to root directory"},
			{Key: ".", Description: "Toggle hidden files"},
			{Key: "S", Description: "Only show directories"},
			{Key: "s", Description: "Only show files"},
			{Key: "y", Description: "Copy file path to clipboard"},
			{Key: "Z", Description: "Zip currently selected tree item"},
			{Key: "U", Description: "Unzip currently selected tree item"},
			{Key: "n", Description: "Create new file"},
			{Key: "N", Description: "Create new directory"},
			{Key: "ctrl+d", Description: "Delete currently selected tree item"},
			{Key: "M", Description: "Move currently selected tree item"},
			{Key: "enter", Description: "Process command"},
			{Key: "E", Description: "Edit currently selected tree item"},
			{Key: "C", Description: "Copy currently selected tree item"},
			{Key: "esc", Description: "Reset FM to initial state"},
			{Key: "O", Description: "Show logs if debugging enabled"},
			{Key: "tab", Description: "Toggle between boxes"},
		})

	return Bubble{
		appConfig:         cfg,
		theme:             theme,
		showHiddenFiles:   true,
		spinner:           s,
		textinput:         t,
		showHelp:          true,
		primaryViewport:   pvp,
		secondaryViewport: svp,
		keyMap:            DefaultKeyMap(),
		help:              h,
	}
}
