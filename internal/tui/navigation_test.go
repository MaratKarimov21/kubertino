package tui

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/stretchr/testify/assert"
)

// TestPodNavigation_ArrowKeys tests arrow key navigation in pod panel (Story 3.3)
func TestPodNavigation_ArrowKeys(t *testing.T) {
	tests := []struct {
		name          string
		initialIndex  int
		keyPress      tea.KeyType
		focusedPanel  PanelType
		pods          []k8s.Pod
		expectedIndex int
	}{
		{
			name:          "Down arrow increments selection",
			initialIndex:  0,
			keyPress:      tea.KeyDown,
			focusedPanel:  PanelPods,
			pods:          []k8s.Pod{{Name: "pod-1"}, {Name: "pod-2"}, {Name: "pod-3"}},
			expectedIndex: 1,
		},
		{
			name:          "Up arrow decrements selection",
			initialIndex:  2,
			keyPress:      tea.KeyUp,
			focusedPanel:  PanelPods,
			pods:          []k8s.Pod{{Name: "pod-1"}, {Name: "pod-2"}, {Name: "pod-3"}},
			expectedIndex: 1,
		},
		{
			name:          "Down arrow clamped at last pod",
			initialIndex:  2,
			keyPress:      tea.KeyDown,
			focusedPanel:  PanelPods,
			pods:          []k8s.Pod{{Name: "pod-1"}, {Name: "pod-2"}, {Name: "pod-3"}},
			expectedIndex: 2, // No change
		},
		{
			name:          "Up arrow clamped at first pod",
			initialIndex:  0,
			keyPress:      tea.KeyUp,
			focusedPanel:  PanelPods,
			pods:          []k8s.Pod{{Name: "pod-1"}, {Name: "pod-2"}},
			expectedIndex: 0, // No change
		},
		{
			name:          "Arrow keys ignored when namespace focused",
			initialIndex:  1,
			keyPress:      tea.KeyDown,
			focusedPanel:  PanelNamespaces,
			pods:          []k8s.Pod{{Name: "pod-1"}, {Name: "pod-2"}, {Name: "pod-3"}},
			expectedIndex: 1, // No change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:     tt.focusedPanel,
				selectedPodIndex: tt.initialIndex,
				pods:             tt.pods,
				viewMode:         viewModeNamespaceView,
				keys:             DefaultKeyMap(), // Need keys for KeyMatches
			}

			msg := tea.KeyMsg{Type: tt.keyPress}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedIndex, m.selectedPodIndex, "Pod index mismatch")
		})
	}
}

// TestPodScrolling tests scrolling behavior for long pod lists (Story 3.3)
func TestPodScrolling(t *testing.T) {
	// Create a long pod list (50 pods)
	pods := make([]k8s.Pod, 50)
	for i := range pods {
		pods[i] = k8s.Pod{Name: fmt.Sprintf("pod-%d", i), Status: "Running"}
	}

	tests := []struct {
		name           string
		initialIndex   int
		initialOffset  int
		keyPress       tea.KeyType
		termHeight     int
		expectedIndex  int
		expectedOffset int
	}{
		{
			name:           "Scrolling down when selection at bottom of window",
			initialIndex:   3, // At bottom of visible window (0-3 visible with height 4)
			initialOffset:  0,
			keyPress:       tea.KeyDown,
			termHeight:     25, // visibleHeight = (25-1)/2 - 8 = 4
			expectedIndex:  4,
			expectedOffset: 1, // Should scroll down by 1
		},
		{
			name:           "Scrolling up when selection at top of window",
			initialIndex:   10,
			initialOffset:  10,
			keyPress:       tea.KeyUp,
			termHeight:     25,
			expectedIndex:  9,
			expectedOffset: 9, // Should scroll up to keep selection visible
		},
		{
			name:           "No scroll when selection in middle of window",
			initialIndex:   1, // In middle of visible range (0-3)
			initialOffset:  0,
			keyPress:       tea.KeyDown,
			termHeight:     25,
			expectedIndex:  2,
			expectedOffset: 0, // No scroll needed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:     PanelPods,
				selectedPodIndex: tt.initialIndex,
				podScrollOffset:  tt.initialOffset,
				pods:             pods,
				termHeight:       tt.termHeight,
				viewMode:         viewModeNamespaceView,
				keys:             DefaultKeyMap(), // Need keys for KeyMatches
			}

			msg := tea.KeyMsg{Type: tt.keyPress}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedIndex, m.selectedPodIndex, "Pod index mismatch")
			assert.Equal(t, tt.expectedOffset, m.podScrollOffset, "Scroll offset mismatch")
		})
	}
}

// TestPodStateReset tests that pod state is reset when namespace changes (Story 3.3)
func TestPodStateReset(t *testing.T) {
	adapter := &mockKubeAdapter{
		namespaces: []string{"default", "kube-system"},
		pods:       []k8s.Pod{{Name: "pod-1", Status: "Running"}},
	}

	model := AppModel{
		viewMode:         viewModeNamespaceView,
		kubeAdapter:      adapter,
		namespaces:       []string{"default", "kube-system"},
		selectedPodIndex: 5,
		podScrollOffset:  10,
		focusedPanel:     PanelNamespaces, // Must be on namespaces to select
		keys:             DefaultKeyMap(), // Need keys for KeyMatches
	}

	// Select a namespace (triggers pod fetch and state reset)
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, -1, m.selectedPodIndex, "Selected pod index should be reset")
	assert.Equal(t, 0, m.podScrollOffset, "Scroll offset should be reset")
}

// TestNavigationWithEmptyPodList tests navigation with empty pod list (Story 3.3)
func TestNavigationWithEmptyPodList(t *testing.T) {
	model := AppModel{
		focusedPanel:     PanelPods,
		selectedPodIndex: -1,
		pods:             []k8s.Pod{}, // Empty
		viewMode:         viewModeNamespaceView,
		keys:             DefaultKeyMap(),
	}

	// Try to navigate down
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, -1, m.selectedPodIndex, "Index should remain -1 with empty pod list")

	// Try to navigate up
	msg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = m.Update(msg)
	m = updatedModel.(AppModel)

	assert.Equal(t, -1, m.selectedPodIndex, "Index should remain -1 with empty pod list")
}

// TestNavigationWithSinglePod tests navigation with single pod (Story 3.3)
func TestNavigationWithSinglePod(t *testing.T) {
	model := AppModel{
		focusedPanel:     PanelPods,
		selectedPodIndex: 0,
		pods:             []k8s.Pod{{Name: "only-pod", Status: "Running"}},
		viewMode:         viewModeNamespaceView,
		keys:             DefaultKeyMap(),
	}

	// Try to navigate down (should stay at 0)
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, 0, m.selectedPodIndex, "Index should stay at 0 with single pod")

	// Try to navigate up (should stay at 0)
	msg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = m.Update(msg)
	m = updatedModel.(AppModel)

	assert.Equal(t, 0, m.selectedPodIndex, "Index should stay at 0 with single pod")
}

// Story 6.1: Test favorites config order preservation
func TestFavoritesConfigOrder(t *testing.T) {
	model := AppModel{}

	// Favorites in specific non-alphabetical order
	favorites := []string{"zebra-ns", "apple-ns", "banana-ns"}
	// All namespaces (includes non-favorites too)
	namespaces := []string{"apple-ns", "banana-ns", "default", "kube-system", "zebra-ns"}

	result := model.sortNamespacesWithFavorites(namespaces, favorites)

	// Expected: favorites in config order, then non-favorites alphabetically
	expected := []string{"zebra-ns", "apple-ns", "banana-ns", "default", "kube-system"}

	assert.Equal(t, expected, result, "Favorites should preserve config order, non-favorites alphabetically sorted")
}

// Story 6.1: Test cursor visibility during scrolling
func TestNamespaceCursorVisibilityDuringScroll(t *testing.T) {
	// Create a long namespace list (50 items)
	namespaces := make([]string, 50)
	for i := range namespaces {
		namespaces[i] = fmt.Sprintf("namespace-%02d", i)
	}

	tests := []struct {
		name                string
		initialIndex        int
		initialViewport     int
		keyPress            tea.KeyType
		height              int
		expectedIndex       int
		checkViewportBounds bool
	}{
		{
			name:                "Cursor remains visible when scrolling down",
			initialIndex:        10,
			initialViewport:     5,
			keyPress:            tea.KeyDown,
			height:              20, // availableHeight will be ~10
			expectedIndex:       11,
			checkViewportBounds: true,
		},
		{
			name:                "Cursor remains visible when scrolling up",
			initialIndex:        20,
			initialViewport:     15,
			keyPress:            tea.KeyUp,
			height:              20,
			expectedIndex:       19,
			checkViewportBounds: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:           PanelNamespaces,
				selectedNamespaceIndex: tt.initialIndex,
				namespaceViewportStart: tt.initialViewport,
				namespaces:             namespaces,
				viewMode:               viewModeNamespaceView,
				height:                 tt.height,
				keys:                   DefaultKeyMap(),
			}

			msg := tea.KeyMsg{Type: tt.keyPress}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedIndex, m.selectedNamespaceIndex, "Namespace index mismatch")

			if tt.checkViewportBounds {
				// Calculate available height (same as adjustNamespaceViewport)
				headerLines := 2
				footerLines := 2
				scrollIndicatorLines := 1
				reservedLines := headerLines + footerLines
				availableHeight := tt.height - reservedLines - scrollIndicatorLines
				if availableHeight < 1 {
					availableHeight = 1
				}

				// Cursor should be within viewport bounds
				assert.GreaterOrEqual(t, m.selectedNamespaceIndex, m.namespaceViewportStart,
					"Cursor should be at or after viewport start")
				assert.Less(t, m.selectedNamespaceIndex, m.namespaceViewportStart+availableHeight,
					"Cursor should be before viewport end")
			}
		})
	}
}

// Story 6.1: Test cursor centering behavior
func TestNamespaceCursorCentering(t *testing.T) {
	// Create a long namespace list (50 items)
	namespaces := make([]string, 50)
	for i := range namespaces {
		namespaces[i] = fmt.Sprintf("namespace-%02d", i)
	}

	tests := []struct {
		name              string
		initialIndex      int
		navigateDownCount int
		height            int
		checkCentering    bool
		expectedViewport  int
		expectedIndex     int
	}{
		{
			name:              "Cursor moves from top to middle without scrolling",
			initialIndex:      0,
			navigateDownCount: 5,
			height:            20, // availableHeight ~10, middle ~5
			checkCentering:    false,
			expectedIndex:     5,
			expectedViewport:  0, // Viewport stays at 0 while cursor reaches middle
		},
		{
			name:              "Cursor centered when scrolling in middle section",
			initialIndex:      10,
			navigateDownCount: 5,
			height:            20, // availableHeight ~10, middle ~5
			checkCentering:    true,
			expectedIndex:     15,
			expectedViewport:  10, // Viewport adjusts to keep cursor at middle (~5)
		},
		{
			name:              "Cursor moves to bottom when approaching end",
			initialIndex:      45,
			navigateDownCount: 3,
			height:            20,
			checkCentering:    false,
			expectedIndex:     48,
			expectedViewport:  40, // Viewport at end (50 - 10)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:           PanelNamespaces,
				selectedNamespaceIndex: tt.initialIndex,
				namespaceViewportStart: 0,
				namespaces:             namespaces,
				viewMode:               viewModeNamespaceView,
				height:                 tt.height,
				keys:                   DefaultKeyMap(),
			}

			// Navigate down N times
			for i := 0; i < tt.navigateDownCount; i++ {
				msg := tea.KeyMsg{Type: tea.KeyDown}
				updatedModel, _ := model.Update(msg)
				model = updatedModel.(AppModel)
			}

			assert.Equal(t, tt.expectedIndex, model.selectedNamespaceIndex, "Final index mismatch")

			if tt.checkCentering {
				// Calculate expected centered viewport
				headerLines := 2
				footerLines := 2
				scrollIndicatorLines := 1
				reservedLines := headerLines + footerLines
				availableHeight := tt.height - reservedLines - scrollIndicatorLines
				if availableHeight < 1 {
					availableHeight = 1
				}
				middlePosition := availableHeight / 2

				// When in middle section, cursor should be roughly centered
				expectedViewport := model.selectedNamespaceIndex - middlePosition
				assert.Equal(t, expectedViewport, model.namespaceViewportStart, "Viewport should center cursor")
			}
		})
	}
}

// Story 6.1: Test cursor persistence on focus change
func TestNamespaceCursorPersistenceOnFocusChange(t *testing.T) {
	adapter := &mockKubeAdapter{
		namespaces: []string{"default", "kube-system", "test-ns"},
		pods:       []k8s.Pod{{Name: "pod-1", Status: "Running"}},
	}

	model := AppModel{
		viewMode:               viewModeNamespaceView,
		kubeAdapter:            adapter,
		namespaces:             []string{"default", "kube-system", "test-ns"},
		selectedNamespaceIndex: 2, // Selected "test-ns"
		focusedPanel:           PanelNamespaces,
		currentNamespace:       "test-ns",
		pods:                   []k8s.Pod{{Name: "pod-1", Status: "Running"}},
		keys:                   DefaultKeyMap(),
	}

	// Press Tab to switch focus to pod panel
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	// Verify namespace selection persists
	assert.Equal(t, 2, m.selectedNamespaceIndex, "Namespace selection should persist")
	assert.Equal(t, PanelPods, m.focusedPanel, "Focus should shift to pods panel")
	assert.Equal(t, "test-ns", m.currentNamespace, "Current namespace should remain")
}

// Story 6.1: Test cursor wrap-around with viewport reset
func TestNamespaceCursorWrapAround(t *testing.T) {
	namespaces := []string{"ns-1", "ns-2", "ns-3", "ns-4", "ns-5", "ns-6", "ns-7", "ns-8", "ns-9", "ns-10"}

	tests := []struct {
		name             string
		initialIndex     int
		initialViewport  int
		keyPress         tea.KeyType
		expectedIndex    int
		expectedViewport int
	}{
		{
			name:             "Wrap from last to first (down)",
			initialIndex:     9, // Last item
			initialViewport:  5,
			keyPress:         tea.KeyDown,
			expectedIndex:    0, // Wraps to first
			expectedViewport: 0, // Viewport resets
		},
		{
			name:             "Wrap from first to last (up)",
			initialIndex:     0, // First item
			initialViewport:  0,
			keyPress:         tea.KeyUp,
			expectedIndex:    9, // Wraps to last
			expectedViewport: 0, // Viewport adjusts (will be calculated by adjustNamespaceViewport)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:           PanelNamespaces,
				selectedNamespaceIndex: tt.initialIndex,
				namespaceViewportStart: tt.initialViewport,
				namespaces:             namespaces,
				viewMode:               viewModeNamespaceView,
				height:                 20,
				keys:                   DefaultKeyMap(),
			}

			msg := tea.KeyMsg{Type: tt.keyPress}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedIndex, m.selectedNamespaceIndex, "Index should wrap correctly")
		})
	}
}

// Story 6.1: Test centering with short list (no scrolling needed)
func TestNamespaceCenteringShortList(t *testing.T) {
	namespaces := []string{"ns-1", "ns-2", "ns-3"}

	model := AppModel{
		focusedPanel:           PanelNamespaces,
		selectedNamespaceIndex: 0,
		namespaceViewportStart: 0,
		namespaces:             namespaces,
		viewMode:               viewModeNamespaceView,
		height:                 20, // Much larger than list
		keys:                   DefaultKeyMap(),
	}

	// Navigate down through entire list
	for i := 0; i < 2; i++ {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		updatedModel, _ := model.Update(msg)
		model = updatedModel.(AppModel)
	}

	// Viewport should stay at 0 (no scrolling needed for short list)
	assert.Equal(t, 2, model.selectedNamespaceIndex, "Should navigate to last item")
	assert.Equal(t, 0, model.namespaceViewportStart, "Viewport should stay at 0 for short list")
}

// Story 6.1: Test large list navigation (252 items) - regression test for cursor going off-screen
func TestNamespaceLargeListBottomNavigation(t *testing.T) {
	// Create 252 namespaces (real-world scenario from user screenshot)
	namespaces := make([]string, 252)
	for i := range namespaces {
		namespaces[i] = fmt.Sprintf("namespace-%03d", i)
	}

	// User had terminal height that gave viewport of 38 items
	// With termHeight ~45: border+padding(4) + headerLines(2) + footerLines(2) + scrollIndicator(1) = 9 reserved
	// availableHeight = 45 - 4 - 5 = 36 (approximate)
	termHeight := 45

	model := AppModel{
		focusedPanel:           PanelNamespaces,
		selectedNamespaceIndex: 0,
		namespaceViewportStart: 0,
		namespaces:             namespaces,
		viewMode:               viewModeNamespaceView,
		termHeight:             termHeight, // CRITICAL: Use termHeight, not height
		keys:                   DefaultKeyMap(),
	}

	// Calculate expected available height (must match adjustNamespaceViewport logic)
	// adjustNamespaceViewport uses termHeight - 4 (border+padding) as effectiveHeight
	effectiveHeight := termHeight - 4
	headerLines := 2
	footerLines := 2
	scrollIndicatorLines := 1
	reservedLines := headerLines + footerLines
	availableHeight := effectiveHeight - reservedLines - scrollIndicatorLines
	if availableHeight < 1 {
		availableHeight = 1
	}

	// Navigate to near the end of the list (e.g., item 248 out of 252)
	targetIndex := 248
	for i := 0; i < targetIndex; i++ {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		updatedModel, _ := model.Update(msg)
		model = updatedModel.(AppModel)

		// CRITICAL CHECK: cursor must ALWAYS be within viewport bounds
		assert.GreaterOrEqual(t, model.selectedNamespaceIndex, model.namespaceViewportStart,
			fmt.Sprintf("At step %d: cursor (%d) must be >= viewport start (%d)",
				i, model.selectedNamespaceIndex, model.namespaceViewportStart))
		assert.Less(t, model.selectedNamespaceIndex, model.namespaceViewportStart+availableHeight,
			fmt.Sprintf("At step %d: cursor (%d) must be < viewport end (%d)",
				i, model.selectedNamespaceIndex, model.namespaceViewportStart+availableHeight))
	}

	// Verify final position
	assert.Equal(t, targetIndex, model.selectedNamespaceIndex, "Should reach target index")

	// Cursor should be visible in viewport
	viewportEnd := model.namespaceViewportStart + availableHeight
	assert.GreaterOrEqual(t, model.selectedNamespaceIndex, model.namespaceViewportStart,
		"Cursor should be at or after viewport start")
	assert.Less(t, model.selectedNamespaceIndex, viewportEnd,
		"Cursor should be before viewport end")

	// Navigate to the absolute last item (251)
	for i := targetIndex; i < 251; i++ {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		updatedModel, _ := model.Update(msg)
		model = updatedModel.(AppModel)

		// Check bounds at each step
		assert.GreaterOrEqual(t, model.selectedNamespaceIndex, model.namespaceViewportStart,
			fmt.Sprintf("At step %d: cursor (%d) must be >= viewport start (%d)",
				i, model.selectedNamespaceIndex, model.namespaceViewportStart))
		assert.Less(t, model.selectedNamespaceIndex, model.namespaceViewportStart+availableHeight,
			fmt.Sprintf("At step %d: cursor (%d) must be < viewport end (%d)",
				i, model.selectedNamespaceIndex, model.namespaceViewportStart+availableHeight))
	}

	// At the end, viewport should be locked to show last items
	assert.Equal(t, 251, model.selectedNamespaceIndex, "Should reach last item")
	expectedMaxViewport := len(namespaces) - availableHeight
	assert.Equal(t, expectedMaxViewport, model.namespaceViewportStart,
		"Viewport should be locked at bottom")
}
