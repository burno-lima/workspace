package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor  = lipgloss.Color("#7D56F4")
	SelectedColor = lipgloss.Color("#F780E2")
	NormalColor   = lipgloss.Color("#FAFAFA")
	DimColor      = lipgloss.Color("#626262")
	BorderColor   = lipgloss.Color("#383838")
	AccentColor   = lipgloss.Color("#00D9FF")
	SuccessColor  = lipgloss.Color("#04B575")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(SelectedColor).
				Background(lipgloss.Color("#2a2a2a")).
				Bold(true).
				Padding(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(NormalColor).
			Padding(0, 1)

	InfoStyle = lipgloss.NewStyle().
			Foreground(BorderColor).
			Bold(true)

	PathStyle = lipgloss.NewStyle().
			Foreground(DimColor).
			Italic(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(DimColor).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(BorderColor).
			Padding(0, 1)

	StatusStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	DimStyle = lipgloss.NewStyle().
			Foreground(DimColor)

	SearchStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(0, 1)

	SearchInputStyle = lipgloss.NewStyle().
				Foreground(NormalColor)
)
