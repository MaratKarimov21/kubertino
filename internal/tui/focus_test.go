package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/stretchr/testify/assert"
)

// TestFocusCycling_Tab tests Tab key focus cycling (Story 3.3)
func TestFocusCycling_Tab(t *testing.T) {
	tests := []struct {
		name          string
		initialFocus  PanelType
		pods          []k8s.Pod
		expectedFocus PanelType
		expectedIndex int // Expected selectedPodIndex after focus switch
	}{
		{
			name:          "Tab from namespaces to pods with pods available",
			initialFocus:  PanelNamespaces,
			pods:          []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			expectedFocus: PanelPods,
			expectedIndex: 0, // Should auto-select first pod
		},
		{
			name:          "Tab from namespaces to pods with no pods",
			initialFocus:  PanelNamespaces,
			pods:          []k8s.Pod{},
			expectedFocus: PanelPods,
			expectedIndex: -1, // No pods to select
		},
		{
			name:          "Tab from pods to namespaces",
			initialFocus:  PanelPods,
			pods:          []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			expectedFocus: PanelNamespaces,
			expectedIndex: -1, // Index unchanged (stays at initial -1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:     tt.initialFocus,
				pods:             tt.pods,
				selectedPodIndex: -1,
				viewMode:         viewModeNamespaceView,
			}

			msg := tea.KeyMsg{Type: tea.KeyTab}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedFocus, m.focusedPanel, "Focus panel mismatch")
			assert.Equal(t, tt.expectedIndex, m.selectedPodIndex, "Selected pod index mismatch")
		})
	}
}

// TestFocusCycling_ShiftTab tests Shift+Tab key backward focus cycling (Story 3.3)
func TestFocusCycling_ShiftTab(t *testing.T) {
	tests := []struct {
		name          string
		initialFocus  PanelType
		pods          []k8s.Pod
		expectedFocus PanelType
		expectedIndex int
	}{
		{
			name:          "Shift+Tab from pods to namespaces",
			initialFocus:  PanelPods,
			pods:          []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			expectedFocus: PanelNamespaces,
			expectedIndex: -1,
		},
		{
			name:          "Shift+Tab from namespaces to pods with pods available",
			initialFocus:  PanelNamespaces,
			pods:          []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			expectedFocus: PanelPods,
			expectedIndex: 0, // Should auto-select first pod
		},
		{
			name:          "Shift+Tab from namespaces to pods with no pods",
			initialFocus:  PanelNamespaces,
			pods:          []k8s.Pod{},
			expectedFocus: PanelPods,
			expectedIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:     tt.initialFocus,
				pods:             tt.pods,
				selectedPodIndex: -1,
				viewMode:         viewModeNamespaceView,
			}

			// Simulate Shift+Tab key
			msg := tea.KeyMsg{Type: tea.KeyShiftTab}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedFocus, m.focusedPanel, "Focus panel mismatch")
			assert.Equal(t, tt.expectedIndex, m.selectedPodIndex, "Selected pod index mismatch")
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
	assert.Equal(t, -1, model.selectedPodIndex, "Should have no pod selected initially")
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
