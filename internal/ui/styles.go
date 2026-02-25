package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			PaddingLeft(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("241")).
			PaddingLeft(1)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238"))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			PaddingLeft(1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			PaddingLeft(3)

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	activeIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			SetString("●")

	windowActiveIndicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				SetString("●")

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			PaddingLeft(1).
			PaddingTop(1)

	confirmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196"))

	inputPromptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170")).
				PaddingLeft(1)
)
