<p align="center">
  <img src="./assets/logo.svg" height="180" width="180" />
  <p align="center">
    Keep those files organized
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

<br />

## Installation

```
go install github.com/knipferrc/fm@latest
```

## Usage

- Start `fm`
- Navigate your files with the following keys
  <br />

  - <kbd>h</kbd> Go back to the previous directory

  - <kbd>j</kbd> Move down in the file tree

  - <kbd>k</kbd> Move up in the file tree

  - <kbd>l</kbd> Opens the currently selected directory

  - <kbd>m</kbd> Move a file or folder. Once pressed you will be prompted in the status bar to type the destination for the currently highlighted file or folder. For example, `test.txt` is currently highlighted, press <kbd>m</kbd>, type `/some/new/destination` and press <kbd>enter</kbd>

  - <kbd>d</kbd> Delete a file or folder. Once you have the file or folder highlighted that you wish to delete, press <kbd>d</kbd>, a prompt will show in the status bar, type <kbd>y</kbd> to delete it or <kbd>n</kbd> to cancel

  - <kbd>r</kbd> Rename a file or folder. Once you have the file or folder highlighted that you wish to rename, press <kbd>r</kbd>, a prompt will show in the status bar, type the new name of the file or folder and then press <kbd>enter</kbd> to confirm those changes

  - <kbd>tab</kbd> Toggle between panes

  - <kbd>esc</kbd> Cancel any current action. Pressing <kbd>esc</kbd> during any action (rename, move, or delete) will cancel that action and return you to file navigation

- A config file will be generated at `.config/fm/config.yml` when you first run `fm`

<br />
<br />

### Sample Config (with defaults)

```yml
colors:
  components:
    spinner: "#F25D94"
  dir_tree:
    selected_dir_item: "#F25D94"
    unselected_dir_item: "#FFFDF5"
  pane:
    active_pane: "#F25D94"
    inactive_pane: "#FFFDF5"
settings:
  show_hidden: true
  show_icons: true
  start_dir: .
```

<br />
<br />

## Local Development

Follow the instructions below to get setup for local development

1. Clone the repo

```sh
git clone https://github.com/knipferrc/fm
```

2. Run

```sh
make run
```

3. Build a binary

```sh
make build
```

<br />

### Credit

- Thank you to this repo https://github.com/Yash-Handa/logo-ls for the icons
