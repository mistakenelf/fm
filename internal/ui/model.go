package ui

import (
	"context"
	"os"

	"github.com/knipferrc/fm/internal/colorimage"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/dirtree"
	"github.com/knipferrc/fm/internal/markdown"
	"github.com/knipferrc/fm/internal/pane"
	"github.com/knipferrc/fm/internal/sourcecode"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/internal/theme"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type directoryItemSizeCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// Model represents the state of the UI.
type Model struct {
	keys                  keyMap
	help                  help.Model
	primaryPane           pane.Model
	secondaryPane         pane.Model
	loader                spinner.Model
	dirTree               dirtree.Model
	dirTreePreview        dirtree.Model
	statusBar             statusbar.Model
	colorimage            colorimage.Model
	markdown              markdown.Model
	sourcecode            sourcecode.Model
	previousKey           tea.KeyMsg
	itemToMove            os.FileInfo
	appConfig             config.Config
	directoryItemSizeCtx  *directoryItemSizeCtx
	theme                 theme.Theme
	previousDirectory     string
	initialMoveDirectory  string
	showCommandBar        bool
	inMoveMode            bool
	inCreateFileMode      bool
	inCreateDirectoryMode bool
	inRenameMode          bool
	ready                 bool
}

// NewModel create an instance of the entire application model.
func NewModel() Model {
	cfg := config.GetConfig()
	keys := getDefaultKeyMap()
	theme := theme.GetTheme(cfg.Settings.Theme)

	// Create a new spinner with some styling based on the config.
	l := spinner.NewModel()
	l.Spinner = spinner.Dot
	l.Style = lipgloss.NewStyle().Foreground(theme.SpinnerColor)

	// Create a new help view.
	h := help.NewModel()
	h.Styles.FullKey.Foreground(theme.DefaultTextColor)
	h.Styles.FullDesc.Foreground(theme.DefaultTextColor)
	h.ShowAll = true

	// Create a new dirtree.
	dirTree := dirtree.NewModel(
		cfg.Settings.ShowIcons,
		theme.SelectedTreeItemColor,
		theme.UnselectedTreeItemColor,
	)

	// Create a new dirtree for previews.
	dirTreePreview := dirtree.NewModel(
		cfg.Settings.ShowIcons,
		theme.UnselectedTreeItemColor,
		theme.UnselectedTreeItemColor,
	)

	// Initialize the primary pane as active and pass in some config values.
	primaryPane := pane.NewModel(
		true,
		cfg.Settings.Borderless,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
	)

	// Initialize the secondary pane as inactive and pass in some config values.
	secondaryPane := pane.NewModel(
		false,
		cfg.Settings.Borderless,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
	)

	// Initialize a status bar passing in config values.
	statusBar := statusbar.NewModel(
		statusbar.Color{
			Background: theme.StatusBarSelectedFileBackgroundColor,
			Foreground: theme.StatusBarSelectedFileForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarBarBackgroundColor,
			Foreground: theme.StatusBarBarForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarTotalFilesBackgroundColor,
			Foreground: theme.StatusBarTotalFilesForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarLogoBackgroundColor,
			Foreground: theme.StatusBarLogoForegroundColor,
		},
		cfg.Settings.ShowIcons,
	)

	return Model{
		keys:           keys,
		help:           h,
		primaryPane:    primaryPane,
		secondaryPane:  secondaryPane,
		loader:         l,
		dirTree:        dirTree,
		dirTreePreview: dirTreePreview,
		statusBar:      statusBar,
		colorimage:     colorimage.Model{},
		markdown:       markdown.Model{},
		sourcecode:     sourcecode.Model{},
		previousKey:    tea.KeyMsg{},
		itemToMove:     nil,
		appConfig:      cfg,
		directoryItemSizeCtx: &directoryItemSizeCtx{
			ctx: context.Background(),
		},
		theme:                 theme,
		previousDirectory:     "",
		initialMoveDirectory:  "",
		showCommandBar:        false,
		inMoveMode:            false,
		inCreateFileMode:      false,
		inCreateDirectoryMode: false,
		inRenameMode:          false,
		ready:                 false,
	}
}
