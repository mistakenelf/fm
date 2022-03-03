package constants

import "github.com/charmbracelet/lipgloss"

const (
	// PrimaryBoxActive represents when the left box is active.
	PrimaryBoxActive = iota

	// SecondaryBoxActive represents when the right box is active.
	SecondaryBoxActive
)
const (
	// StatusBarHeight represents the height of the status bar.
	StatusBarHeight = 1

	// BoxPadding represents the padding of the boxes.
	BoxPadding = 1

	// EllipsisStyle represents the characters displayed when overflowing.
	EllipsisStyle = "..."

	// FileSizeLoadingStyle represents the characters displayed when file sizes are loading.
	FileSizeLoadingStyle = "---"
)

// BoldTextStyle is the style used for bold text.
var BoldTextStyle = lipgloss.NewStyle().Bold(true)

// StarredBorder is the border style used for moving files.
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

// Colors represents the colors used in the application.
var Colors = map[string]lipgloss.Color{
	"black": "#000000",
}
