package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Down                key.Binding
	Up                  key.Binding
	PageUp              key.Binding
	PageDown            key.Binding
	GoToHomeDirectory   key.Binding
	GoToRootDirectory   key.Binding
	ToggleHidden        key.Binding
	OpenDirectory       key.Binding
	PreviousDirectory   key.Binding
	CopyPathToClipboard key.Binding
	CopyDirectoryItem   key.Binding
	DeleteDirectoryItem key.Binding
	ZipDirectoryItem    key.Binding
	UnzipDirectoryItem  key.Binding
	ShowDirectoriesOnly key.Binding
	ShowFilesOnly       key.Binding
	WriteSelectionPath  key.Binding
	OpenInEditor        key.Binding
	CreateFile          key.Binding
	CreateDirectory     key.Binding
	ForceQuit           key.Binding
	Quit                key.Binding
	TogglePane          key.Binding
	OpenFile            key.Binding
	ResetState          key.Binding
	ShowTextInput       key.Binding
	Submit              key.Binding
	GotoTop             key.Binding
	GotoBottom          key.Binding
	MoveDirectoryItem   key.Binding
	RenameDirectoryItem key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ForceQuit:           key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "Force Quit")),
		Quit:                key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit when not performing action")),
		TogglePane:          key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Toggle between l/r panes")),
		OpenFile:            key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "Preview file")),
		ResetState:          key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Reset state")),
		ShowTextInput:       key.NewBinding(key.WithKeys("N", "M", "R"), key.WithHelp("N, M", "Show text input to create file or directory")),
		Submit:              key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Submit text input value")),
		GotoTop:             key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "Go to top of pane")),
		GotoBottom:          key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "Go to bottom of pane")),
		MoveDirectoryItem:   key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "Move directory item")),
		Down:                key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "Go down")),
		Up:                  key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "Go up")),
		PageUp:              key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("pgup", "Page up")),
		PageDown:            key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("pgdown", "Page down")),
		GoToHomeDirectory:   key.NewBinding(key.WithKeys("~"), key.WithHelp("~", "Go to home directory")),
		GoToRootDirectory:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Go to root directory")),
		ToggleHidden:        key.NewBinding(key.WithKeys("."), key.WithHelp(".", "Toggle hidden files/folders")),
		OpenDirectory:       key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "Open directory")),
		PreviousDirectory:   key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "Go to previous directory")),
		CopyPathToClipboard: key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "Copy path to clipboard")),
		CopyDirectoryItem:   key.NewBinding(key.WithKeys("C"), key.WithHelp("C", "Copy directory item")),
		DeleteDirectoryItem: key.NewBinding(key.WithKeys("X"), key.WithHelp("X", "Delete directory item")),
		ZipDirectoryItem:    key.NewBinding(key.WithKeys("Z"), key.WithHelp("Z", "Zip directory item")),
		UnzipDirectoryItem:  key.NewBinding(key.WithKeys("U"), key.WithHelp("U", "Unzip directory item")),
		ShowDirectoriesOnly: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "Show directories only")),
		ShowFilesOnly:       key.NewBinding(key.WithKeys("F"), key.WithHelp("F", "Show files only")),
		WriteSelectionPath:  key.NewBinding(key.WithKeys("W"), key.WithHelp("W", "Write selection path and quit")),
		OpenInEditor:        key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Open in $EDITOR")),
		CreateFile:          key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "Create new file")),
		CreateDirectory:     key.NewBinding(key.WithKeys("M"), key.WithHelp("M", "Create new directory")),
		RenameDirectoryItem: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "Rename directory items")),
	}
}
