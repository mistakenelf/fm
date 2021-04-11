package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) Init() tea.Cmd {
	return m.getInitialDirectoryListing
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()

		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true

			return m, tea.Quit
		}
	}

	switch msg := msg.(type) {
	case fileStatus:
		m.Files = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.Files)-1 {
				m.Cursor++
			}

		case "enter":
			if m.Files[m.Cursor].IsDir() {
				m.Files = getUpdatedDirectoryListing(m.Files[m.Cursor].Name())
			} else {
				dat, err := ioutil.ReadFile(m.Files[m.Cursor].Name())

				if err != nil {
					log.Fatal("Error occured reading file")
				}

				m.FileContent = string(dat)
			}
			m.Cursor = 0
			m.Selected[m.Cursor] = struct{}{}
		}

	default:
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	doc := strings.Builder{}

	historyStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("#FAFAFA")).
		Width(50)

	files := ""

	for i, file := range m.Files {
		files += fmt.Sprintf("%s\n", checkbox(file.Name(), m.Cursor == i))
	}

	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		historyStyle.Copy().Align(lipgloss.Left).Render(files),
		historyStyle.Copy().Align(lipgloss.Center).Render(m.FileContent),
	))

	return doc.String()
}
