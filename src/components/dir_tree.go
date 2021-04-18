package components

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func DirTree(files []fs.FileInfo, cursor, width int) string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range files {
		curFiles += fmt.Sprintf("%s\n", DirItem(file.Name(), cursor == i, file.IsDir(), filepath.Ext(file.Name())))
	}

	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Width(width).
			Align(lipgloss.Left).
			Render(curFiles),
	))

	return doc.String()
}
