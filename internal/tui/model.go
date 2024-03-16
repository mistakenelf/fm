package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/mistakenelf/fm/code"
	"github.com/mistakenelf/fm/filetree"
	"github.com/mistakenelf/fm/help"
	"github.com/mistakenelf/fm/image"
	"github.com/mistakenelf/fm/internal/theme"
	"github.com/mistakenelf/fm/markdown"
	"github.com/mistakenelf/fm/pdf"
	"github.com/mistakenelf/fm/statusbar"
)

type sessionState int

const (
	idleState sessionState = iota
	showCodeState
	showImageState
	showMarkdownState
	showPdfState
	showHelpState
)

type Config struct {
	StartDir       string
	SelectionPath  string
	SyntaxTheme    string
	EnableLogging  bool
	PrettyMarkdown bool
	ShowIcons      bool
	Theme          theme.Theme
}

type model struct {
	filetree      filetree.Model
	help          help.Model
	code          code.Model
	image         image.Model
	markdown      markdown.Model
	pdf           pdf.Model
	statusbar     statusbar.Model
	state         sessionState
	keyMap        keyMap
	activePane    int
	config        Config
	showTextInput bool
	textinput     textinput.Model
}

// New creates a new instance of the UI.
func New(cfg Config) model {
	filetreeModel := filetree.New(cfg.StartDir)
	filetreeModel.SetTheme(cfg.Theme.SelectedTreeItemColor, cfg.Theme.UnselectedTreeItemColor)
	filetreeModel.SetSelectionPath(cfg.SelectionPath)
	filetreeModel.SetShowIcons(cfg.ShowIcons)

	codeModel := code.New()
	codeModel.SetSyntaxTheme(cfg.SyntaxTheme)
	codeModel.SetViewportDisabled(true)

	imageModel := image.New()
	imageModel.SetViewportDisabled(true)

	markdownModel := markdown.New()
	markdownModel.SetViewportDisabled(true)

	pdfModel := pdf.New()
	pdfModel.SetViewportDisabled(true)

	statusbarModel := statusbar.New(
		statusbar.ColorConfig{
			Foreground: cfg.Theme.StatusBarSelectedFileForegroundColor,
			Background: cfg.Theme.StatusBarSelectedFileBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: cfg.Theme.StatusBarBarForegroundColor,
			Background: cfg.Theme.StatusBarBarBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: cfg.Theme.StatusBarTotalFilesForegroundColor,
			Background: cfg.Theme.StatusBarTotalFilesBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: cfg.Theme.StatusBarLogoForegroundColor,
			Background: cfg.Theme.StatusBarLogoBackgroundColor,
		},
	)

	helpModel := help.New(
		"Help",
		help.TitleColor{
			Background: cfg.Theme.TitleBackgroundColor,
			Foreground: cfg.Theme.TitleForegroundColor,
		},
		[]help.Entry{
			{Key: "ctrl+c, q", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
			{Key: "h", Description: "Go to previous directory"},
			{Key: "l", Description: "Open file or directory"},
			{Key: "G", Description: "Jump to bottom"},
			{Key: "g", Description: "Jump to top"},
			{Key: "~", Description: "Go to home directory"},
			{Key: ".", Description: "Toggle hidden files"},
			{Key: "y", Description: "Copy file path to clipboard"},
			{Key: "z", Description: "Zip currently selected tree item"},
			{Key: "u", Description: "Unzip currently selected tree item"},
			{Key: "N", Description: "Create new file"},
			{Key: "M", Description: "Create new directory"},
			{Key: "X", Description: "Delete currently selected tree item"},
			{Key: "m", Description: "Move currently selected tree item"},
			{Key: "enter", Description: "Submit text input"},
			{Key: "e", Description: "Open file in $EDITOR"},
			{Key: "C", Description: "Copy currently selected tree item"},
			{Key: "esc", Description: "Reset state and show help"},
			{Key: "/", Description: "Go to root directory"},
			{Key: "tab", Description: "Toggle between panes"},
			{Key: "K", Description: "Page up in filetree"},
			{Key: "J", Description: "Page down in filetree"},
			{Key: "D", Description: "Only show directories in filetree"},
			{Key: "F", Description: "Only show files in filetree"},
			{Key: "W", Description: "Write selection path to file if its set"},
		},
	)
	helpModel.SetViewportDisabled(true)

	textInput := textinput.New()

	return model{
		filetree:      filetreeModel,
		help:          helpModel,
		code:          codeModel,
		image:         imageModel,
		markdown:      markdownModel,
		pdf:           pdfModel,
		statusbar:     statusbarModel,
		config:        cfg,
		keyMap:        defaultKeyMap(),
		showTextInput: false,
		textinput:     textInput,
	}
}
