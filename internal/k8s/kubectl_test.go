package k8s

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/maratkarimov/kubertino/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContexts(t *testing.T) {
	tests := []struct {
		name           string
		kubeconfigPath string
		want           []string
		wantErr        bool
		errType        error
	}{
		{
			name:           "valid kubeconfig",
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			want:           []string{"minikube", "prod-cluster"},
			wantErr:        false,
		},
		{
			name:           "empty kubeconfig",
			kubeconfigPath: "../testdata/empty-kubeconfig.yml",
			want:           []string{},
			wantErr:        false,
		},
		{
			name:           "missing file",
			kubeconfigPath: "../testdata/nonexistent.yml",
			want:           nil,
			wantErr:        true,
			errType:        ErrKubeconfigNotFound,
		},
		{
			name:           "invalid YAML",
			kubeconfigPath: "../testdata/invalid-kubeconfig.yml",
			want:           nil,
			wantErr:        true,
			errType:        ErrInvalidKubeconfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := NewKubectlAdapter(tt.kubeconfigPath)
			got, err := adapter.GetContexts()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.Is(err, tt.errType), "expected error type %v, got %v", tt.errType, err)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "tilde only",
			path:    "~",
			want:    homeDir,
			wantErr: false,
		},
		{
			name:    "tilde with path",
			path:    "~/.kube/config",
			want:    filepath.Join(homeDir, ".kube/config"),
			wantErr: false,
		},
		{
			name:    "absolute path",
			path:    "/absolute/path",
			want:    "/absolute/path",
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    "relative/path",
			want:    "relative/path",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandPath(tt.path)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMatchContexts(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		kubeconfigPath string
		wantMatched    []string
		wantMissing    []string
		wantErr        bool
	}{
		{
			name: "all contexts match",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "minikube"},
					{Name: "prod-cluster"},
				},
			},
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			wantMatched:    []string{"minikube", "prod-cluster"},
			wantMissing:    nil,
			wantErr:        false,
		},
		{
			name: "partial match",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "minikube"},
					{Name: "nonexistent-context"},
				},
			},
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			wantMatched:    []string{"minikube"},
			wantMissing:    []string{"nonexistent-context"},
			wantErr:        false,
		},
		{
			name: "no contexts match",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "nonexistent-1"},
					{Name: "nonexistent-2"},
				},
			},
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			wantMatched:    nil,
			wantMissing:    []string{"nonexistent-1", "nonexistent-2"},
			wantErr:        false,
		},
		{
			name: "empty config",
			config: &config.Config{
				Contexts: []config.Context{},
			},
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			wantMatched:    nil,
			wantMissing:    nil,
			wantErr:        false,
		},
		{
			name: "kubeconfig read error",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "minikube"},
				},
			},
			kubeconfigPath: "../testdata/nonexistent.yml",
			wantMatched:    nil,
			wantMissing:    nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, missing, err := MatchContexts(tt.config, tt.kubeconfigPath)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantMatched, matched)
			assert.Equal(t, tt.wantMissing, missing)
		})
	}
}

func TestValidateContexts(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		kubeconfigPath string
		wantErr        bool
	}{
		{
			name: "valid contexts",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "minikube"},
				},
			},
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			wantErr:        false,
		},
		{
			name: "no valid contexts",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "nonexistent-context"},
				},
			},
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			wantErr:        true,
		},
		{
			name: "partial match is valid",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "minikube"},
					{Name: "nonexistent-context"},
				},
			},
			kubeconfigPath: "../testdata/valid-kubeconfig.yml",
			wantErr:        false,
		},
		{
			name: "kubeconfig read error",
			config: &config.Config{
				Contexts: []config.Context{
					{Name: "minikube"},
				},
			},
			kubeconfigPath: "../testdata/nonexistent.yml",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContexts(tt.config, tt.kubeconfigPath)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestKubectlAdapter_NotImplemented(t *testing.T) {
	adapter := NewKubectlAdapter("~/.kube/config")

	t.Run("GetNamespaces not implemented", func(t *testing.T) {
		_, err := adapter.GetNamespaces("test-context")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("GetPods not implemented", func(t *testing.T) {
		_, err := adapter.GetPods("test-context", "test-namespace")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})
}
