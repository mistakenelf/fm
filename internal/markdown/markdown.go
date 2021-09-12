package markdown

import (
	"github.com/knipferrc/fm/formatter"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Model contains the properties of markdown.
type Model struct {
	Content string
	Width   int
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

// GetContent returns the markdown content.
func (m Model) GetContent() string {
	return m.Content
}

// SetSize sets the width of the markdown.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// View returns a string representation of the markdown.
func (m Model) View() string {
	return lipgloss.NewStyle().Width(m.Width).Render(
		formatter.ConvertTabsToSpaces(m.Content),
	)
}
