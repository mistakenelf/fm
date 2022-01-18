package tui

import "github.com/charmbracelet/lipgloss"

const StatusBarHeight = 1

var boldTextStyle = lipgloss.NewStyle().Bold(true)
var ellipsisStyle = "..."
var fileSizeLoadingStyle = "---"

var colors = map[string]lipgloss.Color{
	"black": "#000000",
}

const (
	PrimaryBoxActive = iota
	SecondaryBoxActive
)
