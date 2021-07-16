package helpers

import (
	"strings"
)

// Parse command bar commands
func ParseCommand(command string) (string, string) {
	// Split the command string into an array
	cmdString := strings.Split(command, " ")

	// If theres only one item in the array, its a singular
	// command such as rm
	if len(cmdString) == 1 {
		cmdName := cmdString[0]

		return cmdName, ""
	}

	// This command has two values, first one is the name
	// of the command, other is the value to pass back
	// to the UI to update
	if len(cmdString) == 2 {
		cmdName := cmdString[0]
		cmdValue := cmdString[1]

		return cmdName, cmdValue
	}

	return "", ""
}
