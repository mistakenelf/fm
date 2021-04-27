package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
)

func MovePrompt(textInput textinput.Model) string {
	return fmt.Sprintf("%s\n %s", "Where would you like to move this to?", textInput.View())
}

func RenamePrompt(textInput textinput.Model) string {
	return fmt.Sprintf("%s %s", "What would you like to name this file?", textInput.View())
}

func DeletePrompt(textInput textinput.Model, currentFile string) string {
	return fmt.Sprintf("%s %s? [y/n] %s", "Are you sure you want to delete", currentFile, textInput.View())
}
