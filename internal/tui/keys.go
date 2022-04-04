package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the app.
type KeyMap struct {
	Quit      key.Binding
	Exit      key.Binding
	ToggleBox key.Binding
	OpenFile  key.Binding
}

// DefaultKeyMap returns a set of default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
		),
		Exit: key.NewBinding(
			key.WithKeys("q"),
		),
		ToggleBox: key.NewBinding(
			key.WithKeys("tab"),
		),
		OpenFile: key.NewBinding(
			key.WithKeys(" "),
		),
	}
}
