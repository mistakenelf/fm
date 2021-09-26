package ui

import "github.com/charmbracelet/bubbles/key"

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
	OpenHomeDirectory     key.Binding
	OpenPreviousDirectory key.Binding
	ToggleHidden          key.Binding
	Tab                   key.Binding
	EnterMoveMode         key.Binding
	Zip                   key.Binding
	Unzip                 key.Binding
	Copy                  key.Binding
	Escape                key.Binding
	Delete                key.Binding
	CreateFile            key.Binding
	CreateDirectory       key.Binding
	Rename                key.Binding
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
		k.OpenHomeDirectory,
		k.OpenPreviousDirectory,
		k.ToggleHidden,
		k.Tab,
		k.EnterMoveMode,
		k.Zip,
		k.Unzip,
		k.Copy,
		k.Escape,
		k.Delete,
		k.CreateFile,
		k.CreateDirectory,
		k.Rename,
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
			k.OpenHomeDirectory,
			k.OpenPreviousDirectory,
			k.ToggleHidden,
			k.Tab,
			k.EnterMoveMode,
			k.Zip,
			k.Unzip,
			k.Copy,
			k.Escape,
			k.Delete,
			k.CreateFile,
			k.CreateDirectory,
			k.Rename,
		},
	}
}

// getDefaultKeyMap returns the default keybindings.
func getDefaultKeyMap() keyMap {
	return keyMap{
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
		Delete: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "delete the selected file or directory"),
		),
		CreateFile: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "create a new file"),
		),
		CreateDirectory: key.NewBinding(
			key.WithKeys("N"),
			key.WithHelp("N", "create a new directory"),
		),
		Rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename the currently selected file or directory"),
		),
	}
}
