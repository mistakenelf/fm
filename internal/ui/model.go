package ui

import (
	"context"
	"io/fs"

	"github.com/knipferrc/fm/internal/colorimage"
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/constants"
	"github.com/knipferrc/fm/internal/dirtree"
	"github.com/knipferrc/fm/internal/markdown"
	"github.com/knipferrc/fm/internal/pane"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/internal/text"

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
	keys                 keyMap
	help                 help.Model
	primaryPane          pane.Model
	secondaryPane        pane.Model
	loader               spinner.Model
	dirTree              dirtree.Model
	statusBar            statusbar.Model
	colorimage           colorimage.Model
	markdown             markdown.Model
	text                 text.Model
	previousKey          tea.KeyMsg
	itemToMove           fs.DirEntry
	appConfig            config.Config
	previousDirectory    string
	initialMoveDirectory string
	showCommandBar       bool
	inMoveMode           bool
	ready                bool
	directoryItemSizeCtx *directoryItemSizeCtx
}

// NewModel create an instance of the entire application model.
func NewModel() Model {
	cfg := config.GetConfig()
	keys := getDefaultKeyMap()

	// Create a new spinner with some styling based on the config.
	l := spinner.NewModel()
	l.Spinner = spinner.Dot
	l.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Colors.Spinner))

	// Create a new help view.
	h := help.NewModel()
	h.Styles.FullKey.Foreground(lipgloss.Color(constants.Colors.White))
	h.Styles.FullDesc.Foreground(lipgloss.Color(constants.Colors.White))
	h.ShowAll = true

	// Create a new dirtree.
	dirTree := dirtree.NewModel(
		cfg.Settings.ShowIcons,
		cfg.Colors.DirTree.SelectedItem,
		cfg.Colors.DirTree.UnselectedItem,
	)

	// Initialize the primary pane as active and pass in some config values.
	primaryPane := pane.NewModel(
		true,
		cfg.Settings.Borderless,
		cfg.Colors.Pane.ActiveBorderColor,
		cfg.Colors.Pane.InactiveBorderColor,
	)

	// Initialize the secondary pane as inactive and pass in some config values.
	secondaryPane := pane.NewModel(
		false,
		cfg.Settings.Borderless,
		cfg.Colors.Pane.ActiveBorderColor,
		cfg.Colors.Pane.InactiveBorderColor,
	)

	// Set secondary panes initial content to the introText.
	secondaryPane.SetContent(h.View(keys))

	// Initialize a status bar passing in config values.
	statusBar := statusbar.NewModel(
		statusbar.Color{
			Background: cfg.Colors.StatusBar.SelectedFile.Background,
			Foreground: cfg.Colors.StatusBar.SelectedFile.Foreground,
		},
		statusbar.Color{
			Background: cfg.Colors.StatusBar.Bar.Background,
			Foreground: cfg.Colors.StatusBar.Bar.Foreground,
		},
		statusbar.Color{
			Background: cfg.Colors.StatusBar.TotalFiles.Background,
			Foreground: cfg.Colors.StatusBar.TotalFiles.Foreground,
		},
		statusbar.Color{
			Background: cfg.Colors.StatusBar.Logo.Background,
			Foreground: cfg.Colors.StatusBar.Logo.Foreground,
		},
	)

	return Model{
		keys:                 keys,
		help:                 h,
		primaryPane:          primaryPane,
		secondaryPane:        secondaryPane,
		loader:               l,
		dirTree:              dirTree,
		statusBar:            statusBar,
		colorimage:           colorimage.Model{},
		markdown:             markdown.Model{},
		text:                 text.Model{},
		previousKey:          tea.KeyMsg{},
		itemToMove:           nil,
		appConfig:            cfg,
		previousDirectory:    "",
		initialMoveDirectory: "",
		showCommandBar:       false,
		inMoveMode:           false,
		ready:                false,
		directoryItemSizeCtx: &directoryItemSizeCtx{
			ctx: context.Background(),
		},
	}
}
