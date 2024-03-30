package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/mistakenelf/fm/code"
	"github.com/mistakenelf/fm/filetree"
	"github.com/mistakenelf/fm/help"
	"github.com/mistakenelf/fm/icons"
	"github.com/mistakenelf/fm/image"
	"github.com/mistakenelf/fm/internal/theme"
	"github.com/mistakenelf/fm/keys"
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
	showMoveState
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
	filetree              filetree.Model
	secondaryFiletree     filetree.Model
	help                  help.Model
	code                  code.Model
	image                 image.Model
	markdown              markdown.Model
	pdf                   pdf.Model
	statusbar             statusbar.Model
	state                 sessionState
	keyMap                keys.KeyMap
	activePane            int
	height                int
	config                Config
	showTextInput         bool
	textinput             textinput.Model
	statusMessage         string
	directoryBeforeMove   string
	statusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
}

// New creates a new instance of the UI.
func New(cfg Config) model {
	filetreeModel := filetree.New(cfg.StartDir)
	filetreeModel.SetTheme(cfg.Theme.SelectedTreeItemColor, cfg.Theme.UnselectedTreeItemColor)
	filetreeModel.SetSelectionPath(cfg.SelectionPath)
	filetreeModel.SetShowIcons(cfg.ShowIcons)

	secondaryFiletree := filetree.New(cfg.StartDir)
	secondaryFiletree.SetTheme(cfg.Theme.SelectedTreeItemColor, cfg.Theme.UnselectedTreeItemColor)
	secondaryFiletree.SetSelectionPath(cfg.SelectionPath)
	secondaryFiletree.SetShowIcons(cfg.ShowIcons)
	secondaryFiletree.SetDisabled(true)

	if cfg.ShowIcons {
		icons.ParseIcons()
	}

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

	defaultKeyMap := keys.DefaultKeyMap()

	helpModel := help.New(
		"Help",
		help.TitleColor{
			Background: cfg.Theme.TitleBackgroundColor,
			Foreground: cfg.Theme.TitleForegroundColor,
		},
		[]help.Entry{
			{Key: defaultKeyMap.ForceQuit.Help().Key, Description: defaultKeyMap.ForceQuit.Help().Desc},
			{Key: defaultKeyMap.Quit.Help().Key, Description: defaultKeyMap.Quit.Help().Desc},
			{Key: defaultKeyMap.TogglePane.Help().Key, Description: defaultKeyMap.TogglePane.Help().Desc},
			{Key: defaultKeyMap.OpenFile.Help().Key, Description: defaultKeyMap.OpenFile.Help().Desc},
			{Key: defaultKeyMap.ResetState.Help().Key, Description: defaultKeyMap.ResetState.Help().Desc},
			{Key: defaultKeyMap.ShowTextInput.Help().Key, Description: defaultKeyMap.ShowTextInput.Help().Desc},
			{Key: defaultKeyMap.Submit.Help().Key, Description: defaultKeyMap.Submit.Help().Desc},
			{Key: defaultKeyMap.GotoTop.Help().Key, Description: defaultKeyMap.GotoTop.Help().Desc},
			{Key: defaultKeyMap.GotoBottom.Help().Key, Description: defaultKeyMap.GotoBottom.Help().Desc},
			{Key: defaultKeyMap.MoveDirectoryItem.Help().Key, Description: defaultKeyMap.MoveDirectoryItem.Help().Desc},
			{Key: defaultKeyMap.Down.Help().Key, Description: defaultKeyMap.Down.Help().Desc},
			{Key: defaultKeyMap.Up.Help().Key, Description: defaultKeyMap.Up.Help().Desc},
			{Key: defaultKeyMap.PageUp.Help().Key, Description: defaultKeyMap.PageUp.Help().Desc},
			{Key: defaultKeyMap.PageDown.Help().Key, Description: defaultKeyMap.PageDown.Help().Desc},
			{Key: defaultKeyMap.GoToHomeDirectory.Help().Key, Description: defaultKeyMap.GoToHomeDirectory.Help().Desc},
			{Key: defaultKeyMap.GoToRootDirectory.Help().Key, Description: defaultKeyMap.GoToRootDirectory.Help().Desc},
			{Key: defaultKeyMap.ToggleHidden.Help().Key, Description: defaultKeyMap.ToggleHidden.Help().Desc},
			{Key: defaultKeyMap.OpenDirectory.Help().Key, Description: defaultKeyMap.OpenDirectory.Help().Desc},
			{Key: defaultKeyMap.PreviousDirectory.Help().Key, Description: defaultKeyMap.PreviousDirectory.Help().Desc},
			{Key: defaultKeyMap.CopyPathToClipboard.Help().Key, Description: defaultKeyMap.CopyPathToClipboard.Help().Desc},
			{Key: defaultKeyMap.CopyDirectoryItem.Help().Key, Description: defaultKeyMap.CopyDirectoryItem.Help().Desc},
			{Key: defaultKeyMap.DeleteDirectoryItem.Help().Key, Description: defaultKeyMap.DeleteDirectoryItem.Help().Desc},
			{Key: defaultKeyMap.ZipDirectoryItem.Help().Key, Description: defaultKeyMap.ZipDirectoryItem.Help().Desc},
			{Key: defaultKeyMap.UnzipDirectoryItem.Help().Key, Description: defaultKeyMap.UnzipDirectoryItem.Help().Desc},
			{Key: defaultKeyMap.ShowDirectoriesOnly.Help().Key, Description: defaultKeyMap.ShowDirectoriesOnly.Help().Desc},
			{Key: defaultKeyMap.ShowFilesOnly.Help().Key, Description: defaultKeyMap.ShowFilesOnly.Help().Desc},
			{Key: defaultKeyMap.WriteSelectionPath.Help().Key, Description: defaultKeyMap.WriteSelectionPath.Help().Desc},
			{Key: defaultKeyMap.OpenInEditor.Help().Key, Description: defaultKeyMap.OpenInEditor.Help().Desc},
			{Key: defaultKeyMap.CreateFile.Help().Key, Description: defaultKeyMap.CreateFile.Help().Desc},
			{Key: defaultKeyMap.CreateDirectory.Help().Key, Description: defaultKeyMap.CreateDirectory.Help().Desc},
		},
	)
	helpModel.SetViewportDisabled(true)

	textInput := textinput.New()

	return model{
		filetree:              filetreeModel,
		secondaryFiletree:     secondaryFiletree,
		help:                  helpModel,
		code:                  codeModel,
		image:                 imageModel,
		markdown:              markdownModel,
		pdf:                   pdfModel,
		statusbar:             statusbarModel,
		config:                cfg,
		keyMap:                defaultKeyMap,
		showTextInput:         false,
		textinput:             textInput,
		statusMessageLifetime: time.Second,
	}
}
