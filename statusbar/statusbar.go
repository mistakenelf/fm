// Package statusbar provides an statusbar bubble which can render
// four different status sections
package statusbar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Height represents the height of the statusbar.
const Height = 1

// ColorConfig
type ColorConfig struct {
	Foreground lipgloss.AdaptiveColor
	Background lipgloss.AdaptiveColor
}

// Model represents the properties of the statusbar.
type Model struct {
	Width              int
	Height             int
	FirstColumn        string
	SecondColumn       string
	ThirdColumn        string
	FourthColumn       string
	FirstColumnColors  ColorConfig
	SecondColumnColors ColorConfig
	ThirdColumnColors  ColorConfig
	FourthColumnColors ColorConfig
}

// New creates a new instance of the statusbar.
func New(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors ColorConfig) Model {
	return Model{
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
	}
}

// SetSize sets the width of the statusbar.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// Update updates the size of the statusbar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width)
	}

	return m, nil
}

// SetContent sets the content of the statusbar.
func (m *Model) SetContent(firstColumn, secondColumn, thirdColumn, fourthColumn string) {
	m.FirstColumn = firstColumn
	m.SecondColumn = secondColumn
	m.ThirdColumn = thirdColumn
	m.FourthColumn = fourthColumn
}

// SetColors sets the colors of the 4 columns.
func (m *Model) SetColors(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors ColorConfig) {
	m.FirstColumnColors = firstColumnColors
	m.SecondColumnColors = secondColumnColors
	m.ThirdColumnColors = thirdColumnColors
	m.FourthColumnColors = fourthColumnColors
}

// View returns a string representation of a statusbar.
func (m Model) View() string {
	width := lipgloss.Width

	firstColumn := lipgloss.NewStyle().
		Foreground(m.FirstColumnColors.Foreground).
		Background(m.FirstColumnColors.Background).
		Padding(0, 1).
		Height(Height).
		Render(truncate.StringWithTail(m.FirstColumn, 30, "..."))

	thirdColumn := lipgloss.NewStyle().
		Foreground(m.ThirdColumnColors.Foreground).
		Background(m.ThirdColumnColors.Background).
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(Height).
		Render(m.ThirdColumn)

	fourthColumn := lipgloss.NewStyle().
		Foreground(m.FourthColumnColors.Foreground).
		Background(m.FourthColumnColors.Background).
		Padding(0, 1).
		Height(Height).
		Render(m.FourthColumn)

	secondColumn := lipgloss.NewStyle().
		Foreground(m.SecondColumnColors.Foreground).
		Background(m.SecondColumnColors.Background).
		Padding(0, 1).
		Height(Height).
		Width(m.Width - width(firstColumn) - width(thirdColumn) - width(fourthColumn)).
		Render(truncate.StringWithTail(
			m.SecondColumn,
			uint(m.Width-width(firstColumn)-width(thirdColumn)-width(fourthColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		firstColumn,
		secondColumn,
		thirdColumn,
		fourthColumn,
	)
}
