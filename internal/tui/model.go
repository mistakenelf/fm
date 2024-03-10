package tui

import (
	"log"

	"github.com/mistakenelf/fm/internal/config"
	"github.com/mistakenelf/fm/internal/theme"

	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/fm/code"
	"github.com/mistakenelf/fm/filetree"
	"github.com/mistakenelf/fm/help"
	"github.com/mistakenelf/fm/image"
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
)

// model represents the properties of the UI.
type model struct {
	filetree  filetree.Model
	help      help.Model
	code      code.Model
	image     image.Model
	markdown  markdown.Model
	pdf       pdf.Model
	statusbar statusbar.Model
	state     sessionState
	theme     theme.Theme
	config    config.Config
	keys      KeyMap
	activeBox int
}

// New creates a new instance of the UI.
func New(startDir, selectionPath string) model {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	theme := theme.GetTheme(cfg.Theme.AppTheme)

	syntaxTheme := cfg.Theme.SyntaxTheme.Light
	if lipgloss.HasDarkBackground() {
		syntaxTheme = cfg.Theme.SyntaxTheme.Dark
	}

	filetreeModel := filetree.New(true)

	codeModel := code.New(false)
	codeModel.SetSyntaxTheme(syntaxTheme)

	imageModel := image.New(false)
	markdownModel := markdown.New(false)
	pdfModel := pdf.New(false)
	statusbarModel := statusbar.New(
		statusbar.ColorConfig{
			Foreground: theme.StatusBarSelectedFileForegroundColor,
			Background: theme.StatusBarSelectedFileBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusBarBarForegroundColor,
			Background: theme.StatusBarBarBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusBarTotalFilesForegroundColor,
			Background: theme.StatusBarTotalFilesBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusBarLogoForegroundColor,
			Background: theme.StatusBarLogoBackgroundColor,
		},
	)

	helpModel := help.New(
		false,
		"Help",
		help.TitleColor{
			Background: theme.TitleBackgroundColor,
			Foreground: theme.TitleForegroundColor,
		},
		[]help.Entry{
			{Key: "ctrl+c, q", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
			{Key: "h", Description: "Paginate left in current directory"},
			{Key: "l", Description: "Paginate right in current directory"},
			{Key: "space", Description: "Read file or enter directory"},
			{Key: "G", Description: "Jump to bottom"},
			{Key: "g", Description: "Jump to top"},
			{Key: "~", Description: "Go to home directory"},
			{Key: ".", Description: "Toggle hidden files"},
			{Key: "y", Description: "Copy file path to clipboard"},
			{Key: "z", Description: "Zip currently selected tree item"},
			{Key: "u", Description: "Unzip currently selected tree item"},
			{Key: "n", Description: "Create new file"},
			{Key: "N", Description: "Create new directory"},
			{Key: "x", Description: "Delete currently selected tree item"},
			{Key: "m", Description: "Move currently selected tree item"},
			{Key: "enter", Description: "Process command"},
			{Key: "e", Description: "Edit currently selected tree item"},
			{Key: "c", Description: "Copy currently selected tree item"},
			{Key: "esc", Description: "Reset input field"},
			{Key: "R", Description: "Go to root directory"},
			{Key: "tab", Description: "Toggle between boxes"},
		},
	)

	return model{
		filetree:  filetreeModel,
		help:      helpModel,
		code:      codeModel,
		image:     imageModel,
		markdown:  markdownModel,
		pdf:       pdfModel,
		statusbar: statusbarModel,
		theme:     theme,
		config:    cfg,
		keys:      DefaultKeyMap(),
	}
}
