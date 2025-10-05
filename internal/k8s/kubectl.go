package k8s

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/maratkarimov/kubertino/internal/config"

	"gopkg.in/yaml.v3"
)

// Custom error types
var (
	ErrKubeconfigNotFound = errors.New("kubeconfig file not found")
	ErrContextNotFound    = errors.New("context not found")
	ErrInvalidKubeconfig  = errors.New("invalid kubeconfig format")
	ErrKubectlNotFound    = errors.New("kubectl not found in PATH")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrTimeout            = errors.New("operation timeout")
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

// validateContextName validates a context name against allowed characters
// to prevent potential security issues. Allows common Kubernetes context name patterns
// including cluster/user notation (with /), email-based auth (with @), and various separators.
// Prevents only shell metacharacters that could cause issues with exec.CommandContext.
func validateContextName(name string) error {
	if name == "" {
		return fmt.Errorf("context name cannot be empty")
	}

	// Prevent shell metacharacters and control characters while allowing common k8s patterns
	// Blocked: semicolon, pipe, ampersand, backtick, dollar, parentheses, angle brackets, newlines, tabs
	// Allowed: alphanumeric, dash, underscore, dot, slash, at-sign, plus, colon
	matched, err := regexp.MatchString(`^[a-zA-Z0-9._\-/@+:]+$`, name)
	if err != nil {
		return fmt.Errorf("failed to validate context name: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid context name '%s': contains forbidden characters (use only alphanumeric, .-_/@+:)", name)
	}

	return nil
}

// GetNamespaces fetches namespaces for the specified context using kubectl
func (k *KubectlAdapter) GetNamespaces(ctxName string) ([]string, error) {
	// Validate context name for security
	if err := validateContextName(ctxName); err != nil {
		return nil, err
	}

	// Find kubectl in PATH
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrKubectlNotFound, err)
	}

	// Expand kubeconfig path
	kubeconfigPath, err := expandPath(k.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand kubeconfig path: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute kubectl command
	cmd := exec.CommandContext(ctx, kubectlPath, "--kubeconfig", kubeconfigPath, "--context", ctxName, "get", "namespaces", "-o", "json")
	output, err := cmd.Output()

	if err != nil {
		// Check for specific error types
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("%w: kubectl command timed out after 10s", ErrTimeout)
		}

		// Check for exit error to extract stderr
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)

			// Check for permission denied
			if strings.Contains(stderr, "forbidden") || strings.Contains(stderr, "Forbidden") {
				return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, stderr)
			}

			// Check for context not found
			if strings.Contains(stderr, "context") && (strings.Contains(stderr, "not found") || strings.Contains(stderr, "does not exist")) {
				return nil, fmt.Errorf("%w: %s", ErrContextNotFound, ctxName)
			}

			return nil, fmt.Errorf("kubectl command failed: %s", stderr)
		}

		return nil, fmt.Errorf("failed to execute kubectl: %w", err)
	}

	// Parse JSON response
	var response NamespaceList
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse kubectl output: %w", err)
	}

	// Extract namespace names
	namespaces := make([]string, 0, len(response.Items))
	for _, ns := range response.Items {
		namespaces = append(namespaces, ns.Metadata.Name)
	}

	return namespaces, nil
}

// GetPods fetches pods for the specified context and namespace using kubectl
func (k *KubectlAdapter) GetPods(ctxName, namespace string) ([]Pod, error) {
	// Validate context name for security
	if err := validateContextName(ctxName); err != nil {
		return nil, err
	}

	// Validate namespace name for security
	if err := validateNamespaceName(namespace); err != nil {
		return nil, err
	}

	// Find kubectl in PATH
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrKubectlNotFound, err)
	}

	// Expand kubeconfig path
	kubeconfigPath, err := expandPath(k.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand kubeconfig path: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute kubectl command
	cmd := exec.CommandContext(ctx, kubectlPath, "--kubeconfig", kubeconfigPath, "--context", ctxName, "get", "pods", "-n", namespace, "-o", "json")
	output, err := cmd.Output()

	if err != nil {
		// Check for specific error types
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("%w: kubectl command timed out after 10s", ErrTimeout)
		}

		// Check for exit error to extract stderr
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)

			// Check for permission denied
			if strings.Contains(stderr, "forbidden") || strings.Contains(stderr, "Forbidden") {
				return nil, fmt.Errorf("%w: %s", ErrPermissionDenied, stderr)
			}

			// Check for namespace not found
			if strings.Contains(stderr, "namespace") && (strings.Contains(stderr, "not found") || strings.Contains(stderr, "does not exist")) {
				return nil, fmt.Errorf("namespace not found: %s", namespace)
			}

			return nil, fmt.Errorf("kubectl command failed: %s", stderr)
		}

		return nil, fmt.Errorf("failed to execute kubectl: %w", err)
	}

	// Parse JSON response
	var response PodList
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("failed to parse kubectl output: %w", err)
	}

	// Convert to []Pod
	pods := make([]Pod, 0, len(response.Items))
	for _, item := range response.Items {
		pods = append(pods, Pod{
			Name:   item.Metadata.Name,
			Status: item.Status.Phase,
		})
	}

	return pods, nil
}

// validateNamespaceName validates a namespace name against allowed characters
func validateNamespaceName(name string) error {
	if name == "" {
		return fmt.Errorf("namespace name cannot be empty")
	}

	// Kubernetes namespace names must be valid DNS labels
	// RFC 1123: lowercase alphanumeric characters or '-', start/end with alphanumeric
	matched, err := regexp.MatchString(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, name)
	if err != nil {
		return fmt.Errorf("failed to validate namespace name: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid namespace name '%s': must be lowercase alphanumeric with optional hyphens", name)
	}

	return nil
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

// ExecInPod returns a command configured to execute in a pod with interactive TTY
func (k *KubectlAdapter) ExecInPod(ctxName, namespace, pod, container, command string) (*exec.Cmd, error) {
	// Validate inputs for security
	if err := validateContextName(ctxName); err != nil {
		return nil, err
	}
	if err := validateNamespaceName(namespace); err != nil {
		return nil, err
	}
	if err := validatePodName(pod); err != nil {
		return nil, err
	}
	if container != "" {
		if err := validateContainerName(container); err != nil {
			return nil, err
		}
	}

	// Find kubectl in PATH
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrKubectlNotFound, err)
	}

	// Expand kubeconfig path
	kubeconfigPath, err := expandPath(k.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to expand kubeconfig path: %w", err)
	}

	// Build kubectl exec command with interactive TTY
	args := []string{
		"--kubeconfig", kubeconfigPath,
		"--context", ctxName,
		"-n", namespace,
		"exec",
		"-it",
		pod,
	}

	// Add container flag if specified (for multi-container pods)
	if container != "" {
		args = append(args, "-c", container)
	}

	// Add command
	args = append(args, "--", "sh", "-c", command)

	cmd := exec.Command(kubectlPath, args...)

	// Set stdin/stdout/stderr for interactive mode
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}

// validatePodName validates a pod name against allowed characters
func validatePodName(name string) error {
	if name == "" {
		return fmt.Errorf("pod name cannot be empty")
	}

	// Kubernetes pod names must be valid DNS subdomain names
	// RFC 1123: lowercase alphanumeric characters, '-' or '.', start/end with alphanumeric
	matched, err := regexp.MatchString(`^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$`, name)
	if err != nil {
		return fmt.Errorf("failed to validate pod name: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid pod name '%s': must be lowercase alphanumeric with optional hyphens or dots", name)
	}

	return nil
}

// validateContainerName validates a container name against allowed characters
func validateContainerName(name string) error {
	if name == "" {
		return nil // Empty is OK - kubectl will use default
	}

	// Container names follow same rules as pod names (RFC 1123)
	matched, err := regexp.MatchString(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, name)
	if err != nil {
		return fmt.Errorf("failed to validate container name: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid container name '%s': must be lowercase alphanumeric with optional hyphens", name)
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
