package statusbar

import (
	"github.com/charmbracelet/lipgloss"
)

type Color struct {
	Background string
	Foreground string
}

type Model struct {
	Width               int
	FirstColumnContent  string
	SecondColumnContent string
	ThirdColumnContent  string
	FourthColumnContent string
	FirstColumnColors   Color
	SecondColumnColors  Color
	ThirdColumnColors   Color
	FourthColumnColors  Color
}

func NewModel(
	width int,
	firstColumnContent, secondColumnContent, thirdColumnContent, fourthColumnContent string,
	firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color,
) Model {
	return Model{
		Width:               width,
		FirstColumnContent:  firstColumnContent,
		SecondColumnContent: secondColumnContent,
		ThirdColumnContent:  thirdColumnContent,
		FourthColumnContent: fourthColumnContent,
		FirstColumnColors:   firstColumnColors,
		SecondColumnColors:  secondColumnColors,
		ThirdColumnColors:   thirdColumnColors,
		FourthColumnColors:  fourthColumnColors,
	}
}

func (m *Model) SetContent(firstColumnContent, secondColumnContent, thirdColumnContent, fourthColumnContent string) {
	m.FirstColumnContent = firstColumnContent
	m.SecondColumnContent = secondColumnContent
	m.ThirdColumnContent = thirdColumnContent
	m.FourthColumnContent = fourthColumnContent
}

func (m *Model) SetSize(width int) {
	m.Width = width
}

func (m Model) View() string {
	width := lipgloss.Width

	firstColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FirstColumnColors.Foreground)).
		Background(lipgloss.Color(m.FirstColumnColors.Background)).
		Padding(0, 1).
		Render(m.FirstColumnContent)

	thirdColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.ThirdColumnColors.Foreground)).
		Background(lipgloss.Color(m.ThirdColumnColors.Background)).
		Align(lipgloss.Right).
		Padding(0, 1).
		Render(m.ThirdColumnContent)

	fourthColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FourthColumnColors.Foreground)).
		Background(lipgloss.Color(m.FourthColumnColors.Background)).
		Padding(0, 1).
		Render(m.FourthColumnContent)

	secondColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.SecondColumnColors.Foreground)).
		Background(lipgloss.Color(m.SecondColumnColors.Background)).
		Padding(0, 1).
		Width(m.Width - width(firstColumn) - width(thirdColumn) - width(fourthColumn)).
		Render(m.SecondColumnContent)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		firstColumn,
		secondColumn,
		thirdColumn,
		fourthColumn,
	)
}
