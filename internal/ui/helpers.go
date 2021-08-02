package ui

import (
	"fmt"

	"github.com/knipferrc/fm/internal/config"
	"github.com/knipferrc/fm/internal/constants"
	"github.com/knipferrc/fm/internal/helpers"
	"github.com/knipferrc/fm/pkg/icons"
)

// scrollPrimaryPane handles the scrolling of the primary pane which will handle
// infinite scroll on the dirtree and the scrolling of the viewport.
func (m *model) scrollPrimaryPane() {
	top := m.primaryPane.GetYOffset()
	bottom := m.primaryPane.GetHeight() + m.primaryPane.GetYOffset() - 1

	// If the cursor is above the top of the viewport scroll up on the viewport
	// else were at the bottom and need to scroll the viewport down.
	if m.dirTree.GetCursor() < top {
		m.primaryPane.LineUp(1)
	} else if m.dirTree.GetCursor() > bottom {
		m.primaryPane.LineDown(1)
	}

	// If the cursor of the dirtree is at the bottom of the files
	// set the cursor to 0 to go to the top of the dirtree and
	// scroll the pane to the top else, were at the top of the dirtree and pane so
	// scroll the pane to the bottom and set the cursor to the bottom.
	if m.dirTree.GetCursor() > m.dirTree.GetTotalFiles()-1 {
		m.dirTree.GotoTop()
		m.primaryPane.GotoTop()
	} else if m.dirTree.GetCursor() < top {
		m.dirTree.GotoBottom()
		m.primaryPane.GotoBottom()
	}
}

// getStatusBarContent returns the content of the status bar (current file name, file info, file count, logo).
func (m model) getStatusBarContent() (string, string, string, string) {
	cfg := config.GetConfig()
	currentPath, err := helpers.GetWorkingDirectory()
	if err != nil {
		currentPath = constants.CurrentDirectory
	}

	if m.dirTree.GetTotalFiles() == 0 {
		return "", "", "", ""
	}

	logo := ""

	// If icons are enabled, show the directory icon next to the logo text
	// else just show the text of the logo.
	if cfg.Settings.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	// Display some information about the currently seleted file including
	// its size, the mode and the current path.
	status := fmt.Sprintf("%s %s %s",
		helpers.ConvertBytesToSizeString(m.dirTree.GetSelectedFile().Size()),
		m.dirTree.GetSelectedFile().Mode().String(),
		currentPath,
	)

	// If the command bar is shown, show the text input.
	if m.showCommandBar {
		status = m.textInput.View()
	}

	// If in move mode, update the status text to indicate move mode is enabled
	// and the name of the file or directory being moved.
	if m.inMoveMode {
		status = fmt.Sprintf("Currently moving %s", m.itemToMove.Name())
	}

	return m.dirTree.GetSelectedFile().Name(),
		status,
		fmt.Sprintf("%d/%d", m.dirTree.GetCursor()+1, m.dirTree.GetTotalFiles()),
		logo
}
