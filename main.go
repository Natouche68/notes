package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.getNotesSpinner.Init(), getNotes)
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
			m.saveNoteSpinner = saveNoteSpinner()
			cmds = append(cmds, m.saveNoteSpinner.Init(), saveNote(m))
		}

	case ErrorMsg:
		m.error = msg

	case GetNotesMsg:
		m = Model{
			notes:          msg.notes,
			currentState:   "home",
			selectNoteForm: homeForm(msg.notes),
			db:             msg.db,
		}
		cmds = append(cmds, m.selectNoteForm.Init())

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
		m.saveNoteSpinner = nil
		cmds = append(cmds, m.selectNoteForm.Init())
	}

	if m.getNotesSpinner != nil {
		spinnerUpdated, cmd := m.getNotesSpinner.Update(msg)
		if s, ok := spinnerUpdated.(*spinner.Spinner); ok {
			m.getNotesSpinner = s
		}

		cmds = append(cmds, cmd)
	} else if m.saveNoteSpinner != nil {
		spinnerUpdated, cmd := m.saveNoteSpinner.Update(msg)
		if s, ok := spinnerUpdated.(*spinner.Spinner); ok {
			m.saveNoteSpinner = s
		}

		cmds = append(cmds, cmd)
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

	if m.error != nil {
		return errorStyle.Render(m.error.Error())
	}

	if m.getNotesSpinner != nil {
		return m.getNotesSpinner.View()
	}

	if m.currentState == "home" {
		title := titleStyle.Render("Notes")

		if m.selectNoteForm.State == huh.StateNormal {
			s = lipgloss.JoinVertical(lipgloss.Left, title, m.selectNoteForm.View())
		}
	} else if m.currentState == "note" {
		editingForm := m.editingNoteForm.View()

		if m.saveNoteSpinner != nil {
			spinnerView := m.saveNoteSpinner.View()
			s = lipgloss.JoinVertical(lipgloss.Left, editingForm, spinnerView)
		} else {
			help := helpStyle.Render("crtl+c - quit • esc - go home")
			s = lipgloss.JoinVertical(lipgloss.Left, editingForm, help)
		}
	} else if m.currentState == "create" {
		s = m.createNoteForm.View()
	}

	return s
}

func main() {
	p := tea.NewProgram(Model{
		getNotesSpinner: getNotesSpinner(),
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println(errorStyle.Render("There was an error while instantiating the program : " + err.Error()))
		os.Exit(1)
	}
}
