package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/tui/styles"
)

const (
	// View mode constants
	viewModeContextSelection = "context_selection"
	viewModeNamespaceView    = "namespace_view"
)

// KubeAdapter is an interface for Kubernetes operations
type KubeAdapter interface {
	GetNamespaces(context string) ([]string, error)
}

// namespaceFetchedMsg is sent when namespaces are fetched
type namespaceFetchedMsg struct {
	namespaces []string
	err        error
}

// AppModel is the main Bubble Tea model for the Kubertino TUI
type AppModel struct {
	config                 *config.Config
	currentContext         *config.Context
	contexts               []config.Context
	selectedContextIndex   int
	viewMode               string // "context_selection" or "namespace_view"
	err                    error
	width                  int
	height                 int
	keys                   KeyMap
	namespaces             []string
	selectedNamespaceIndex int
	namespaceViewportStart int // Starting index for namespace viewport
	namespacesLoading      bool
	namespacesError        error
	kubeAdapter            KubeAdapter
}

// NewAppModel creates a new AppModel with the provided configuration and KubeAdapter
func NewAppModel(cfg *config.Config, adapter KubeAdapter) AppModel {
	model := AppModel{
		config:      cfg,
		contexts:    cfg.Contexts,
		keys:        DefaultKeyMap(),
		kubeAdapter: adapter,
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
	// If single context auto-selected, fetch namespaces immediately
	if m.viewMode == viewModeNamespaceView && m.currentContext != nil {
		return m.fetchNamespacesCmd()
	}
	return nil
}

// fetchNamespacesCmd returns a command that fetches namespaces asynchronously
func (m AppModel) fetchNamespacesCmd() tea.Cmd {
	return func() tea.Msg {
		if m.currentContext == nil {
			return namespaceFetchedMsg{err: fmt.Errorf("no context selected")}
		}

		slog.Info("fetching namespaces", "context", m.currentContext.Name)
		namespaces, err := m.kubeAdapter.GetNamespaces(m.currentContext.Name)

		if err != nil {
			slog.Error("namespace fetch failed", "context", m.currentContext.Name, "error", err)
			return namespaceFetchedMsg{err: err}
		}

		slog.Info("namespaces fetched", "context", m.currentContext.Name, "count", len(namespaces))
		return namespaceFetchedMsg{namespaces: namespaces}
	}
}

// Update handles incoming messages and returns an updated model and optional command
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case namespaceFetchedMsg:
		// Handle namespace fetch results
		m.namespacesLoading = false
		if msg.err != nil {
			m.namespacesError = msg.err
			return m, nil
		}
		m.namespacesError = nil
		m.namespaces = msg.namespaces

		// Sort namespaces: favorites first
		if m.currentContext != nil {
			m.namespaces = m.sortNamespacesWithFavorites(m.namespaces, m.currentContext.FavoriteNamespaces)
		}

		return m, nil

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
				m.namespacesLoading = true
				m.selectedNamespaceIndex = 0 // Reset namespace selection
				m.namespaceViewportStart = 0 // Reset viewport position
				m.namespaces = nil           // Clear previous namespaces
				return m, m.fetchNamespacesCmd()
			}
		}

		// Handle namespace view navigation
		if m.viewMode == viewModeNamespaceView {
			if KeyMatches(msg, m.keys.Up) {
				// Navigate up with wrap-around
				if len(m.namespaces) > 0 {
					m.selectedNamespaceIndex--
					if m.selectedNamespaceIndex < 0 {
						m.selectedNamespaceIndex = len(m.namespaces) - 1
						// Wrap to end: adjust viewport to show last item
						availableHeight := m.height - 4
						if availableHeight < 1 {
							availableHeight = 10
						}
						if len(m.namespaces) > availableHeight {
							m.namespaceViewportStart = len(m.namespaces) - availableHeight
						}
					} else if m.selectedNamespaceIndex < m.namespaceViewportStart {
						// Scroll up: selected item went above viewport
						m.namespaceViewportStart = m.selectedNamespaceIndex
					}
				}
				return m, nil
			}

			if KeyMatches(msg, m.keys.Down) {
				// Navigate down with wrap-around
				if len(m.namespaces) > 0 {
					m.selectedNamespaceIndex++
					if m.selectedNamespaceIndex >= len(m.namespaces) {
						m.selectedNamespaceIndex = 0
						// Wrap to start: reset viewport to beginning
						m.namespaceViewportStart = 0
					} else {
						// Check if we need to scroll down
						availableHeight := m.height - 4
						if availableHeight < 1 {
							availableHeight = 10
						}
						viewportEnd := m.namespaceViewportStart + availableHeight
						if m.selectedNamespaceIndex >= viewportEnd {
							// Scroll down: selected item went below viewport
							m.namespaceViewportStart = m.selectedNamespaceIndex - availableHeight + 1
						}
					}
				}
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

// sortNamespacesWithFavorites sorts namespaces with favorites first
func (m AppModel) sortNamespacesWithFavorites(namespaces []string, favorites []string) []string {
	// Create a set of favorites for quick lookup
	favSet := make(map[string]bool)
	for _, fav := range favorites {
		favSet[fav] = true
	}

	// Separate favorites and non-favorites
	var favs, nonFavs []string
	for _, ns := range namespaces {
		if favSet[ns] {
			favs = append(favs, ns)
		} else {
			nonFavs = append(nonFavs, ns)
		}
	}

	// Combine: favorites first, then rest
	result := make([]string, 0, len(namespaces))
	result = append(result, favs...)
	result = append(result, nonFavs...)

	return result
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

	if m.viewMode == viewModeNamespaceView {
		return m.renderNamespaceList()
	}

	// Default to basic layout (shouldn't reach here)
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

// renderNamespaceList renders the namespace list panel
func (m AppModel) renderNamespaceList() string {
	var s string

	// Header with namespace count
	header := styles.TitleStyle.Render(fmt.Sprintf("Namespaces (%d)", len(m.namespaces)))
	if m.currentContext != nil {
		header += styles.DimStyle.Render(fmt.Sprintf(" - %s", m.currentContext.Name))
	}
	s += header + "\n\n"

	// Show loading indicator
	if m.namespacesLoading {
		s += styles.DimStyle.Render("Loading namespaces...") + "\n"
		s += "\n"
		s += styles.DimStyle.Render("ESC/q: Quit")
		return s
	}

	// Show error if present
	if m.namespacesError != nil {
		errorMsg := fmt.Sprintf("Error fetching namespaces: %v", m.namespacesError)
		s += styles.ErrorStyle.Render(errorMsg) + "\n"
		s += "\n"
		s += styles.DimStyle.Render("ESC/q: Quit")
		return s
	}

	// Render namespace list
	if len(m.namespaces) == 0 {
		s += styles.DimStyle.Render("No namespaces found") + "\n"
	} else {
		// Create favorite set for lookup
		favSet := make(map[string]bool)
		if m.currentContext != nil {
			for _, fav := range m.currentContext.FavoriteNamespaces {
				favSet[fav] = true
			}
		}

		// Calculate viewport (visible window)
		// Reserve space for: header (2 lines) + footer (2 lines) + margins
		availableHeight := m.height - 4
		if availableHeight < 1 {
			availableHeight = 10 // Minimum reasonable height
		}

		// Use stored viewport position (updated during navigation)
		start := m.namespaceViewportStart
		end := len(m.namespaces)

		if len(m.namespaces) > availableHeight {
			// List is longer than screen - use viewport
			end = start + availableHeight
			if end > len(m.namespaces) {
				end = len(m.namespaces)
			}
		}

		// Render visible namespaces
		for i := start; i < end; i++ {
			ns := m.namespaces[i]
			prefix := "  "
			if i == m.selectedNamespaceIndex {
				prefix = "> "
			}

			// Add star for favorites
			if favSet[ns] {
				prefix += "★ "
			}

			// Render with appropriate styling
			if i == m.selectedNamespaceIndex {
				s += styles.SelectedStyle.Render(prefix+ns) + "\n"
			} else {
				s += styles.NormalStyle.Render(prefix+ns) + "\n"
			}
		}

		// Show scroll indicators if list is longer than viewport
		if len(m.namespaces) > availableHeight {
			indicator := fmt.Sprintf(" [%d-%d of %d]", start+1, end, len(m.namespaces))
			s += styles.DimStyle.Render(indicator) + "\n"
		}
	}

	// Footer with key hints
	s += "\n"
	footer := styles.DimStyle.Render("↑/↓ Navigate | /: Search | ESC/q: Quit")
	s += footer

	return s
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
