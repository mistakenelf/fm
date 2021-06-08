package constants

const (
	CurrentDirectory  = "."
	PreviousDirectory = ".."
	HomeDirectory     = "~"

	StatusBarHeight = 2

	White       = "#FFFDF5"
	Pink        = "#F25D94"
	LightPurple = "#A550DF"
	DarkPurple  = "#6124DF"
	DarkGray    = "#353533"

	HelpText = `# FM (File Manager)
- h or left arrow     | go back a directory
- j or down arrow     | move cursor down
- k or up arrow       | move cursor up
- l or right arrow    | open selected folder / view file
- gg                  | go to top of pane
- G                   | go to botom of pane
- ~                   | switch to home directory
- .                   | toggle hidden files and directories
- (-)                 | Go To previous directory
- :                   | open command bar
- mkdir dirname       | create directory in current directory
- touch filename.txt  | create file in current directory
- mv newname.txt      | rename currently selected file or directory
- cp /dir/to/move/to  | move file or directory
- rm                  | remove file or directory
- tab                 | toggle between panes
	`
)
