package filetree

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down       key.Binding
	Up         key.Binding
	GoToTop    key.Binding
	GoToBottom key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Down:       key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "down")),
		Up:         key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "up")),
		GoToTop:    key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "go to top")),
		GoToBottom: key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "go to bottom")),
		PageUp:     key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("pgup", "page up")),
		PageDown:   key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("pgdown", "page down")),
	}
}
