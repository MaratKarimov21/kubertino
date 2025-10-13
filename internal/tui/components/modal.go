package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ErrorModal represents an error dialog overlay
type ErrorModal struct {
	Message    string
	Operation  string
	Suggestion string
	RetryFunc  func() tea.Cmd
	IsVisible  bool
	termWidth  int
	termHeight int
}

// Error modal styles
var (
	errorModalStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")). // Red
			Padding(1, 2).
			Width(60)

	modalTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true)

	modalFooterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")) // Dim gray
)

// NewErrorModal creates a new error modal
func NewErrorModal() *ErrorModal {
	return &ErrorModal{
		IsVisible: false,
	}
}

// Show displays the error modal with the given message and operation context
func (e *ErrorModal) Show(message, operation string, retryFunc func() tea.Cmd) {
	e.Message = message
	e.Operation = operation
	e.RetryFunc = retryFunc
	e.IsVisible = true
}

// ShowWithSuggestion displays the error modal with a suggestion for the user
func (e *ErrorModal) ShowWithSuggestion(message, operation, suggestion string, retryFunc func() tea.Cmd) {
	e.Message = message
	e.Operation = operation
	e.Suggestion = suggestion
	e.RetryFunc = retryFunc
	e.IsVisible = true
}

// Hide dismisses the error modal
func (e *ErrorModal) Hide() {
	e.IsVisible = false
	e.Message = ""
	e.Operation = ""
	e.Suggestion = ""
	e.RetryFunc = nil
}

// SetSize updates the terminal dimensions for proper centering
func (e *ErrorModal) SetSize(width, height int) {
	e.termWidth = width
	e.termHeight = height
}

// View renders the error modal overlay
func (e *ErrorModal) View() string {
	if !e.IsVisible {
		return ""
	}

	// Build modal content
	var content string

	// Title
	title := modalTitleStyle.Render("Error: " + e.Operation)
	content = title + "\n\n"

	// Message
	content += e.Message + "\n"

	// Suggestion (if provided)
	if e.Suggestion != "" {
		content += "\n" + e.Suggestion + "\n"
	}

	// Footer with instructions
	content += "\n"
	var footer string
	if e.RetryFunc != nil {
		footer = modalFooterStyle.Render("[Press Enter to retry] [Press ESC to exit]")
	} else {
		footer = modalFooterStyle.Render("[Press ESC to exit]")
	}
	content += footer

	// Apply modal styling
	modal := errorModalStyle.Render(content)

	// Center modal on screen
	if e.termWidth > 0 && e.termHeight > 0 {
		return lipgloss.Place(
			e.termWidth, e.termHeight,
			lipgloss.Center, lipgloss.Center,
			modal,
		)
	}

	// Fallback if terminal size not set
	return modal
}

// HandleKeyPress processes keyboard input for the modal
// Returns true if the key was handled, false otherwise
func (e *ErrorModal) HandleKeyPress(key string) (bool, tea.Cmd) {
	if !e.IsVisible {
		return false, nil
	}

	switch key {
	case "enter":
		// Retry if retry function is available
		if e.RetryFunc != nil {
			retryCmd := e.RetryFunc
			e.Hide()
			return true, retryCmd()
		}
		// No retry func - just hide modal
		e.Hide()
		return true, nil

	case "esc":
		// Dismiss modal
		e.Hide()
		return true, nil
	}

	// Block all other input when modal is visible
	return true, nil
}
