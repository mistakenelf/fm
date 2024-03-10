package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mistakenelf/fm/statusbar"
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

// resetViewports goes to the top of all bubbles viewports.
func (m *model) resetViewports() {
	m.code.GotoTop()
	m.pdf.GotoTop()
	m.markdown.GotoTop()
	m.help.GotoTop()
	m.image.GotoTop()
}

// deactivateALlBubbles sets all bubbles to inactive.
func (m *model) deactivateAllBubbles() {
	m.filetree.SetIsActive(false)
	m.code.SetIsActive(false)
	m.markdown.SetIsActive(false)
	m.image.SetIsActive(false)
	m.pdf.SetIsActive(false)
	m.help.SetIsActive(false)
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
		case selectedFile.Extension == ".md" && m.config.Settings.PrettyMarkdown:
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

// togglePane toggles between the two boxes.
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

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.filetree, cmd = m.filetree.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		halfSize := msg.Width / 2
		bubbleHeight := msg.Height - statusbar.Height - 2

		cmds = append(cmds, m.image.SetSize(halfSize, bubbleHeight))
		cmds = append(cmds, m.markdown.SetSize(halfSize, bubbleHeight))

		m.filetree.SetSize(halfSize, bubbleHeight)
		m.help.SetSize(halfSize, bubbleHeight)
		m.code.SetSize(halfSize, bubbleHeight)
		m.pdf.SetSize(halfSize, bubbleHeight)
		m.statusbar.SetSize(msg.Width)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.OpenFile):
			cmds = append(cmds, tea.Batch(m.openFile()...))
		case key.Matches(msg, m.keys.TogglePane):
			m.togglePane()
		}
	}

	if m.filetree.GetSelectedItem().Name != "" {
		m.statusbar.SetContent(
			m.filetree.GetSelectedItem().Name,
			m.filetree.GetSelectedItem().CurrentDirectory,
			fmt.Sprintf("%d/%d", m.filetree.Cursor, m.filetree.GetTotalItems()),
			fmt.Sprintf("%s %s", "🗀", "FM"),
		)
	}

	m.code, cmd = m.code.Update(msg)
	cmds = append(cmds, cmd)

	m.markdown, cmd = m.markdown.Update(msg)
	cmds = append(cmds, cmd)

	m.image, cmd = m.image.Update(msg)
	cmds = append(cmds, cmd)

	m.pdf, cmd = m.pdf.Update(msg)
	cmds = append(cmds, cmd)

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
