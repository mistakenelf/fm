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
	White              string
	DarkGray           string
	Red                string
	Black              string
	DefaultPink        string
	DefaultLightPurple string
	DefaultDarkPurple  string
	GruvGreen          string
	GruvBlue           string
	GruvYellow         string
	GruvOrange         string
	SpookyPurple       string
	SpookyOrange       string
	SpookyYellow       string
}

// Colors contains the different kinds of colors and their values.
var colors = appColors{
	White:              "#FFFDF5",
	DarkGray:           "#3c3836",
	Red:                "#cc241d",
	Black:              "#000000",
	DefaultPink:        "#F25D94",
	DefaultLightPurple: "#A550DF",
	DefaultDarkPurple:  "#6124DF",
	GruvGreen:          "#b8bb26",
	GruvBlue:           "#458588",
	GruvYellow:         "#d79921",
	GruvOrange:         "#d65d0e",
	SpookyPurple:       "#881EE4 ",
	SpookyOrange:       "#F75F1C ",
	SpookyYellow:       "#FF9A00 ",
}

var defaultTheme = Theme{
	SelectedTreeItemColor:                colors.DefaultPink,
	UnselectedTreeItemColor:              colors.White,
	ActivePaneBorderColor:                colors.DefaultPink,
	InactivePaneBorderColor:              colors.White,
	SpinnerColor:                         colors.DefaultPink,
	StatusBarSelectedFileForegroundColor: colors.White,
	StatusBarSelectedFileBackgroundColor: colors.DefaultPink,
	StatusBarBarForegroundColor:          colors.White,
	StatusBarBarBackgroundColor:          colors.DarkGray,
	StatusBarTotalFilesForegroundColor:   colors.White,
	StatusBarTotalFilesBackgroundColor:   colors.DefaultLightPurple,
	StatusBarLogoForegroundColor:         colors.White,
	StatusBarLogoBackgroundColor:         colors.DefaultDarkPurple,
	ErrorColor:                           colors.Red,
	DefaultTextColor:                     colors.White,
}

var gruvboxTheme = Theme{
	SelectedTreeItemColor:                colors.GruvOrange,
	UnselectedTreeItemColor:              colors.White,
	ActivePaneBorderColor:                colors.GruvGreen,
	InactivePaneBorderColor:              colors.White,
	SpinnerColor:                         colors.Red,
	StatusBarSelectedFileForegroundColor: colors.White,
	StatusBarSelectedFileBackgroundColor: colors.Red,
	StatusBarBarForegroundColor:          colors.White,
	StatusBarBarBackgroundColor:          colors.DarkGray,
	StatusBarTotalFilesForegroundColor:   colors.White,
	StatusBarTotalFilesBackgroundColor:   colors.GruvYellow,
	StatusBarLogoForegroundColor:         colors.White,
	StatusBarLogoBackgroundColor:         colors.GruvBlue,
	ErrorColor:                           colors.Red,
	DefaultTextColor:                     colors.White,
}

var spookyTheme = Theme{
	SelectedTreeItemColor:                colors.SpookyOrange,
	UnselectedTreeItemColor:              colors.White,
	ActivePaneBorderColor:                colors.SpookyOrange,
	InactivePaneBorderColor:              colors.White,
	SpinnerColor:                         colors.Red,
	StatusBarSelectedFileForegroundColor: colors.White,
	StatusBarSelectedFileBackgroundColor: colors.SpookyPurple,
	StatusBarBarForegroundColor:          colors.White,
	StatusBarBarBackgroundColor:          colors.Black,
	StatusBarTotalFilesForegroundColor:   colors.White,
	StatusBarTotalFilesBackgroundColor:   colors.SpookyYellow,
	StatusBarLogoForegroundColor:         colors.White,
	StatusBarLogoBackgroundColor:         colors.SpookyOrange,
	ErrorColor:                           colors.Red,
	DefaultTextColor:                     colors.White,
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
