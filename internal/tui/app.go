package tui

import (
	"fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/maratkarimov/kubertino/internal/search"
	"github.com/maratkarimov/kubertino/internal/tui/styles"
)

const (
	// View mode constants
	viewModeContextSelection = "context_selection"
	viewModeNamespaceView    = "namespace_view"

	// Terminal size constraints
	MinTerminalWidth  = 80
	MinTerminalHeight = 24
	HeaderHeight      = 1
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
	// Search mode fields
	searchMode         bool
	searchQuery        string
	filteredNamespaces []string
	// Terminal size fields
	termWidth        int
	termHeight       int
	terminalTooSmall bool
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

		return namespaceFetchedMsg{namespaces: namespaces}
	}
}

// Update handles incoming messages and returns an updated model and optional command
// nolint:gocyclo // Bubble Tea Update pattern inherently has high complexity due to message routing
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
		// Handle search mode ESC first (before quit keys)
		if m.searchMode && msg.Type == tea.KeyEsc {
			m.deactivateSearch()
			return m, nil
		}

		// Handle quit keys (but not in search mode where ESC is handled above)
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
						// TODO: Actual namespace selection logic (Story 2.4+)
						// For now, just deactivate search
						m.deactivateSearch()
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

			// Navigation (works in both normal and search mode)
			// Get the list we're navigating (filtered or full)
			navList := m.namespaces
			if m.searchMode && m.filteredNamespaces != nil {
				navList = m.filteredNamespaces
			}

			if KeyMatches(msg, m.keys.Up) {
				// Navigate up with wrap-around
				if len(navList) > 0 {
					m.selectedNamespaceIndex--
					if m.selectedNamespaceIndex < 0 {
						m.selectedNamespaceIndex = len(navList) - 1
						// Wrap to end: adjust viewport to show last item
						availableHeight := m.height - 4
						if availableHeight < 1 {
							availableHeight = 10
						}
						if len(navList) > availableHeight {
							m.namespaceViewportStart = len(navList) - availableHeight
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
				if len(navList) > 0 {
					m.selectedNamespaceIndex++
					if m.selectedNamespaceIndex >= len(navList) {
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
		m.termWidth = msg.Width
		m.termHeight = msg.Height

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
	} else if m.searchMode && len(m.filteredNamespaces) == 0 {
		// Show "no matches" message when search returns empty
		s += styles.DimStyle.Render("No matches found") + "\n"
	} else {
		// Create favorite set for lookup
		favSet := make(map[string]bool)
		if m.currentContext != nil {
			for _, fav := range m.currentContext.FavoriteNamespaces {
				favSet[fav] = true
			}
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

			// Add star for favorites
			if favSet[ns] {
				prefix += "★ "
			}

			// Render namespace name with highlighting if in search mode
			var renderedName string
			if m.searchMode && m.searchQuery != "" {
				renderedName = m.renderNamespaceWithHighlight(ns, prefix)
			} else {
				renderedName = prefix + ns
			}

			// Apply selection styling
			if i == m.selectedNamespaceIndex {
				s += styles.SelectedStyle.Render(renderedName) + "\n"
			} else {
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

// renderSplitLayout renders the split-pane layout with header, namespace, pods, and actions panels
func (m AppModel) renderSplitLayout() string {
	// Render header
	header := m.renderHeader()

	// Calculate dimensions
	availableHeight := m.termHeight - HeaderHeight
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

	// Combine left and right panels horizontally
	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, namespacePanel, rightSide)

	// Add header at top
	return lipgloss.JoinVertical(lipgloss.Left, header, mainLayout)
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
	contentHeight := height - 4 // Subtract border (2) + padding (2)

	// Get the namespace list content with correct panel height
	content := m.renderNamespaceList(contentHeight)

	// Width calculation: Lip Gloss adds border (2) + padding (4) = 6 chars total
	contentWidth := width - 6

	// Don't use MaxHeight - it cuts off borders
	// Instead, renderNamespaceList already limits content to contentHeight
	// Just apply the border and let Lip Gloss add padding + border
	return styles.PanelBorderStyle.
		Width(contentWidth).
		Render(content)
}

// renderPodPanel renders the placeholder pod panel
func (m AppModel) renderPodPanel(width, height int) string {
	title := styles.PanelTitleStyle.Render("Pods")
	placeholder := styles.PlaceholderStyle.Render("Select a namespace to view pods")
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", placeholder)

	// Apply border style with calculated dimensions
	contentWidth := width - 4   // 2 for border + 2*2 for padding
	contentHeight := height - 2 // 2 for border

	return styles.PanelBorderStyle.
		Width(contentWidth).
		Height(contentHeight).
		Render(content)
}

// renderActionsPanel renders the placeholder actions panel
func (m AppModel) renderActionsPanel(width, height int) string {
	title := styles.PanelTitleStyle.Render("Actions")
	placeholder := styles.PlaceholderStyle.Render("Available actions will appear here")
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", placeholder)

	// Apply border style with calculated dimensions
	contentWidth := width - 4   // 2 for border + 2*2 for padding
	contentHeight := height - 2 // 2 for border

	return styles.PanelBorderStyle.
		Width(contentWidth).
		Height(contentHeight).
		Render(content)
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
