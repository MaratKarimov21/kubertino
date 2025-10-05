package executor

import (
	"os/exec"
	"testing"

	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockKubeAdapter is a test mock for KubeAdapter
type mockKubeAdapter struct {
	execInPodFunc func(context, namespace, pod, container, command string) (*exec.Cmd, error)
}

func (m *mockKubeAdapter) ExecInPod(context, namespace, pod, container, command string) (*exec.Cmd, error) {
	if m.execInPodFunc != nil {
		return m.execInPodFunc(context, namespace, pod, container, command)
	}
	// Default: return a simple echo command
	return exec.Command("echo", "test"), nil
}

// TestMatchPod tests the matchPod function with various scenarios
func TestMatchPod(t *testing.T) {
	pods := []k8s.Pod{
		{Name: "rails-web-abc123", Status: "Running"},
		{Name: "rails-web-def456", Status: "Running"},
		{Name: "worker-xyz789", Status: "Running"},
	}

	tests := []struct {
		name        string
		pattern     string
		wantPodName string
		wantErr     bool
		errContains string
	}{
		{
			name:        "exact match",
			pattern:     "^rails-web-abc123$",
			wantPodName: "rails-web-abc123",
			wantErr:     false,
		},
		{
			name:        "regex match single pod",
			pattern:     "worker-.*",
			wantPodName: "worker-xyz789",
			wantErr:     false,
		},
		{
			name:        "no match",
			pattern:     "notfound-.*",
			wantErr:     true,
			errContains: "no pods match pattern",
		},
		{
			name:        "multiple matches",
			pattern:     "rails-web-.*",
			wantErr:     true,
			errContains: "multiple pods match pattern",
		},
		{
			name:        "invalid regex",
			pattern:     "[invalid",
			wantErr:     true,
			errContains: "invalid pod pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod, err := matchPod(pods, tt.pattern)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPodName, pod.Name)
			}
		})
	}
}

// TestMatchPod_EmptyPodList tests matchPod with empty pod list
func TestMatchPod_EmptyPodList(t *testing.T) {
	pods := []k8s.Pod{}

	pod, err := matchPod(pods, ".*")

	require.Error(t, err)
	assert.Nil(t, pod)
	assert.Contains(t, err.Error(), "no pods in namespace")
}

// TestGetPodPattern tests the getPodPattern function
func TestGetPodPattern(t *testing.T) {
	tests := []struct {
		name        string
		action      config.Action
		context     config.Context
		wantPattern string
	}{
		{
			name: "action pattern takes precedence",
			action: config.Action{
				PodPattern: "action-pattern-.*",
			},
			context: config.Context{
				DefaultPodPattern: "context-pattern-.*",
			},
			wantPattern: "action-pattern-.*",
		},
		{
			name: "context pattern used when action pattern empty",
			action: config.Action{
				PodPattern: "",
			},
			context: config.Context{
				DefaultPodPattern: "context-pattern-.*",
			},
			wantPattern: "context-pattern-.*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := getPodPattern(tt.action, tt.context)
			assert.Equal(t, tt.wantPattern, pattern)
		})
	}
}

// TestNewExecutor tests executor initialization
func TestNewExecutor(t *testing.T) {
	adapter := &mockKubeAdapter{}
	executor := NewExecutor(adapter)

	assert.NotNil(t, executor)
	assert.Equal(t, adapter, executor.kubeAdapter)
}

// TestPreparePodExec tests preparing a pod exec command
func TestPreparePodExec(t *testing.T) {
	tests := []struct {
		name      string
		action    config.Action
		context   config.Context
		namespace string
		pods      []k8s.Pod
		wantErr   bool
	}{
		{
			name: "successful preparation",
			action: config.Action{
				Name:       "Console",
				Type:       "pod_exec",
				Command:    "rails console",
				PodPattern: "rails-web-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "default",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
			},
			wantErr: false,
		},
		{
			name: "no pods match",
			action: config.Action{
				Name:       "Console",
				Type:       "pod_exec",
				Command:    "rails console",
				PodPattern: "rails-web-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "default",
			pods: []k8s.Pod{
				{Name: "worker-abc123", Status: "Running"},
			},
			wantErr: true,
		},
		{
			name: "multiple pods match",
			action: config.Action{
				Name:       "Console",
				Type:       "pod_exec",
				Command:    "rails console",
				PodPattern: "rails-web-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "default",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
				{Name: "rails-web-def456", Status: "Running"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &mockKubeAdapter{}
			executor := NewExecutor(adapter)

			cmd, err := executor.PreparePodExec(tt.action, tt.context, tt.namespace, tt.pods)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cmd)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cmd)
			}
		})
	}
}

// TestPreparePodExec_UsesDefaultPattern tests that default pattern is used when action pattern is empty
func TestPreparePodExec_UsesDefaultPattern(t *testing.T) {
	action := config.Action{
		Name:       "Console",
		Type:       "pod_exec",
		Command:    "rails console",
		PodPattern: "", // Empty pattern
	}
	context := config.Context{
		Name:              "production",
		DefaultPodPattern: "rails-web-.*",
	}
	pods := []k8s.Pod{
		{Name: "rails-web-abc123", Status: "Running"},
	}

	adapter := &mockKubeAdapter{}
	executor := NewExecutor(adapter)

	cmd, err := executor.PreparePodExec(action, context, "default", pods)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
}
