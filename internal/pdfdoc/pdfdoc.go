package pdfdoc

import (
	"bytes"

	"github.com/knipferrc/fm/strfmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ledongthuc/pdf"
)

// Model represents the properties of text.
type Model struct {
	Content string
	Width   int
}

// ReadPdf reads a PDF file given a name.
func ReadPdf(name string) (string, error) {
	f, r, err := pdf.Open(name)
	if err != nil {
		return "", err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	b, err := r.GetPlainText()

	if err != nil {
		return "", err
	}

	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SetSize sets the size of the pdfdoc.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// SetContent sets the content of the pdfdoc.
func (m *Model) SetContent(content string) {
	m.Content = content
}

// GetContent returns the content of the pdfdoc.
func (m Model) GetContent() string {
	return m.Content
}

// View returns a string representation of pdfdoc.
func (m *Model) View() string {
	return lipgloss.NewStyle().
		Width(m.Width).
		Render(strfmt.ConvertTabsToSpaces(m.Content))
}
