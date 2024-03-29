package main

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/charm/kv"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"golang.org/x/term"
)

func getNotes() tea.Msg {
	db, err := kv.OpenWithDefaults("notes-db")
	if err != nil {
		return ErrorMsg(err)
	}

	if err := db.Sync(); err != nil {
		return ErrorMsg(err)
	}

	notesFromDb, err := db.Get([]byte("notes"))
	if err != nil {
		if err.Error() == "Key not found" {
			notesFromDb = []byte("[]")
		} else {
			return ErrorMsg(err)
		}
	}

	var notes []Note
	error := json.Unmarshal(notesFromDb, &notes)
	if error != nil {
		return ErrorMsg(error)
	}

	return GetNotesMsg{notes, db}
}

func homeForm(notes []Note) *huh.Form {
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].lastEdited > notes[j].lastEdited
	})

	notesTitles := []string{
		"Create a new note",
	}
	for _, note := range notes {
		notesTitles = append(notesTitles, note.title)
	}
	notesTitles = append(notesTitles, "Delete a note")

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("note").
				Title("Select a note or create a new one").
				Options(huh.NewOptions(notesTitles...)...),
		),
	).WithShowHelp(false).WithTheme(huh.ThemeBase16())
}

func openNote(notes []Note, noteToOpen string) tea.Cmd {
	return func() tea.Msg {
		openedNoteIndex := slices.IndexFunc(notes, func(note Note) bool {
			return note.title == noteToOpen
		})

		notes[openedNoteIndex].content = strings.ReplaceAll(notes[openedNoteIndex].content, "\\n", "\n")

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
					CharLimit(3200).
					WithTheme(&huh.Theme{
						Focused: huh.FieldStyles{
							Title: titleStyle,
						},
					}),
			),
		).
			WithShowHelp(false).
			WithTheme(huh.ThemeBase16()).
			WithKeyMap(&huh.KeyMap{
				Text: huh.TextKeyMap{
					NewLine: key.NewBinding(key.WithKeys("enter")),
				},
			}).
			WithWidth(width - 1).
			WithTheme(&huh.Theme{
				Focused: huh.FieldStyles{
					Title: titleStyle,
				},
			}),
		)
	}
}

func saveNote(m Model) tea.Cmd {
	var error error

	m.notes[m.currentNote].lastEdited = time.Now().Unix()

	jsonNotes, err := json.Marshal(m.notes)
	if err != nil {
		error = err
	}
	if m.isQuitting {
		if err := m.db.Set([]byte("notes"), jsonNotes); err != nil {
			error = err
		}
	}

	return func() tea.Msg {
		if !m.isQuitting {
			if err := m.db.Set([]byte("notes"), jsonNotes); err != nil {
				error = err
			}
		}

		if error != nil {
			return ErrorMsg(error)
		} else {
			return NoteSavedMsg(m)
		}
	}
}

func getNotesSpinner() *spinner.Spinner {
	return spinner.New().
		Title("Loading notes...").
		Type(spinner.Points).
		Style(spinnerStyle).
		TitleStyle(spinnerTitleStyle)
}

func saveNoteSpinner() *spinner.Spinner {
	return spinner.New().
		Title("Saving note...").
		Type(spinner.Points).
		Style(spinnerStyle).
		TitleStyle(spinnerTitleStyle)
}

func quit(m Model) tea.Cmd {
	m.db.Close()

	return tea.Quit
}

func deleteNoteForm(notes []Note) *huh.Form {
	notesTitle := []string{}
	for _, note := range notes {
		notesTitle = append(notesTitle, note.title)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("noteToDelete").
				Title("Select a note to delete").
				Options(huh.NewOptions(notesTitle...)...),
		),
	).WithShowHelp(false).WithTheme(huh.ThemeBase16())
}

func deleteNote(notes []Note, noteToDelete string) tea.Cmd {
	return func() tea.Msg {
		return NoteDeletedMsg(slices.DeleteFunc(notes, func(note Note) bool {
			return note.title == noteToDelete
		}))
	}
}

func saveAfterDeletion(m Model) tea.Cmd {
	var error error

	jsonNotes, err := json.Marshal(m.notes)
	if err != nil {
		error = err
	}

	return func() tea.Msg {
		if err := m.db.Set([]byte("notes"), jsonNotes); err != nil {
			error = err
		}

		if error != nil {
			return ErrorMsg(error)
		} else {
			return NoteSavedAfterDeletionMsg(m)
		}
	}
}
