package tui

import "github.com/charmbracelet/lipgloss"

const (
	PrimaryBoxActive = iota
	SecondaryBoxActive
)
const (
	StatusBarHeight      = 1
	BoxPadding           = 1
	ellipsisStyle        = "..."
	fileSizeLoadingStyle = "---"
)

var boldTextStyle = lipgloss.NewStyle().Bold(true)
var starredBorder = lipgloss.Border{
	Top:         "-",
	Bottom:      "-",
	Left:        "|",
	Right:       "|",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}

var colors = map[string]lipgloss.Color{
	"black": "#000000",
}
