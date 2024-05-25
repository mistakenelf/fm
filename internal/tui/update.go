package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/mistakenelf/fm/filetree"
	"github.com/mistakenelf/fm/statusbar"
)

// Update handles all UI interactions.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case statusMessageTimeoutMsg:
		m.statusMessage = ""

		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		halfSize := msg.Width / 2
		height := msg.Height - statusbar.Height

		cmds = append(cmds, m.image.SetSizeCmd(halfSize, height))
		cmds = append(cmds, m.markdown.SetSizeCmd(halfSize, height))
		cmds = append(cmds, m.csv.SetSizeCmd(halfSize, height))

		m.filetree[m.activeWorkspace].SetSize(halfSize, height-3)
		m.secondaryFiletree.SetSize(halfSize, height-3)
		m.code.SetSize(halfSize, height)
		m.pdf.SetSize(halfSize, height)
		m.statusbar.SetSize(msg.Width)
		m.help.SetSize(halfSize, height)

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Quit):
			if m.filetree[m.activeWorkspace].State == filetree.IdleState {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keyMap.OpenFile):
			if !m.showTextInput && m.activePane == 0 {
				cmds = append(cmds, m.openFileCmd())
			}
		case key.Matches(msg, m.keyMap.ResetState):
			if m.state == showMoveState {
				cmds = append(cmds, m.filetree[m.activeWorkspace].GetDirectoryListingCmd(m.directoryBeforeMove))
			}

			m.state = idleState
			m.showTextInput = false
			m.disableAllViewports()
			m.resetViewports()
			m.filetree[m.activeWorkspace].SetDisabled(false)
			m.textinput.Blur()
			m.filetree[m.activeWorkspace].State = filetree.IdleState
			m.secondaryFiletree.SetDisabled(true)
			m.activePane = 0

			m.textinput, cmd = m.textinput.Update(msg)
			cmds = append(cmds, cmd)

			m.textinput.Reset()
		case key.Matches(msg, m.keyMap.MoveDirectoryItem):
			if m.activePane == 0 && m.filetree[m.activeWorkspace].State == filetree.IdleState {
				m.activePane = (m.activePane + 1) % 2
				m.directoryBeforeMove = m.filetree[m.activeWorkspace].CurrentDirectory
				m.state = showMoveState
				m.filetree[m.activeWorkspace].State = filetree.MoveState
				m.filetree[m.activeWorkspace].SetDisabled(true)
				m.secondaryFiletree.SetDisabled(false)
				cmds = append(cmds, m.secondaryFiletree.GetDirectoryListingCmd(m.filetree[m.activeWorkspace].CurrentDirectory))
			}
		case key.Matches(msg, m.keyMap.ShowTextInput):
			if m.activePane == 0 && m.filetree[m.activeWorkspace].State == filetree.IdleState {
				m.showTextInput = true
				m.textinput.Focus()
				m.disableAllViewports()

				m.textinput, cmd = m.textinput.Update(msg)
				cmds = append(cmds, cmd)

				m.textinput.Reset()
			}
		case key.Matches(msg, m.keyMap.Submit):
			switch {
			case m.filetree[m.activeWorkspace].State == filetree.CreateFileState:
				cmds = append(cmds, m.filetree[m.activeWorkspace].CreateFileCmd(m.textinput.Value()))
			case m.filetree[m.activeWorkspace].State == filetree.CreateDirectoryState:
				cmds = append(cmds, m.filetree[m.activeWorkspace].CreateDirectoryCmd(m.textinput.Value()))
			case m.filetree[m.activeWorkspace].State == filetree.MoveState:
				cmds = append(
					cmds,
					m.filetree[m.activeWorkspace].MoveDirectoryItemCmd(
						m.filetree[m.activeWorkspace].GetSelectedItem().Path,
						m.secondaryFiletree.CurrentDirectory+"/"+m.filetree[m.activeWorkspace].GetSelectedItem().Name,
					),
				)
			case m.filetree[m.activeWorkspace].State == filetree.RenameState:
				cmds = append(cmds,
					m.filetree[m.activeWorkspace].RenameDirectoryItemCmd(
						m.filetree[m.activeWorkspace].GetSelectedItem().Path,
						m.filetree[m.activeWorkspace].CurrentDirectory+"/"+m.textinput.Value(),
					),
				)
			default:
				return m, nil
			}

			m.resetViewports()
			m.textinput.Blur()
			m.textinput.Reset()
			m.showTextInput = false
			m.activePane = 0
		case key.Matches(msg, m.keyMap.TogglePane):
			if !m.showTextInput {
				m.activePane = (m.activePane + 1) % 2

				if m.activePane == 0 {
					m.filetree[m.activeWorkspace].SetDisabled(false)
					m.secondaryFiletree.SetDisabled(true)
					m.disableAllViewports()
				} else {
					m.filetree[m.activeWorkspace].SetDisabled(true)

					switch m.state {
					case idleState:
						m.disableAllViewports()
						m.help.SetViewportDisabled(false)
					case showCodeState:
						m.disableAllViewports()
						m.code.SetViewportDisabled(false)
					case showImageState:
						m.disableAllViewports()
						m.image.SetViewportDisabled(false)
					case showPdfState:
						m.disableAllViewports()
						m.pdf.SetViewportDisabled(false)
					case showMarkdownState:
						m.disableAllViewports()
						m.markdown.SetViewportDisabled(false)
					case showMoveState:
						m.secondaryFiletree.SetDisabled(false)
						m.disableAllViewports()
					}
				}
			}
		case key.Matches(msg, m.keyMap.GotoTop):
			if m.activePane != 0 {
				m.resetViewports()
			}
		case key.Matches(msg, m.keyMap.GotoBottom):
			if m.activePane != 0 {
				m.code.GotoBottom()
				m.pdf.GotoBottom()
				m.markdown.GotoBottom()
				m.help.GotoBottom()
				m.image.GotoBottom()
			}
		case key.Matches(msg, m.keyMap.WorkspaceOne):
			m.workspaces = []int{1}
			m.activeWorkspace = 0

			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keyMap.WorkspaceTwo):
			filetreeModel := filetree.New(m.config.StartDir)
			filetreeModel.SetTheme(m.config.Theme.SelectedTreeItemColor, m.config.Theme.UnselectedTreeItemColor)
			filetreeModel.SetSelectionPath(m.config.SelectionPath)
			filetreeModel.SetShowIcons(m.config.ShowIcons)

			m.workspaces = append(m.workspaces, 2)
			m.activeWorkspace = 1
			m.filetree = append(m.filetree, filetreeModel)

			cmds = append(cmds, m.filetree[m.activeWorkspace].Init())
			m.filetree[m.activeWorkspace].SetSize(m.width/2, m.height-3)
			return m, tea.Batch(cmds...)
		}
	}

	if m.filetree[m.activeWorkspace].State == filetree.CreateDirectoryState ||
		m.filetree[m.activeWorkspace].State == filetree.CreateFileState ||
		m.filetree[m.activeWorkspace].State == filetree.RenameState {
		m.textinput, cmd = m.textinput.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.filetree[m.activeWorkspace], cmd = m.filetree[m.activeWorkspace].Update(msg)
	cmds = append(cmds, cmd)

	m.secondaryFiletree, cmd = m.secondaryFiletree.Update(msg)
	cmds = append(cmds, cmd)

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

	m.csv, cmd = m.csv.Update(msg)
	cmds = append(cmds, cmd)

	m.updateStatusBar()

	return m, tea.Batch(cmds...)
}
