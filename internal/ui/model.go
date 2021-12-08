package ui

import (
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/filetree"
	"github.com/knipferrc/fm/internal/previewer"
	"github.com/knipferrc/fm/internal/renderer"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/internal/theme"
)

// Bubble represents the state of the UI.
type Bubble struct {
	fileTree  filetree.Bubble
	previewer previewer.Bubble
	statusBar statusbar.Bubble
	renderer  renderer.Bubble
	appConfig config.Config
	theme     theme.Theme
}

// NewBubble create an instance of the entire application model.
func NewBubble() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.Settings.Theme)

	fileTree := filetree.NewBubble(
		!cfg.Settings.SimpleMode && cfg.Settings.ShowIcons,
		cfg.Settings.Borderless,
		true,
		true,
		theme.SelectedTreeItemColor,
		theme.UnselectedTreeItemColor,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
		cfg,
	)

	previewer := previewer.NewBubble(
		!cfg.Settings.SimpleMode && cfg.Settings.ShowIcons,
		cfg.Settings.Borderless,
		false,
		true,
		theme.UnselectedTreeItemColor,
		theme.UnselectedTreeItemColor,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
		cfg,
	)

	renderer := renderer.NewBubble(
		cfg.Settings.Borderless,
		false,
		theme.ActivePaneBorderColor,
		theme.InactivePaneBorderColor,
	)

	statusBar := statusbar.NewBubble(
		statusbar.Color{
			Background: theme.StatusBarSelectedFileBackgroundColor,
			Foreground: theme.StatusBarSelectedFileForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarBarBackgroundColor,
			Foreground: theme.StatusBarBarForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarTotalFilesBackgroundColor,
			Foreground: theme.StatusBarTotalFilesForegroundColor,
		},
		statusbar.Color{
			Background: theme.StatusBarLogoBackgroundColor,
			Foreground: theme.StatusBarLogoForegroundColor,
		},
		!cfg.Settings.SimpleMode && cfg.Settings.ShowIcons,
		cfg.Settings.SimpleMode,
	)

	return Bubble{
		fileTree:  fileTree,
		previewer: previewer,
		statusBar: statusBar,
		renderer:  renderer,
		appConfig: cfg,
		theme:     theme,
	}
}
