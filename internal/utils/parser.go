package utils

import (
	"strings"
)

func ParseCommand(command string) (string, string) {
	cmdString := strings.Split(command, " ")

	if len(cmdString) == 1 {
		cmdName := cmdString[0]

		return cmdName, ""
	}

	if len(cmdString) == 2 {
		cmdName := cmdString[0]
		cmdValue := cmdString[1]

		return cmdName, cmdValue
	}

	return "", ""
}
