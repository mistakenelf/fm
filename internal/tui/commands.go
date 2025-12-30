package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/fm/filesystem"
	"github.com/mistakenelf/fm/polish"
)

type statusMessageTimeoutMsg struct{}

func (m *model) openFileCmd() tea.Cmd {
	selectedFile := m.filetree.GetSelectedItem()

	if !selectedFile.IsDirectory {
		m.resetViewports()

		isBinary, err := filesystem.IsBinary(selectedFile.Path)
		if err != nil {
			return m.newStatusMessageCmd(lipgloss.NewStyle().
				Foreground(polish.Colors.Red600).
				Bold(true).
				Render("Selected file type is not supported"))
		}

		switch {
		case selectedFile.Extension == ".csv":
			m.state = showCsvState

			return m.csv.SetFileNameCmd(selectedFile.Path)
		case selectedFile.Extension == ".png" || selectedFile.Extension == ".jpg" || selectedFile.Extension == ".jpeg" || selectedFile.Extension == ".gif":
			m.state = showImageState

			return m.image.SetFileNameCmd(selectedFile.Path)
		case selectedFile.Extension == ".md" && m.config.PrettyMarkdown:
			m.state = showMarkdownState

			return m.markdown.SetFileNameCmd(selectedFile.Path)
		case selectedFile.Extension == ".pdf":
			m.state = showPdfState

			return m.pdf.SetFileNameCmd(selectedFile.Path)
		case isBinary:
			return m.newStatusMessageCmd(
				lipgloss.NewStyle().
					Foreground(polish.Colors.Red600).
					Bold(true).
					Render("Selected file cannot be opened"))
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
