package text

import (
	"bytes"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the properties of text.
type Model struct {
	Content string
	Width   int
}

// Highlight returns a highlighted string of text.
func Highlight(content, extension string) (string, error) {
	buf := new(bytes.Buffer)
	if err := quick.Highlight(buf, content, extension, "terminal256", "dracula"); err != nil {
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

// View returns a string representation of text.
func (m *Model) View() string {
	return lipgloss.NewStyle().Width(m.Width).Render(m.Content)
}
