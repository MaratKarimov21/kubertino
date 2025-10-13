package tui

import (
	"fmt"
	"log/slog"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/executor"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/maratkarimov/kubertino/internal/search"
	"github.com/maratkarimov/kubertino/internal/tui/components"
	"github.com/maratkarimov/kubertino/internal/tui/styles"
)

const (
	// View mode constants
	viewModeContextSelection = "context_selection"
	viewModeNamespaceView    = "namespace_view"

	// Terminal size constraints
	MinTerminalWidth  = 80
	MinTerminalHeight = 24
	HeaderHeight      = 0 // No header displayed (Story 6.1)
)

// KubeAdapter is an interface for Kubernetes operations
type KubeAdapter interface {
	GetNamespaces(context string) ([]string, error)
	GetPods(context, namespace string) ([]k8s.Pod, error)
}

// namespaceFetchedMsg is sent when namespaces are fetched
type namespaceFetchedMsg struct {
	namespaces []string
	err        error
}

// podsFetchedMsg is sent when pods are fetched
type podsFetchedMsg struct {
	pods []k8s.Pod
	err  error
}

// execFinishedMsg is sent when an external command execution finishes
type execFinishedMsg struct {
	err error
}

// PanelType represents which panel has keyboard focus
type PanelType int

const (
	PanelNamespaces PanelType = iota
	PanelPods
	PanelActions // Story 4.1
)

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
	// Search mode fields
	searchMode         bool
	searchQuery        string
	filteredNamespaces []string
	// Pod state fields
	pods             []k8s.Pod
	podsLoading      bool
	podsError        error
	currentNamespace string
	// Terminal size fields
	termWidth        int
	termHeight       int
	terminalTooSmall bool
	// Focus and navigation state (Story 3.3)
	focusedPanel     PanelType // Which panel has keyboard focus
	selectedPodIndex int       // Index of selected pod in pods slice (-1 if none, Story 6.2: cursor position = selection)
	podScrollOffset  int       // Scroll offset for long pod lists
	// Actions state (Story 4.1)
	actions []config.Action // Actions for current context
	// Executor and error state (Story 4.2)
	executor     *executor.Executor
	errorMessage string // Error message to display in TUI (deprecated in Story 6.3, use errorModal)
	// Favorites state (Story 5.3)
	favoriteNamespaces []string // Favorite namespaces for current context
	// Error modal and spinners (Story 6.3)
	errorModal        *components.ErrorModal
	namespacesSpinner *components.Spinner
	podsSpinner       *components.Spinner
	actionSpinner     *components.Spinner
}

// NewAppModel creates a new AppModel with the provided configuration and KubeAdapter
func NewAppModel(cfg *config.Config, adapter KubeAdapter) AppModel {
	model := AppModel{
		config:            cfg,
		contexts:          cfg.Contexts,
		keys:              DefaultKeyMap(),
		kubeAdapter:       adapter,
		focusedPanel:      PanelNamespaces,            // Story 3.3: Start with namespace panel focused
		selectedPodIndex:  -1,                         // Story 6.2: No pod selected initially (cursor = selection)
		podScrollOffset:   0,                          // Story 3.3: No scroll offset initially
		executor:          executor.NewExecutor(),     // Story 6.2: Initialize executor (no adapter needed)
		errorModal:        components.NewErrorModal(), // Story 6.3: Initialize error modal
		namespacesSpinner: components.NewSpinner(),    // Story 6.3: Initialize namespace spinner
		podsSpinner:       components.NewSpinner(),    // Story 6.3: Initialize pod spinner
		actionSpinner:     components.NewSpinner(),    // Story 6.3: Initialize action spinner
	}

	// Initialize viewMode based on number of contexts
	if len(cfg.Contexts) > 1 {
		model.viewMode = viewModeContextSelection
		model.selectedContextIndex = 0
	} else if len(cfg.Contexts) == 1 {
		// Auto-select single context
		model.currentContext = &cfg.Contexts[0]
		model.viewMode = viewModeNamespaceView
		// Load actions from context (Story 4.1)
		model.actions = cfg.Contexts[0].Actions
	}

	return model
}

// Init initializes the model. Returns nil as no initial commands are needed
func (m AppModel) Init() tea.Cmd {
	// If single context auto-selected, fetch namespaces immediately
	if m.viewMode == viewModeNamespaceView && m.currentContext != nil {
		// Story 6.3: Start namespace spinner
		m.namespacesSpinner.Start("Loading namespaces...")
		return tea.Batch(m.fetchNamespacesCmd(), components.TickCmd())
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

		return namespaceFetchedMsg{namespaces: namespaces}
	}
}

// fetchPodsCmd returns a command that fetches pods asynchronously
func (m AppModel) fetchPodsCmd() tea.Cmd {
	return func() tea.Msg {
		if m.currentContext == nil {
			return podsFetchedMsg{err: fmt.Errorf("no context selected")}
		}

		if m.currentNamespace == "" {
			return podsFetchedMsg{err: fmt.Errorf("no namespace selected")}
		}

		slog.Info("fetching pods", "context", m.currentContext.Name, "namespace", m.currentNamespace)
		pods, err := m.kubeAdapter.GetPods(m.currentContext.Name, m.currentNamespace)

		if err != nil {
			slog.Error("pod fetch failed", "context", m.currentContext.Name, "namespace", m.currentNamespace, "error", err)
			return podsFetchedMsg{err: err}
		}

		return podsFetchedMsg{pods: pods}
	}
}

// Update handles incoming messages and returns an updated model and optional command
// nolint:gocyclo // Bubble Tea Update pattern inherently has high complexity due to message routing
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case components.SpinnerTickMsg:
		// Story 6.3: Handle spinner animation tick
		m.namespacesSpinner.Tick()
		m.podsSpinner.Tick()
		m.actionSpinner.Tick()

		// Re-subscribe if any spinner is active
		if m.namespacesSpinner.IsActive || m.podsSpinner.IsActive || m.actionSpinner.IsActive {
			return m, components.TickCmd()
		}
		return m, nil

	case namespaceFetchedMsg:
		// Handle namespace fetch results (Story 6.3: use spinners and modal)
		m.namespacesLoading = false
		m.namespacesSpinner.Stop()

		if msg.err != nil {
			m.namespacesError = msg.err
			// Story 6.3: Show error modal with retry capability
			m.errorModal.ShowWithSuggestion(
				msg.err.Error(),
				"Fetch Namespaces",
				"Check your network connection and cluster access",
				func() tea.Cmd { return m.fetchNamespacesCmd() },
			)
			return m, nil
		}
		m.namespacesError = nil

		// Story 5.3: Get favorites for current context
		if m.currentContext != nil {
			favorites, err := config.GetFavorites(m.config, m.currentContext.Name)
			if err != nil {
				slog.Warn("failed to get favorites", "context", m.currentContext.Name, "error", err)
				m.favoriteNamespaces = []string{}
			} else {
				m.favoriteNamespaces = favorites
			}
		}

		// Story 5.3: Sort namespaces with favorites first
		m.namespaces = m.sortNamespacesWithFavorites(msg.namespaces, m.favoriteNamespaces)

		return m, nil

	case podsFetchedMsg:
		// Handle pod fetch results (Story 6.3: use spinners and modal)
		m.podsLoading = false
		m.podsSpinner.Stop()

		if msg.err != nil {
			m.podsError = msg.err
			// Story 6.3: Show error modal with retry capability
			m.errorModal.ShowWithSuggestion(
				msg.err.Error(),
				"Fetch Pods",
				"Check your network connection and cluster access",
				func() tea.Cmd { return m.fetchPodsCmd() },
			)
			return m, nil
		}
		m.podsError = nil
		m.pods = msg.pods

		// Story 6.2: No default pod pattern matching
		// User must manually select a pod

		return m, nil

	case execFinishedMsg:
		// Handle command execution completion (Story 6.3: use modal for errors)
		m.actionSpinner.Stop()

		if msg.err != nil {
			m.errorModal.Show(
				fmt.Sprintf("Command failed: %s", msg.err.Error()),
				"Action Execution",
				nil,
			)
		}
		return m, nil

	case tea.KeyMsg:
		// Story 6.3: Handle error modal key presses first (blocks other input)
		if m.errorModal.IsVisible {
			// Bug Fix: Capture operation BEFORE HandleKeyPress clears it
			operation := m.errorModal.Operation

			handled, cmd := m.errorModal.HandleKeyPress(msg.String())
			if handled {
				// QA Fix: ESC should exit app, not just dismiss (user testing feedback)
				if msg.String() == "esc" {
					return m, tea.Quit
				}
				// Modal handled the key - start spinner if retrying
				if cmd != nil && msg.String() == "enter" {
					// QA Fix: Restart appropriate spinner based on operation being retried
					// Bug Fix: Clear error state and set loading state when retrying
					slog.Debug("retry operation", "operation", operation)
					switch operation {
					case "Fetch Namespaces":
						slog.Debug("clearing namespace error and starting spinner")
						m.namespacesError = nil
						m.namespacesLoading = true
						m.namespacesSpinner.Start("Loading namespaces...")
					case "Fetch Pods":
						slog.Debug("clearing pod error and starting spinner")
						m.podsError = nil
						m.podsLoading = true
						m.podsSpinner.Start("Loading pods...")
					case "Action Execution":
						// Action spinner already handled in handleActionExecution
					}
					return m, tea.Batch(cmd, components.TickCmd())
				}
				return m, cmd
			}
		}

		// Clear error message on any key press (Story 4.2)
		if m.errorMessage != "" {
			m.errorMessage = ""
			return m, nil
		}

		// Handle search mode ESC first (before quit keys)
		if m.searchMode && msg.Type == tea.KeyEsc {
			m.deactivateSearch()
			return m, nil
		}

		// Handle quit keys (but not in search mode where ESC is handled above)
		if KeyMatches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		// Check for action shortcut key presses (Story 4.2)
		// Only in namespace view mode and not in search mode
		if m.viewMode == viewModeNamespaceView && !m.searchMode {
			keyStr := msg.String()
			for _, action := range m.actions {
				if keyStr == action.Shortcut {
					return m.handleActionExecution(action)
				}
			}
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
				// Story 5.3: Clear favorites (will be loaded with namespaces)
				m.favoriteNamespaces = nil
				// Load actions from context (Story 6.2)
				m.actions = m.currentContext.Actions
				// Story 6.3: Start namespace spinner
				m.namespacesSpinner.Start("Loading namespaces...")
				return m, tea.Batch(m.fetchNamespacesCmd(), components.TickCmd())
			}
		}

		// Handle namespace view navigation
		if m.viewMode == viewModeNamespaceView {
			// Handle Tab key for focus switching (Story 6.2: Skip actions panel)
			if msg.String() == "tab" {
				switch m.focusedPanel {
				case PanelNamespaces:
					m.focusedPanel = PanelPods
					// Auto-select first pod when focusing pod panel
					if len(m.pods) > 0 && m.selectedPodIndex == -1 {
						m.selectedPodIndex = 0
					}
				case PanelPods:
					m.focusedPanel = PanelNamespaces
				}
				return m, nil
			}

			// Handle Shift+Tab key for backward focus switching (Story 6.2: Skip actions panel)
			if msg.String() == "shift+tab" {
				switch m.focusedPanel {
				case PanelPods:
					m.focusedPanel = PanelNamespaces
				case PanelNamespaces:
					m.focusedPanel = PanelPods
					// Auto-select first pod when focusing pod panel
					if len(m.pods) > 0 && m.selectedPodIndex == -1 {
						m.selectedPodIndex = 0
					}
				}
				return m, nil
			}

			// Handle search mode activation
			if !m.searchMode && msg.String() == "/" {
				m.activateSearch()
				return m, nil
			}

			// Handle search mode input
			if m.searchMode {
				// Enter selects current filtered namespace and exits search
				if KeyMatches(msg, m.keys.Enter) {
					if len(m.filteredNamespaces) > 0 && m.selectedNamespaceIndex < len(m.filteredNamespaces) {
						// Select namespace and fetch pods
						m.currentNamespace = m.filteredNamespaces[m.selectedNamespaceIndex]
						m.deactivateSearch()
						m.podsLoading = true
						m.podsError = nil
						m.pods = nil
						// Reset pod panel state (Story 6.2)
						m.selectedPodIndex = -1
						m.podScrollOffset = 0
						// QA Fix: Auto-switch focus to pods panel after namespace selection
						m.focusedPanel = PanelPods
						// Story 6.3: Start pod spinner
						m.podsSpinner.Start("Loading pods...")
						return m, tea.Batch(m.fetchPodsCmd(), components.TickCmd())
					}
					return m, nil
				}

				// Backspace removes last character
				if msg.Type == tea.KeyBackspace {
					if len(m.searchQuery) > 0 {
						m.updateSearchQuery(m.searchQuery[:len(m.searchQuery)-1])
					}
					return m, nil
				}

				// Handle regular character input (alphanumeric, dash, dot, underscore)
				if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 {
					r := msg.Runes[0]
					// Accept alphanumeric, dash, dot, underscore
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '.' || r == '_' {
						m.updateSearchQuery(m.searchQuery + string(r))
					}
					return m, nil
				}
			}

			// Handle Enter key in normal mode (namespace selection)
			if !m.searchMode && KeyMatches(msg, m.keys.Enter) {
				// Story 6.2: Cursor position = pod selection (no Enter confirmation needed)
				// Enter only used for namespace selection
				if m.focusedPanel == PanelNamespaces && len(m.namespaces) > 0 && m.selectedNamespaceIndex < len(m.namespaces) {
					// Select namespace and fetch pods
					m.currentNamespace = m.namespaces[m.selectedNamespaceIndex]
					m.podsLoading = true
					m.podsError = nil
					m.pods = nil
					// Reset pod panel state (Story 6.2)
					m.selectedPodIndex = -1
					m.podScrollOffset = 0
					// QA Fix: Auto-switch focus to pods panel after namespace selection
					m.focusedPanel = PanelPods
					// Story 6.3: Start pod spinner
					m.podsSpinner.Start("Loading pods...")
					return m, tea.Batch(m.fetchPodsCmd(), components.TickCmd())
				}
				return m, nil
			}

			// Navigation (works in both normal and search mode)
			// Handle arrow keys based on focused panel (Story 6.2: Actions panel removed)
			if KeyMatches(msg, m.keys.Up) {
				switch m.focusedPanel {
				case PanelNamespaces:
					// Navigate namespace panel with cursor centering (Story 6.1)
					navList := m.namespaces
					if m.searchMode && m.filteredNamespaces != nil {
						navList = m.filteredNamespaces
					}

					if len(navList) > 0 {
						m.selectedNamespaceIndex--
						if m.selectedNamespaceIndex < 0 {
							// Wrap to end: adjust viewport to show last item
							m.selectedNamespaceIndex = len(navList) - 1
							// Story 6.1: Adjust viewport to show last item with proper centering
							m.adjustNamespaceViewport(len(navList))
						} else {
							// Story 6.1: Implement cursor centering
							m.adjustNamespaceViewport(len(navList))
						}
					}
				case PanelPods:
					// Navigate pod panel (Story 3.3)
					if len(m.pods) > 0 && m.selectedPodIndex > 0 {
						m.selectedPodIndex--
						m.adjustPodScrollOffset()
					}
				}
				return m, nil
			}

			if KeyMatches(msg, m.keys.Down) {
				switch m.focusedPanel {
				case PanelNamespaces:
					// Navigate namespace panel with cursor centering (Story 6.1)
					navList := m.namespaces
					if m.searchMode && m.filteredNamespaces != nil {
						navList = m.filteredNamespaces
					}

					if len(navList) > 0 {
						m.selectedNamespaceIndex++
						if m.selectedNamespaceIndex >= len(navList) {
							// Wrap to start: reset viewport to beginning
							m.selectedNamespaceIndex = 0
							m.namespaceViewportStart = 0
						} else {
							// Story 6.1: Implement cursor centering
							m.adjustNamespaceViewport(len(navList))
						}
					}
				case PanelPods:
					// Navigate pod panel (Story 3.3)
					if len(m.pods) > 0 && m.selectedPodIndex < len(m.pods)-1 {
						m.selectedPodIndex++
						m.adjustPodScrollOffset()
					}
				}
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		// Handle terminal resize
		m.width = msg.Width
		m.height = msg.Height
		m.termWidth = msg.Width
		m.termHeight = msg.Height

		// Story 6.3: Update error modal size for proper centering
		m.errorModal.SetSize(msg.Width, msg.Height)

		// Check minimum size
		if msg.Width < MinTerminalWidth || msg.Height < MinTerminalHeight {
			m.terminalTooSmall = true
		} else {
			m.terminalTooSmall = false
		}
		return m, nil
	}

	return m, nil
}

// sortNamespacesWithFavorites sorts namespaces with favorites first
// Story 6.1: Favorites preserve config order (not sorted), regular namespaces sorted alphabetically
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

	// Story 6.1: Preserve config order for favorites (don't sort)
	// Sort favorites by their position in the original favorites slice
	favOrder := make(map[string]int)
	for i, fav := range favorites {
		favOrder[fav] = i
	}
	sort.SliceStable(favs, func(i, j int) bool {
		return favOrder[favs[i]] < favOrder[favs[j]]
	})

	// Sort only non-favorites alphabetically
	sort.Strings(nonFavs)

	// Combine: favorites first (in config order), then rest (alphabetically)
	result := make([]string, 0, len(namespaces))
	result = append(result, favs...)
	result = append(result, nonFavs...)

	return result
}

// activateSearch enables search mode and initializes filtered list
func (m *AppModel) activateSearch() {
	slog.Debug("search mode activated")
	m.searchMode = true
	m.searchQuery = ""
	m.filteredNamespaces = m.namespaces
	m.selectedNamespaceIndex = 0
	m.namespaceViewportStart = 0
}

// deactivateSearch disables search mode and clears search state
func (m *AppModel) deactivateSearch() {
	slog.Debug("search mode deactivated")
	m.searchMode = false
	m.searchQuery = ""
	m.filteredNamespaces = nil
	// Reset selected index to 0 (don't try to preserve position)
	m.selectedNamespaceIndex = 0
	m.namespaceViewportStart = 0
}

// updateSearchQuery updates the search query and filters namespaces
func (m *AppModel) updateSearchQuery(query string) {
	m.searchQuery = query

	if query == "" {
		// Empty query shows all namespaces
		m.filteredNamespaces = m.namespaces
	} else {
		// Perform fuzzy search (will be implemented next)
		m.filteredNamespaces = m.performFuzzySearch(query)
	}

	// Reset selection to first result
	m.selectedNamespaceIndex = 0
	m.namespaceViewportStart = 0

	slog.Debug("search query updated", "query", query, "results", len(m.filteredNamespaces))
}

// performFuzzySearch performs fuzzy search and returns filtered namespace list
func (m *AppModel) performFuzzySearch(query string) []string {
	// Convert string slice to Namespace slice for fuzzy matching
	namespaces := make([]k8s.Namespace, len(m.namespaces))
	for i, ns := range m.namespaces {
		namespaces[i] = k8s.Namespace{Name: ns}
	}

	// Perform fuzzy search
	matches := search.FuzzyMatch(query, namespaces)

	// Convert matches back to string slice
	results := make([]string, len(matches))
	for i, match := range matches {
		results[i] = match.Namespace.Name
	}

	return results
}

// adjustPodScrollOffset adjusts the pod scroll offset based on selected pod index (Story 3.3)
func (m *AppModel) adjustPodScrollOffset() {
	// Calculate visible window size based on pod panel height
	// Pod panel height is roughly half the available height
	availableHeight := m.termHeight - HeaderHeight
	podPanelHeight := availableHeight / 2

	// Reserve space for: border (2) + padding (2) + title (1) + blank (1) + help text (2) = 8 lines
	visibleHeight := podPanelHeight - 8
	if visibleHeight < 1 {
		visibleHeight = 5 // Minimum visible items
	}

	// Scroll down if selection below visible window
	if m.selectedPodIndex >= m.podScrollOffset+visibleHeight {
		m.podScrollOffset = m.selectedPodIndex - visibleHeight + 1
	}

	// Scroll up if selection above visible window
	if m.selectedPodIndex < m.podScrollOffset {
		m.podScrollOffset = m.selectedPodIndex
	}
}

// adjustNamespaceViewport adjusts the namespace viewport with cursor centering (Story 6.1)
func (m *AppModel) adjustNamespaceViewport(listLength int) {
	// Calculate available height for namespace list
	// CRITICAL: This must match the calculation in renderNamespacePanel!
	// renderNamespacePanel passes: contentHeight = termHeight - 4 (border+padding)
	// to renderNamespaceList

	// Get the panel height (full terminal height for namespace panel)
	panelHeight := m.termHeight
	if panelHeight == 0 {
		panelHeight = 24 // Default for tests
	}

	// Subtract border+padding that renderNamespacePanel removes (4 lines total)
	effectiveHeight := panelHeight - 4
	if effectiveHeight < 1 {
		effectiveHeight = 20 // Minimum
	}

	headerLines := 2
	footerLines := 2
	scrollIndicatorLines := 1
	searchBoxLines := 0
	if m.searchMode {
		searchBoxLines = 2
	}

	reservedLines := headerLines + footerLines + searchBoxLines
	availableHeight := effectiveHeight - reservedLines - scrollIndicatorLines
	if availableHeight < 1 {
		availableHeight = 1
	}

	// If list fits in viewport, no scrolling needed
	if listLength <= availableHeight {
		m.namespaceViewportStart = 0
		return
	}

	// Story 6.1: Implement cursor centering with strict bounds checking
	// Calculate middle position
	middlePosition := availableHeight / 2

	// Calculate the maximum valid viewport start position
	maxViewportStart := listLength - availableHeight
	if maxViewportStart < 0 {
		maxViewportStart = 0
	}

	// Determine viewport position based on cursor position
	if m.selectedNamespaceIndex < middlePosition {
		// Near top: cursor moves freely, viewport stays at top
		m.namespaceViewportStart = 0
	} else if m.selectedNamespaceIndex >= listLength-middlePosition {
		// Near bottom: lock viewport to show last items, cursor moves within bottom section
		m.namespaceViewportStart = maxViewportStart
	} else {
		// Middle section: keep cursor centered in viewport
		m.namespaceViewportStart = m.selectedNamespaceIndex - middlePosition
	}

	// Final safety check: absolutely ensure cursor is within viewport bounds
	// This handles any edge cases where the above logic might fail
	viewportEnd := m.namespaceViewportStart + availableHeight

	// If cursor is above viewport, scroll up to show it
	if m.selectedNamespaceIndex < m.namespaceViewportStart {
		m.namespaceViewportStart = m.selectedNamespaceIndex
	}

	// If cursor is below viewport, scroll down to show it
	// CRITICAL: cursor must be strictly within [start, end), not >= end
	if m.selectedNamespaceIndex >= viewportEnd {
		m.namespaceViewportStart = m.selectedNamespaceIndex - availableHeight + 1
	}

	// Ensure viewport doesn't go negative
	if m.namespaceViewportStart < 0 {
		m.namespaceViewportStart = 0
	}

	// Ensure viewport doesn't exceed maximum
	if m.namespaceViewportStart > maxViewportStart {
		m.namespaceViewportStart = maxViewportStart
	}
}

// getMatchIndices returns the match indices for a given namespace in the current search
// Returns nil if not in search mode or namespace not found
func (m *AppModel) getMatchIndices(namespaceName string) []int {
	if !m.searchMode || m.searchQuery == "" {
		return nil
	}

	// Convert string slice to Namespace slice
	namespaces := make([]k8s.Namespace, len(m.namespaces))
	for i, ns := range m.namespaces {
		namespaces[i] = k8s.Namespace{Name: ns}
	}

	// Perform fuzzy search
	matches := search.FuzzyMatch(m.searchQuery, namespaces)

	// Find matching indices for this namespace
	for _, match := range matches {
		if match.Namespace.Name == namespaceName {
			return match.MatchIndices
		}
	}

	return nil
}

// handleActionExecution executes an action using tea.ExecProcess (Story 6.2)
func (m AppModel) handleActionExecution(action config.Action) (tea.Model, tea.Cmd) {
	// Story 6.2: ALL actions execute locally with template substitution

	// Ensure we have a current context and namespace (Story 6.3: use modal for errors)
	if m.currentContext == nil {
		m.errorModal.Show("No context selected", "Execute Action", nil)
		return m, nil
	}

	if m.currentNamespace == "" {
		m.errorModal.Show("No namespace selected", "Execute Action", nil)
		return m, nil
	}

	if len(m.pods) == 0 {
		m.errorModal.Show("No pods available in namespace", "Execute Action", nil)
		return m, nil
	}

	// Story 6.2: Check if pod is selected (cursor position = selection)
	if m.selectedPodIndex == -1 {
		m.errorModal.ShowWithSuggestion(
			"No pod selected",
			"Execute Action",
			"Press Tab to focus pod panel, then use arrow keys to select a pod",
			nil,
		)
		return m, nil
	}

	// Get the selected pod
	selectedPod := m.pods[m.selectedPodIndex]

	// Prepare the local command using executor (Story 6.2: all actions are local now)
	cmd, err := m.executor.PrepareLocal(action, *m.currentContext, m.currentNamespace, selectedPod, m.config.Kubeconfig)
	if err != nil {
		m.errorModal.Show(fmt.Sprintf("Action failed: %s", err.Error()), "Execute Action", nil)
		return m, nil
	}

	// Story 6.3: Start action spinner before executing
	m.actionSpinner.Start(fmt.Sprintf("Executing %s...", action.Name))

	// Use tea.ExecProcess to suspend TUI and run command
	// This gives full terminal control to the command
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return execFinishedMsg{err: err}
	})
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
		// Check terminal size before rendering
		if m.terminalTooSmall {
			return m.renderTerminalTooSmallWarning()
		}
		return m.renderSplitLayout()
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
// panelHeight is the actual height available for the namespace list (0 = use m.height)
func (m AppModel) renderNamespaceList(panelHeight int) string {
	var s string

	// Use provided panel height or fall back to full terminal height
	effectiveHeight := panelHeight
	if effectiveHeight == 0 {
		effectiveHeight = m.height
	}
	// If still 0 (tests don't set terminal size), use a reasonable default
	if effectiveHeight == 0 {
		effectiveHeight = 20
	}

	// Determine which list to render
	renderList := m.namespaces
	listCount := len(m.namespaces)
	if m.searchMode && m.filteredNamespaces != nil {
		renderList = m.filteredNamespaces
		listCount = len(m.filteredNamespaces)
	}

	// Header with namespace count
	header := styles.TitleStyle.Render(fmt.Sprintf("Namespaces (%d)", len(m.namespaces)))
	if m.currentContext != nil {
		header += styles.DimStyle.Render(fmt.Sprintf(" - %s", m.currentContext.Name))
	}
	s += header + "\n\n"

	// Show loading indicator (Story 6.3: use spinner)
	if m.namespacesLoading {
		spinnerView := m.namespacesSpinner.View()
		if spinnerView != "" {
			s += spinnerView + "\n"
		} else {
			s += styles.DimStyle.Render("Loading namespaces...") + "\n"
		}
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
	} else if m.searchMode && len(m.filteredNamespaces) == 0 {
		// Show "no matches" message when search returns empty
		s += styles.DimStyle.Render("No matches found") + "\n"
	} else {
		// Story 5.3: Build favorites set for marking
		favSet := make(map[string]bool)
		for _, fav := range m.favoriteNamespaces {
			favSet[fav] = true
		}

		// Calculate viewport (visible window)
		// Count exact lines used by UI elements:
		// - Header: "Namespaces (N)" + blank line = 2 lines
		// - Each namespace item: 1 line
		// - Scroll indicator (if needed): 1 line
		// - Blank line before footer: 1 line
		// - Search box (if active): blank + search line = 2 lines
		// - Footer: 1 line

		headerLines := 2
		footerLines := 2          // blank line + footer text
		scrollIndicatorLines := 1 // will be shown if list > available
		searchBoxLines := 0
		if m.searchMode {
			searchBoxLines = 2 // blank + search input
		}

		// Reserve space for fixed UI elements
		reservedLines := headerLines + footerLines + searchBoxLines

		// Calculate available space for namespace items
		// We'll add scroll indicator line if needed, so reserve 1 more
		availableHeight := effectiveHeight - reservedLines - scrollIndicatorLines
		if availableHeight < 1 {
			availableHeight = 1 // Minimum: show at least 1 namespace
		}

		// DEBUG: Log to see actual values
		slog.Debug("viewport calculation",
			"effectiveHeight", effectiveHeight,
			"reservedLines", reservedLines,
			"scrollIndicatorLines", scrollIndicatorLines,
			"availableHeight", availableHeight,
			"listCount", listCount)

		// Use stored viewport position (updated during navigation)
		start := m.namespaceViewportStart
		end := listCount

		if listCount > availableHeight {
			// List is longer than screen - use viewport
			end = start + availableHeight
			if end > listCount {
				end = listCount
			}
		}

		// Render visible namespaces
		for i := start; i < end; i++ {
			ns := renderList[i]
			prefix := "  "
			if i == m.selectedNamespaceIndex {
				prefix = "> "
			}

			// Story 6.1: No star icon for favorites (removed)

			// Render namespace name with highlighting if in search mode
			var renderedName string
			if m.searchMode && m.searchQuery != "" {
				renderedName = m.renderNamespaceWithHighlight(ns, prefix)
			} else {
				renderedName = prefix + ns
			}

			// Story 6.1: Apply selection or favorite styling
			if i == m.selectedNamespaceIndex {
				// Selected item gets selection style (highest priority)
				s += styles.SelectedStyle.Render(renderedName) + "\n"
			} else if favSet[ns] {
				// Favorite namespace gets color highlight (Story 6.1)
				s += styles.FavoriteNamespaceStyle.Render(renderedName) + "\n"
			} else {
				// Regular namespace - no special styling
				s += renderedName + "\n"
			}
		}

		// Show scroll indicators if list is longer than viewport
		if listCount > availableHeight {
			indicator := fmt.Sprintf(" [%d-%d of %d]", start+1, end, listCount)
			s += styles.DimStyle.Render(indicator) + "\n"
		}
	}

	// Search input box (if search mode active)
	if m.searchMode {
		s += "\n"
		searchBox := "Search: " + m.searchQuery + "_"
		s += styles.NormalStyle.Render(searchBox) + "\n"
	}

	// Footer with key hints
	s += "\n"
	if m.searchMode {
		footer := styles.DimStyle.Render("Type to search | ESC: Cancel | Enter: Select")
		s += footer
	} else {
		footer := styles.DimStyle.Render("↑/↓ Navigate | /: Search | ESC/q: Quit")
		s += footer
	}

	return s
}

// renderNamespaceWithHighlight renders a namespace name with matching characters highlighted
func (m AppModel) renderNamespaceWithHighlight(namespace string, prefix string) string {
	matchIndices := m.getMatchIndices(namespace)
	if len(matchIndices) == 0 {
		return prefix + namespace
	}

	// Create a set of match indices for quick lookup
	matchSet := make(map[int]bool)
	for _, idx := range matchIndices {
		matchSet[idx] = true
	}

	// Build the highlighted string
	result := prefix
	for i, char := range namespace {
		if matchSet[i] {
			// Highlight matching character
			result += styles.HighlightStyle.Render(string(char))
		} else {
			// Normal character
			result += string(char)
		}
	}

	return result
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

		// After Story 5.1: FavoriteNamespaces moved to Config.Favorites
		// Namespace count display temporarily disabled
		namespaceCount := ""

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

// renderSplitLayout renders the split-pane layout with header, namespace, pods, and actions panels
func (m AppModel) renderSplitLayout() string {
	// Calculate dimensions (no header, use full height)
	availableHeight := m.termHeight
	leftWidth := m.termWidth / 2
	rightWidth := m.termWidth - leftWidth
	rightTopHeight := availableHeight / 2
	rightBottomHeight := availableHeight - rightTopHeight

	// Render panels
	namespacePanel := m.renderNamespacePanel(leftWidth, availableHeight)
	podPanel := m.renderPodPanel(rightWidth, rightTopHeight)
	actionsPanel := m.renderActionsPanel(rightWidth, rightBottomHeight)

	// Compose layout: combine right panels vertically
	rightSide := lipgloss.JoinVertical(lipgloss.Left, podPanel, actionsPanel)

	// Combine left and right panels horizontally (no header)
	fullLayout := lipgloss.JoinHorizontal(lipgloss.Top, namespacePanel, rightSide)

	// Story 6.3: Removed error bar at bottom - errors now shown via modal

	// Story 6.3: Render error modal overlay on top of everything
	if m.errorModal.IsVisible {
		modalView := m.errorModal.View()
		if modalView != "" {
			// Modal is rendered as an overlay - it will center itself
			return modalView
		}
	}

	return fullLayout
}

// renderHeader renders the header bar with context name
func (m AppModel) renderHeader() string {
	contextName := "None"
	if m.currentContext != nil {
		contextName = m.currentContext.Name
	}
	text := fmt.Sprintf("Context: %s", contextName)
	return styles.HeaderStyle.Width(m.termWidth).Render(text)
}

// renderNamespacePanel renders the namespace list panel with borders
func (m AppModel) renderNamespacePanel(width, height int) string {
	// PanelBorderStyle has Padding(1, 2) and Border (adds 2 lines vertically, 2 chars horizontally)
	// Total additions: vertical = 2 (border) + 2 (padding) = 4 lines
	//                  horizontal = 2 (border) + 4 (padding) = 6 chars

	// Calculate available height for content (what renderNamespaceList will receive)
	// We need to account for what Lip Gloss will add
	// CRITICAL: This must match the effectiveHeight used in adjustNamespaceViewport!
	contentHeight := height - 4 // Subtract border (2) + padding (2)

	// Get the namespace list content with correct panel height
	content := m.renderNamespaceList(contentHeight)

	// Width calculation: Lip Gloss adds border (2) + padding (4) = 6 chars total
	contentWidth := width - 6

	// Select border style based on focus (Story 3.3)
	borderStyle := styles.UnfocusedPanelBorderStyle
	if m.focusedPanel == PanelNamespaces {
		borderStyle = styles.FocusedPanelBorderStyle
	}

	// Don't use MaxHeight - it cuts off borders
	// Instead, renderNamespaceList already limits content to contentHeight
	// Just apply the border and let Lip Gloss add padding + border
	return borderStyle.
		Width(contentWidth).
		Render(content)
}

// renderPodPanel renders the pod panel with real data, loading, or error states
func (m AppModel) renderPodPanel(width, height int) string {
	title := styles.PanelTitleStyle.Render("Pods")

	var content string

	// Loading state (Story 6.3: use spinner)
	if m.podsLoading {
		spinnerView := m.podsSpinner.View()
		if spinnerView != "" {
			content = spinnerView
		} else {
			content = styles.LoadingStyle.Render("Loading pods...")
		}
	} else if m.podsError != nil {
		// Error state
		errorMsg := fmt.Sprintf("Error: %v", m.podsError)
		content = styles.ErrorStyle.Render(errorMsg)
	} else if m.currentNamespace == "" {
		// No namespace selected
		content = styles.PlaceholderStyle.Render("Select a namespace to view pods")
	} else if len(m.pods) == 0 {
		// Empty state
		content = styles.PlaceholderStyle.Render("No pods in this namespace")
	} else {
		// Calculate visible window for scrolling (Story 6.2)
		// Reserve space for: border (2) + padding (2) + title (1) + blank (1) + help text (2) = 8 lines
		visibleHeight := height - 8
		if visibleHeight < 1 {
			visibleHeight = 5 // Minimum visible items
		}

		// Calculate visible pod range
		visiblePods := m.pods[m.podScrollOffset:]
		if len(visiblePods) > visibleHeight {
			visiblePods = visiblePods[:visibleHeight]
		}

		// Render visible pods (Story 6.2: manual selection only)
		var podLines []string
		for i, pod := range visiblePods {
			actualIndex := i + m.podScrollOffset
			statusStyle := m.getPodStatusStyle(pod.Status)
			statusText := statusStyle.Render(pod.Status)

			// Build selection marker (Story 6.2: cursor position = pod selection)
			var marker string
			if actualIndex == m.selectedPodIndex {
				marker = "> "
			} else {
				marker = "  "
			}

			// Apply styling
			podName := pod.Name
			if actualIndex == m.selectedPodIndex {
				// Selected pod gets special highlighting (Story 6.2: cursor = selection)
				podName = styles.SelectedPodStyle.Render(podName)
			}

			line := fmt.Sprintf("%s%-12s %s", marker, statusText, podName)
			podLines = append(podLines, line)
		}
		content = lipgloss.JoinVertical(lipgloss.Left, podLines...)

		// Show scroll indicators (Story 3.3)
		if m.podScrollOffset > 0 {
			scrollUp := styles.HelpTextStyle.Render("↑ More above")
			content = lipgloss.JoinVertical(lipgloss.Left, scrollUp, content)
		}
		if m.podScrollOffset+visibleHeight < len(m.pods) {
			remaining := len(m.pods) - (m.podScrollOffset + visibleHeight)
			scrollDown := styles.HelpTextStyle.Render(fmt.Sprintf("↓ %d more", remaining))
			content = lipgloss.JoinVertical(lipgloss.Left, content, scrollDown)
		}

		// Add help text (Story 6.2)
		helpText := styles.HelpTextStyle.Render("↑/↓: Navigate | Enter: Select | Tab: Switch panel")
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", helpText)
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left, title, "", content)

	// Select border style based on focus (Story 3.3)
	borderStyle := styles.UnfocusedPanelBorderStyle
	if m.focusedPanel == PanelPods {
		borderStyle = styles.FocusedPanelBorderStyle
	}

	// Apply border style with calculated dimensions
	contentWidth := width - 4   // 2 for border + 2*2 for padding
	contentHeight := height - 2 // 2 for border

	return borderStyle.
		Width(contentWidth).
		Height(contentHeight).
		Render(fullContent)
}

// getPodStatusStyle returns the appropriate style for a given pod status
func (m AppModel) getPodStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Running", "Succeeded":
		return styles.RunningStyle
	case "Pending":
		return styles.PendingStyle
	case "Failed":
		return styles.FailedStyle
	default:
		return styles.DimStyle
	}
}

// renderActionsPanel renders the actions panel with multi-column layout (Story 6.2)
func (m AppModel) renderActionsPanel(width, height int) string {
	title := styles.PanelTitleStyle.Render("Actions")

	var content string

	if len(m.actions) == 0 {
		// Empty state
		content = styles.PlaceholderStyle.Render("No actions configured")
	} else {
		// Story 6.2: Multi-column layout (no scrolling, no focus)
		// 2 columns for width < 80, 3 columns for width >= 80
		columnCount := 2
		if width >= 80 {
			columnCount = 3
		}

		itemsPerColumn := (len(m.actions) + columnCount - 1) / columnCount

		var columns []string
		for col := 0; col < columnCount; col++ {
			var columnLines []string
			start := col * itemsPerColumn
			end := start + itemsPerColumn
			if end > len(m.actions) {
				end = len(m.actions)
			}

			for i := start; i < end; i++ {
				action := m.actions[i]
				shortcut := styles.ShortcutStyle.Render(fmt.Sprintf("[%s]", action.Shortcut))
				actionName := styles.ActionStyle.Render(action.Name)
				line := fmt.Sprintf("%s %s", shortcut, actionName)
				columnLines = append(columnLines, line)
			}

			if len(columnLines) > 0 {
				columns = append(columns, lipgloss.JoinVertical(lipgloss.Left, columnLines...))
			}
		}

		content = lipgloss.JoinHorizontal(lipgloss.Top, columns...)

		// Add help text (Story 6.2)
		helpText := styles.HelpTextStyle.Render("[key]: Execute action (works from any panel)")
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", helpText)
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left, title, "", content)

	// Story 6.2: Actions panel is never focused (always use unfocused style)
	borderStyle := styles.UnfocusedPanelBorderStyle

	// Apply border style with calculated dimensions
	contentWidth := width - 4   // 2 for border + 2*2 for padding
	contentHeight := height - 2 // 2 for border

	return borderStyle.
		Width(contentWidth).
		Height(contentHeight).
		Render(fullContent)
}

// renderTerminalTooSmallWarning renders a warning message when terminal is too small
func (m AppModel) renderTerminalTooSmallWarning() string {
	message := fmt.Sprintf(
		"Terminal too small.\nMinimum size: %dx%d\nCurrent: %dx%d",
		MinTerminalWidth, MinTerminalHeight,
		m.termWidth, m.termHeight,
	)
	return styles.PlaceholderStyle.Render(message)
}
