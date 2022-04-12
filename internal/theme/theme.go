package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents the properties that make up a theme.
type Theme struct {
	SelectedTreeItemColor                lipgloss.AdaptiveColor
	UnselectedTreeItemColor              lipgloss.AdaptiveColor
	ActiveBoxBorderColor                 lipgloss.AdaptiveColor
	InactiveBoxBorderColor               lipgloss.AdaptiveColor
	StatusBarSelectedFileForegroundColor lipgloss.AdaptiveColor
	StatusBarSelectedFileBackgroundColor lipgloss.AdaptiveColor
	StatusBarBarForegroundColor          lipgloss.AdaptiveColor
	StatusBarBarBackgroundColor          lipgloss.AdaptiveColor
	StatusBarTotalFilesForegroundColor   lipgloss.AdaptiveColor
	StatusBarTotalFilesBackgroundColor   lipgloss.AdaptiveColor
	StatusBarLogoForegroundColor         lipgloss.AdaptiveColor
	StatusBarLogoBackgroundColor         lipgloss.AdaptiveColor
	TitleBackgroundColor                 lipgloss.AdaptiveColor
	TitleForegroundColor                 lipgloss.AdaptiveColor
}

// themeMap represents the mapping of different themes.
var themeMap = map[string]Theme{
	"default": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: "63", Light: "63"},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: "#F25D94", Light: "#F25D94"},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: "#F25D94", Light: "#F25D94"},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: "#3c3836", Light: "#3c3836"},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: "#A550DF", Light: "#A550DF"},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: "#6124DF", Light: "#6124DF"},
		TitleBackgroundColor:                 lipgloss.AdaptiveColor{Dark: "63", Light: "63"},
		TitleForegroundColor:                 lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
	},
	"gruvbox": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: "#d65d0e", Light: "#d65d0e"},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: "#b8bb26", Light: "#b8bb26"},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: "#cc241d", Light: "#cc241d"},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: "#3c3836", Light: "#3c3836"},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: "#ebcb8b", Light: "#ebcb8b"},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: "#458588", Light: "#458588"},
		TitleBackgroundColor:                 lipgloss.AdaptiveColor{Dark: "#d65d0e", Light: "#d65d0e"},
		TitleForegroundColor:                 lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
	},
	"nord": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: "#d08770", Light: "#d08770"},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: "#e5e9f0", Light: "#3b4252"},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: "#a3be8c", Light: "#a3be8c"},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: "#e5e9f0", Light: "#3b4252"},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: "#e5e9f0", Light: "#e5e9f0"},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: "#bf616a", Light: "#bf616a"},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: "#e5e9f0", Light: "#e5e9f0"},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: "#4c566a", Light: "#4c566a"},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: "#e5e9f0", Light: "#e5e9f0"},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: "#ebcb8b", Light: "#ebcb8b"},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: "#e5e9f0", Light: "#e5e9f0"},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: "#81a1c1", Light: "#81a1c1"},
		TitleBackgroundColor:                 lipgloss.AdaptiveColor{Dark: "#d08770", Light: "#d08770"},
		TitleForegroundColor:                 lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
	},
}

// GetTheme returns a theme based on the given name.
func GetTheme(theme string) Theme {
	switch theme {
	case "default":
		return themeMap["default"]
	case "gruvbox":
		return themeMap["gruvbox"]
	case "nord":
		return themeMap["nord"]
	default:
		return themeMap["default"]
	}
}
