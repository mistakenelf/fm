package ui

import (
	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/filetree"
	"github.com/knipferrc/fm/internal/renderer"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/knipferrc/fm/internal/theme"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the state of the UI.
type Model struct {
	loader    spinner.Model
	fileTree  filetree.Model
	statusBar statusbar.Model
	renderer  renderer.Model
	appConfig config.Config
	theme     theme.Theme
	ready     bool
}

// NewModel create an instance of the entire application model.
func NewModel() Model {
	cfg := config.GetConfig()
	theme := theme.GetTheme(cfg.Settings.Theme)

	// Create a new spinner with some styling based on the config.
	l := spinner.NewModel()
	l.Spinner = spinner.Dot
	l.Style = lipgloss.NewStyle().Foreground(theme.SpinnerColor)

	fileTree := filetree.NewModel(
		!cfg.Settings.SimpleMode && cfg.Settings.ShowIcons,
		cfg.Settings.Borderless,
		theme.SelectedTreeItemColor,
		theme.UnselectedTreeItemColor,
		cfg,
	)

	// Initialize a status bar passing in config values.
	statusBar := statusbar.NewModel(
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

	return Model{
		loader:    l,
		fileTree:  fileTree,
		statusBar: statusBar,
		renderer:  renderer.Model{},
		appConfig: cfg,
		theme:     theme,
		ready:     false,
	}
}
