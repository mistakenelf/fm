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

		if i == m.cursor {
			fileList.WriteString(selectedItemStyle.Render(file.name) + "\n")
			// fileList.WriteString(selectedItemStyle.Render(file.details) + "\n\n")
		} else {
			fileList.WriteString(file.name + "\n")
			// fileList.WriteString(file.details + "\n\n")
		}
	}

	for i := lipgloss.Height(fileList.String()); i <= m.height; i++ {
		fileList.WriteRune('\n')
	}

	return fileList.String()
}
