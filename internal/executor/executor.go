package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
)

// Executor manages action execution
type Executor struct {
}

// NewExecutor creates a new Executor instance
func NewExecutor() *Executor {
	return &Executor{}
}

// ExecuteLocal executes a local command action with template variable substitution
func (e *Executor) ExecuteLocal(action config.Action, context config.Context, namespace string, pod k8s.Pod, kubeconfigPath string) error {
	// 1. Parse command template
	tmpl, err := template.New("action").Parse(action.Command)
	if err != nil {
		return fmt.Errorf("invalid command template: %w", err)
	}

	// 2. Substitute template variables
	data := map[string]string{
		"context":   context.Name,
		"namespace": namespace,
		"pod":       pod.Name,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	command := buf.String()

	// 3. Build context box and compound command
	contextBox := renderContextBox(context.Name, namespace, pod.Name, action.Name, command)
	compoundCommand := buildCompoundCommand(contextBox, command)

	// 4. Execute compound command with shell
	cmd := exec.Command("sh", "-c", compoundCommand)
	cmd.Env = os.Environ() // Preserve parent environment

	// Expand kubeconfig path (tilde expansion)
	if kubeconfigPath != "" {
		expandedPath := kubeconfigPath
		if strings.HasPrefix(kubeconfigPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				expandedPath = filepath.Join(homeDir, kubeconfigPath[2:])
			}
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("KUBECONFIG=%s", expandedPath))
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

// PrepareLocal prepares a local command without executing it (for TUI integration)
// Returns the command ready to be executed by tea.ExecProcess
func (e *Executor) PrepareLocal(action config.Action, context config.Context, namespace string, pod k8s.Pod, kubeconfigPath string) (*exec.Cmd, error) {
	// 1. Parse command template
	tmpl, err := template.New("action").Parse(action.Command)
	if err != nil {
		return nil, fmt.Errorf("invalid command template: %w", err)
	}

	// 2. Substitute template variables
	data := map[string]string{
		"context":   context.Name,
		"namespace": namespace,
		"pod":       pod.Name,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("template execution failed: %w", err)
	}

	command := buf.String()

	// 3. Build context box and compound command
	contextBox := renderContextBox(context.Name, namespace, pod.Name, action.Name, command)
	compoundCommand := buildCompoundCommand(contextBox, command)

	// 4. Build command with shell
	cmd := exec.Command("sh", "-c", compoundCommand)
	cmd.Env = os.Environ() // Preserve parent environment

	// Expand kubeconfig path (tilde expansion)
	if kubeconfigPath != "" {
		expandedPath := kubeconfigPath
		if strings.HasPrefix(kubeconfigPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				expandedPath = filepath.Join(homeDir, kubeconfigPath[2:])
			}
		}
		cmd.Env = append(cmd.Env, fmt.Sprintf("KUBECONFIG=%s", expandedPath))
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}
