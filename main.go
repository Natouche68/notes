package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

func newModel() Model {
	notes := []Note{
		{
			title:      "Note 1",
			content:    "This is the content of note 1",
			lastEdited: time.Now().Unix(),
		},
		{
			title:      "Note 2",
			content:    "This is the content of note 2",
			lastEdited: time.Now().Unix(),
		},
	}

	return Model{
		notes:          notes,
		currentState:   "home",
		selectNoteForm: homeForm(notes),
	}
}

func (m Model) Init() tea.Cmd {
	return m.selectNoteForm.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.currentState == "home" {
				return m, tea.Quit
			}

		case "esc":
			cmds = append(cmds, saveNote(m))
		}

	case OpenedNoteMsg:
		m.currentNote = int(msg)
		cmds = append(cmds, initNoteForm(m))

	case CreatingFormMsg:
		m.currentState = "create"
		m.createNoteForm = msg
		cmds = append(cmds, m.createNoteForm.Init())

	case EditingFormMsg:
		m.currentState = "note"
		m.editingNoteForm = msg
		m.createNoteForm = nil
		m.selectNoteForm = nil
		cmds = append(cmds, m.editingNoteForm.Init())

	case NoteSavedMsg:
		m.currentState = "home"
		m.selectNoteForm = homeForm(m.notes)
		m.editingNoteForm = nil
		cmds = append(cmds, m.selectNoteForm.Init())
	}

	if m.currentState == "home" {
		form, cmd := m.selectNoteForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.selectNoteForm = f
		}

		cmds = append(cmds, cmd)

		if m.selectNoteForm.State == huh.StateCompleted {
			if m.selectNoteForm.GetString("note") == "Create a new note" {
				cmds = append(cmds, createNote(m.notes))
			} else {
				cmds = append(cmds, openNote(m.notes, m.selectNoteForm.GetString("note")))
			}
		}
	} else if m.currentState == "note" {
		form, cmd := m.editingNoteForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.selectNoteForm = f
		}

		cmds = append(cmds, cmd)
	} else if m.currentState == "create" {
		form, cmd := m.createNoteForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.selectNoteForm = f
		}

		cmds = append(cmds, cmd)

		if m.createNoteForm.State == huh.StateCompleted {
			m.currentNote = len(m.notes)
			m.notes = append(m.notes, Note{
				title:      m.createNoteForm.GetString("title"),
				content:    "",
				lastEdited: time.Now().Unix(),
			})

			cmds = append(cmds, initNoteForm(m))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := ""

	if m.currentState == "home" {
		title := titleStyle.Render("Notes")

		if m.selectNoteForm.State == huh.StateNormal {
			s = lipgloss.JoinVertical(lipgloss.Left, title, m.selectNoteForm.View())
		}
	} else if m.currentState == "note" {
		editingForm := m.editingNoteForm.View()
		help := helpStyle.Render("crtl+c - quit • esc - go home")
		s = lipgloss.JoinVertical(lipgloss.Left, editingForm, help)
	} else if m.currentState == "create" {
		s = m.createNoteForm.View()
	}

	return s
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(errorStyle.Render("There was an error while instantiating the program : " + err.Error()))
		os.Exit(1)
	}
}
