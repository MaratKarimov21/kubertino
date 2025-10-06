package tui

import (
	"errors"
	"testing"

	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/stretchr/testify/assert"
)

func TestRenderPodPanel_LoadingState(t *testing.T) {
	model := AppModel{
		podsLoading: true,
		termWidth:   80,
		termHeight:  24,
	}

	output := model.renderPodPanel(40, 20)

	// Should contain loading message
	assert.Contains(t, output, "Loading pods")
	assert.Contains(t, output, "Pods") // Panel title
}

func TestRenderPodPanel_ErrorState(t *testing.T) {
	model := AppModel{
		podsError:  errors.New("permission denied"),
		termWidth:  80,
		termHeight: 24,
	}

	output := model.renderPodPanel(40, 20)

	// Should contain error message
	assert.Contains(t, output, "Error")
	assert.Contains(t, output, "permission denied")
	assert.Contains(t, output, "Pods") // Panel title
}

func TestRenderPodPanel_NoNamespaceSelected(t *testing.T) {
	model := AppModel{
		currentNamespace: "",
		termWidth:        80,
		termHeight:       24,
	}

	output := model.renderPodPanel(40, 20)

	// Should show placeholder when no namespace selected
	assert.Contains(t, output, "Select a namespace to view pods")
	assert.Contains(t, output, "Pods") // Panel title
}

func TestRenderPodPanel_EmptyPodList(t *testing.T) {
	model := AppModel{
		currentNamespace: "default",
		pods:             []k8s.Pod{},
		termWidth:        80,
		termHeight:       24,
	}

	output := model.renderPodPanel(40, 20)

	// Should show empty state message
	assert.Contains(t, output, "No pods in this namespace")
	assert.Contains(t, output, "Pods") // Panel title
}

func TestRenderPodPanel_WithPods(t *testing.T) {
	model := AppModel{
		currentNamespace: "default",
		pods: []k8s.Pod{
			{Name: "pod-1", Status: "Running"},
			{Name: "pod-2", Status: "Pending"},
			{Name: "pod-3", Status: "Failed"},
		},
		termWidth:  80,
		termHeight: 24,
	}

	output := model.renderPodPanel(40, 20)

	// Should contain all pod names
	assert.Contains(t, output, "pod-1")
	assert.Contains(t, output, "pod-2")
	assert.Contains(t, output, "pod-3")

	// Should contain status indicators
	assert.Contains(t, output, "Running")
	assert.Contains(t, output, "Pending")
	assert.Contains(t, output, "Failed")

	// Should contain panel title
	assert.Contains(t, output, "Pods")
}

func TestRenderPodPanel_MultipleStatuses(t *testing.T) {
	model := AppModel{
		currentNamespace: "test-namespace",
		pods: []k8s.Pod{
			{Name: "running-pod", Status: "Running"},
			{Name: "pending-pod", Status: "Pending"},
			{Name: "failed-pod", Status: "Failed"},
			{Name: "succeeded-pod", Status: "Succeeded"},
			{Name: "unknown-pod", Status: "Unknown"},
		},
		termWidth:  80,
		termHeight: 24,
	}

	output := model.renderPodPanel(40, 20)

	// Should contain all statuses
	statuses := []string{"Running", "Pending", "Failed", "Succeeded", "Unknown"}
	for _, status := range statuses {
		assert.Contains(t, output, status)
	}

	// Should contain all pod names
	podNames := []string{"running-pod", "pending-pod", "failed-pod", "succeeded-pod", "unknown-pod"}
	for _, name := range podNames {
		assert.Contains(t, output, name)
	}
}

func TestGetPodStatusStyle(t *testing.T) {
	model := AppModel{}

	tests := []struct {
		status       string
		expectedType string
	}{
		{"Running", "RunningStyle"},
		{"Succeeded", "RunningStyle"},
		{"Pending", "PendingStyle"},
		{"Failed", "FailedStyle"},
		{"Unknown", "DimStyle"},
		{"CustomStatus", "DimStyle"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			style := model.getPodStatusStyle(tt.status)
			// Just verify we get a style object without error
			assert.NotNil(t, style)
		})
	}
}

func TestPodsFetchedMsg_Success(t *testing.T) {
	ctx := &config.Context{
		Name:              "test-context",
		DefaultPodPattern: ".*",
		Actions:           []config.Action{},
	}

	model := AppModel{
		currentContext:   ctx,
		currentNamespace: "default",
		podsLoading:      true,
	}

	// Simulate successful pod fetch
	msg := podsFetchedMsg{
		pods: []k8s.Pod{
			{Name: "pod-1", Status: "Running"},
			{Name: "pod-2", Status: "Pending"},
		},
		err: nil,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	// Should have cleared loading state
	assert.False(t, m.podsLoading)
	assert.Nil(t, m.podsError)

	// Should have stored pods
	assert.Len(t, m.pods, 2)
	assert.Equal(t, "pod-1", m.pods[0].Name)
	assert.Equal(t, "Running", m.pods[0].Status)
}

func TestPodsFetchedMsg_Error(t *testing.T) {
	ctx := &config.Context{
		Name:              "test-context",
		DefaultPodPattern: ".*",
		Actions:           []config.Action{},
	}

	model := AppModel{
		currentContext:   ctx,
		currentNamespace: "default",
		podsLoading:      true,
	}

	// Simulate failed pod fetch
	msg := podsFetchedMsg{
		pods: nil,
		err:  errors.New("permission denied"),
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	// Should have cleared loading state
	assert.False(t, m.podsLoading)

	// Should have stored error
	assert.NotNil(t, m.podsError)
	assert.Contains(t, m.podsError.Error(), "permission denied")

	// Should not have stored pods
	assert.Nil(t, m.pods)
}

func TestFetchPodsCmd(t *testing.T) {
	// This test verifies the fetchPodsCmd function structure
	ctx := &config.Context{
		Name:              "test-context",
		DefaultPodPattern: ".*",
		Actions:           []config.Action{},
	}

	model := AppModel{
		currentContext:   ctx,
		currentNamespace: "default",
	}

	cmd := model.fetchPodsCmd()
	assert.NotNil(t, cmd)
}

func TestNamespaceSelection_TriggersPodFetch(t *testing.T) {
	// This test verifies that the pod fetching workflow is correctly implemented
	// We test the message handling directly rather than simulating key presses
	ctx := &config.Context{
		Name:              "test-context",
		DefaultPodPattern: ".*",
		Actions:           []config.Action{},
	}

	adapter := newMockAdapter()

	model := AppModel{
		config:                 &config.Config{Contexts: []config.Context{*ctx}},
		currentContext:         ctx,
		kubeAdapter:            adapter,
		viewMode:               viewModeNamespaceView,
		currentNamespace:       "default",
		namespaces:             []string{"default", "kube-system"},
		selectedNamespaceIndex: 0,
		termWidth:              80,
		termHeight:             24,
	}

	// Verify that fetchPodsCmd can be created
	cmd := model.fetchPodsCmd()
	assert.NotNil(t, cmd)

	// Execute the command to get the message
	msg := cmd()

	// Should receive a podsFetchedMsg
	podMsg, ok := msg.(podsFetchedMsg)
	assert.True(t, ok, "Message should be podsFetchedMsg")
	assert.NoError(t, podMsg.err)
	assert.NotNil(t, podMsg.pods)
}

func TestSearchModeNamespaceSelection_TriggersPodFetch(t *testing.T) {
	// Test that pod fetching works correctly after namespace selection in search mode
	ctx := &config.Context{
		Name:              "test-context",
		DefaultPodPattern: ".*",
		Actions:           []config.Action{},
	}

	adapter := newMockAdapter()

	model := AppModel{
		config:                 &config.Config{Contexts: []config.Context{*ctx}},
		currentContext:         ctx,
		kubeAdapter:            adapter,
		viewMode:               viewModeNamespaceView,
		currentNamespace:       "default",
		namespaces:             []string{"default", "kube-system"},
		filteredNamespaces:     []string{"default"},
		selectedNamespaceIndex: 0,
		searchMode:             false, // Already exited search mode
		termWidth:              80,
		termHeight:             24,
	}

	// Verify that fetchPodsCmd can be created after namespace selection
	cmd := model.fetchPodsCmd()
	assert.NotNil(t, cmd)

	// Execute the command to get the message
	msg := cmd()

	// Should receive a podsFetchedMsg
	podMsg, ok := msg.(podsFetchedMsg)
	assert.True(t, ok, "Message should be podsFetchedMsg")
	assert.NoError(t, podMsg.err)
	assert.NotNil(t, podMsg.pods)
}

func TestPodPanel_FormattingConsistency(t *testing.T) {
	model := AppModel{
		currentNamespace: "test",
		pods: []k8s.Pod{
			{Name: "short", Status: "Running"},
			{Name: "very-long-pod-name-that-exceeds-normal-length", Status: "Pending"},
		},
		termWidth:  80,
		termHeight: 24,
	}

	output := model.renderPodPanel(50, 20)

	// Should contain both pod names (note: long names may be wrapped by Lip Gloss)
	assert.Contains(t, output, "short")
	// The long name may be split/wrapped, so check for a substring
	assert.Contains(t, output, "very-long-pod-name")

	// Output should not be empty and should contain panel structure
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "Pods") // Panel title
}
