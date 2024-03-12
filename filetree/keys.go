package filetree

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down              key.Binding
	Up                key.Binding
	GoToTop           key.Binding
	GoToBottom        key.Binding
	PageUp            key.Binding
	PageDown          key.Binding
	GoToHomeDirectory key.Binding
	GoToRootDirectory key.Binding
	ToggleHidden      key.Binding
	OpenDirectory     key.Binding
	PreviousDirectory key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Down:              key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "down")),
		Up:                key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "up")),
		GoToTop:           key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "go to top")),
		GoToBottom:        key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "go to bottom")),
		PageUp:            key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("pgup", "page up")),
		PageDown:          key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("pgdown", "page down")),
		GoToHomeDirectory: key.NewBinding(key.WithKeys("~"), key.WithHelp("~", "go to home directory")),
		GoToRootDirectory: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "go to root directory")),
		ToggleHidden:      key.NewBinding(key.WithKeys("."), key.WithHelp(".", "toggle hidden")),
		OpenDirectory:     key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "open directory")),
		PreviousDirectory: key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "previous directory")),
	}
}
