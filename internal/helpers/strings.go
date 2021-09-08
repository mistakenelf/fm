package helpers

import (
	"strings"
)

// ConvertTabsToSpaces converts tabs to spaces.
func ConvertTabsToSpaces(input string) string {
	return strings.Replace(input, "\t", "    ", -1)
}
