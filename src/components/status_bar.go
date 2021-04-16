package components

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/knipferrc/fm/src/icons"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	selectedFileStyle = lipgloss.NewStyle().
				Inherit(statusBarStyle).
				Foreground(lipgloss.Color("#FFFDF5")).
				Background(lipgloss.Color("#FF5F87")).
				Padding(0, 1).
				MarginRight(1)

	fileEncodingStyle = statusItemStyle.Copy().
				Background(lipgloss.Color("#A550DF")).
				Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	logoStyle = statusItemStyle.Copy().Background(lipgloss.Color("#6124DF"))
)

func getMovingPrompt(textInput *textinput.Model, width func(string) int, screenWidth int, selectedFileName string, fileEncoding string, logo string) string {
	prompt := fmt.Sprintf("%s %s", "Where would you like to move this to?", textInput.View())

	return statusText.Copy().
		Width(screenWidth - width(selectedFileName) - width(fileEncoding) - width(logo)).
		Render(prompt)
}

func getRenamingPrompt(textInput *textinput.Model, width func(string) int, screenWidth int, selectedFileName string, fileEncoding string, logo string) string {
	prompt := fmt.Sprintf("%s %s", "What would you like to name this file?", textInput.View())

	return statusText.Copy().
		Width(screenWidth - width(selectedFileName) - width(fileEncoding) - width(logo)).
		Render(prompt)
}

func getDeletingPrompt(textInput *textinput.Model, width func(string) int, screenWidth int, currentFile, selectedFileName string, fileEncoding string, logo string) string {
	prompt := fmt.Sprintf("%s %s? [y/n] %s", "Are you sure you want to delete", currentFile, textInput.View())

	return statusText.Copy().
		Width(screenWidth - width(selectedFileName) - width(fileEncoding) - width(logo)).
		Render(prompt)
}

func StatusBar(screenWidth int, currentFile fs.FileInfo, isMoving, isRenaming, isDeleting bool, textInput *textinput.Model) string {
	doc := strings.Builder{}
	width := lipgloss.Width
	selectedFileName := ""

	if currentFile != nil {
		selectedFileName = selectedFileStyle.Render(currentFile.Name())
	}

	fileEncoding := fileEncodingStyle.Render("UTF-8")

	logo := logoStyle.Render(fmt.Sprintf("%s %s", icons.Icon_Def["dir"].GetGlyph(), "FM"))

	status := statusText.Copy().
		Width(screenWidth - width(selectedFileName) - width(fileEncoding) - width(logo)).
		Render("m - move, d - delete, r - rename, i - help")

	if isMoving {
		status = getMovingPrompt(textInput, width, screenWidth, selectedFileName, fileEncoding, logo)
	}

	if isRenaming {
		status = getRenamingPrompt(textInput, width, screenWidth, selectedFileName, fileEncoding, logo)
	}

	if isDeleting {
		status = getDeletingPrompt(textInput, width, screenWidth, currentFile.Name(), selectedFileName, fileEncoding, logo)
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileName,
		status,
		fileEncoding,
		logo,
	)

	doc.WriteString(statusBarStyle.Render(bar))

	return doc.String()
}
