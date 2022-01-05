package tui

import "github.com/charmbracelet/lipgloss"

const StatusBarHeight = 1

var boxStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).PaddingLeft(1).PaddingRight(1)
var boldTextStyle = lipgloss.NewStyle().Bold(true)
var ellipsisStyle = "..."
var fileSizeLoadingStyle = "---"
var colors = map[string]lipgloss.Color{
	"black": "#000000",
}
