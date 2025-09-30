package k8s

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/maratkarimov/kubertino/internal/config"

	"gopkg.in/yaml.v3"
)

// Custom error types
var (
	ErrKubeconfigNotFound = errors.New("kubeconfig file not found")
	ErrContextNotFound    = errors.New("context not found")
	ErrInvalidKubeconfig  = errors.New("invalid kubeconfig format")
)

// KubectlAdapter implements KubeAdapter by reading kubeconfig files
type KubectlAdapter struct {
	kubeconfigPath string
}

// NewKubectlAdapter creates a new KubectlAdapter with the specified kubeconfig path
func NewKubectlAdapter(kubeconfigPath string) *KubectlAdapter {
	return &KubectlAdapter{
		kubeconfigPath: kubeconfigPath,
	}
}

// GetContexts reads the kubeconfig file and returns available context names
func (k *KubectlAdapter) GetContexts() ([]string, error) {
	kubeconfigPath, err := expandPath(k.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand kubeconfig path: %w", err)
	}

	data, err := os.ReadFile(kubeconfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrKubeconfigNotFound, kubeconfigPath)
		}
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	var config KubeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKubeconfig, err)
	}

	contexts := make([]string, 0, len(config.Contexts))
	for _, ctx := range config.Contexts {
		contexts = append(contexts, ctx.Name)
	}

	return contexts, nil
}

// GetNamespaces is a placeholder for future implementation
func (k *KubectlAdapter) GetNamespaces(context string) ([]string, error) {
	return nil, errors.New("not implemented")
}

// GetPods is a placeholder for future implementation
func (k *KubectlAdapter) GetPods(context, namespace string) ([]Pod, error) {
	return nil, errors.New("not implemented")
}

// MatchContexts compares configured contexts with available kubectl contexts
// Returns matched contexts, missing contexts, and any error encountered
func MatchContexts(cfg *config.Config, kubeconfigPath string) (matched, missing []string, err error) {
	adapter := NewKubectlAdapter(kubeconfigPath)
	availableContexts, err := adapter.GetContexts()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read kubectl contexts: %w", err)
	}

	contextSet := make(map[string]bool)
	for _, ctx := range availableContexts {
		contextSet[ctx] = true
	}

	for _, configCtx := range cfg.Contexts {
		if contextSet[configCtx.Name] {
			matched = append(matched, configCtx.Name)
		} else {
			missing = append(missing, configCtx.Name)
		}
	}

	return matched, missing, nil
}

// ValidateContexts validates that at least one configured context exists in kubeconfig
// Returns an error if no valid contexts are found
func ValidateContexts(cfg *config.Config, kubeconfigPath string) error {
	matched, missing, err := MatchContexts(cfg, kubeconfigPath)
	if err != nil {
		return fmt.Errorf("failed to match contexts: %w", err)
	}

	if len(matched) == 0 {
		return fmt.Errorf("no valid contexts found: configured contexts %v not found in %s. Check your kubeconfig and ~/.kubertino.yml", missing, kubeconfigPath)
	}

	return nil
}

// expandPath expands ~ to the user's home directory
func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	if path == "~" {
		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}
