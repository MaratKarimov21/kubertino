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
		defaultPodIndex:  3,
		focusedPanel:     PanelNamespaces, // Must be on namespaces to select
		keys:             DefaultKeyMap(), // Need keys for KeyMatches
	}

	// Select a namespace (triggers pod fetch and state reset)
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, -1, m.selectedPodIndex, "Selected pod index should be reset")
	assert.Equal(t, 0, m.podScrollOffset, "Scroll offset should be reset")
	assert.Equal(t, -1, m.defaultPodIndex, "Default pod index should be reset")
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
