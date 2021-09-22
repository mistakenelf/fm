package theme

import "github.com/knipferrc/fm/internal/constants"

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
}

func GetCurrentTheme(theme string) Theme {
	switch theme {
	case "default":
		return Theme{
			SelectedTreeItemColor:                constants.Colors.Pink,
			UnselectedTreeItemColor:              constants.Colors.White,
			ActivePaneBorderColor:                constants.Colors.Pink,
			InactivePaneBorderColor:              constants.Colors.White,
			SpinnerColor:                         constants.Colors.Pink,
			StatusBarSelectedFileForegroundColor: constants.Colors.White,
			StatusBarSelectedFileBackgroundColor: constants.Colors.Pink,
			StatusBarBarForegroundColor:          constants.Colors.White,
			StatusBarBarBackgroundColor:          constants.Colors.DarkGray,
			StatusBarTotalFilesForegroundColor:   constants.Colors.White,
			StatusBarTotalFilesBackgroundColor:   constants.Colors.LightPurple,
			StatusBarLogoForegroundColor:         constants.Colors.White,
			StatusBarLogoBackgroundColor:         constants.Colors.DarkPurple,
		}
	case "gruvbox":
		return Theme{
			SelectedTreeItemColor:                constants.Colors.Orange,
			UnselectedTreeItemColor:              constants.Colors.White,
			ActivePaneBorderColor:                constants.Colors.Green,
			InactivePaneBorderColor:              constants.Colors.White,
			SpinnerColor:                         constants.Colors.Red,
			StatusBarSelectedFileForegroundColor: constants.Colors.White,
			StatusBarSelectedFileBackgroundColor: constants.Colors.Red,
			StatusBarBarForegroundColor:          constants.Colors.White,
			StatusBarBarBackgroundColor:          constants.Colors.DarkGray,
			StatusBarTotalFilesForegroundColor:   constants.Colors.White,
			StatusBarTotalFilesBackgroundColor:   constants.Colors.Yellow,
			StatusBarLogoForegroundColor:         constants.Colors.White,
			StatusBarLogoBackgroundColor:         constants.Colors.Blue,
		}
	default:
		return Theme{
			SelectedTreeItemColor:                constants.Colors.Pink,
			UnselectedTreeItemColor:              constants.Colors.White,
			ActivePaneBorderColor:                constants.Colors.Pink,
			InactivePaneBorderColor:              constants.Colors.White,
			SpinnerColor:                         constants.Colors.Pink,
			StatusBarSelectedFileForegroundColor: constants.Colors.White,
			StatusBarSelectedFileBackgroundColor: constants.Colors.Pink,
			StatusBarBarForegroundColor:          constants.Colors.White,
			StatusBarBarBackgroundColor:          constants.Colors.DarkGray,
			StatusBarTotalFilesForegroundColor:   constants.Colors.White,
			StatusBarTotalFilesBackgroundColor:   constants.Colors.LightPurple,
			StatusBarLogoForegroundColor:         constants.Colors.White,
			StatusBarLogoBackgroundColor:         constants.Colors.DarkPurple,
		}
	}
}
