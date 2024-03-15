package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit            key.Binding
	TogglePane      key.Binding
	OpenFile        key.Binding
	ResetState      key.Binding
	ShowTextInput   key.Binding
	SubmitTextInput key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit:            key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		TogglePane:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "toggle pane")),
		OpenFile:        key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "open file")),
		ResetState:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "reset state")),
		ShowTextInput:   key.NewBinding(key.WithKeys("ctrl+f"), key.WithHelp("ctrl+f", "show text input")),
		SubmitTextInput: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit text input")),
	}
}
