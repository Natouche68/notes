package main

import (
	"github.com/charmbracelet/huh"
)

type Model struct {
	notes          []Note
	currentState   string
	currentNote    int
	selectNoteForm *huh.Form
	createNoteForm *huh.Form
}

type Note struct {
	title      string
	content    string
	lastEdited int64
}

type OpenedNoteMsg int
type CreatingFormMsg *huh.Form
