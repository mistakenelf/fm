package icons

import (
	"os"
	"strings"
)

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

func GetIcon(name, ext, indicator string) (icon, color string) {
	var i *Icon_Info
	var ok bool
	const DOT = '.'

	switch indicator {
	case "/":
		i, ok = Icon_Dir[strings.ToLower(name+ext)]
		if ok {
			break
		}
		if len(name) == 0 || DOT == name[0] {
			i = Icon_Def["hiddendir"]
			break
		}
		i = Icon_Def["dir"]
	default:
		i, ok = Icon_FileName[strings.ToLower(name+ext)]
		if ok {
			break
		}

		if ext == ".go" && strings.HasSuffix(name, "_test") {
			i = Icon_Set["go-test"]
			break
		}

		t := strings.Split(name, ".")
		if len(t) > 1 && t[0] != "" {
			i, ok = Icon_SubExt[strings.ToLower(t[len(t)-1]+ext)]
			if ok {
				break
			}
		}

		i, ok = Icon_Ext[strings.ToLower(strings.TrimPrefix(ext, "."))]
		if ok {
			break
		}

		if len(name) == 0 || DOT == name[0] {
			i = Icon_Def["hiddenfile"]
			break
		}
		i = Icon_Def["file"]
	}

	if indicator == "*" {
		if i.GetGlyph() == "\uf723" {
			i = Icon_Def["exe"]
		}

		i.MakeExe()
	}

	return i.GetGlyph(), i.GetColor(1)
}
