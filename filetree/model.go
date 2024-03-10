package filetree

type DirectoryItem struct {
	Name             string
	Details          string
	Path             string
	Extension        string
	IsDirectory      bool
	CurrentDirectory string
}

type Model struct {
	Cursor int
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
		Cursor: 0,
		active: true,
		keyMap: DefaultKeyMap(),
		min:    0,
		max:    0,
	}
}
