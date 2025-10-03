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

	// SelectedStyle is used for highlighting selected items
	// Bright cyan background with black text and bold
	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).  // Black
			Background(lipgloss.Color("14")). // Bright cyan
			Bold(true)

	// DimStyle is used for secondary/dim text
	// Gray color
	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // Gray

	// HighlightStyle is used for highlighting matching characters in search results
	// Cyan color with bold
	HighlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")). // Bright cyan
			Bold(true)

	// PanelBorderStyle is used for panel borders
	// Purple border with rounded corners and padding
	PanelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")). // Purple
				Padding(1, 2)

	// PanelTitleStyle is used for panel headers/titles within panels
	// Cyan color with bold
	PanelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")). // Cyan
			Padding(0, 1)

	// HeaderStyle is used for the top header bar showing context
	// Yellow text on dark gray background
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")). // Yellow
			Background(lipgloss.Color("235")). // Dark gray
			Padding(0, 2)

	// PlaceholderStyle is used for placeholder text in empty panels
	// Gray color with italic
	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")). // Gray
				Italic(true)

	// Pod status styles
	// RunningStyle is used for running/succeeded pods (green)
	RunningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")). // Bright green
			Bold(true)

	// PendingStyle is used for pending pods (yellow)
	PendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")). // Bright yellow
			Bold(true)

	// FailedStyle is used for failed pods (red)
	FailedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Bright red
			Bold(true)

	// LoadingStyle is used for loading indicators (gray, italic)
	LoadingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Gray
			Italic(true)
)
