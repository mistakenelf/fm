package tui

import (
	"fmt"

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

// contains returns true if the slice contains the string.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (m *model) updateStatusbarContent() {
	if m.filetree.GetSelectedItem().Name != "" {
		if m.filetree.StatusMessage != "" {
			m.statusbar.SetContent(
				m.filetree.GetSelectedItem().Name,
				m.filetree.StatusMessage,
				fmt.Sprintf("%d/%d", m.filetree.Cursor, m.filetree.GetTotalItems()),
				fmt.Sprintf("%s %s", "🗀", "FM"),
			)
		} else {
			m.statusbar.SetContent(
				m.filetree.GetSelectedItem().Name,
				m.filetree.GetSelectedItem().CurrentDirectory,
				fmt.Sprintf("%d/%d", m.filetree.Cursor, m.filetree.GetTotalItems()),
				fmt.Sprintf("%s %s", "🗀", "FM"),
			)
		}
	}
}

// togglePane toggles between the two panes.
func (m *model) togglePane() {
	m.activeBox = (m.activeBox + 1) % 2

	if m.activeBox == 0 {
		m.deactivateAllBubbles()
		m.filetree.SetIsActive(true)
	} else {
		switch m.state {
		case idleState:
			m.deactivateAllBubbles()
			m.help.SetIsActive(true)
		case showCodeState:
			m.deactivateAllBubbles()
			m.code.SetIsActive(true)
		case showImageState:
			m.deactivateAllBubbles()
			m.image.SetIsActive(true)
		case showMarkdownState:
			m.deactivateAllBubbles()
			m.markdown.SetIsActive(true)
		case showPdfState:
			m.deactivateAllBubbles()
			m.markdown.SetIsActive(true)
		}
	}
}

// openFile opens the currently selected file.
func (m *model) openFile() []tea.Cmd {
	var cmds []tea.Cmd

	selectedFile := m.filetree.GetSelectedItem()
	if !selectedFile.IsDirectory {
		m.resetViewports()

		switch {
		case selectedFile.Extension == ".png" || selectedFile.Extension == ".jpg" || selectedFile.Extension == ".jpeg":
			m.state = showImageState
			readFileCmd := m.image.SetFileName(selectedFile.Name)
			cmds = append(cmds, readFileCmd)
		case selectedFile.Extension == ".md" && m.config.PrettyMarkdown:
			m.state = showMarkdownState
			markdownCmd := m.markdown.SetFileName(selectedFile.Name)
			cmds = append(cmds, markdownCmd)
		case selectedFile.Extension == ".pdf":
			m.state = showPdfState
			pdfCmd := m.pdf.SetFileName(selectedFile.Name)
			cmds = append(cmds, pdfCmd)
		case contains(forbiddenExtensions, selectedFile.Extension):
			return nil
		default:
			m.state = showCodeState
			readFileCmd := m.code.SetFileName(selectedFile.Name)
			cmds = append(cmds, readFileCmd)
		}
	}

	return cmds
}

// resetViewports goes to the top of all bubbles viewports.
func (m *model) resetViewports() {
	m.code.GotoTop()
	m.pdf.GotoTop()
	m.markdown.GotoTop()
	m.help.GotoTop()
	m.image.GotoTop()
}

// deactivateAllBubbles sets all bubbles to inactive.
func (m *model) deactivateAllBubbles() {
	m.filetree.SetIsActive(false)
	m.code.SetIsActive(false)
	m.markdown.SetIsActive(false)
	m.image.SetIsActive(false)
	m.pdf.SetIsActive(false)
	m.help.SetIsActive(false)
}
