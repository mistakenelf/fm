package tui

import (
	"fmt"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/theme"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/help"
	"github.com/knipferrc/teacup/icons"
	"github.com/knipferrc/teacup/statusbar"
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

// resetBorderColors resets all bubble border colors to default.
func (m *model) resetBorderColors() {
	m.filetree.SetBorderColor(m.theme.InactiveBoxBorderColor)
	m.help.SetBorderColor(m.theme.InactiveBoxBorderColor)
	m.code.SetBorderColor(m.theme.InactiveBoxBorderColor)
	m.image.SetBorderColor(m.theme.InactiveBoxBorderColor)
	m.markdown.SetBorderColor(m.theme.InactiveBoxBorderColor)
	m.pdf.SetBorderColor(m.theme.InactiveBoxBorderColor)
}

// reloadConfig reloads the config file and updates the UI.
func (m *model) reloadConfig() []tea.Cmd {
	var cmds []tea.Cmd

	cfg, err := config.ParseConfig()
	if err != nil {
		return nil
	}

	m.config = cfg
	syntaxTheme := cfg.Theme.SyntaxTheme.Light
	if lipgloss.HasDarkBackground() {
		syntaxTheme = cfg.Theme.SyntaxTheme.Dark
	}

	m.code.SetSyntaxTheme(syntaxTheme)

	theme := theme.GetTheme(cfg.Theme.AppTheme)
	m.theme = theme
	m.statusbar.SetColors(
		statusbar.ColorConfig{
			Foreground: theme.StatusBarSelectedFileForegroundColor,
			Background: theme.StatusBarSelectedFileBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusBarBarForegroundColor,
			Background: theme.StatusBarBarBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusBarTotalFilesForegroundColor,
			Background: theme.StatusBarTotalFilesBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: theme.StatusBarLogoForegroundColor,
			Background: theme.StatusBarLogoBackgroundColor,
		},
	)

	m.help.SetTitleColor(
		help.TitleColor{
			Background: theme.TitleBackgroundColor,
			Foreground: theme.TitleForegroundColor,
		},
	)

	m.filetree.SetTitleColors(theme.TitleForegroundColor, theme.TitleBackgroundColor)
	m.filetree.SetSelectedItemColors(theme.SelectedTreeItemColor)
	cmds = append(cmds, m.filetree.ToggleShowIcons(cfg.Settings.ShowIcons))

	m.filetree.SetBorderless(cfg.Settings.Borderless)
	m.code.SetBorderless(cfg.Settings.Borderless)
	m.help.SetBorderless(cfg.Settings.Borderless)
	m.markdown.SetBorderless(cfg.Settings.Borderless)
	m.pdf.SetBorderless(cfg.Settings.Borderless)
	m.image.SetBorderless(cfg.Settings.Borderless)

	if m.activeBox == 0 {
		m.deactivateAllBubbles()
		m.filetree.SetIsActive(true)
		m.resetBorderColors()
		m.filetree.SetBorderColor(theme.ActiveBoxBorderColor)
	} else {
		switch m.state {
		case idleState:
			m.deactivateAllBubbles()
			m.help.SetIsActive(true)
			m.resetBorderColors()
			m.help.SetBorderColor(theme.ActiveBoxBorderColor)
		case showCodeState:
			m.deactivateAllBubbles()
			m.code.SetIsActive(true)
			m.resetBorderColors()
			m.code.SetBorderColor(theme.ActiveBoxBorderColor)
		case showImageState:
			m.deactivateAllBubbles()
			m.image.SetIsActive(true)
			m.resetBorderColors()
			m.image.SetBorderColor(theme.ActiveBoxBorderColor)
		case showMarkdownState:
			m.deactivateAllBubbles()
			m.markdown.SetIsActive(true)
			m.resetBorderColors()
			m.markdown.SetBorderColor(theme.ActiveBoxBorderColor)
		case showPdfState:
			m.deactivateAllBubbles()
			m.markdown.SetIsActive(true)
			m.resetBorderColors()
			m.pdf.SetBorderColor(theme.ActiveBoxBorderColor)
		}
	}

	return cmds
}

// openFile opens the currently selected file.
func (m *model) openFile() []tea.Cmd {
	var cmds []tea.Cmd

	selectedFile := m.filetree.GetSelectedItem()
	if !selectedFile.IsDirectory() {
		m.resetViewports()

		switch {
		case selectedFile.FileExtension() == ".png" || selectedFile.FileExtension() == ".jpg" || selectedFile.FileExtension() == ".jpeg":
			m.state = showImageState
			readFileCmd := m.image.SetFileName(selectedFile.FileName())
			cmds = append(cmds, readFileCmd)
		case selectedFile.FileExtension() == ".md" && m.config.Settings.PrettyMarkdown:
			m.state = showMarkdownState
			markdownCmd := m.markdown.SetFileName(selectedFile.FileName())
			cmds = append(cmds, markdownCmd)
		case selectedFile.FileExtension() == ".pdf":
			m.state = showPdfState
			pdfCmd := m.pdf.SetFileName(selectedFile.FileName())
			cmds = append(cmds, pdfCmd)
		case contains(forbiddenExtensions, selectedFile.FileExtension()):
			return nil
		default:
			m.state = showCodeState
			readFileCmd := m.code.SetFileName(selectedFile.FileName())
			cmds = append(cmds, readFileCmd)
		}
	}

	return cmds
}

// toggleBox toggles between the two boxes.
func (m *model) toggleBox() {
	m.activeBox = (m.activeBox + 1) % 2
	if m.activeBox == 0 {
		m.deactivateAllBubbles()
		m.filetree.SetIsActive(true)
		m.resetBorderColors()
		m.filetree.SetBorderColor(m.theme.ActiveBoxBorderColor)
	} else {
		switch m.state {
		case idleState:
			m.deactivateAllBubbles()
			m.help.SetIsActive(true)
			m.resetBorderColors()
			m.help.SetBorderColor(m.theme.ActiveBoxBorderColor)
		case showCodeState:
			m.deactivateAllBubbles()
			m.code.SetIsActive(true)
			m.resetBorderColors()
			m.code.SetBorderColor(m.theme.ActiveBoxBorderColor)
		case showImageState:
			m.deactivateAllBubbles()
			m.image.SetIsActive(true)
			m.resetBorderColors()
			m.image.SetBorderColor(m.theme.ActiveBoxBorderColor)
		case showMarkdownState:
			m.deactivateAllBubbles()
			m.markdown.SetIsActive(true)
			m.resetBorderColors()
			m.markdown.SetBorderColor(m.theme.ActiveBoxBorderColor)
		case showPdfState:
			m.deactivateAllBubbles()
			m.markdown.SetIsActive(true)
			m.resetBorderColors()
			m.pdf.SetBorderColor(m.theme.ActiveBoxBorderColor)
		}
	}
}

// updateStatusbar updates the content of the statusbar.
func (m *model) updateStatusbar() {
	logoText := fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM")
	if !m.config.Settings.ShowIcons {
		logoText = "FM"
	}

	m.statusbar.SetContent(
		m.filetree.GetSelectedItem().ShortName(),
		m.filetree.GetSelectedItem().CurrentDirectory(),
		fmt.Sprintf("%d/%d", m.filetree.Cursor(), m.filetree.TotalItems()),
		logoText,
	)
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
		bubbleHeight := msg.Height - statusbar.Height

		resizeImgCmd := m.image.SetSize(halfSize, bubbleHeight)
		markdownCmd := m.markdown.SetSize(halfSize, bubbleHeight)
		m.filetree.SetSize(halfSize, bubbleHeight)
		m.help.SetSize(halfSize, bubbleHeight)
		m.code.SetSize(halfSize, bubbleHeight)
		m.pdf.SetSize(halfSize, bubbleHeight)
		m.statusbar.SetSize(msg.Width)

		cmds = append(cmds, m.filetree.ToggleShowIcons(m.config.Settings.ShowIcons))
		cmds = append(cmds, resizeImgCmd, markdownCmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Exit):
			if !m.filetree.IsFiltering() {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keys.ReloadConfig):
			if !m.filetree.IsFiltering() {
				cmds = append(cmds, tea.Batch(m.reloadConfig()...))
			}
		case key.Matches(msg, m.keys.OpenFile):
			cmds = append(cmds, tea.Batch(m.openFile()...))
		case key.Matches(msg, m.keys.ToggleBox):
			m.toggleBox()
		}
	}

	m.updateStatusbar()

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
