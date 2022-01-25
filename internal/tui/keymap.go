package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the app.
type KeyMap struct {
	Quit                key.Binding
	Down                key.Binding
	Up                  key.Binding
	Left                key.Binding
	Right               key.Binding
	Preview             key.Binding
	GotoBottom          key.Binding
	HomeShortcut        key.Binding
	RootShortcut        key.Binding
	ToggleHidden        key.Binding
	ShowDirectoriesOnly key.Binding
	ShowFilesOnly       key.Binding
	CopyPathToClipboard key.Binding
	Zip                 key.Binding
	Unzip               key.Binding
	NewFile             key.Binding
	NewDirectory        key.Binding
	Delete              key.Binding
	Move                key.Binding
	Enter               key.Binding
	Edit                key.Binding
	Copy                key.Binding
	Find                key.Binding
	Rename              key.Binding
	Escape              key.Binding
	ShowLogs            key.Binding
	ToggleBox           key.Binding
}

// DefaultKeyMap returns a set of default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
		),
		Preview: key.NewBinding(
			key.WithKeys("p"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("G"),
		),
		HomeShortcut: key.NewBinding(
			key.WithKeys("~"),
		),
		RootShortcut: key.NewBinding(
			key.WithKeys("/"),
		),
		ToggleHidden: key.NewBinding(
			key.WithKeys("."),
		),
		ShowDirectoriesOnly: key.NewBinding(
			key.WithKeys("S"),
		),
		ShowFilesOnly: key.NewBinding(
			key.WithKeys("s"),
		),
		CopyPathToClipboard: key.NewBinding(
			key.WithKeys("y"),
		),
		Zip: key.NewBinding(
			key.WithKeys("Z"),
		),
		Unzip: key.NewBinding(
			key.WithKeys("U"),
		),
		NewFile: key.NewBinding(
			key.WithKeys("n"),
		),
		NewDirectory: key.NewBinding(
			key.WithKeys("N"),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
		),
		Move: key.NewBinding(
			key.WithKeys("M"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
		),
		Edit: key.NewBinding(
			key.WithKeys("E"),
		),
		Copy: key.NewBinding(
			key.WithKeys("C"),
		),
		Find: key.NewBinding(
			key.WithKeys("ctrl+f"),
		),
		Rename: key.NewBinding(
			key.WithKeys("R"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
		),
		ShowLogs: key.NewBinding(
			key.WithKeys("O"),
		),
		ToggleBox: key.NewBinding(
			key.WithKeys("tab"),
		),
	}
}
