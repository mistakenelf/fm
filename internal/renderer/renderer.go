package renderer

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/ledongthuc/pdf"
	"github.com/lucasb-eyer/go-colorful"
)

// Model is a struct that contains all the properties of renderer.
type Model struct {
	Image   image.Image
	Content string
	Width   int
}

// ConvertTabsToSpaces converts tabs to spaces.
func ConvertTabsToSpaces(input string) string {
	return strings.ReplaceAll(input, "\t", "    ")
}

// ConvertByesToSizeString converts a byte count to a human readable string.
func ConvertBytesToSizeString(size int64) string {
	if size < 1000 {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / 1000
	for _, s := range suffix {
		if curr < 10 {
			return fmt.Sprintf("%.1f%s", curr-0.0499, s)
		} else if curr < 1000 {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= 1000
	}

	return ""
}

// ImageToString converts an image to a string representation of an image.
func ImageToString(width int, img image.Image) string {
	img = imaging.Resize(img, width, 0, imaging.Lanczos)
	b := img.Bounds()
	w := b.Max.X
	h := b.Max.Y
	str := strings.Builder{}

	for y := 0; y < h; y += 2 {
		for x := w; x < width; x += 2 {
			str.WriteString(" ")
		}

		for x := 0; x < w; x++ {
			c1, _ := colorful.MakeColor(img.At(x, y))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, y+1))
			color2 := lipgloss.Color(c2.Hex())
			str.WriteString(lipgloss.NewStyle().Foreground(color1).
				Background(color2).Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
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

// Highlight returns a syntax highlighted string of text.
func Highlight(content, extension, syntaxTheme string) (string, error) {
	buf := new(bytes.Buffer)
	if err := quick.Highlight(buf, content, extension, "terminal256", syntaxTheme); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SetSize sets the size of the renderer.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// GetImage returns the currently set image.
func (m Model) GetImage() image.Image {
	return m.Image
}

// SetContent sets the content of the renderer.
func (m *Model) SetContent(content string) {
	m.Content = content
}

// SetImage sets the image of the renderer.
func (m *Model) SetImage(img image.Image) {
	m.Image = img
}

// GetWidth returns the width of the renderer.
func (m Model) GetWidth() int {
	return m.Width
}

// GetContent returns the content of the renderer.
func (m Model) GetContent() string {
	return m.Content
}

// View returns a string representation of a renderer.
func (m Model) View() string {
	return lipgloss.NewStyle().Width(m.Width).Render(m.Content)
}
