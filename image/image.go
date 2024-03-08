// Package image provides an image bubble which can render
// images as strings.
package image

import (
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
)

type convertImageToStringMsg string
type errorMsg error

// ToString converts an image to a string representation of an image.
func ToString(width int, img image.Image) string {
	img = imaging.Resize(img, width, 0, imaging.Lanczos)
	b := img.Bounds()
	imageWidth := b.Max.X
	h := b.Max.Y
	str := strings.Builder{}

	for heightCounter := 0; heightCounter < h; heightCounter += 2 {
		for x := imageWidth; x < width; x += 2 {
			str.WriteString(" ")
		}

		for x := 0; x < imageWidth; x++ {
			c1, _ := colorful.MakeColor(img.At(x, heightCounter))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
			color2 := lipgloss.Color(c2.Hex())
			str.WriteString(lipgloss.NewStyle().Foreground(color1).
				Background(color2).Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
}

// convertImageToStringCmd redraws the image based on the width provided.
func convertImageToStringCmd(width int, filename string) tea.Cmd {
	return func() tea.Msg {
		imageContent, err := os.Open(filepath.Clean(filename))
		if err != nil {
			return errorMsg(err)
		}

		img, _, err := image.Decode(imageContent)
		if err != nil {
			return errorMsg(err)
		}

		imageString := ToString(width, img)

		return convertImageToStringMsg(imageString)
	}
}

// Model represents the properties of a code bubble.
type Model struct {
	Viewport    viewport.Model
	Active      bool
	Borderless  bool
	FileName    string
	ImageString string
}

// New creates a new instance of code.
func New(active, borderless bool, borderColor lipgloss.AdaptiveColor) Model {
	viewPort := viewport.New(0, 0)

	return Model{
		Viewport: viewPort,
		Active:   active,
	}
}

// Init initializes the code bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetFileName sets current file to highlight, this
// returns a cmd which will highlight the text.
func (m *Model) SetFileName(filename string) tea.Cmd {
	m.FileName = filename

	return convertImageToStringCmd(m.Viewport.Width, filename)
}

// SetSize sets the size of the bubble.
func (m *Model) SetSize(w, h int) tea.Cmd {
	m.Viewport.Width = w
	m.Viewport.Height = h

	if m.FileName != "" {
		return convertImageToStringCmd(m.Viewport.Width, m.FileName)
	}

	return nil
}

// SetIsActive sets if the bubble is currently active
func (m *Model) SetIsActive(active bool) {
	m.Active = active
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// SetBorderless sets weather or not to show the border.
func (m *Model) SetBorderless(borderless bool) {
	m.Borderless = borderless
}

// Update handles updating the UI of a code bubble.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case convertImageToStringMsg:
		m.ImageString = lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render(string(msg))

		m.Viewport.SetContent(m.ImageString)

		return m, nil
	case errorMsg:
		m.FileName = ""
		m.ImageString = lipgloss.NewStyle().
			Width(m.Viewport.Width).
			Height(m.Viewport.Height).
			Render("Error: " + msg.Error())

		m.Viewport.SetContent(m.ImageString)

		return m, nil
	}

	if m.Active {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the code bubble.
func (m Model) View() string {
	return m.Viewport.View()
}
