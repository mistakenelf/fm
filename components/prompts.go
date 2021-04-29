package components

import (
	"fmt"
)

func MovePrompt(textInputValue string) string {
	return fmt.Sprintf("%s\n %s", "Where would you like to move this to?", textInputValue)
}

func RenamePrompt(textInputValue string) string {
	return fmt.Sprintf("%s %s", "What would you like to name this file?", textInputValue)
}

func DeletePrompt(textInputValue, currentFile string) string {
	return fmt.Sprintf("%s %s? [y/n] %s", "Are you sure you want to delete", currentFile, textInputValue)
}
