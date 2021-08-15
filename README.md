<p align="center">
  <img src="./assets/logo.svg" height="180" width="180" />
  <p align="center">
    Keep those files organized
  </p>
  <p align="center">
    <a href="https://github.com/knipferrc/fm/releases"><img src="https://img.shields.io/github/v/release/knipferrc/fm" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/knipferrc/fm?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/knipferrc/fm/actions"><img src="https://img.shields.io/github/workflow/status/knipferrc/fm/Release" alt="Build Status"></a>
  </p>
</p>

<p align="center" style="margin-top: 30px; margin-bottom: 20px;">
  <img src="./assets/screenshot.png" width="450" title="hover text">
</p>

## About The Project

A terminal based file manager

### Built With

- [Go](https://golang.org/)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [bubbles](https://github.com/charmbracelet/bubbles)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [Glamour](https://github.com/charmbracelet/glamour)
- [Chroma](https://github.com/alecthomas/chroma)
- [Viper](https://github.com/spf13/viper)
- [Cobra](https://github.com/spf13/cobra)

## Installation

```
go install github.com/knipferrc/fm@latest
```

## Features

- Double pane layout
- File icons
- Layout adjusts to terminal resize
- Syntax highlighting for source code
- Mouse support
- Customizable colors
- Render PNG, JPG and JPEG as ASCII

## Usage

- Run `fm` or `fm /some/dir`

## Navigation

| Key                  | Description                                                                                                                                                                                                                                                      |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| h or left            | Go back to previous directory                                                                                                                                                                                                                                    |
| j or down            | Move down in the file tree or scroll pane down                                                                                                                                                                                                                   |
| k or up              | Move up in the file tree or scroll pane up                                                                                                                                                                                                                       |
| l or right           | Opens the currently selected directory or file                                                                                                                                                                                                                   |
| gg                   | Jump to bottom of file tree or pane                                                                                                                                                                                                                              |
| G                    | Jump to top of file tree or pane                                                                                                                                                                                                                                 |
| ~                    | Go to home directory                                                                                                                                                                                                                                             |
| .                    | Toggle hide files and directories                                                                                                                                                                                                                                |
| -                    | Go to previous directory                                                                                                                                                                                                                                         |
| ctrl+c               | Exit                                                                                                                                                                                                                                                             |
| q                    | Exit if command bar is not open                                                                                                                                                                                                                                  |
| m                    | Move the currently selected file or directory. Once pressed, the file manager enters move mode. Navigate the tree as usual and press enter in the desired destination directory. It will navigate back to the starting direcotry in which the move was initiated |
| :                    | Open command bar                                                                                                                                                                                                                                                 |
| mkdir dirname        | Create a new directory in the current directory                                                                                                                                                                                                                  |
| touch filename.txt   | Create a new file in the current directory                                                                                                                                                                                                                       |
| rename or mv newname | Rename currently selected file or directory                                                                                                                                                                                                                      |
| delete or rm         | Delete the currently selected file or directory                                                                                                                                                                                                                  |
| tab                  | Toggle between panes                                                                                                                                                                                                                                             |
| esc                  | Cancel any current action. Pressing escape during any action (rename, move, delete) will cancel that operation and return back to file navigation                                                                                                                |
| z                    | Create a zip file of the currently selected directory item                                                                                                                                                                                                       |
| u                    | Unzip a zip file                                                                                                                                                                                                                                                 |
| c                    | Create a copy of a file or directory                                                                                                                                                                                                                             |

## Configuration

- A config file will be generated at `~/.fm.yml` when you first run `fm`

```yml
colors:
  dir_tree:
    selected_item: "#F25D94"
    unselected_item: "#FFFDF5"
  pane:
    active_border_color: "#F25D94"
    inactive_border_color: "#FFFDF5"
  spinner: "#F25D94"
  status_bar:
    bar:
      background: "#353533"
      foreground: "#FFFDF5"
    logo:
      background: "#6124DF"
      foreground: "#FFFDF5"
    selected_file:
      background: "#F25D94"
      foreground: "#FFFDF5"
    total_files:
      background: "#A550DF"
      foreground: "#FFFDF5"
settings:
  enable_logging: false
  enable_mousewheel: true
  pretty_markdown: true
  rounded_panes: false
  show_icons: true
  start_dir: .
```

## Local Development

Follow the instructions below to get setup for local development

1. Clone the repo

```sh
git clone https://github.com/knipferrc/fm
```

2. Run

```sh
make
```

3. Build a binary

```sh
make build
```

## Credit

- Thank you to this repo https://github.com/Yash-Handa/logo-ls for the icons
