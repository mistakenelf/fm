package polish

import "github.com/charmbracelet/lipgloss"

type ColorMap struct {
	Red600    lipgloss.Color
	Yellow500 lipgloss.Color
}

var Colors = ColorMap{
	Red600:    lipgloss.Color("#dc2626"),
	Yellow500: lipgloss.Color("#eab308"),
}

type AdaptiveColorMap struct {
	DefaultText lipgloss.AdaptiveColor
}

var AdaptiveColors = AdaptiveColorMap{
	DefaultText: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#000000"},
}
