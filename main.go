package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct{}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "Hello, World!"

	return s
}

func main() {
	p := tea.NewProgram(Model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("There was an error while instantiating the program : ", err)
		os.Exit(1)
	}
}
