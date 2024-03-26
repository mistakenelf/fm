package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit              key.Binding
	TogglePane        key.Binding
	OpenFile          key.Binding
	ResetState        key.Binding
	ShowTextInput     key.Binding
	Submit            key.Binding
	GotoTop           key.Binding
	GotoBottom        key.Binding
	MoveDirectoryItem key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit:              key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		TogglePane:        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle pane")),
		OpenFile:          key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "open file")),
		ResetState:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "reset state")),
		ShowTextInput:     key.NewBinding(key.WithKeys("N", "M"), key.WithHelp("N, M", "show text input")),
		Submit:            key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit text input")),
		GotoTop:           key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "go to top")),
		GotoBottom:        key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "go to bottom")),
		MoveDirectoryItem: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "move directory item")),
	}
}
