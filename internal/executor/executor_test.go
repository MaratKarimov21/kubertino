package executor

import (
	"fmt"
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

// TestExecuteLocal tests the ExecuteLocal method with various scenarios
func TestExecuteLocal(t *testing.T) {
	tests := []struct {
		name           string
		action         config.Action
		context        config.Context
		namespace      string
		pods           []k8s.Pod
		kubeconfigPath string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful template substitution with all variables",
			action: config.Action{
				Name:       "Logs",
				Command:    "kubectl logs -n {{.namespace}} {{.pod}}",
				PodPattern: "rails-web-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
			},
			kubeconfigPath: "/path/to/kubeconfig",
			wantErr:        false,
		},
		{
			name: "template with context variable",
			action: config.Action{
				Name:       "Custom",
				Command:    "echo {{.context}}/{{.namespace}}/{{.pod}}",
				PodPattern: ".*",
			},
			context: config.Context{
				Name:              "staging",
				DefaultPodPattern: ".*",
			},
			namespace: "default",
			pods: []k8s.Pod{
				{Name: "test-pod", Status: "Running"},
			},
			kubeconfigPath: "",
			wantErr:        false,
		},
		{
			name: "invalid template syntax",
			action: config.Action{
				Name:       "Bad",
				Command:    "kubectl logs {{.namespace",
				PodPattern: ".*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "test-pod", Status: "Running"},
			},
			wantErr:     true,
			errContains: "invalid command template",
		},
		{
			name: "no pods match pattern",
			action: config.Action{
				Name:       "Logs",
				Command:    "kubectl logs {{.pod}}",
				PodPattern: "nonexistent-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
			},
			wantErr:     true,
			errContains: "no pods match pattern",
		},
		{
			name: "multiple pods match pattern",
			action: config.Action{
				Name:       "Logs",
				Command:    "kubectl logs {{.pod}}",
				PodPattern: "rails-web-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
				{Name: "rails-web-def456", Status: "Running"},
			},
			wantErr:     true,
			errContains: "multiple pods match pattern",
		},
		{
			name: "uses context default pattern when action pattern empty",
			action: config.Action{
				Name:       "Logs",
				Command:    "kubectl logs {{.pod}}",
				PodPattern: "",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: "rails-web-.*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &mockKubeAdapter{}
			executor := NewExecutor(adapter)

			// Note: We can't actually run the command in tests, so we skip this for now
			// In a real test environment, we would mock exec.Command
			err := executor.ExecuteLocal(tt.action, tt.context, tt.namespace, tt.pods, tt.kubeconfigPath)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				// Commands will fail in test environment since we don't have kubectl
				// But we can verify template parsing and pod matching worked
				if err != nil {
					// Allow command execution failures in tests
					assert.Contains(t, err.Error(), "command execution failed")
				}
			}
		})
	}
}

// TestPrepareLocal tests the PrepareLocal method
func TestPrepareLocal(t *testing.T) {
	tests := []struct {
		name           string
		action         config.Action
		context        config.Context
		namespace      string
		pods           []k8s.Pod
		kubeconfigPath string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful preparation with all variables",
			action: config.Action{
				Name:       "Logs",
				Command:    "kubectl logs -n {{.namespace}} {{.pod}}",
				PodPattern: "rails-web-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
			},
			kubeconfigPath: "/path/to/kubeconfig",
			wantErr:        false,
		},
		{
			name: "preparation without kubeconfig path",
			action: config.Action{
				Name:       "Echo",
				Command:    "echo {{.namespace}}",
				PodPattern: ".*",
			},
			context: config.Context{
				Name:              "staging",
				DefaultPodPattern: ".*",
			},
			namespace: "default",
			pods: []k8s.Pod{
				{Name: "test-pod", Status: "Running"},
			},
			kubeconfigPath: "",
			wantErr:        false,
		},
		{
			name: "invalid template syntax",
			action: config.Action{
				Name:       "Bad",
				Command:    "kubectl logs {{.namespace",
				PodPattern: ".*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "test-pod", Status: "Running"},
			},
			wantErr:     true,
			errContains: "invalid command template",
		},
		{
			name: "no pods match",
			action: config.Action{
				Name:       "Logs",
				Command:    "kubectl logs {{.pod}}",
				PodPattern: "nonexistent-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
			},
			wantErr:     true,
			errContains: "no pods match pattern",
		},
		{
			name: "multiple pods match",
			action: config.Action{
				Name:       "Logs",
				Command:    "kubectl logs {{.pod}}",
				PodPattern: "rails-web-.*",
			},
			context: config.Context{
				Name:              "production",
				DefaultPodPattern: ".*",
			},
			namespace: "app",
			pods: []k8s.Pod{
				{Name: "rails-web-abc123", Status: "Running"},
				{Name: "rails-web-def456", Status: "Running"},
			},
			wantErr:     true,
			errContains: "multiple pods match pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := &mockKubeAdapter{}
			executor := NewExecutor(adapter)

			cmd, err := executor.PrepareLocal(tt.action, tt.context, tt.namespace, tt.pods, tt.kubeconfigPath)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, cmd)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.NotNil(t, cmd)
				assert.NotNil(t, cmd.Stdin)
				assert.NotNil(t, cmd.Stdout)
				assert.NotNil(t, cmd.Stderr)

				// Verify environment variables are set
				assert.NotEmpty(t, cmd.Env)

				// If kubeconfigPath is set, verify it's in the environment
				if tt.kubeconfigPath != "" {
					found := false
					for _, env := range cmd.Env {
						if env == fmt.Sprintf("KUBECONFIG=%s", tt.kubeconfigPath) {
							found = true
							break
						}
					}
					assert.True(t, found, "KUBECONFIG environment variable should be set")
				}
			}
		})
	}
}

// TestPrepareLocal_EnvironmentPreservation tests that parent environment is preserved
func TestPrepareLocal_EnvironmentPreservation(t *testing.T) {
	action := config.Action{
		Name:       "Echo",
		Command:    "echo test",
		PodPattern: ".*",
	}
	context := config.Context{
		Name:              "production",
		DefaultPodPattern: ".*",
	}
	pods := []k8s.Pod{
		{Name: "test-pod", Status: "Running"},
	}

	adapter := &mockKubeAdapter{}
	executor := NewExecutor(adapter)

	cmd, err := executor.PrepareLocal(action, context, "default", pods, "")

	require.NoError(t, err)
	assert.NotNil(t, cmd)

	// Environment should contain parent environment variables
	assert.Greater(t, len(cmd.Env), 0, "Environment should be preserved from parent")
}
