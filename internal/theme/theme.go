package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	SelectedTreeItemColor                lipgloss.AdaptiveColor
	UnselectedTreeItemColor              lipgloss.AdaptiveColor
	ActivePaneBorderColor                lipgloss.AdaptiveColor
	InactivePaneBorderColor              lipgloss.AdaptiveColor
	SpinnerColor                         lipgloss.AdaptiveColor
	StatusBarSelectedFileForegroundColor lipgloss.AdaptiveColor
	StatusBarSelectedFileBackgroundColor lipgloss.AdaptiveColor
	StatusBarBarForegroundColor          lipgloss.AdaptiveColor
	StatusBarBarBackgroundColor          lipgloss.AdaptiveColor
	StatusBarTotalFilesForegroundColor   lipgloss.AdaptiveColor
	StatusBarTotalFilesBackgroundColor   lipgloss.AdaptiveColor
	StatusBarLogoForegroundColor         lipgloss.AdaptiveColor
	StatusBarLogoBackgroundColor         lipgloss.AdaptiveColor
	ErrorColor                           lipgloss.AdaptiveColor
	DefaultTextColor                     lipgloss.AdaptiveColor
}

// appColors contains the different types of colors.
type appColors struct {
	white              string
	darkGray           string
	red                string
	black              string
	defaultPink        string
	defaultLightPurple string
	defaultDarkPurple  string
	gruvGreen          string
	gruvBlue           string
	gruvYellow         string
	gruvOrange         string
	spookyPurple       string
	spookyOrange       string
	spookyYellow       string
}

// Colors contains the different kinds of colors and their values.
var colors = appColors{
	white:              "#FFFDF5",
	darkGray:           "#3c3836",
	red:                "#cc241d",
	black:              "#000000",
	defaultPink:        "#F25D94",
	defaultLightPurple: "#A550DF",
	defaultDarkPurple:  "#6124DF",
	gruvGreen:          "#b8bb26",
	gruvBlue:           "#458588",
	gruvYellow:         "#d79921",
	gruvOrange:         "#d65d0e",
	spookyPurple:       "#881EE4 ",
	spookyOrange:       "#F75F1C ",
	spookyYellow:       "#FF9A00 ",
}

var themeMap = map[string]Theme{
	"default": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		ActivePaneBorderColor:                lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		InactivePaneBorderColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.darkGray, Light: colors.darkGray},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.defaultLightPurple, Light: colors.defaultLightPurple},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.defaultDarkPurple, Light: colors.defaultDarkPurple},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
	},
	"gruvbox": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.gruvOrange, Light: colors.gruvOrange},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		ActivePaneBorderColor:                lipgloss.AdaptiveColor{Dark: colors.gruvGreen, Light: colors.gruvGreen},
		InactivePaneBorderColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.darkGray, Light: colors.darkGray},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.gruvYellow, Light: colors.gruvYellow},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.gruvBlue, Light: colors.gruvBlue},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
	},
	"spooky": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.spookyOrange, Light: colors.spookyOrange},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		ActivePaneBorderColor:                lipgloss.AdaptiveColor{Dark: colors.spookyOrange, Light: colors.spookyOrange},
		InactivePaneBorderColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.spookyPurple, Light: colors.spookyPurple},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.black, Light: colors.black},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.spookyYellow, Light: colors.spookyYellow},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.spookyOrange, Light: colors.spookyOrange},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
	},
}

func GetCurrentTheme(theme string) Theme {
	switch theme {
	case "default":
		return themeMap["default"]
	case "gruvbox":
		return themeMap["gruvbox"]
	case "spooky":
		return themeMap["spooky"]
	default:
		return themeMap["default"]
	}
}
