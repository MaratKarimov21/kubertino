package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/tui/styles"
)

// AppModel is the main Bubble Tea model for the Kubertino TUI
type AppModel struct {
	config         *config.Config
	currentContext *config.Context
	err            error
	width          int
	height         int
	keys           KeyMap
}

// NewAppModel creates a new AppModel with the provided configuration
func NewAppModel(cfg *config.Config) AppModel {
	return AppModel{
		config: cfg,
		keys:   DefaultKeyMap(),
	}
}

// Init initializes the model. Returns nil as no initial commands are needed for Story 1.4
func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and returns an updated model and optional command
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit keys
		if KeyMatches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Handle terminal resize
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View renders the UI based on the current model state
func (m AppModel) View() string {
	// If there's an error, display it
	if m.err != nil {
		return styles.ErrorStyle.Render(
			"Error: " + m.err.Error() + "\n\n" +
				"Press 'q' or ESC to quit",
		)
	}

	// For Story 1.4, render a basic placeholder UI
	// Detailed panel rendering will be implemented in Epic 2+
	return m.renderBasicLayout()
}

// renderBasicLayout renders a basic split-pane layout scaffold
func (m AppModel) renderBasicLayout() string {
	// Calculate panel dimensions (50% left, 50% right split)
	leftPanelWidth := m.width / 2
	rightPanelWidth := m.width - leftPanelWidth
	rightTopHeight := m.height / 2
	rightBottomHeight := m.height - rightTopHeight

	// Create placeholder content using fmt package for proper integer formatting
	content := "Kubertino TUI Framework Initialized\n\n"
	content += fmt.Sprintf("Terminal Size: %dx%d\n", m.width, m.height)
	content += fmt.Sprintf("Left Panel: %dw\n", leftPanelWidth)
	content += fmt.Sprintf("Right Panel: %dw x %dh / %dh\n\n", rightPanelWidth, rightTopHeight, rightBottomHeight)
	content += "Press 'q', ESC, or Ctrl+C to quit"

	return content
}
