package text

import (
	"log"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Width   int
	Content string
}

func NewModel(width int, content string) Model {
	return Model{
		Width:   width,
		Content: content,
	}
}

func (m Model) View() string {
	bg := "light"
	if lipgloss.HasDarkBackground() {
		bg = "dark"
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithWordWrap(m.Width),
		glamour.WithStandardStyle(bg),
	)

	out, err := r.Render(m.Content)
	if err != nil {
		log.Fatal(err)
	}

	return out
}
