package helpers

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// RenderMarkdown renders the markdown content with glamour.
func RenderMarkdown(width int, content string) (string, error) {
	bg := "light"

	if lipgloss.HasDarkBackground() {
		bg = "dark"
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithStandardStyle(bg),
	)

	out, err := r.Render(content)
	if err != nil {
		return "", err
	}

	return out, nil
}
