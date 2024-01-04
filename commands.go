package main

import (
	"errors"
	"os"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"golang.org/x/term"
)

func openNote(notes []Note, noteToOpen string) tea.Cmd {
	return func() tea.Msg {
		openedNoteIndex := slices.IndexFunc(notes, func(note Note) bool {
			return note.title == noteToOpen
		})

		return OpenedNoteMsg(openedNoteIndex)
	}
}

func createNote(notes []Note) tea.Cmd {
	return func() tea.Msg {
		return CreatingFormMsg(huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("title").
					Title("Create a new note").
					Placeholder("Title").
					Validate(func(s string) error {
						if s == "" {
							return errors.New("Title cannot be empty")
						}
						for _, note := range notes {
							if s == note.title {
								return errors.New("Title already exists")
							}
						}
						return nil
					}),
			),
		).WithShowHelp(false).WithTheme(huh.ThemeBase16()))
	}
}

func initNoteForm(m Model) tea.Cmd {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))

	return func() tea.Msg {
		return EditingFormMsg(huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Key("content").
					Title(m.notes[m.currentNote].title).
					Value(&m.notes[m.currentNote].content).
					Lines(height - 6).
					CharLimit(3200),
			),
		).WithShowHelp(false).WithTheme(huh.ThemeBase16()).WithKeyMap(&huh.KeyMap{
			Text: huh.TextKeyMap{
				NewLine: key.NewBinding(key.WithKeys("enter")),
			},
		}).WithWidth(width - 1))
	}
}
