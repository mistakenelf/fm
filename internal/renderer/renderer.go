package renderer

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/knipferrc/fm/internal/statusbar"
	"github.com/ledongthuc/pdf"
	"github.com/lucasb-eyer/go-colorful"
)

// Model is a struct that contains all the properties of renderer.
type Model struct {
	Image               image.Image
	Style               lipgloss.Style
	Viewport            viewport.Model
	ActiveBorderColor   lipgloss.AdaptiveColor
	InactiveBorderColor lipgloss.AdaptiveColor
	Borderless          bool
	IsActive            bool
}

// NewModel creates a new instance of a renderer.
func NewModel(borderless, isActive bool, activeBorderColor, inactiveBorderColor lipgloss.AdaptiveColor) Model {
	border := lipgloss.NormalBorder()
	padding := 1

	if borderless {
		border = lipgloss.HiddenBorder()
	}

	style := lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border)

	return Model{
		Borderless:          borderless,
		IsActive:            isActive,
		ActiveBorderColor:   activeBorderColor,
		InactiveBorderColor: inactiveBorderColor,
		Style:               style,
	}
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
func (m *Model) SetSize(width, height int) {
	m.Viewport.Width = (width / 2) - m.Style.GetHorizontalBorderSize()
	m.Viewport.Height = height - m.Style.GetVerticalBorderSize() - statusbar.StatusbarHeight
}

// GetImage returns the currently set image.
func (m Model) GetImage() image.Image {
	return m.Image
}

// SetContent sets the content of the renderer.
func (m *Model) SetContent(content string) {
	m.Viewport.SetContent(content)
}

// SetImage sets the image of the renderer.
func (m *Model) SetImage(img image.Image) {
	m.Image = img
}

// GetWidth returns the width of the renderer.
func (m Model) GetWidth() int {
	return m.Viewport.Width
}

// GetIsActive returns the active state of the renderer.
func (m Model) GetIsActive() bool {
	return m.IsActive
}

// SetIsActive sets the active state of the renderer.
func (m *Model) SetIsActive(isActive bool) {
	m.IsActive = isActive
}

// Update updates the statusbar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		m.SetContent("Welcome to FM")
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			if m.IsActive {
				m.Viewport.LineDown(1)
			}
		case "up", "k":
			if m.IsActive {
				m.Viewport.LineUp(1)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of a renderer.
func (m Model) View() string {
	borderColor := m.InactiveBorderColor
	border := lipgloss.NormalBorder()
	padding := 1

	if m.Borderless {
		border = lipgloss.HiddenBorder()
	}

	if m.IsActive {
		borderColor = m.ActiveBorderColor
	}

	return m.Style.Copy().
		BorderForeground(borderColor).
		PaddingLeft(padding).
		PaddingRight(padding).
		Border(border).
		Width(m.Viewport.Width).
		Height(m.Viewport.Height).
		Render(m.Viewport.View())
}
