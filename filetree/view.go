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
			if file.IsDirectory {
				fileList.WriteString(inactiveStyle.Render("🗀 "))
			} else {
				fileList.WriteString(inactiveStyle.Render("🗎 "))
			}

			fileList.WriteString(inactiveStyle.Render(file.Name) + "\n")
		case i == m.Cursor && !m.Disabled:
			if file.IsDirectory {
				fileList.WriteString(selectedItemStyle.Render("🗀 "))
			} else {
				fileList.WriteString(selectedItemStyle.Render("🗎 "))
			}

			fileList.WriteString(selectedItemStyle.Render(file.Name) + "\n")
		case i != m.Cursor && !m.Disabled:
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
