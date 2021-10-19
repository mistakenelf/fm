package text

import (
	"bytes"

	"github.com/knipferrc/fm/strfmt"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the properties of text.
type Model struct {
	Content string
	Width   int
}

// Highlight returns a syntax highlighted string of text.
func Highlight(content, extension, syntaxTheme string) (string, error) {
	buf := new(bytes.Buffer)
	if err := quick.Highlight(buf, content, extension, "terminal256", syntaxTheme); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SetSize sets the size of the text.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// SetContent sets the content of the text.
func (m *Model) SetContent(content string) {
	m.Content = content
}

// GetContent returns the content of the text.
func (m Model) GetContent() string {
	return m.Content
}

// View returns a string representation of text.
func (m *Model) View() string {
	return lipgloss.NewStyle().
		Width(m.Width).
		Render(strfmt.ConvertTabsToSpaces(m.Content))
}
