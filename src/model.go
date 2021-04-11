package main

import "io/fs"

type model struct {
	Files       []fs.FileInfo
	Selected    map[int]struct{}
	Cursor      int
	Quitting    bool
	FileContent string
}
