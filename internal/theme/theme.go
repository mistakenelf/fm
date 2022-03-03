package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents the properties that make up a theme.
type Theme struct {
	SelectedTreeItemColor                lipgloss.AdaptiveColor
	UnselectedTreeItemColor              lipgloss.AdaptiveColor
	ActiveBoxBorderColor                 lipgloss.AdaptiveColor
	InactiveBoxBorderColor               lipgloss.AdaptiveColor
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
	nordRed            string
	nordGreen          string
	nordBlue           string
	nordYellow         string
	nordWhite          string
	nordBlack          string
	nordGray           string
	nordOrange         string
	spookyPurple       string
	spookyOrange       string
	spookyYellow       string
	holidayRed         string
	holidayGreen       string
	holidayGold        string
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
	nordRed:            "#bf616a",
	nordGreen:          "#a3be8c",
	nordBlue:           "#81a1c1",
	nordYellow:         "#ebcb8b",
	nordWhite:          "#e5e9f0",
	nordBlack:          "#3b4252",
	nordGray:           "#4c566a",
	nordOrange:         "#d08770",
	spookyPurple:       "#881EE4",
	spookyOrange:       "#F75F1C",
	spookyYellow:       "#FF9A00",
	holidayRed:         "#B70D00",
	holidayGreen:       "#005C01",
	holidayGold:        "#CC9901",
}

// themeMap represents the mapping of different themes.
var themeMap = map[string]Theme{
	"default": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: colors.defaultPink, Light: colors.defaultPink},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
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
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: colors.gruvGreen, Light: colors.gruvGreen},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
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
	"nord": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.nordOrange, Light: colors.nordOrange},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.nordWhite, Light: colors.nordBlack},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: colors.nordGreen, Light: colors.nordGreen},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: colors.nordWhite, Light: colors.nordBlack},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.nordRed, Light: colors.nordRed},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.nordWhite, Light: colors.nordWhite},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.nordRed, Light: colors.nordRed},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.nordWhite, Light: colors.nordWhite},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.nordGray, Light: colors.nordGray},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.nordWhite, Light: colors.nordWhite},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.nordYellow, Light: colors.nordYellow},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.nordWhite, Light: colors.nordWhite},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.nordBlue, Light: colors.nordBlue},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.nordRed, Light: colors.nordRed},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.nordWhite, Light: colors.nordBlack},
	},
	"spooky": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.spookyOrange, Light: colors.spookyOrange},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: colors.spookyOrange, Light: colors.spookyOrange},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
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
	"holiday": {
		SelectedTreeItemColor:                lipgloss.AdaptiveColor{Dark: colors.holidayRed, Light: colors.holidayRed},
		UnselectedTreeItemColor:              lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		ActiveBoxBorderColor:                 lipgloss.AdaptiveColor{Dark: colors.holidayGold, Light: colors.holidayGold},
		InactiveBoxBorderColor:               lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
		SpinnerColor:                         lipgloss.AdaptiveColor{Dark: colors.holidayGreen, Light: colors.holidayGreen},
		StatusBarSelectedFileForegroundColor: lipgloss.AdaptiveColor{Dark: colors.holidayGold, Light: colors.holidayGold},
		StatusBarSelectedFileBackgroundColor: lipgloss.AdaptiveColor{Dark: colors.holidayGreen, Light: colors.holidayGreen},
		StatusBarBarForegroundColor:          lipgloss.AdaptiveColor{Dark: colors.holidayGreen, Light: colors.holidayGreen},
		StatusBarBarBackgroundColor:          lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.white},
		StatusBarTotalFilesForegroundColor:   lipgloss.AdaptiveColor{Dark: colors.holidayGold, Light: colors.holidayGold},
		StatusBarTotalFilesBackgroundColor:   lipgloss.AdaptiveColor{Dark: colors.holidayRed, Light: colors.holidayRed},
		StatusBarLogoForegroundColor:         lipgloss.AdaptiveColor{Dark: colors.holidayGreen, Light: colors.holidayGreen},
		StatusBarLogoBackgroundColor:         lipgloss.AdaptiveColor{Dark: colors.holidayGold, Light: colors.holidayGold},
		ErrorColor:                           lipgloss.AdaptiveColor{Dark: colors.red, Light: colors.red},
		DefaultTextColor:                     lipgloss.AdaptiveColor{Dark: colors.white, Light: colors.black},
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
	case "spooky":
		return themeMap["spooky"]
	case "holiday":
		return themeMap["holiday"]
	default:
		return themeMap["default"]
	}
}
