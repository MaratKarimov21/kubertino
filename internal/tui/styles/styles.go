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

	// DefaultPodStyle is used for default pod indicator (Story 3.2)
	// Bright blue/cyan color with bold
	DefaultPodStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")). // Bright blue/cyan
			Bold(true)

	// WarningStyle is used for warning messages (Story 3.2)
	// Orange color with bold
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")). // Orange
			Bold(true)

	// HelpTextStyle is used for help text (Story 3.2)
	// Gray color with italic
	HelpTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")). // Gray
			Italic(true)

	// FocusedPanelBorderStyle is used for focused panel borders (Story 3.3)
	// Bright cyan border with rounded corners and padding
	FocusedPanelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("39")). // Bright cyan
				Padding(1, 2)

	// UnfocusedPanelBorderStyle is used for unfocused panel borders (Story 3.3)
	// Gray border with rounded corners and padding
	UnfocusedPanelBorderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("240")). // Gray
					Padding(1, 2)

	// SelectedPodStyle is used for selected pod highlighting (Story 3.3)
	// Black text on bright cyan background with bold
	SelectedPodStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).  // Black text
				Background(lipgloss.Color("39")). // Bright cyan background
				Bold(true)

	// Action display styles (Story 4.1)

	// ActionStyle is used for normal action display
	// Light gray color
	ActionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")) // Light gray

	// SelectedActionStyle is used for selected action highlighting
	// Black text on bright cyan background with bold
	SelectedActionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).  // Black text
				Background(lipgloss.Color("39")). // Bright cyan background
				Bold(true)

	// GroupHeaderStyle is used for action type group headers
	// Gray color with italic
	GroupHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")). // Gray
				Italic(true)

	// ShortcutStyle is used for shortcut key highlighting
	// Orange/yellow color with bold
	ShortcutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")). // Orange/yellow
			Bold(true)

	// FavoriteNamespaceStyle is used for favorite namespace highlighting (Story 6.1)
	// Yellow/gold color with bold - no icon, just color highlight
	FavoriteNamespaceStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("226")). // Yellow/gold
				Bold(true)
)
