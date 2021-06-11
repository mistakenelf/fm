package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/knipferrc/fm/config"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/utils"
)

func (m *model) scrollPrimaryPane() {
	top := m.primaryPane.Viewport.YOffset
	bottom := m.primaryPane.Height + m.primaryPane.YOffset - 1

	if m.dirTree.GetCursor() < top {
		m.primaryPane.LineUp(1)
	} else if m.dirTree.GetCursor() > bottom {
		m.primaryPane.LineDown(1)
	}

	if m.dirTree.GetCursor() > m.dirTree.GetTotalFiles()-1 {
		m.dirTree.GotoTop()
		m.primaryPane.GotoTop()
	} else if m.dirTree.GetCursor() < top {
		m.dirTree.GotoBottom()
		m.primaryPane.GotoBottom()
	}
}

func (m model) getStatusBarContent() (string, string, string, string) {
	cfg := config.GetConfig()
	currentPath, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	logo := ""
	if cfg.Settings.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	status := fmt.Sprintf("%s %s %s",
		utils.ConvertBytesToSizeString(m.dirTree.GetSelectedFile().Size()),
		m.dirTree.GetSelectedFile().Mode().String(),
		currentPath,
	)

	if m.showCommandBar {
		status = m.textInput.View()
	}

	return m.dirTree.GetSelectedFile().Name(), status, fmt.Sprintf("%d/%d", m.dirTree.GetCursor()+1, m.dirTree.GetTotalFiles()), logo
}
