package filetree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/mistakenelf/fm/icons"
)

func (m Model) View() string {
	var fileList strings.Builder

	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}

	for i, file := range m.files {
		if i < m.min || i > m.max {
			continue
		}

		switch {
		case m.Disabled:
			fallthrough
		case i == m.Cursor && !m.Disabled:
			textColor := m.inactiveItemColor

			if i == m.Cursor && !m.Disabled {
				textColor = m.selectedItemColor
			}

			if m.showIcons {
				icon := icons.GetElementIcon(file.Name, file.IsDirectory)

				fileList.WriteString(
					lipgloss.NewStyle().
						Bold(true).
						Foreground(textColor).
						Render(icon.Icon) + " ",
				)

				fileList.WriteString(
					lipgloss.NewStyle().
						Bold(true).
						Foreground(textColor).
						Render(file.Name) + "\n",
				)
			} else {
				fileList.WriteString(
					lipgloss.NewStyle().
						Bold(true).
						Foreground(textColor).
						Render(file.Name) + "\n",
				)
			}

		case i != m.Cursor && !m.Disabled:
			textColor := m.unselectedItemColor

			if m.showIcons {
				icon := icons.GetElementIcon(file.Name, file.IsDirectory)

				fileList.WriteString(
					lipgloss.NewStyle().
						Bold(true).
						Foreground(lipgloss.Color(icon.Color)).
						Render(icon.Icon) + " ",
				)

				fileList.WriteString(
					lipgloss.NewStyle().
						Bold(true).
						Foreground(textColor).
						Render(file.Name) + "\n",
				)
			} else {
				fileList.WriteString(
					lipgloss.NewStyle().
						Bold(true).
						Foreground(textColor).
						Render(file.Name) + "\n",
				)
			}
		}
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Render(fileList.String())
}
