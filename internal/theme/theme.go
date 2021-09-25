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

var defaultTheme = Theme{
	SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
	UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
	ActivePaneBorderColor:                lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
	InactivePaneBorderColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
	SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
	StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.White},
	StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.DefaultPink},
	StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.White},
	StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.darkGray, Light: colors.DarkGray},
	StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.White},
	StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.defaultLightPurple, Light: colors.DefaultLightPurple},
	StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.White},
	StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.defaultDarkPurple, Light: colors.DefaultDarkPurple},
	ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.Red},
	DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.Black},
}

var gruvboxTheme = Theme{
	SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.GruvOrange, Light: colors.GruvOrange},
	UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.Black},
	ActivePaneBorderColor:                lipgloss.AdaptiveColor{Dark: colors.GruvGreen, Light: colors.GruvGreen},
	InactivePaneBorderColor:              lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.Black},
	SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.Red, Light: colors.Red},
	StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.Red, Light: colors.Red},
	StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.DarkGray, Light: colors.DarkGray},
	StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.GruvYellow, Light: colors.GruvYellow},
	StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.GruvBlue, Light: colors.GruvBlue},
	ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.Red, Light: colors.Red},
	DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.Black},
}

var spookyTheme = Theme{
	SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.SpookyOrange, Light: colors.SpookyOrange},
	UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.Black},
	ActivePaneBorderColor:                lipgloss.AdaptiveColor{Dark: colors.SpookyOrange, Light: colors.SpookyOrange},
	InactivePaneBorderColor:              lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.Black},
	SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.Red, Light: colors.Red},
	StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.SpookyPurple, Light: colors.SpookyPurple},
	StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.Black, Light: colors.Black},
	StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.SpookyYellow, Light: colors.SpookyYellow},
	StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.White},
	StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.SpookyOrange, Light: colors.SpookyOrange},
	ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.Red, Light: colors.Red},
	DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.White, Light: colors.Black},
}

func GetCurrentTheme(theme string) Theme {
	switch theme {
	case "default":
		return defaultTheme
	case "gruvbox":
		return gruvboxTheme
	case "spooky":
		return spookyTheme
	default:
		return defaultTheme
	}
}
