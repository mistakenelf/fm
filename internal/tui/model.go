package tui

import (
	"log"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/theme"

	"github.com/knipferrc/teacup/code"
	"github.com/knipferrc/teacup/filetree"
	"github.com/knipferrc/teacup/help"
	"github.com/knipferrc/teacup/image"
	"github.com/knipferrc/teacup/markdown"
	"github.com/knipferrc/teacup/pdf"
	"github.com/knipferrc/teacup/statusbar"
)

type sessionState int

const (
	idleState sessionState = iota
	showCodeState
	showImageState
	showMarkdownState
	showPdfState
)

// Bubble represents the properties of the UI.
type Bubble struct {
	filetree  filetree.Bubble
	help      help.Bubble
	code      code.Bubble
	image     image.Bubble
	markdown  markdown.Bubble
	pdf       pdf.Bubble
	statusbar statusbar.Bubble
	state     sessionState
	theme     theme.Theme
	config    config.Config
	keys      KeyMap
	activeBox int
}

// New creates a new instance of the UI.
func New() Bubble {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	theme := theme.GetTheme(cfg.Theme.AppTheme)

	filetreeModel := filetree.New(true, cfg.Settings.Borderless, theme.ActiveBoxBorderColor, theme.SelectedTreeItemColor)
	codeModel := code.New(false, cfg.Settings.Borderless, theme.InactiveBoxBorderColor)
	imageModel := image.New(false, cfg.Settings.Borderless, theme.InactiveBoxBorderColor)
	markdownModel := markdown.New(false, cfg.Settings.Borderless, theme.InactiveBoxBorderColor)
	pdfModel := pdf.New(false, cfg.Settings.Borderless, theme.InactiveBoxBorderColor)
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
		theme.InactiveBoxBorderColor,
		"Help",
		[]help.Entry{
			{Key: "ctrl+c, q", Description: "Exit FM"},
			{Key: "j/up", Description: "Move up"},
			{Key: "k/down", Description: "Move down"},
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
			{Key: "enter", Description: "Process command"},
			{Key: "E", Description: "Edit currently selected tree item"},
			{Key: "c", Description: "Copy currently selected tree item"},
			{Key: "esc", Description: "Reset input field"},
			{Key: "tab", Description: "Toggle between boxes"},
		},
		false,
		cfg.Settings.Borderless,
	)

	return Bubble{
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
