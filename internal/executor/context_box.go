package executor

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	colorCyan = "39"  // Bright cyan for box border
	colorGray = "240" // Gray for help hint
)

// renderContextBox creates a decorated box showing execution context
func renderContextBox(context, namespace, pod, action, command string) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colorCyan)).
		Padding(1, 2).
		Bold(true)

	// Dim style for the help hint
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray))

	// Build main content
	mainContent := fmt.Sprintf(
		"Context:   %s\nNamespace: %s\nPod:       %s\nAction:    %s\nCommand:   %s",
		context, namespace, pod, action, command,
	)

	// Add help hint at the bottom
	helpHint := dimStyle.Render("Press Ctrl+D to return")
	content := mainContent + "\n\n" + helpHint

	return boxStyle.Render(content)
}

// buildWaitPrompt creates a styled prompt message for wait-on-exit functionality
func buildWaitPrompt() string {
	dimStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray))

	return dimStyle.Render("Command finished. Press Ctrl+D to return to Kubertino.")
}

// buildCompoundCommand creates a shell command that displays the context box before executing the actual command.
// This ensures the box is displayed AFTER the TUI releases terminal control.
// The context box is properly escaped for shell safety.
// If waitOnExit is true, the command will wait for Ctrl+D (EOF) before returning.
func buildCompoundCommand(contextBox, command string, waitOnExit bool) string {
	// Escape single quotes for shell safety using POSIX shell escaping:
	// Replace each ' with '\'' (end quote, escaped quote, start quote)
	escapedBox := strings.ReplaceAll(contextBox, "'", "'\\''")

	// Build base compound command: print box, then run actual command
	compoundCommand := fmt.Sprintf("printf '%%s\\n\\n' '%s' && %s", escapedBox, command)

	// If wait-on-exit is enabled, append wait logic
	if waitOnExit {
		waitPrompt := buildWaitPrompt()
		escapedPrompt := strings.ReplaceAll(waitPrompt, "'", "'\\''")
		compoundCommand = fmt.Sprintf("%s && { printf '%%s\\n' '%s' && cat > /dev/null; }", compoundCommand, escapedPrompt)
	}

	return compoundCommand
}
