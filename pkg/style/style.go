package style

import "github.com/charmbracelet/lipgloss"

var (
	// Focused styles a focused item
	Focused = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

	// Placeholder styles a faded item
	Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// None provides default styling
	None = lipgloss.NewStyle()
)
