// Package image provides an image bubble which can render
// images as strings.
package image

import (
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/mistakenelf/fm/polish"
)

type convertImageToStringMsg string
type errorMsg string
type statusMessageTimeoutMsg struct{}

// Model represents the properties of a image bubble.
type Model struct {
	Viewport              viewport.Model
	ViewportDisabled      bool
	FileName              string
	ImageString           string
	StatusMessage         string
	StatusMessageLifetime time.Duration
	statusMessageTimer    *time.Timer
}

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

// NewStatusMessage sets a new status message, which will show for a limited
// amount of time.
func (m *Model) NewStatusMessageCmd(s string) tea.Cmd {
	m.StatusMessage = s

	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}

	m.statusMessageTimer = time.NewTimer(m.StatusMessageLifetime)

	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}

func convertImageToStringCmd(width int, filename string) tea.Cmd {
	return func() tea.Msg {
		imageContent, err := os.Open(filepath.Clean(filename))
		if err != nil {
			return errorMsg(err.Error())
		}

		img, _, err := image.Decode(imageContent)
		if err != nil {
			return errorMsg(err.Error())
		}

		imageString := ToString(width, img)

		return convertImageToStringMsg(imageString)
	}
}

// SetFileName sets the image file and converts it to a string.
func (m *Model) SetFileNameCmd(filename string) tea.Cmd {
	m.FileName = filename

	return convertImageToStringCmd(m.Viewport.Width, filename)
}

// SetSize sets the size of the bubble.
func (m *Model) SetSizeCmd(w, h int) tea.Cmd {
	m.Viewport.Width = w
	m.Viewport.Height = h

	if m.FileName != "" {
		return convertImageToStringCmd(m.Viewport.Width, m.FileName)
	}

	return nil
}

// New creates a new instance of an image.
func New() Model {
	viewPort := viewport.New(0, 0)

	return Model{
		Viewport:              viewPort,
		ViewportDisabled:      false,
		StatusMessageLifetime: time.Second,
	}
}

// Init initializes the image bubble.
func (m Model) Init() tea.Cmd {
	return nil
}

// SetViewportDisabled toggles the state of the viewport.
func (m *Model) SetViewportDisabled(disabled bool) {
	m.ViewportDisabled = disabled
}

// GotoTop jumps to the top of the viewport.
func (m *Model) GotoTop() {
	m.Viewport.GotoTop()
}

// GotoBottom jumps to the bottom of the viewport.
func (m *Model) GotoBottom() {
	m.Viewport.GotoBottom()
}

// Update handles updating the UI of the image bubble.
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
		cmds = append(cmds, m.NewStatusMessageCmd(
			lipgloss.NewStyle().
				Foreground(polish.Colors.Red600).
				Bold(true).
				Render(string(msg)),
		))
	}

	if !m.ViewportDisabled {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View returns a string representation of the image bubble.
func (m Model) View() string {
	return m.Viewport.View()
}
