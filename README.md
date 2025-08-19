<p align="center">
  <img src="./assets/logo.svg" height="180" width="180" />
  <p align="center">
    Keep those files organized
  </p>
  <p align="center">
    <a href="https://github.com/mistakenelf/fm/releases"><img src="https://img.shields.io/github/v/release/mistakenelf/fm" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/mistakenelf/fm?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/mistakenelf/fm/actions"><img src="https://img.shields.io/github/actions/workflow/status/mistakenelf/fm/build.yml?branch=main" alt="Build Status"></a>
  </p>
</p>

<p align="center" style="margin-top: 30px; margin-bottom: 20px;">
  <img src="./assets/default.png" width="450" alt="default screenshot">
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
- [Cobra](https://github.com/spf13/cobra)

## Installation

### Curl

```sh
curl -sfL https://raw.githubusercontent.com/mistakenelf/fm/main/install.sh | sh
```

### Go

```
go install github.com/mistakenelf/fm@latest
```

### AUR

Install through the Arch User Repository with your favorite AUR helper.
There are currently two possible packages:

- [fm-git](https://aur.archlinux.org/packages/fm-git/): Builds the package from the main branch

```sh
paru -S fm-git
```

- [fm-bin](https://aur.archlinux.org/packages/fm-bin/): Uses the github release package

```sh
paru -S fm-bin
```

## Features

- File icons (requires nerd font)
- Layout adjusts to terminal resize
- Syntax highlighting for source code with customizable themes using styles from [chroma](https://swapoff.org/chroma/playground/) (dracula, monokai etc.)
- Render pretty markdown
- Mouse support
- Themes (`default`, `gruvbox`, `nord`)
- Render PNG, JPG and JPEG as strings
- Colors adapt to terminal background, for syntax highlighting to work properly on light/dark terminals, set the appropriate theme flags
- Open selected file in editor set in EDITOR environment variable
- Copy selected directory items path to the clipboard
- Read PDF files

## Themes

### Default

<img src="./assets/default.png" width="350" alt="default">

### Gruvbox

<img src="./assets/gruvbox.png" width="350" alt="gruvbox">

### Nord

<img src="./assets/nord.png" width="350" alt="nord">

## Usage

- `fm` will start fm in the current directory
- `fm update` will update fm to the latest version
- `fm --start-dir=/some/start/dir` will start fm in the specified directory
- `fm --selection-path=/tmp/tmpfile` will write the selected items path to the selection path when pressing <kbd>E</kbd> and exit fm
- `fm --enable-logging=true` start fm with logging enabled
- `fm --pretty-markdown=true` render markdown using glamour to make it look nice
- `fm --theme=default` set the theme of fm
- `fm --show-icons=false` set whether to show icons or not
- `fm --syntax-theme=dracula` sets the syntax theme to render code with

## Local Development

Follow the instructions below to get setup for local development

1. Clone the repo

```sh
git clone https://github.com/mistakenelf/fm
```

2. Run

```sh
make
```

3. Build a binary

```sh
make build
```
