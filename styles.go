package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("4")).
			Bold(true).
			MarginBottom(1).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true).
			Margin(1, 4)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("4")).
			Width(3).
			MarginRight(1)

	spinnerTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("7")).
				Italic(true)
)
