package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, cfg *Config)
	}{
		{
			name:     "valid config",
			filename: "../testdata/valid-config.yml",
			wantErr:  false,
			validate: func(t *testing.T, cfg *Config) {
				assert.Equal(t, "1.0", cfg.Version)
				require.Len(t, cfg.Contexts, 1)

				ctx := cfg.Contexts[0]
				assert.Equal(t, "test-context", ctx.Name)
				assert.Equal(t, "~/.kube/config", ctx.Kubeconfig)
				assert.Equal(t, "https://k8s.test.example.com", ctx.ClusterURL)
				assert.Equal(t, "^api-.*", ctx.DefaultPodPattern)
				assert.Equal(t, []string{"default", "api"}, ctx.FavoriteNamespaces)

				require.Len(t, ctx.Actions, 3)

				// Check pod_exec action
				assert.Equal(t, "Console", ctx.Actions[0].Name)
				assert.Equal(t, "c", ctx.Actions[0].Shortcut)
				assert.Equal(t, "pod_exec", ctx.Actions[0].Type)
				assert.Equal(t, "/bin/sh", ctx.Actions[0].Command)
				assert.Equal(t, "^api-.*", ctx.Actions[0].PodPattern)

				// Check url action
				assert.Equal(t, "Grafana", ctx.Actions[1].Name)
				assert.Equal(t, "g", ctx.Actions[1].Shortcut)
				assert.Equal(t, "url", ctx.Actions[1].Type)
				assert.Equal(t, "https://grafana.example.com", ctx.Actions[1].URL)

				// Check local action
				assert.Equal(t, "Port Forward", ctx.Actions[2].Name)
				assert.Equal(t, "p", ctx.Actions[2].Shortcut)
				assert.Equal(t, "local", ctx.Actions[2].Type)
				assert.Equal(t, "kubectl port-forward pod 8080:8080", ctx.Actions[2].Command)
			},
		},
		{
			name:        "file not found",
			filename:    "../testdata/nonexistent.yml",
			wantErr:     true,
			errContains: "failed to read config file",
		},
		{
			name:        "invalid yaml",
			filename:    "../testdata/invalid-yaml.yml",
			wantErr:     true,
			errContains: "failed to parse YAML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Parse(tt.filename)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

func TestParse_TildeExpansion(t *testing.T) {
	// Create temp config file
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tempFile := filepath.Join(homeDir, ".kubertino-test.yml")
	defer func() {
		_ = os.Remove(tempFile)
	}()

	content := `version: "1.0"
contexts:
  - name: test
    default_pod_pattern: ".*"`

	err = os.WriteFile(tempFile, []byte(content), 0644)
	require.NoError(t, err)

	// Test parsing with tilde path
	cfg, err := Parse("~/.kubertino-test.yml")
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.Contexts[0].Name)
}

// TestParse_RegexCompilation tests that default_pod_pattern is compiled during parsing (Story 3.2)
func TestParse_RegexCompilation(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		wantPattern string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid regex pattern",
			yamlContent: `version: "1.0"
contexts:
  - name: test
    default_pod_pattern: "^api-.*"`,
			wantPattern: "^api-.*",
			wantErr:     false,
		},
		{
			name: "invalid regex pattern",
			yamlContent: `version: "1.0"
contexts:
  - name: test
    default_pod_pattern: "api-["`,
			wantErr:     true,
			errContains: "invalid default_pod_pattern",
		},
		{
			name: "empty pattern (valid)",
			yamlContent: `version: "1.0"
contexts:
  - name: test
    default_pod_pattern: ""`,
			wantPattern: "",
			wantErr:     false,
		},
		{
			name: "complex regex pattern",
			yamlContent: `version: "1.0"
contexts:
  - name: test
    default_pod_pattern: "^(api|web)-.*-v[0-9]+$"`,
			wantPattern: "^(api|web)-.*-v[0-9]+$",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tempFile := filepath.Join(t.TempDir(), "config.yml")
			err := os.WriteFile(tempFile, []byte(tt.yamlContent), 0644)
			require.NoError(t, err)

			cfg, err := Parse(tempFile)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, cfg)
			assert.Equal(t, tt.wantPattern, cfg.Contexts[0].DefaultPodPattern)
			if tt.wantPattern != "" {
				assert.NotNil(t, cfg.Contexts[0].CompiledPattern)
			}
		})
	}
}

// Helper functions for creating test data
func NewTestConfig() *Config {
	return &Config{
		Version: "1.0",
		Contexts: []Context{
			NewTestContext("test-context"),
		},
	}
}

func NewTestContext(name string) Context {
	return Context{
		Name:               name,
		Kubeconfig:         "~/.kube/config",
		DefaultPodPattern:  ".*",
		FavoriteNamespaces: []string{"default"},
		Actions: []Action{
			NewTestAction("test", "t", "pod_exec", "/bin/sh"),
		},
	}
}

func NewTestAction(name, shortcut, actionType, command string) Action {
	return Action{
		Name:     name,
		Shortcut: shortcut,
		Type:     actionType,
		Command:  command,
	}
}
