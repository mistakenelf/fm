package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var forbiddenExtensions = []string{
	".FCStd",
	".gif",
	".zip",
	".rar",
	".webm",
	".sqlite",
	".sqlite-shm",
	".sqlite-wal",
	".DS_Store",
	".db",
	".data",
	".plist",
	".webp",
	".img",
}

type statusMessageTimeoutMsg struct{}

func (m *model) openFileCmd() tea.Cmd {
	selectedFile := m.filetree.GetSelectedItem()

	if !selectedFile.IsDirectory {
		m.resetViewports()

		switch {
		case selectedFile.Extension == ".png" || selectedFile.Extension == ".jpg" || selectedFile.Extension == ".jpeg":
			m.state = showImageState

			return m.image.SetFileNameCmd(selectedFile.Path)
		case selectedFile.Extension == ".md" && m.config.PrettyMarkdown:
			m.state = showMarkdownState

			return m.markdown.SetFileNameCmd(selectedFile.Path)
		case selectedFile.Extension == ".pdf":
			m.state = showPdfState

			return m.pdf.SetFileNameCmd(selectedFile.Path)
		case contains(forbiddenExtensions, selectedFile.Extension):
			return m.newStatusMessageCmd(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#cc241d")).
				Bold(true).
				Render("Selected file type is not supported"))
		default:
			m.state = showCodeState

			return m.code.SetFileNameCmd(selectedFile.Path)
		}
	}

	return nil
}

// newStatusMessage sets a new status message, which will show for a limited
// amount of time.
func (m *model) newStatusMessageCmd(s string) tea.Cmd {
	m.statusMessage = s

	if m.statusMessageTimer != nil {
		m.statusMessageTimer.Stop()
	}

	m.statusMessageTimer = time.NewTimer(m.statusMessageLifetime)

	return func() tea.Msg {
		<-m.statusMessageTimer.C
		return statusMessageTimeoutMsg{}
	}
}
