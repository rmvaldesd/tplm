package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorAccent)).
			PaddingLeft(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorMuted)).
			PaddingLeft(1)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorDim))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorAccent)).
			PaddingLeft(1)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorText)).
			PaddingLeft(3)

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted))

	activeIndicator = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorGreen)).
			SetString(SymbolActive)

	windowActiveIndicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorGreen)).
				SetString(SymbolActive)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMuted)).
			PaddingLeft(1).
			PaddingTop(1)

	confirmStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorRed))

	inputPromptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(ColorAccent)).
				PaddingLeft(1)
)
