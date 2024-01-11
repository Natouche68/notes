package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/charm/kv"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type Model struct {
	notes           []Note
	currentState    string
	currentNote     int
	isSaving        bool
	isQuitting      bool
	db              *kv.KV
	error           error
	selectNoteForm  *huh.Form
	createNoteForm  *huh.Form
	editingNoteForm *huh.Form
	deleteNoteForm  *huh.Form
	getNotesSpinner *spinner.Spinner
	saveNoteSpinner *spinner.Spinner
}

type Note struct {
	title      string
	content    string
	lastEdited int64
}

type ErrorMsg error
type GetNotesMsg struct {
	notes []Note
	db    *kv.KV
}
type OpenedNoteMsg int
type CreatingFormMsg *huh.Form
type EditingFormMsg *huh.Form
type NoteSavedMsg Model
type NoteDeletedMsg []Note
type NoteSavedAfterDeletionMsg Model

func (note Note) MarshalJSON() ([]byte, error) {
	return []byte(
		fmt.Sprintf(
			`{"title": "%s", "content": "%s", "lastEdited": %d}`,
			note.title, strings.ReplaceAll(note.content, "\n", "\\n"), note.lastEdited,
		),
	), nil
}

func (note *Note) UnmarshalJSON(data []byte) error {
	var raw map[string]*json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		switch k {
		case "title":
			if err := json.Unmarshal(*v, &note.title); err != nil {
				return err
			}
		case "content":
			if err := json.Unmarshal(*v, &note.content); err != nil {
				return err
			}
		case "lastEdited":
			if err := json.Unmarshal(*v, &note.lastEdited); err != nil {
				return err
			}
		}
	}

	return nil
}
