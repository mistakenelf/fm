package tui

import (
	tea "github.com/charmbracelet/bubbletea"
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (m *model) openFile() tea.Cmd {
	selectedFile := m.filetree.GetSelectedItem()

	if !selectedFile.IsDirectory {
		m.resetViewports()

		switch {
		case selectedFile.Extension == ".png" || selectedFile.Extension == ".jpg" || selectedFile.Extension == ".jpeg":
			m.state = showImageState

			return m.image.SetFileName(selectedFile.Name)
		case selectedFile.Extension == ".md" && m.config.PrettyMarkdown:
			m.state = showMarkdownState

			return m.markdown.SetFileName(selectedFile.Name)
		case selectedFile.Extension == ".pdf":
			m.state = showPdfState

			return m.pdf.SetFileName(selectedFile.Name)
		case contains(forbiddenExtensions, selectedFile.Extension):
			//TODO: Display an error status message in the statusbar.
			return nil
		default:
			m.state = showCodeState

			return m.code.SetFileName(selectedFile.Name)
		}
	}

	return nil
}

func (m *model) disableAllViewports() {
	m.code.SetViewportDisabled(true)
	m.pdf.SetViewportDisabled(true)
	m.markdown.SetViewportDisabled(true)
	m.help.SetViewportDisabled(true)
	m.image.SetViewportDisabled(true)
}

func (m *model) resetViewports() {
	m.code.GotoTop()
	m.pdf.GotoTop()
	m.markdown.GotoTop()
	m.help.GotoTop()
	m.image.GotoTop()
}
