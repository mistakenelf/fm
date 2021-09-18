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
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap struct contains all keybindings.
type keyMap struct {
	Exit                  key.Binding
	Quit                  key.Binding
	Left                  key.Binding
	Down                  key.Binding
	Up                    key.Binding
	Right                 key.Binding
	GotoBottom            key.Binding
	Enter                 key.Binding
	OpenCommandBar        key.Binding
	OpenHomeDirectory     key.Binding
	OpenPreviousDirectory key.Binding
	ToggleHidden          key.Binding
	Tab                   key.Binding
	EnterMoveMode         key.Binding
	Zip                   key.Binding
	Unzip                 key.Binding
	Copy                  key.Binding
	Escape                key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Exit,
		k.Quit,
		k.Left,
		k.Down,
		k.Up,
		k.Right,
		k.GotoBottom,
		k.Enter,
		k.OpenCommandBar,
		k.OpenHomeDirectory,
		k.OpenPreviousDirectory,
		k.ToggleHidden,
		k.Tab,
		k.EnterMoveMode,
		k.Zip,
		k.Unzip,
		k.Copy,
		k.Escape,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Exit,
			k.Quit,
			k.Left,
			k.Down,
			k.Up,
			k.Right,
			k.GotoBottom,
			k.Enter,
			k.OpenCommandBar,
			k.OpenHomeDirectory,
			k.OpenPreviousDirectory,
			k.ToggleHidden,
			k.Tab,
			k.EnterMoveMode,
			k.Zip,
			k.Unzip,
			k.Copy,
			k.Escape,
		},
	}
}

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

	// Create a new spinner with some styling based on the config.
	l := spinner.NewModel()
	l.Spinner = spinner.Dot
	l.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Colors.Spinner))

	// Create a new help view.
	h := help.NewModel()
	h.Styles.FullKey.Foreground(lipgloss.Color(constants.Colors.White))
	h.Styles.FullDesc.Foreground(lipgloss.Color(constants.Colors.White))
	h.ShowAll = true

	// keys represents the key bindings in the app along with the help text.
	var keys = keyMap{
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "exit"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "go back a directory"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "scroll active pane down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "scroll active pane up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "go to bottom of active pane"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "handle move mode and command parsing"),
		),
		OpenCommandBar: key.NewBinding(
			key.WithKeys(":"),
			key.WithHelp(":", "open command bar in the status bar"),
		),
		OpenHomeDirectory: key.NewBinding(
			key.WithKeys("~"),
			key.WithHelp("~", "go to home directory"),
		),
		OpenPreviousDirectory: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "go to previous directory"),
		),
		ToggleHidden: key.NewBinding(
			key.WithKeys("."),
			key.WithHelp(".", "toggle hidden files and directories"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle between panes"),
		),
		EnterMoveMode: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "enter move mode to move files or directories"),
		),
		Zip: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp("z", "zip the currently selected file or directory"),
		),
		Unzip: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "un-zip the currently selected file or directory"),
		),
		Copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy the currently selected file or directory"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "reset to initial state"),
		),
	}

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
