package statusbar

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Color is a struct that contains the foreground and background colors of the statusbar.
type Color struct {
	Background string
	Foreground string
}

// Model is a struct that contains all the properties of the statusbar.
type Model struct {
	Width               int
	Height              int
	FirstColumnContent  string
	SecondColumnContent string
	ThirdColumnContent  string
	FourthColumnContent string
	FirstColumnColors   Color
	SecondColumnColors  Color
	ThirdColumnColors   Color
	FourthColumnColors  Color
}

// NewModel creates an instance of a statusbar.
func NewModel(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color) Model {
	return Model{
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
	}
}

// SetContent sets the content of the statusbar.
func (m *Model) SetContent(firstColumnContent, secondColumnContent, thirdColumnContent, fourthColumnContent string) {
	m.FirstColumnContent = firstColumnContent
	m.SecondColumnContent = secondColumnContent
	m.ThirdColumnContent = thirdColumnContent
	m.FourthColumnContent = fourthColumnContent
}

// SetSize sets the size of the statusbar, useful when the terminal is resized.
func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}

// View returns a string representation of the statusbar.
func (m Model) View() string {
	width := lipgloss.Width

	firstColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FirstColumnColors.Foreground)).
		Background(lipgloss.Color(m.FirstColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Render(truncate.StringWithTail(m.FirstColumnContent, 30, "..."))

	thirdColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.ThirdColumnColors.Foreground)).
		Background(lipgloss.Color(m.ThirdColumnColors.Background)).
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(m.Height).
		Render(m.ThirdColumnContent)

	fourthColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FourthColumnColors.Foreground)).
		Background(lipgloss.Color(m.FourthColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Render(m.FourthColumnContent)

	// Second column of the status bar displayed in the center with configurable
	// foreground and background colors and some padding. Also calculate the
	// width of the other three columns so that this one can take up the rest of the space.
	secondColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.SecondColumnColors.Foreground)).
		Background(lipgloss.Color(m.SecondColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Width(m.Width - width(firstColumn) - width(thirdColumn) - width(fourthColumn)).
		Render(truncate.StringWithTail(m.SecondColumnContent, uint(m.Width-width(firstColumn)-width(thirdColumn)-width(fourthColumn)-3), "..."))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		firstColumn,
		secondColumn,
		thirdColumn,
		fourthColumn,
	)
}
