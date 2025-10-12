package executor

import (
	"fmt"
	"testing"

	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewExecutor tests executor initialization
func TestNewExecutor(t *testing.T) {
	executor := NewExecutor()

	assert.NotNil(t, executor)
}

// TestPrepareLocal tests the PrepareLocal method with template substitution
func TestPrepareLocal(t *testing.T) {
	tests := []struct {
		name           string
		action         config.Action
		context        config.Context
		namespace      string
		pod            k8s.Pod
		kubeconfigPath string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful preparation with all variables",
			action: config.Action{
				Name:    "Logs",
				Command: "kubectl logs -n {{.namespace}} {{.pod}}",
			},
			context: config.Context{
				Name: "production",
			},
			namespace:      "app",
			pod:            k8s.Pod{Name: "rails-web-abc123", Status: "Running"},
			kubeconfigPath: "/path/to/kubeconfig",
			wantErr:        false,
		},
		{
			name: "preparation without kubeconfig path",
			action: config.Action{
				Name:    "Echo",
				Command: "echo {{.namespace}}",
			},
			context: config.Context{
				Name: "staging",
			},
			namespace:      "default",
			pod:            k8s.Pod{Name: "test-pod", Status: "Running"},
			kubeconfigPath: "",
			wantErr:        false,
		},
		{
			name: "template with context variable",
			action: config.Action{
				Name:    "Custom",
				Command: "echo {{.context}}/{{.namespace}}/{{.pod}}",
			},
			context: config.Context{
				Name: "staging",
			},
			namespace:      "default",
			pod:            k8s.Pod{Name: "test-pod", Status: "Running"},
			kubeconfigPath: "",
			wantErr:        false,
		},
		{
			name: "invalid template syntax",
			action: config.Action{
				Name:    "Bad",
				Command: "kubectl logs {{.namespace",
			},
			context: config.Context{
				Name: "production",
			},
			namespace:   "app",
			pod:         k8s.Pod{Name: "test-pod", Status: "Running"},
			wantErr:     true,
			errContains: "invalid command template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor()

			cmd, err := executor.PrepareLocal(tt.action, tt.context, tt.namespace, tt.pod, tt.kubeconfigPath)

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
		Name:    "Echo",
		Command: "echo test",
	}
	context := config.Context{
		Name: "production",
	}
	pod := k8s.Pod{Name: "test-pod", Status: "Running"}

	executor := NewExecutor()

	cmd, err := executor.PrepareLocal(action, context, "default", pod, "")

	require.NoError(t, err)
	assert.NotNil(t, cmd)

	// Environment should contain parent environment variables
	assert.Greater(t, len(cmd.Env), 0, "Environment should be preserved from parent")
}

// TestExecuteLocal tests the ExecuteLocal method with various scenarios
func TestExecuteLocal(t *testing.T) {
	tests := []struct {
		name           string
		action         config.Action
		context        config.Context
		namespace      string
		pod            k8s.Pod
		kubeconfigPath string
		wantErr        bool
		errContains    string
	}{
		{
			name: "successful template substitution with all variables",
			action: config.Action{
				Name:    "Logs",
				Command: "kubectl logs -n {{.namespace}} {{.pod}}",
			},
			context: config.Context{
				Name: "production",
			},
			namespace:      "app",
			pod:            k8s.Pod{Name: "rails-web-abc123", Status: "Running"},
			kubeconfigPath: "/path/to/kubeconfig",
			wantErr:        false,
		},
		{
			name: "template with context variable",
			action: config.Action{
				Name:    "Custom",
				Command: "echo {{.context}}/{{.namespace}}/{{.pod}}",
			},
			context: config.Context{
				Name: "staging",
			},
			namespace:      "default",
			pod:            k8s.Pod{Name: "test-pod", Status: "Running"},
			kubeconfigPath: "",
			wantErr:        false,
		},
		{
			name: "invalid template syntax",
			action: config.Action{
				Name:    "Bad",
				Command: "kubectl logs {{.namespace",
			},
			context: config.Context{
				Name: "production",
			},
			namespace:   "app",
			pod:         k8s.Pod{Name: "test-pod", Status: "Running"},
			wantErr:     true,
			errContains: "invalid command template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewExecutor()

			// Note: We can't actually run the command in tests, so we skip this for now
			// In a real test environment, we would mock exec.Command
			err := executor.ExecuteLocal(tt.action, tt.context, tt.namespace, tt.pod, tt.kubeconfigPath)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				// Commands will fail in test environment since we don't have kubectl
				// But we can verify template parsing worked
				if err != nil {
					// Allow command execution failures in tests
					assert.Contains(t, err.Error(), "command execution failed")
				}
			}
		})
	}
}
