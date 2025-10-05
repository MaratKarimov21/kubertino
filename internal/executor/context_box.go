package executor

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// renderContextBox creates a decorated box showing execution context
func renderContextBox(context, namespace, pod, action string) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")). // Bright cyan
		Padding(1, 2).
		Bold(true)

	content := fmt.Sprintf(
		"Context:   %s\nNamespace: %s\nPod:       %s\nAction:    %s",
		context, namespace, pod, action,
	)

	return boxStyle.Render(content)
}
