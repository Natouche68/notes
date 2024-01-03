package main

import (
	"fmt"
	"os"
	"sort"
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

	sort.Slice(notes, func(i, j int) bool {
		return notes[i].lastEdited > notes[j].lastEdited
	})

	notesTitles := []string{
		"Create a new note",
	}
	for _, note := range notes {
		notesTitles = append(notesTitles, note.title)
	}

	return Model{
		notes:        notes,
		currentState: "home",
		selectNoteForm: huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Key("note").
					Title("Select a note or create a new one").
					Options(huh.NewOptions(notesTitles...)...),
			),
		).WithTheme(huh.ThemeBase16()),
	}
}

func (m Model) Init() tea.Cmd {
	return m.selectNoteForm.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	if m.currentState == "home" {
		form, cmd := m.selectNoteForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.selectNoteForm = f
		}

		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	s := ""

	if m.currentState == "home" {
		s += titleStyle.Render("Notes")

		if m.selectNoteForm.State == huh.StateCompleted {
			s = lipgloss.JoinVertical(lipgloss.Left, m.selectNoteForm.GetString("note"))
		} else {
			s = lipgloss.JoinVertical(lipgloss.Left, s, m.selectNoteForm.View())
		}
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
