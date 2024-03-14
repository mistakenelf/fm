package filetree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var fileList strings.Builder

	for i, file := range m.files {
		if i < m.min || i > m.max {
			continue
		}

		switch {
		case m.Disabled:
			if m.showIcons {
				if file.IsDirectory {
					fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.inactiveItemColor).Render("🗀 "))
				} else {
					fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.inactiveItemColor).Render("🗎 "))
				}
			}

			fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.inactiveItemColor).Render(file.Name) + "\n")
		case i == m.Cursor && !m.Disabled:
			if m.showIcons {
				if file.IsDirectory {
					fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.selectedItemColor).Render("🗀 "))
				} else {
					fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.selectedItemColor).Render("🗎 "))
				}
			}

			fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.selectedItemColor).Render(file.Name) + "\n")
		case i != m.Cursor && !m.Disabled:
			if m.showIcons {
				if file.IsDirectory {
					fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.unselectedItemColor).Render("🗀 "))
				} else {
					fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.unselectedItemColor).Render("🗎 "))
				}
			}

			fileList.WriteString(lipgloss.NewStyle().Bold(true).Foreground(m.unselectedItemColor).Render(file.Name) + "\n")
		}
	}

	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(fileList.String())
}
