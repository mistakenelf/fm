package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/fm/internal/config"
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
}

// resetViewports goes to the top of all bubbles viewports.
func (b *Bubble) resetViewports() {
	b.code.GotoTop()
	b.pdf.GotoTop()
	b.markdown.GotoTop()
	b.help.GotoTop()
	b.image.GotoTop()
}

// deactivateALlBubbles sets all bubbles to inactive.
func (b *Bubble) deactivateAllBubbles() {
	b.filetree.SetIsActive(false)
	b.code.SetIsActive(false)
	b.markdown.SetIsActive(false)
	b.image.SetIsActive(false)
	b.pdf.SetIsActive(false)
	b.help.SetIsActive(false)
}

// resetBorderColors resets all bubble border colors to default.
func (b *Bubble) resetBorderColors() {
	b.filetree.SetBorderColor(b.theme.InactiveBoxBorderColor)
	b.help.SetBorderColor(b.theme.InactiveBoxBorderColor)
	b.code.SetBorderColor(b.theme.InactiveBoxBorderColor)
	b.image.SetBorderColor(b.theme.InactiveBoxBorderColor)
	b.markdown.SetBorderColor(b.theme.InactiveBoxBorderColor)
	b.pdf.SetBorderColor(b.theme.InactiveBoxBorderColor)
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
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	b.filetree, cmd = b.filetree.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		resizeImgCmd := b.image.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		markdownCmd := b.markdown.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.filetree.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.help.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.code.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.pdf.SetSize(msg.Width/2, msg.Height-statusbar.Height)
		b.statusbar.SetSize(msg.Width)

		cmds = append(cmds, resizeImgCmd, markdownCmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, b.keys.Quit):
			return b, tea.Quit
		case key.Matches(msg, b.keys.Exit):
			if !b.filetree.IsFiltering() {
				return b, tea.Quit
			}
		case key.Matches(msg, b.keys.ReloadConfig):
			if !b.filetree.IsFiltering() {
				cfg, err := config.ParseConfig()
				if err != nil {
					return b, nil
				}

				b.config = cfg
				syntaxTheme := cfg.Theme.SyntaxTheme.Light
				if lipgloss.HasDarkBackground() {
					syntaxTheme = cfg.Theme.SyntaxTheme.Dark
				}

				b.code.SetSyntaxTheme(syntaxTheme)
			}
		case key.Matches(msg, b.keys.OpenFile):
			selectedFile := b.filetree.GetSelectedItem()
			if !selectedFile.IsDirectory() {
				b.resetViewports()

				switch {
				case selectedFile.FileExtension() == ".png" || selectedFile.FileExtension() == ".jpg" || selectedFile.FileExtension() == ".jpeg":
					b.state = showImageState
					readFileCmd := b.image.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				case selectedFile.FileExtension() == ".md" && b.config.Settings.PrettyMarkdown:
					b.state = showMarkdownState
					markdownCmd := b.markdown.SetFileName(selectedFile.FileName())
					cmds = append(cmds, markdownCmd)
				case selectedFile.FileExtension() == ".pdf":
					b.state = showPdfState
					pdfCmd := b.pdf.SetFileName(selectedFile.FileName())
					cmds = append(cmds, pdfCmd)
				case contains(forbiddenExtensions, selectedFile.FileExtension()):
					return b, nil
				default:
					b.state = showCodeState
					readFileCmd := b.code.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				}
			}
		case key.Matches(msg, b.keys.ToggleBox):
			b.activeBox = (b.activeBox + 1) % 2
			if b.activeBox == 0 {
				b.deactivateAllBubbles()
				b.filetree.SetIsActive(true)
				b.resetBorderColors()
				b.filetree.SetBorderColor(b.theme.ActiveBoxBorderColor)
			} else {
				switch b.state {
				case idleState:
					b.deactivateAllBubbles()
					b.help.SetIsActive(true)
					b.resetBorderColors()
					b.help.SetBorderColor(b.theme.ActiveBoxBorderColor)
				case showCodeState:
					b.deactivateAllBubbles()
					b.code.SetIsActive(true)
					b.resetBorderColors()
					b.code.SetBorderColor(b.theme.ActiveBoxBorderColor)
				case showImageState:
					b.deactivateAllBubbles()
					b.image.SetIsActive(true)
					b.resetBorderColors()
					b.image.SetBorderColor(b.theme.ActiveBoxBorderColor)
				case showMarkdownState:
					b.deactivateAllBubbles()
					b.markdown.SetIsActive(true)
					b.resetBorderColors()
					b.markdown.SetBorderColor(b.theme.ActiveBoxBorderColor)
				case showPdfState:
					b.deactivateAllBubbles()
					b.markdown.SetIsActive(true)
					b.resetBorderColors()
					b.pdf.SetBorderColor(b.theme.ActiveBoxBorderColor)
				}
			}
		}
	}

	b.statusbar.SetContent(
		b.filetree.GetSelectedItem().ShortName(),
		b.filetree.GetSelectedItem().CurrentDirectory(),
		fmt.Sprintf("%d/%d", b.filetree.Cursor(), b.filetree.TotalItems()),
		fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM"),
	)

	b.code, cmd = b.code.Update(msg)
	cmds = append(cmds, cmd)

	b.markdown, cmd = b.markdown.Update(msg)
	cmds = append(cmds, cmd)

	b.image, cmd = b.image.Update(msg)
	cmds = append(cmds, cmd)

	b.pdf, cmd = b.pdf.Update(msg)
	cmds = append(cmds, cmd)

	b.help, cmd = b.help.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}
