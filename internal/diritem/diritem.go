package diritem

import "io/fs"

type Model struct {
	Item fs.DirEntry
}

func (m *Model) SetContent(dirItem fs.DirEntry) {
	m.Item = dirItem
}

func (m Model) View() string {
	return ""
}
