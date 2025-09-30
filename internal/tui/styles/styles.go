package styles

import "github.com/charmbracelet/lipgloss"

var (
	// ErrorStyle is used for displaying error messages
	// Red, bold, with rounded border and padding
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	// TitleStyle is used for panel headers and titles
	// Bright blue and bold
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")) // Bright blue

	// NormalStyle is the default style for normal text
	// White color
	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")) // White
)
