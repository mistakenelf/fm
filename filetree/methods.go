package filetree

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	thousand    = 1000
	ten         = 10
	fivePercent = 0.0499
)

// ConvertBytesToSizeString converts a byte count to a human readable string.
func ConvertBytesToSizeString(size int64) string {
	if size < thousand {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / thousand
	for _, s := range suffix {
		if curr < ten {
			return fmt.Sprintf("%.1f%s", curr-fivePercent, s)
		} else if curr < thousand {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= thousand
	}

	return ""
}

// SetDisabled sets if the bubble is currently active.
func (m *Model) SetDisabled(disabled bool) {
	m.Disabled = disabled
}

// GetSelectedItem returns the currently selected file/dir.
func (m Model) GetSelectedItem() DirectoryItem {
	if len(m.files) > 0 {
		return m.files[m.Cursor]
	}

	return DirectoryItem{}
}

// GetTotalItems returns total number of tree items.
func (m Model) GetTotalItems() int {
	return len(m.files)
}

// SetSize Sets the size of the filetree.
func (m *Model) SetSize(width, height int) {
	m.height = height
	m.width = width
	m.max = m.height - 1
}

// SetTheme sets the theme of the tree.
func (m *Model) SetTheme(selectedItemColor, unselectedItemColor lipgloss.AdaptiveColor) {
	m.selectedItemColor = selectedItemColor
	m.unselectedItemColor = unselectedItemColor
}

// SetSelectionPath sets the selection path to be written.
func (m *Model) SetSelectionPath(path string) {
	m.selectionPath = path
}

// SetShowIcons sets whether icons will show or not.
func (m *Model) SetShowIcons(show bool) {
	m.showIcons = show
}
