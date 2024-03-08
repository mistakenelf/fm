// Package icons provides a set of unicode icons
// based on type and extension.
package icons

import (
	"os"
	"strings"
)

// GetIndicator returns the indicator for the given file.
func GetIndicator(modebit os.FileMode) (i string) {
	switch {
	case modebit&os.ModeDir > 0:
		i = "/"
	case modebit&os.ModeNamedPipe > 0:
		i = "|"
	case modebit&os.ModeSymlink > 0:
		i = "@"
	case modebit&os.ModeSocket > 0:
		i = "="
	case modebit&1000000 > 0:
		i = "*"
	}

	return i
}

// GetIcon returns the icon based on its name, extension and indicator.
func GetIcon(name, ext, indicator string) (icon, color string) {
	var i *IconInfo
	var ok bool
	const DOT = '.'

	switch indicator {
	case "/":
		i, ok = IconDir[strings.ToLower(name+ext)]
		if ok {
			break
		}
		if len(name) == 0 || DOT == name[0] {
			i = IconDef["hiddendir"]
			break
		}
		i = IconDef["dir"]
	default:
		i, ok = IconFileName[strings.ToLower(name+ext)]
		if ok {
			break
		}

		if ext == ".go" && strings.HasSuffix(name, "_test") {
			i = IconSet["go-test"]

			break
		}

		t := strings.Split(name, ".")
		if len(t) > 1 && t[0] != "" {
			i, ok = IconSubExt[strings.ToLower(t[len(t)-1]+ext)]
			if ok {
				break
			}
		}

		i, ok = IconExt[strings.ToLower(strings.TrimPrefix(ext, "."))]
		if ok {
			break
		}

		if len(name) == 0 || DOT == name[0] {
			i = IconDef["hiddenfile"]

			break
		}
		i = IconDef["file"]
	}

	if indicator == "*" {
		if i.GetGlyph() == "\uf723" {
			i = IconDef["exe"]
		}

		i.MakeExe()
	}

	return i.GetGlyph(), i.GetColor(1)
}
