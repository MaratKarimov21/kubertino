package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/tui/styles"
)

const (
	// View mode constants
	viewModeContextSelection = "context_selection"
	viewModeNamespaceView    = "namespace_view"
)

// AppModel is the main Bubble Tea model for the Kubertino TUI
type AppModel struct {
	config               *config.Config
	currentContext       *config.Context
	contexts             []config.Context
	selectedContextIndex int
	viewMode             string // "context_selection" or "namespace_view"
	err                  error
	width                int
	height               int
	keys                 KeyMap
}

// NewAppModel creates a new AppModel with the provided configuration
func NewAppModel(cfg *config.Config) AppModel {
	model := AppModel{
		config:   cfg,
		contexts: cfg.Contexts,
		keys:     DefaultKeyMap(),
	}

	// Initialize viewMode based on number of contexts
	if len(cfg.Contexts) > 1 {
		model.viewMode = viewModeContextSelection
		model.selectedContextIndex = 0
	} else if len(cfg.Contexts) == 1 {
		// Auto-select single context
		model.currentContext = &cfg.Contexts[0]
		model.viewMode = viewModeNamespaceView
	}

	return model
}

// Init initializes the model. Returns nil as no initial commands are needed
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

		// Handle context selection mode navigation
		if m.viewMode == viewModeContextSelection {
			if KeyMatches(msg, m.keys.Up) {
				// Navigate up with wrap-around
				m.selectedContextIndex--
				if m.selectedContextIndex < 0 {
					m.selectedContextIndex = len(m.contexts) - 1
				}
				return m, nil
			}

			if KeyMatches(msg, m.keys.Down) {
				// Navigate down with wrap-around
				m.selectedContextIndex++
				if m.selectedContextIndex >= len(m.contexts) {
					m.selectedContextIndex = 0
				}
				return m, nil
			}

			if KeyMatches(msg, m.keys.Enter) {
				// Select context and transition to namespace view
				m.currentContext = &m.contexts[m.selectedContextIndex]
				m.viewMode = viewModeNamespaceView
				return m, nil
			}
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

	// Render based on current view mode
	if m.viewMode == viewModeContextSelection {
		return m.renderContextList()
	}

	// Default to basic layout (namespace view for Story 1.4)
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

// renderContextList renders the full-screen context selection list
func (m AppModel) renderContextList() string {
	var s string

	// Header
	header := styles.TitleStyle.Render("Select Kubernetes Context")
	s += header + "\n\n"

	// Context list
	for i, ctx := range m.contexts {
		// Determine if this context is selected
		prefix := "  "
		if i == m.selectedContextIndex {
			prefix = "> "
		}

		// Build namespace count suffix
		namespaceCount := ""
		if len(ctx.FavoriteNamespaces) > 0 {
			namespaceCount = styles.DimStyle.Render(fmt.Sprintf(" (%d namespaces)", len(ctx.FavoriteNamespaces)))
		}

		// Render context line with appropriate styling
		if i == m.selectedContextIndex {
			s += styles.SelectedStyle.Render(prefix+ctx.Name) + namespaceCount + "\n"
		} else {
			s += styles.NormalStyle.Render(prefix+ctx.Name) + namespaceCount + "\n"
		}
	}

	// Footer with key hints
	s += "\n"
	footer := styles.DimStyle.Render("↑/↓ Navigate | Enter: Select | ESC/q: Quit")
	s += footer

	return s
}
