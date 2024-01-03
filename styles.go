package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true).MarginBottom(1)

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Margin(1, 4)
)
