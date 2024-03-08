package filetree

type DirectoryItem struct {
	name             string
	details          string
	path             string
	extension        string
	isDirectory      bool
	currentDirectory string
}

type Model struct {
	cursor int
	files  []DirectoryItem
	active bool
	keyMap KeyMap
	min    int
	max    int
	height int
	width  int
}

func New() Model {
	return Model{
		cursor: 0,
		active: true,
		keyMap: DefaultKeyMap(),
		min:    0,
		max:    0,
	}
}
