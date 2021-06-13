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
	Blue        = "#1D4ED8"

	HelpText = `# FM (File Manager)
| Key                | Description                                                                                                                                                                                                                                                      |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| h or left          | Go back to previous directory                                                                                                                                                                                                                                    |
| j or down          | Move down in the file tree or scroll pane down                                                                                                                                                                                                                   |
| k or up            | Move up in the file tree or scroll pane up                                                                                                                                                                                                                       |
| l or right         | Opens the currently selected directory or file                                                                                                                                                                                                                   |
| gg                 | Jump to bottom of file tree or pane                                                                                                                                                                                                                              |
| G                  | Jump to top of file tree or pane                                                                                                                                                                                                                                 |
| ~                  | Go to home directory                                                                                                                                                                                                                                             |
| .                  | Toggle hide files and directories                                                                                                                                                                                                                                |
| -                  | Go to previous directory                                                                                                                                                                                                                                         |
| ctrl+c             | Exit                                                                                                                                                                                                                                                             |
| q                  | Exit if command bar is not open                                                                                                                                                                                                                                  |
| m                  | Move the currently selected file or directory. Once pressed, the file manager enters move mode. Navigate the tree as usual and press enter in the desired destination directory. It will navigate back to the starting direcotry in which the move was initiated |
| :                  | Open command bar                                                                                                                                                                                                                                                 |
| mkdir dirname      | Create a new directory in the current directory                                                                                                                                                                                                                  |
| touch filename.txt | Create a new file in the current directory                                                                                                                                                                                                                       |
| rename or mv       | Rename currently selected file or directory                                                                                                                                                                                                                      |
| delete or rm       | Delete the currently selected file or directory                                                                                                                                                                                                                  |
| tab                | Toggle between panes                                                                                                                                                                                                                                             |
| esc                | Cancel any current action. Pressing escape during any action (rename, move, delete) will cancel that operation and return back to file navigation                                                                                                                |                          
| tab                | toggle between panes
	`
)
