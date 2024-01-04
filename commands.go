package main

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func openNote(notes []Note, noteToOpen string) tea.Cmd {
	return func() tea.Msg {
		openedNoteIndex := slices.IndexFunc(notes, func(note Note) bool {
			return note.title == noteToOpen
		})

		return OpenedNoteMsg(openedNoteIndex)
	}
}

func createNote() tea.Msg {
	return CreatingFormMsg(huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Create a new note").
				Placeholder("Title").
				Key("title"),
		),
	).WithTheme(huh.ThemeBase16()))
}
