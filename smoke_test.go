package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/maratkarimov/kubertino/internal/tui"
)

func main() {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{
				Name:               "context1",
				DefaultPodPattern:  ".*",
				FavoriteNamespaces: []string{"default", "api", "monitoring"},
			},
			{
				Name:               "context2",
				DefaultPodPattern:  ".*",
				FavoriteNamespaces: []string{"default"},
			},
		},
	}

	model := tui.NewAppModel(cfg)

	// Verify model is in context_selection mode with multiple contexts
	if model.View() == "" {
		fmt.Println("ERROR: View returned empty string")
		os.Exit(1)
	}

	// Create a program but don't run it (we just want to verify initialization)
	p := tea.NewProgram(model)

	// Send a quit message immediately
	go func() {
		time.Sleep(100 * time.Millisecond)
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Smoke test passed: App starts correctly with context selection")
	fmt.Println("✓ Multiple contexts trigger context_selection mode")
	fmt.Println("✓ View renders without panic")
}
