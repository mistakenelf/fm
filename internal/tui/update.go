package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

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
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		case "q":
			if !b.filetree.IsFiltering() {
				return b, tea.Quit
			}
		case " ":
			selectedFile := b.filetree.GetSelectedItem()
			if !selectedFile.IsDirectory() {
				if selectedFile.FileExtension() == ".png" || selectedFile.FileExtension() == ".jpg" {
					b.state = showImageState
					readFileCmd := b.image.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				} else if selectedFile.FileExtension() == ".md" {
					b.state = showMarkdownState
					markdownCmd := b.markdown.SetFileName(selectedFile.FileName())
					cmds = append(cmds, markdownCmd)
				} else if selectedFile.FileExtension() == ".pdf" {
					b.state = showPdfState
					pdfCmd := b.pdf.SetFileName(selectedFile.FileName())
					cmds = append(cmds, pdfCmd)
				} else {
					b.state = showCodeState
					readFileCmd := b.code.SetFileName(selectedFile.FileName())
					cmds = append(cmds, readFileCmd)
				}
			}
		case "tab":
			b.activeBox = (b.activeBox + 1) % 2
			if b.activeBox == 0 {
				b.filetree.SetIsActive(true)
				b.code.SetIsActive(false)
				b.markdown.SetIsActive(false)
				b.image.SetIsActive(false)
				b.pdf.SetIsActive(false)
				b.help.SetIsActive(false)
				b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
			} else {
				switch b.state {
				case idleState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(false)
					b.markdown.SetIsActive(false)
					b.image.SetIsActive(false)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(true)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				case showCodeState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(true)
					b.markdown.SetIsActive(false)
					b.image.SetIsActive(false)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(false)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
				case showImageState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(false)
					b.markdown.SetIsActive(false)
					b.image.SetIsActive(true)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(false)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				case showMarkdownState:
					b.filetree.SetIsActive(false)
					b.code.SetIsActive(false)
					b.markdown.SetIsActive(true)
					b.image.SetIsActive(false)
					b.pdf.SetIsActive(false)
					b.help.SetIsActive(false)
					b.filetree.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.help.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.code.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.image.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"})
					b.markdown.SetBorderColor(lipgloss.AdaptiveColor{Light: "#000000", Dark: "63"})
				}
			}
		}
	}

	b.statusbar.SetContent(
		b.filetree.GetSelectedItem().ShortName(),
		b.filetree.GetSelectedItem().CurrentDirectory(),
		fmt.Sprintf("%d/%d", b.filetree.Cursor(), b.filetree.TotalItems()),
		"FM",
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
