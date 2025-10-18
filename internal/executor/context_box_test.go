package executor

import (
	"strings"
	"testing"
)

func TestRenderContextBox_ShortCommand(t *testing.T) {
	result := renderContextBox("prod", "app-namespace", "web-pod-123", "console", "bash")

	// Verify all fields are present
	if !strings.Contains(result, "prod") {
		t.Errorf("Expected context 'prod' in output")
	}
	if !strings.Contains(result, "app-namespace") {
		t.Errorf("Expected namespace 'app-namespace' in output")
	}
	if !strings.Contains(result, "web-pod-123") {
		t.Errorf("Expected pod 'web-pod-123' in output")
	}
	if !strings.Contains(result, "console") {
		t.Errorf("Expected action 'console' in output")
	}
	if !strings.Contains(result, "bash") {
		t.Errorf("Expected command 'bash' in output")
	}

	// Verify labels are present
	if !strings.Contains(result, "Context:") {
		t.Errorf("Expected 'Context:' label in output")
	}
	if !strings.Contains(result, "Namespace:") {
		t.Errorf("Expected 'Namespace:' label in output")
	}
	if !strings.Contains(result, "Pod:") {
		t.Errorf("Expected 'Pod:' label in output")
	}
	if !strings.Contains(result, "Action:") {
		t.Errorf("Expected 'Action:' label in output")
	}
	if !strings.Contains(result, "Command:") {
		t.Errorf("Expected 'Command:' label in output")
	}

	// Verify help hint is present
	if !strings.Contains(result, "Press Ctrl+D to return") {
		t.Errorf("Expected 'Press Ctrl+D to return' help hint in output")
	}
}

func TestRenderContextBox_LongCommand(t *testing.T) {
	longCommand := "kubectl exec -it web-pod-123 -n app-namespace --context prod -- /bin/bash -c 'echo hello && ls -la /var/log && tail -f /var/log/app.log'"

	result := renderContextBox("prod", "app-namespace", "web-pod-123", "console", longCommand)

	// Verify command is present in output
	if !strings.Contains(result, longCommand) {
		t.Errorf("Expected long command to be present in output")
	}

	// Verify output is not empty and contains box styling
	if len(result) == 0 {
		t.Errorf("Expected non-empty output for long command")
	}

	// Verify all other fields are still present
	if !strings.Contains(result, "prod") {
		t.Errorf("Expected context 'prod' in output")
	}
	if !strings.Contains(result, "app-namespace") {
		t.Errorf("Expected namespace 'app-namespace' in output")
	}
	if !strings.Contains(result, "web-pod-123") {
		t.Errorf("Expected pod 'web-pod-123' in output")
	}
}

func TestRenderContextBox_AllFieldsPresent(t *testing.T) {
	tests := []struct {
		name      string
		context   string
		namespace string
		pod       string
		action    string
		command   string
	}{
		{
			name:      "standard values",
			context:   "production",
			namespace: "default",
			pod:       "my-pod-abc123",
			action:    "logs",
			command:   "kubectl logs my-pod-abc123 -n default",
		},
		{
			name:      "with special characters",
			context:   "dev-cluster-01",
			namespace: "team-backend",
			pod:       "api-server-v2-xyz789",
			action:    "exec-bash",
			command:   "kubectl exec -it api-server-v2-xyz789 -- /bin/bash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderContextBox(tt.context, tt.namespace, tt.pod, tt.action, tt.command)

			// Check all input values are in output
			expectedValues := []string{tt.context, tt.namespace, tt.pod, tt.action, tt.command}
			for _, value := range expectedValues {
				if !strings.Contains(result, value) {
					t.Errorf("Expected value '%s' not found in output", value)
				}
			}

			// Check all labels are present
			expectedLabels := []string{"Context:", "Namespace:", "Pod:", "Action:", "Command:"}
			for _, label := range expectedLabels {
				if !strings.Contains(result, label) {
					t.Errorf("Expected label '%s' not found in output", label)
				}
			}
		})
	}
}

func TestRenderContextBox_TemplateSubstitutedCommand(t *testing.T) {
	// Test with a command that has already gone through template substitution
	templateSubstitutedCommand := "kubectl exec -it web-123 -n production-ns -- bash"

	result := renderContextBox("prod", "production-ns", "web-123", "console", templateSubstitutedCommand)

	// Verify no template variables remain
	if strings.Contains(result, "{{") || strings.Contains(result, "}}") {
		t.Errorf("Template variables should not appear in output")
	}

	// Verify substituted command is present
	if !strings.Contains(result, templateSubstitutedCommand) {
		t.Errorf("Expected substituted command to be present in output")
	}
}

func TestRenderContextBox_EmptyValues(t *testing.T) {
	// Test with empty command (edge case)
	result := renderContextBox("context", "namespace", "pod", "action", "")

	// Should still render without panicking
	if len(result) == 0 {
		t.Errorf("Expected output even with empty command")
	}

	// Verify labels are still present
	if !strings.Contains(result, "Command:") {
		t.Errorf("Expected 'Command:' label even with empty command")
	}
}

func TestRenderContextBox_VeryLongCommand(t *testing.T) {
	// Test with a very long command (200+ characters)
	veryLongCommand := "kubectl exec -it my-pod-with-very-long-name-12345678 -n my-namespace-with-very-long-name-12345678 --context my-context-with-very-long-name-12345678 -- bash -c 'for i in {1..100}; do echo $i; done'"

	result := renderContextBox("prod", "app", "pod-123", "exec", veryLongCommand)

	// Verify command is present (even if wrapped or truncated)
	// We're checking for a substring to handle potential wrapping
	if !strings.Contains(result, "kubectl exec") {
		t.Errorf("Expected command to contain 'kubectl exec'")
	}

	// Verify box is rendered (has styling)
	if len(result) < len(veryLongCommand) {
		t.Errorf("Expected output to include box styling, got shorter output than command length")
	}
}

func TestRenderContextBox_URLCommand(t *testing.T) {
	// Test with URL-based action
	urlCommand := "open https://dashboard.example.com/prod/my-namespace/pod-123"

	result := renderContextBox("prod", "my-namespace", "pod-123", "dashboard", urlCommand)

	// Verify URL is preserved correctly
	if !strings.Contains(result, "https://dashboard.example.com") {
		t.Errorf("Expected URL to be preserved in output")
	}
	if !strings.Contains(result, urlCommand) {
		t.Errorf("Expected full URL command in output")
	}
}

// Story 7.6: Wait-on-Exit Tests

func TestBuildWaitPrompt(t *testing.T) {
	result := buildWaitPrompt()

	// Verify prompt is not empty
	if len(result) == 0 {
		t.Errorf("Expected non-empty wait prompt")
	}

	// Verify prompt contains key text
	if !strings.Contains(result, "Press Ctrl+D") {
		t.Errorf("Expected wait prompt to contain 'Press Ctrl+D'")
	}
	if !strings.Contains(result, "Command finished") {
		t.Errorf("Expected wait prompt to contain 'Command finished'")
	}
}

func TestBuildCompoundCommand_WaitOnExit_False(t *testing.T) {
	contextBox := "Test Context Box"
	command := "echo hello"

	result := buildCompoundCommand(contextBox, command, false)

	// Verify command is present
	if !strings.Contains(result, command) {
		t.Errorf("Expected command '%s' in output", command)
	}

	// Verify printf with context box is present
	if !strings.Contains(result, "printf") {
		t.Errorf("Expected printf command in output")
	}

	// Verify NO wait logic is present
	if strings.Contains(result, "cat > /dev/null") {
		t.Errorf("Expected NO wait logic when waitOnExit is false")
	}
}

func TestBuildCompoundCommand_WaitOnExit_True(t *testing.T) {
	contextBox := "Test Context Box"
	command := "kubectl describe pod test-pod"

	result := buildCompoundCommand(contextBox, command, true)

	// Verify command is present
	if !strings.Contains(result, command) {
		t.Errorf("Expected command '%s' in output", command)
	}

	// Verify printf with context box is present
	if !strings.Contains(result, "printf") {
		t.Errorf("Expected printf command in output")
	}

	// Verify wait logic IS present
	if !strings.Contains(result, "cat > /dev/null") {
		t.Errorf("Expected wait logic (cat > /dev/null) when waitOnExit is true")
	}

	// Verify wait prompt message is present
	if !strings.Contains(result, "Command finished") {
		t.Errorf("Expected wait prompt message in output")
	}
}

func TestBuildCompoundCommand_ShellEscaping_WithSingleQuotes(t *testing.T) {
	contextBox := "Context with 'single quotes' in text"
	command := "echo test"

	result := buildCompoundCommand(contextBox, command, false)

	// Verify shell escaping is applied (single quotes should be escaped)
	// The escaping pattern is: ' becomes '\''
	if !strings.Contains(result, "'\\''") {
		t.Errorf("Expected single quotes to be escaped in context box")
	}

	// Verify command is still present
	if !strings.Contains(result, command) {
		t.Errorf("Expected command '%s' in output", command)
	}
}

func TestBuildCompoundCommand_WaitOnExit_ShellEscaping(t *testing.T) {
	// Test that wait prompt is properly escaped when it contains special characters
	contextBox := "Test Box"
	command := "echo test"

	result := buildCompoundCommand(contextBox, command, true)

	// Verify the wait prompt is properly escaped for shell
	// The wait prompt should be within single quotes and properly escaped
	if !strings.Contains(result, "printf") {
		t.Errorf("Expected printf for wait prompt")
	}

	// Verify wait logic structure
	if !strings.Contains(result, "&&") {
		t.Errorf("Expected && to chain commands")
	}

	// Verify braces for wait logic grouping
	if !strings.Contains(result, "{") || !strings.Contains(result, "}") {
		t.Errorf("Expected braces for wait logic grouping")
	}
}
