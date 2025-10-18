package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/maratkarimov/kubertino/internal/tui/components"
	"github.com/stretchr/testify/assert"
)

// TestFocusCycling_Tab tests Tab key focus cycling (Story 6.2 - Two panels only, actions not focusable)
func TestFocusCycling_Tab(t *testing.T) {
	tests := []struct {
		name             string
		initialFocus     PanelType
		pods             []k8s.Pod
		actions          []config.Action
		expectedFocus    PanelType
		expectedPodIndex int // Expected selectedPodIndex after focus switch
	}{
		{
			name:             "Tab from namespaces to pods with pods available",
			initialFocus:     PanelNamespaces,
			pods:             []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:          []config.Action{{Name: "Test", Shortcut: "t", Command: "echo test"}},
			expectedFocus:    PanelPods,
			expectedPodIndex: 0, // Should auto-select first pod
		},
		{
			name:             "Tab from namespaces to pods with no pods",
			initialFocus:     PanelNamespaces,
			pods:             []k8s.Pod{},
			actions:          []config.Action{{Name: "Test", Shortcut: "t", Command: "echo test"}},
			expectedFocus:    PanelPods,
			expectedPodIndex: -1, // No pods to select
		},
		{
			name:             "Tab from pods back to namespaces (Story 6.2 - skip actions)",
			initialFocus:     PanelPods,
			pods:             []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:          []config.Action{{Name: "Test", Shortcut: "t", Command: "echo test"}},
			expectedFocus:    PanelNamespaces,
			expectedPodIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:      tt.initialFocus,
				pods:              tt.pods,
				actions:           tt.actions,
				selectedPodIndex:  -1,
				viewMode:          viewModeNamespaceView,
				errorModal:        components.NewErrorModal(), // Story 6.3: Initialize components
				namespacesSpinner: components.NewSpinner(),    // Story 6.3
				podsSpinner:       components.NewSpinner(),    // Story 6.3
				actionSpinner:     components.NewSpinner(),    // Story 6.3
			}

			msg := tea.KeyMsg{Type: tea.KeyTab}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedFocus, m.focusedPanel, "Focus panel mismatch")
			assert.Equal(t, tt.expectedPodIndex, m.selectedPodIndex, "Selected pod index mismatch")
		})
	}
}

// TestFocusCycling_ShiftTab tests Shift+Tab key backward focus cycling (Story 6.2)
func TestFocusCycling_ShiftTab(t *testing.T) {
	tests := []struct {
		name             string
		initialFocus     PanelType
		pods             []k8s.Pod
		actions          []config.Action
		expectedFocus    PanelType
		expectedPodIndex int
	}{
		{
			name:             "Shift+Tab from pods to namespaces",
			initialFocus:     PanelPods,
			pods:             []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:          []config.Action{{Name: "Test", Shortcut: "t", Command: "echo test"}},
			expectedFocus:    PanelNamespaces,
			expectedPodIndex: -1,
		},
		{
			name:             "Shift+Tab from namespaces to pods (Story 6.2 - skip actions)",
			initialFocus:     PanelNamespaces,
			pods:             []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:          []config.Action{{Name: "Test", Shortcut: "t", Command: "echo test"}},
			expectedFocus:    PanelPods,
			expectedPodIndex: 0, // Should auto-select first pod
		},
		{
			name:             "Shift+Tab from namespaces to pods with no pods",
			initialFocus:     PanelNamespaces,
			pods:             []k8s.Pod{},
			actions:          []config.Action{},
			expectedFocus:    PanelPods,
			expectedPodIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:      tt.initialFocus,
				pods:              tt.pods,
				actions:           tt.actions,
				selectedPodIndex:  -1,
				viewMode:          viewModeNamespaceView,
				errorModal:        components.NewErrorModal(), // Story 6.3
				namespacesSpinner: components.NewSpinner(),    // Story 6.3
				podsSpinner:       components.NewSpinner(),    // Story 6.3
				actionSpinner:     components.NewSpinner(),    // Story 6.3
			}

			// Simulate Shift+Tab key
			msg := tea.KeyMsg{Type: tea.KeyShiftTab}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedFocus, m.focusedPanel, "Focus panel mismatch")
			assert.Equal(t, tt.expectedPodIndex, m.selectedPodIndex, "Selected pod index mismatch")
		})
	}
}

// TestFocusInitialization tests that focus state is initialized correctly (Story 3.3)
func TestFocusInitialization(t *testing.T) {
	adapter := &mockKubeAdapter{}
	cfg := &config.Config{
		Version:  "1.0",
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, adapter)

	assert.Equal(t, PanelNamespaces, model.focusedPanel, "Should start with namespace panel focused")
	assert.Equal(t, -1, model.selectedPodIndex, "Should have no pod selected initially (Story 6.2: cursor = selection)")
	assert.Equal(t, 0, model.podScrollOffset, "Should have no scroll offset initially")
}

// TestAutoSelectFirstPod tests that first pod is auto-selected when focusing pod panel (Story 3.3)
func TestAutoSelectFirstPod(t *testing.T) {
	model := AppModel{
		focusedPanel:      PanelNamespaces,
		pods:              []k8s.Pod{{Name: "pod-1", Status: "Running"}, {Name: "pod-2", Status: "Running"}},
		selectedPodIndex:  -1,
		viewMode:          viewModeNamespaceView,
		errorModal:        components.NewErrorModal(), // Story 6.3
		namespacesSpinner: components.NewSpinner(),    // Story 6.3
		podsSpinner:       components.NewSpinner(),    // Story 6.3
		actionSpinner:     components.NewSpinner(),    // Story 6.3
	}

	// Tab to pod panel
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, PanelPods, m.focusedPanel)
	assert.Equal(t, 0, m.selectedPodIndex, "First pod should be auto-selected")
}

// TestNoAutoSelectWhenAlreadySelected tests that auto-select doesn't override existing selection (Story 3.3)
func TestNoAutoSelectWhenAlreadySelected(t *testing.T) {
	model := AppModel{
		focusedPanel:      PanelNamespaces,
		pods:              []k8s.Pod{{Name: "pod-1", Status: "Running"}, {Name: "pod-2", Status: "Running"}},
		selectedPodIndex:  1, // Already selected
		viewMode:          viewModeNamespaceView,
		errorModal:        components.NewErrorModal(), // Story 6.3
		namespacesSpinner: components.NewSpinner(),    // Story 6.3
		podsSpinner:       components.NewSpinner(),    // Story 6.3
		actionSpinner:     components.NewSpinner(),    // Story 6.3
	}

	// Tab to pod panel
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, PanelPods, m.focusedPanel)
	assert.Equal(t, 1, m.selectedPodIndex, "Selection should remain at pod 1")
}

// TestActionsPanelNeverFocused tests that actions panel is never focusable (Story 6.2)
func TestActionsPanelNeverFocused(t *testing.T) {
	adapter := &mockKubeAdapter{}
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{{
			Name: "test",
			Actions: []config.Action{
				{Name: "Console", Shortcut: "c", Command: "echo test"},
			},
		}},
	}

	model := NewAppModel(cfg, adapter)

	assert.Equal(t, 1, len(model.actions), "Actions should be loaded from context")
	assert.Equal(t, PanelNamespaces, model.focusedPanel, "Should start on namespace panel")

	// Tab should cycle between namespaces and pods only, never landing on actions
	// Simulate multiple Tab presses to verify cycling
	model.viewMode = viewModeNamespaceView
	model.pods = []k8s.Pod{{Name: "test", Status: "Running"}}

	// Tab 1: Namespaces -> Pods
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)
	assert.Equal(t, PanelPods, m.focusedPanel, "First tab should go to pods")

	// Tab 2: Pods -> Namespaces (skipping actions)
	msg = tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ = m.Update(msg)
	m = updatedModel.(AppModel)
	assert.Equal(t, PanelNamespaces, m.focusedPanel, "Second tab should go back to namespaces, skipping actions")
}

// Story 6.2 Update: Pod confirmation test removed - cursor position now equals selection (no Enter needed)

// Story 7.5: Test namespace cursor preservation on context switch
func TestContextSwitch_PreservesNamespaceCursor(t *testing.T) {
	tests := []struct {
		name                  string
		initialNamespaceIndex int
		newContextNamespaces  []string
		expectedIndex         int
	}{
		{
			name:                  "Cursor preserved when new context has same namespace count",
			initialNamespaceIndex: 5,
			newContextNamespaces:  []string{"ns-1", "ns-2", "ns-3", "ns-4", "ns-5", "ns-6", "ns-7"},
			expectedIndex:         5,
		},
		{
			name:                  "Cursor preserved when new context has more namespaces",
			initialNamespaceIndex: 3,
			newContextNamespaces:  []string{"ns-1", "ns-2", "ns-3", "ns-4", "ns-5", "ns-6", "ns-7", "ns-8"},
			expectedIndex:         3,
		},
		{
			name:                  "Cursor clamped when new context has fewer namespaces",
			initialNamespaceIndex: 5,
			newContextNamespaces:  []string{"ns-1", "ns-2", "ns-3"},
			expectedIndex:         2, // Clamped to last item (index 2)
		},
		{
			name:                  "Cursor at 0 preserved when new context has namespaces",
			initialNamespaceIndex: 0,
			newContextNamespaces:  []string{"new-ns-1", "new-ns-2"},
			expectedIndex:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &mockKubeAdapter{
				namespaces: tt.newContextNamespaces,
			}

			// Start with context selection mode
			cfg := &config.Config{
				Version: "1.0",
				Contexts: []config.Context{
					{Name: "context-1"},
					{Name: "context-2"},
				},
			}

			model := NewAppModel(cfg, adapter)
			model.viewMode = viewModeContextSelection
			model.selectedContextIndex = 0
			model.selectedNamespaceIndex = tt.initialNamespaceIndex

			// Select context (triggers namespace fetch)
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			// Verify namespace index is preserved during context switch
			assert.Equal(t, tt.initialNamespaceIndex, m.selectedNamespaceIndex, "Cursor position should be preserved")

			// Now simulate namespace fetch completion
			namespaceMsg := namespaceFetchedMsg{namespaces: tt.newContextNamespaces}
			updatedModel, _ = m.Update(namespaceMsg)
			m = updatedModel.(AppModel)

			// Verify cursor is clamped to valid range
			assert.Equal(t, tt.expectedIndex, m.selectedNamespaceIndex, "Cursor should be adjusted to valid range")
		})
	}
}

// Story 7.5: Test auto-select first pod on namespace selection
func TestNamespaceSelection_AutoSelectsFirstPod(t *testing.T) {
	tests := []struct {
		name             string
		pods             []k8s.Pod
		expectedPodIndex int
	}{
		{
			name: "Auto-select first pod when pods available",
			pods: []k8s.Pod{
				{Name: "pod-1", Status: "Running"},
				{Name: "pod-2", Status: "Running"},
			},
			expectedPodIndex: 0,
		},
		{
			name:             "No pod selected when no pods available",
			pods:             []k8s.Pod{},
			expectedPodIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &mockKubeAdapter{
				namespaces: []string{"default", "kube-system"},
				pods:       tt.pods,
			}

			model := AppModel{
				errorModal:             components.NewErrorModal(),
				namespacesSpinner:      components.NewSpinner(),
				podsSpinner:            components.NewSpinner(),
				actionSpinner:          components.NewSpinner(),
				viewMode:               viewModeNamespaceView,
				kubeAdapter:            adapter,
				namespaces:             []string{"default", "kube-system"},
				selectedNamespaceIndex: 0,
				focusedPanel:           PanelNamespaces,
				selectedPodIndex:       -1,
				keys:                   DefaultKeyMap(),
			}

			// Select namespace (triggers pod fetch and focus switch)
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			// Verify focus switched to pods panel
			assert.Equal(t, PanelPods, m.focusedPanel, "Focus should switch to pods panel")
			assert.Equal(t, -1, m.selectedPodIndex, "Pod index should be reset before fetch")

			// Simulate pod fetch completion
			podMsg := podsFetchedMsg{pods: tt.pods}
			updatedModel, _ = m.Update(podMsg)
			m = updatedModel.(AppModel)

			// Verify first pod is auto-selected
			assert.Equal(t, tt.expectedPodIndex, m.selectedPodIndex, "First pod should be auto-selected")
		})
	}
}

// Story 7.5: Test auto-select consistency between Tab and namespace selection
func TestAutoSelectFirstPod_ConsistentBehavior(t *testing.T) {
	pods := []k8s.Pod{
		{Name: "pod-1", Status: "Running"},
		{Name: "pod-2", Status: "Running"},
	}

	adapter := &mockKubeAdapter{
		namespaces: []string{"default"},
		pods:       pods,
	}

	t.Run("Tab key auto-selects first pod", func(t *testing.T) {
		model := AppModel{
			errorModal:        components.NewErrorModal(),
			namespacesSpinner: components.NewSpinner(),
			podsSpinner:       components.NewSpinner(),
			actionSpinner:     components.NewSpinner(),
			viewMode:          viewModeNamespaceView,
			focusedPanel:      PanelNamespaces,
			pods:              pods,
			selectedPodIndex:  -1,
		}

		// Press Tab to switch to pods
		msg := tea.KeyMsg{Type: tea.KeyTab}
		updatedModel, _ := model.Update(msg)
		m := updatedModel.(AppModel)

		assert.Equal(t, PanelPods, m.focusedPanel)
		assert.Equal(t, 0, m.selectedPodIndex, "Tab should auto-select first pod")
	})

	t.Run("Namespace selection auto-selects first pod", func(t *testing.T) {
		model := AppModel{
			errorModal:             components.NewErrorModal(),
			namespacesSpinner:      components.NewSpinner(),
			podsSpinner:            components.NewSpinner(),
			actionSpinner:          components.NewSpinner(),
			viewMode:               viewModeNamespaceView,
			kubeAdapter:            adapter,
			namespaces:             []string{"default"},
			selectedNamespaceIndex: 0,
			focusedPanel:           PanelNamespaces,
			selectedPodIndex:       -1,
			keys:                   DefaultKeyMap(),
		}

		// Select namespace
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		updatedModel, _ := model.Update(msg)
		m := updatedModel.(AppModel)

		// Simulate pod fetch
		podMsg := podsFetchedMsg{pods: pods}
		updatedModel, _ = m.Update(podMsg)
		m = updatedModel.(AppModel)

		assert.Equal(t, PanelPods, m.focusedPanel)
		assert.Equal(t, 0, m.selectedPodIndex, "Namespace selection should auto-select first pod")
	})
}
