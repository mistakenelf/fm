package components

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/icons"

	"github.com/charmbracelet/lipgloss"
)

func DirTree(files []fs.FileInfo, cursor, width int) string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range files {
		curFiles += fmt.Sprintf("%s\n", DirItem(cursor == i, file.IsDir(), file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode())))
	}

	doc.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Left).
			Render(curFiles),
	))

	return doc.String()
}
