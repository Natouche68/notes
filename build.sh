#!/bin/sh

gum spin --spinner points --title "Building for Windows..." "go build -o build/notes-windows.exe"
gum spin --spinner points --title "Building for Linux..." "go build -o build/notes-linux"
gum spin --spinner points --title "Building for macOS..." "go build -o build/notes-macos"

gum style --foreground 2 --bold --margin "1 2" "All builds completed !"
