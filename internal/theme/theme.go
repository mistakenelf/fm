package theme

type Theme struct {
	SelectedTreeItemColor                string
	UnselectedTreeItemColor              string
	ActivePaneBorderColor                string
	InactivePaneBorderColor              string
	SpinnerColor                         string
	StatusBarSelectedFileForegroundColor string
	StatusBarSelectedFileBackgroundColor string
	StatusBarBarForegroundColor          string
	StatusBarBarBackgroundColor          string
	StatusBarTotalFilesForegroundColor   string
	StatusBarTotalFilesBackgroundColor   string
	StatusBarLogoForegroundColor         string
	StatusBarLogoBackgroundColor         string
	ErrorColor                           string
	DefaultTextColor                     string
}

// appColors contains the different types of colors.
type appColors struct {
	White       string
	Pink        string
	LightPurple string
	DarkPurple  string
	DarkGray    string
	Red         string
	Green       string
	Blue        string
	Yellow      string
	Orange      string
}

// Colors contains the different kinds of colors and their values.
var colors = appColors{
	White:       "#FFFDF5",
	Pink:        "#F25D94",
	LightPurple: "#A550DF",
	DarkPurple:  "#6124DF",
	DarkGray:    "#3c3836",
	Red:         "#cc241d",
	Green:       "#b8bb26",
	Blue:        "#458588",
	Yellow:      "#d79921",
	Orange:      "#d65d0e",
}

var defaultTheme = Theme{
	SelectedTreeItemColor:                colors.Pink,
	UnselectedTreeItemColor:              colors.White,
	ActivePaneBorderColor:                colors.Pink,
	InactivePaneBorderColor:              colors.White,
	SpinnerColor:                         colors.Pink,
	StatusBarSelectedFileForegroundColor: colors.White,
	StatusBarSelectedFileBackgroundColor: colors.Pink,
	StatusBarBarForegroundColor:          colors.White,
	StatusBarBarBackgroundColor:          colors.DarkGray,
	StatusBarTotalFilesForegroundColor:   colors.White,
	StatusBarTotalFilesBackgroundColor:   colors.LightPurple,
	StatusBarLogoForegroundColor:         colors.White,
	StatusBarLogoBackgroundColor:         colors.DarkPurple,
	ErrorColor:                           colors.Red,
	DefaultTextColor:                     colors.White,
}

var gruvboxTheme = Theme{
	SelectedTreeItemColor:                colors.Orange,
	UnselectedTreeItemColor:              colors.White,
	ActivePaneBorderColor:                colors.Green,
	InactivePaneBorderColor:              colors.White,
	SpinnerColor:                         colors.Red,
	StatusBarSelectedFileForegroundColor: colors.White,
	StatusBarSelectedFileBackgroundColor: colors.Red,
	StatusBarBarForegroundColor:          colors.White,
	StatusBarBarBackgroundColor:          colors.DarkGray,
	StatusBarTotalFilesForegroundColor:   colors.White,
	StatusBarTotalFilesBackgroundColor:   colors.Yellow,
	StatusBarLogoForegroundColor:         colors.White,
	StatusBarLogoBackgroundColor:         colors.Blue,
	ErrorColor:                           colors.Red,
	DefaultTextColor:                     colors.White,
}

func GetCurrentTheme(theme string) Theme {
	switch theme {
	case "default":
		return defaultTheme
	case "gruvbox":
		return gruvboxTheme
	default:
		return defaultTheme
	}
}
