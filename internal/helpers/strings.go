package helpers

import (
	"fmt"
	"strings"
)

// Convert bytes to readable strings such as 5MB, 2kB
func ConvertBytesToSizeString(b int64) string {
	const unit = 1000

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// Remove tabs and convert them to spaces
func ConvertTabsToSpaces(input string) string {
	return strings.Replace(input, "\t", "    ", -1)
}
