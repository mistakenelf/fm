package constants

import "github.com/charmbracelet/lipgloss"

const (
	PrimaryBoxActive = iota
	SecondaryBoxActive
)
const (
	StatusBarHeight      = 1
	BoxPadding           = 1
	EllipsisStyle        = "..."
	FileSizeLoadingStyle = "---"
)

var BoldTextStyle = lipgloss.NewStyle().Bold(true)
var StarredBorder = lipgloss.Border{
	Top:         "-",
	Bottom:      "-",
	Left:        "|",
	Right:       "|",
	TopLeft:     "*",
	TopRight:    "*",
	BottomLeft:  "*",
	BottomRight: "*",
}

var Colors = map[string]lipgloss.Color{
	"black": "#000000",
}
