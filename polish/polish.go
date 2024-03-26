package polish

import "github.com/charmbracelet/lipgloss"

type ColorMap struct {
	Red600 lipgloss.Color
}

var Colors = ColorMap{
	Red600: lipgloss.Color("#dc2626"),
}
