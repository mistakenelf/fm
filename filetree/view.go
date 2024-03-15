package filetree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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
			iconColor := m.inactiveItemColor
			textColor := m.inactiveItemColor
			if i == m.Cursor && !m.Disabled {
				iconColor = m.selectedItemColor
				textColor = m.selectedItemColor
			}

			if m.showIcons {
				icon := "🗀 "
				if !file.IsDirectory {
					icon = "🗎 "
				}
				fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(iconColor).Render(icon))
			}

			fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(textColor).Render(file.Name) + "\n")
		case i != m.Cursor && !m.Disabled:
			iconColor := m.unselectedItemColor
			textColor := m.unselectedItemColor

			if m.showIcons {
				icon := "🗀 "
				if !file.IsDirectory {
					icon = "🗎 "
				}
				fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(iconColor).Render(icon))
			}

			fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(textColor).Render(file.Name) + "\n")
		}
	}

	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(fileList.String())
}
