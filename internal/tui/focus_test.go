package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/stretchr/testify/assert"
)

// TestFocusCycling_Tab tests Tab key focus cycling (Story 3.3, 4.1 - Three panels)
func TestFocusCycling_Tab(t *testing.T) {
	tests := []struct {
		name                string
		initialFocus        PanelType
		pods                []k8s.Pod
		actions             []config.Action
		expectedFocus       PanelType
		expectedPodIndex    int // Expected selectedPodIndex after focus switch
		expectedActionIndex int // Expected selectedActionIndex after focus switch
	}{
		{
			name:                "Tab from namespaces to pods with pods available",
			initialFocus:        PanelNamespaces,
			pods:                []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:             []config.Action{{Name: "Test", Shortcut: "t", Type: "pod_exec"}},
			expectedFocus:       PanelPods,
			expectedPodIndex:    0,  // Should auto-select first pod
			expectedActionIndex: -1, // No change
		},
		{
			name:                "Tab from namespaces to pods with no pods",
			initialFocus:        PanelNamespaces,
			pods:                []k8s.Pod{},
			actions:             []config.Action{{Name: "Test", Shortcut: "t", Type: "pod_exec"}},
			expectedFocus:       PanelPods,
			expectedPodIndex:    -1, // No pods to select
			expectedActionIndex: -1,
		},
		{
			name:                "Tab from pods to actions with actions available (Story 4.1)",
			initialFocus:        PanelPods,
			pods:                []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:             []config.Action{{Name: "Test", Shortcut: "t", Type: "pod_exec"}},
			expectedFocus:       PanelActions,
			expectedPodIndex:    -1,
			expectedActionIndex: 0, // Should auto-select first action
		},
		{
			name:                "Tab from pods to actions with no actions (Story 4.1)",
			initialFocus:        PanelPods,
			pods:                []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:             []config.Action{},
			expectedFocus:       PanelActions,
			expectedPodIndex:    -1,
			expectedActionIndex: -1, // No actions to select
		},
		{
			name:                "Tab from actions to namespaces (Story 4.1)",
			initialFocus:        PanelActions,
			pods:                []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:             []config.Action{{Name: "Test", Shortcut: "t", Type: "pod_exec"}},
			expectedFocus:       PanelNamespaces,
			expectedPodIndex:    -1,
			expectedActionIndex: -1, // Index unchanged (stays at initial -1)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:        tt.initialFocus,
				pods:                tt.pods,
				actions:             tt.actions,
				selectedPodIndex:    -1,
				selectedActionIndex: -1,
				viewMode:            viewModeNamespaceView,
			}

			msg := tea.KeyMsg{Type: tea.KeyTab}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedFocus, m.focusedPanel, "Focus panel mismatch")
			assert.Equal(t, tt.expectedPodIndex, m.selectedPodIndex, "Selected pod index mismatch")
			assert.Equal(t, tt.expectedActionIndex, m.selectedActionIndex, "Selected action index mismatch")
		})
	}
}

// TestFocusCycling_ShiftTab tests Shift+Tab key backward focus cycling (Story 3.3, 4.1)
func TestFocusCycling_ShiftTab(t *testing.T) {
	tests := []struct {
		name                string
		initialFocus        PanelType
		pods                []k8s.Pod
		actions             []config.Action
		expectedFocus       PanelType
		expectedPodIndex    int
		expectedActionIndex int
	}{
		{
			name:                "Shift+Tab from pods to namespaces",
			initialFocus:        PanelPods,
			pods:                []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:             []config.Action{{Name: "Test", Shortcut: "t", Type: "pod_exec"}},
			expectedFocus:       PanelNamespaces,
			expectedPodIndex:    -1,
			expectedActionIndex: -1,
		},
		{
			name:                "Shift+Tab from namespaces to actions with actions available (Story 4.1)",
			initialFocus:        PanelNamespaces,
			pods:                []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:             []config.Action{{Name: "Test", Shortcut: "t", Type: "pod_exec"}},
			expectedFocus:       PanelActions,
			expectedPodIndex:    -1,
			expectedActionIndex: 0, // Should auto-select first action
		},
		{
			name:                "Shift+Tab from namespaces to actions with no actions (Story 4.1)",
			initialFocus:        PanelNamespaces,
			pods:                []k8s.Pod{},
			actions:             []config.Action{},
			expectedFocus:       PanelActions,
			expectedPodIndex:    -1,
			expectedActionIndex: -1,
		},
		{
			name:                "Shift+Tab from actions to pods with pods available (Story 4.1)",
			initialFocus:        PanelActions,
			pods:                []k8s.Pod{{Name: "test-pod", Status: "Running"}},
			actions:             []config.Action{{Name: "Test", Shortcut: "t", Type: "pod_exec"}},
			expectedFocus:       PanelPods,
			expectedPodIndex:    0, // Should auto-select first pod
			expectedActionIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:        tt.initialFocus,
				pods:                tt.pods,
				actions:             tt.actions,
				selectedPodIndex:    -1,
				selectedActionIndex: -1,
				viewMode:            viewModeNamespaceView,
			}

			// Simulate Shift+Tab key
			msg := tea.KeyMsg{Type: tea.KeyShiftTab}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedFocus, m.focusedPanel, "Focus panel mismatch")
			assert.Equal(t, tt.expectedPodIndex, m.selectedPodIndex, "Selected pod index mismatch")
			assert.Equal(t, tt.expectedActionIndex, m.selectedActionIndex, "Selected action index mismatch")
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

// TestActionsPanelFocusInitialization tests that actions panel state is initialized correctly (Story 4.1)
func TestActionsPanelFocusInitialization(t *testing.T) {
	adapter := &mockKubeAdapter{}
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{{
			Name: "test",
			Actions: []config.Action{
				{Name: "Console", Shortcut: "c", Type: "pod_exec"},
			},
		}},
	}

	model := NewAppModel(cfg, adapter)

	assert.Equal(t, -1, model.selectedActionIndex, "Should have no action selected initially")
	assert.Equal(t, 0, model.actionsScrollOffset, "Should have no scroll offset initially")
	assert.Equal(t, 1, len(model.actions), "Actions should be loaded from context")
}

// TestAutoSelectFirstAction tests that first action is auto-selected when focusing actions panel (Story 4.1)
func TestAutoSelectFirstAction(t *testing.T) {
	model := AppModel{
		focusedPanel: PanelPods,
		actions: []config.Action{
			{Name: "Console", Shortcut: "c", Type: "pod_exec"},
			{Name: "Logs", Shortcut: "l", Type: "local"},
		},
		selectedActionIndex: -1,
		viewMode:            viewModeNamespaceView,
	}

	// Tab to actions panel
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, PanelActions, m.focusedPanel)
	assert.Equal(t, 0, m.selectedActionIndex, "First action should be auto-selected")
}

// TestActionsNavigation_ArrowKeys tests arrow key navigation in actions panel (Story 4.1)
func TestActionsNavigation_ArrowKeys(t *testing.T) {
	actions := []config.Action{
		{Name: "Console", Shortcut: "c", Type: "pod_exec"},
		{Name: "Logs", Shortcut: "l", Type: "local"},
		{Name: "Dashboard", Shortcut: "d", Type: "url"},
	}

	tests := []struct {
		name          string
		initialIndex  int
		keyPress      tea.KeyType
		focusedPanel  PanelType
		expectedIndex int
	}{
		{
			name:          "Down arrow increments selection",
			initialIndex:  0,
			keyPress:      tea.KeyDown,
			focusedPanel:  PanelActions,
			expectedIndex: 1,
		},
		{
			name:          "Up arrow decrements selection",
			initialIndex:  2,
			keyPress:      tea.KeyUp,
			focusedPanel:  PanelActions,
			expectedIndex: 1,
		},
		{
			name:          "Down arrow clamped at last action",
			initialIndex:  2,
			keyPress:      tea.KeyDown,
			focusedPanel:  PanelActions,
			expectedIndex: 2,
		},
		{
			name:          "Up arrow clamped at first action",
			initialIndex:  0,
			keyPress:      tea.KeyUp,
			focusedPanel:  PanelActions,
			expectedIndex: 0,
		},
		{
			name:          "Arrow keys ignored when not focused on actions panel",
			initialIndex:  1,
			keyPress:      tea.KeyDown,
			focusedPanel:  PanelPods,
			expectedIndex: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := AppModel{
				focusedPanel:        tt.focusedPanel,
				selectedActionIndex: tt.initialIndex,
				actions:             actions,
				viewMode:            viewModeNamespaceView,
				termHeight:          30,              // Set terminal height for scroll calculations
				keys:                DefaultKeyMap(), // Set key mappings
			}

			msg := tea.KeyMsg{Type: tt.keyPress}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectedIndex, m.selectedActionIndex)
		})
	}
}
