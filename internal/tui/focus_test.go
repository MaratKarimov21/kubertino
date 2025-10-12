package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
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
				focusedPanel:     tt.initialFocus,
				pods:             tt.pods,
				actions:          tt.actions,
				selectedPodIndex: -1,
				viewMode:         viewModeNamespaceView,
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
				focusedPanel:     tt.initialFocus,
				pods:             tt.pods,
				actions:          tt.actions,
				selectedPodIndex: -1,
				viewMode:         viewModeNamespaceView,
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
		focusedPanel:     PanelNamespaces,
		pods:             []k8s.Pod{{Name: "pod-1", Status: "Running"}, {Name: "pod-2", Status: "Running"}},
		selectedPodIndex: -1,
		viewMode:         viewModeNamespaceView,
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
		focusedPanel:     PanelNamespaces,
		pods:             []k8s.Pod{{Name: "pod-1", Status: "Running"}, {Name: "pod-2", Status: "Running"}},
		selectedPodIndex: 1, // Already selected
		viewMode:         viewModeNamespaceView,
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
