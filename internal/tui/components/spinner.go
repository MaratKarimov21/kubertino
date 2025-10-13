package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerTickMsg is sent when the spinner should advance to the next frame
type SpinnerTickMsg struct{}

// Spinner represents a loading animation
type Spinner struct {
	FrameIndex int
	Message    string
	IsActive   bool
	Frames     []string
}

var (
	// Spinner styling
	spinnerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")) // Pink/magenta color for visibility
)

// NewSpinner creates a new spinner with default frames
func NewSpinner() *Spinner {
	return &Spinner{
		FrameIndex: 0,
		IsActive:   false,
		// Unicode dot spinner frames
		Frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Start activates the spinner with a loading message
func (s *Spinner) Start(message string) {
	s.Message = message
	s.IsActive = true
	s.FrameIndex = 0
}

// Stop deactivates the spinner
func (s *Spinner) Stop() {
	s.IsActive = false
	s.Message = ""
}

// Tick advances the spinner to the next frame
func (s *Spinner) Tick() {
	if s.IsActive && len(s.Frames) > 0 {
		s.FrameIndex = (s.FrameIndex + 1) % len(s.Frames)
	}
}

// View renders the current spinner frame with message
func (s *Spinner) View() string {
	if !s.IsActive {
		return ""
	}

	if len(s.Frames) == 0 {
		return s.Message
	}

	frame := s.Frames[s.FrameIndex]
	return spinnerStyle.Render(frame) + " " + s.Message
}

// TickCmd returns a command that sends a spinner tick message after a delay
func TickCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(100 * time.Millisecond)
		return SpinnerTickMsg{}
	}
}
