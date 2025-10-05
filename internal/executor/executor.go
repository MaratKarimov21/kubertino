package executor

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/k8s"
)

// KubeAdapter defines the interface for Kubernetes operations
type KubeAdapter interface {
	ExecInPod(context, namespace, pod, container, command string) (*exec.Cmd, error)
}

// Executor manages action execution
type Executor struct {
	kubeAdapter KubeAdapter
}

// NewExecutor creates a new Executor instance
func NewExecutor(kubeAdapter KubeAdapter) *Executor {
	return &Executor{
		kubeAdapter: kubeAdapter,
	}
}

// ExecutePodExec executes a pod_exec action
func (e *Executor) ExecutePodExec(action config.Action, context config.Context, namespace string, pods []k8s.Pod) error {
	// 1. Match pod using pattern
	pattern := getPodPattern(action, context)
	pod, err := matchPod(pods, pattern)
	if err != nil {
		return err
	}

	// 2. Render context box
	fmt.Println(renderContextBox(context.Name, namespace, pod.Name, action.Name))

	// 3. Build kubectl exec command
	cmd, err := e.kubeAdapter.ExecInPod(context.Name, namespace, pod.Name, action.Container, action.Command)
	if err != nil {
		return fmt.Errorf("failed to build exec command: %w", err)
	}

	// 4. Execute command with full terminal control
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kubectl exec failed: %w", err)
	}

	return nil
}

// PreparePodExec prepares a pod_exec command without executing it (for TUI integration)
// Returns the command ready to be executed by tea.ExecProcess
func (e *Executor) PreparePodExec(action config.Action, context config.Context, namespace string, pods []k8s.Pod) (*exec.Cmd, error) {
	// 1. Match pod using pattern
	pattern := getPodPattern(action, context)
	pod, err := matchPod(pods, pattern)
	if err != nil {
		return nil, err
	}

	// 2. Render context box (shows before TUI suspends)
	fmt.Println(renderContextBox(context.Name, namespace, pod.Name, action.Name))

	// 3. Build kubectl exec command
	cmd, err := e.kubeAdapter.ExecInPod(context.Name, namespace, pod.Name, action.Container, action.Command)
	if err != nil {
		return nil, fmt.Errorf("failed to build exec command: %w", err)
	}

	return cmd, nil
}

// getPodPattern returns the pod pattern from action or context default
func getPodPattern(action config.Action, context config.Context) string {
	if action.PodPattern != "" {
		return action.PodPattern
	}
	return context.DefaultPodPattern
}

// matchPod matches a single pod by pattern (regex)
func matchPod(pods []k8s.Pod, pattern string) (*k8s.Pod, error) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pod pattern '%s': %w", pattern, err)
	}

	var matches []k8s.Pod
	for _, pod := range pods {
		if regex.MatchString(pod.Name) {
			matches = append(matches, pod)
		}
	}

	if len(matches) == 0 {
		if len(pods) == 0 {
			return nil, fmt.Errorf("no pods match pattern '%s' - no pods in namespace", pattern)
		}
		// List available pods for debugging
		names := make([]string, len(pods))
		for i, p := range pods {
			names[i] = p.Name
		}
		return nil, fmt.Errorf("no pods match pattern '%s' - available pods: %s", pattern, strings.Join(names, ", "))
	}

	if len(matches) > 1 {
		names := make([]string, len(matches))
		for i, p := range matches {
			names[i] = p.Name
		}
		return nil, fmt.Errorf("multiple pods match pattern '%s': %s", pattern, strings.Join(names, ", "))
	}

	return &matches[0], nil
}
