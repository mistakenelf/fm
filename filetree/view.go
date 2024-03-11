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

		if i == m.Cursor {
			if file.IsDirectory {
				fileList.WriteString(selectedItemStyle.Render("🗀 "))
			} else {
				fileList.WriteString(selectedItemStyle.Render("🗎 "))
			}

			fileList.WriteString(selectedItemStyle.Render(file.Name) + "\n")
		} else {
			if file.IsDirectory {
				fileList.WriteString(unselectedItemStyle.Render("🗀 "))
			} else {
				fileList.WriteString(unselectedItemStyle.Render("🗎 "))
			}

			fileList.WriteString(unselectedItemStyle.Render(file.Name) + "\n")
		}
	}

	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(fileList.String())
}
