package markdown

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Model contains the properties of markdown.
type Model struct {
	Content string
}

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

// SetContent sets the markdown content.
func (m *Model) SetContent(content string) {
	m.Content = content
}

// View returns a string representation of the markdown.
func (m Model) View() string {
	return m.Content
}
