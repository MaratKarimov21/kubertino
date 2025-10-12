package config

import (
	"bytes"
	"fmt"
	"text/template"
)

// Validate validates the configuration and returns an error if invalid
func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if cfg.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(cfg.Contexts) == 0 {
		return fmt.Errorf("no contexts defined")
	}

	// Validate favorites format if present
	if cfg.Favorites != nil {
		if err := validateFavorites(cfg.Favorites); err != nil {
			return fmt.Errorf("invalid favorites: %w", err)
		}
	}

	// Validate global actions
	if len(cfg.Actions) > 0 {
		globalShortcuts := make(map[string]string)
		for i, action := range cfg.Actions {
			if err := validateAction(&action, i, "global"); err != nil {
				return err
			}
			if existingAction, exists := globalShortcuts[action.Shortcut]; exists {
				return fmt.Errorf("global actions: duplicate shortcut '%s' found in actions '%s' and '%s'",
					action.Shortcut, existingAction, action.Name)
			}
			globalShortcuts[action.Shortcut] = action.Name
		}
	}

	for i, ctx := range cfg.Contexts {
		if err := validateContext(&ctx, i, cfg.Actions); err != nil {
			return err
		}
	}

	return nil
}

// validateContext validates a single context
func validateContext(ctx *Context, index int, globalActions []Action) error {
	// Validate required fields
	if ctx.Name == "" {
		return fmt.Errorf("context[%d]: name is required", index)
	}

	// Story 6.2: default_pod_pattern removed - no validation needed

	// Validate per-context actions
	shortcuts := make(map[string]string)
	for j, action := range ctx.Actions {
		if err := validateAction(&action, j, ctx.Name); err != nil {
			return err
		}
		// Check for duplicate shortcuts within per-context actions
		if existingAction, exists := shortcuts[action.Shortcut]; exists {
			return fmt.Errorf("context[%d] (%s): duplicate shortcut '%s' found in actions '%s' and '%s'",
				index, ctx.Name, action.Shortcut, existingAction, action.Name)
		}
		shortcuts[action.Shortcut] = action.Name
	}

	return nil
}

// validateAction validates a single action
func validateAction(action *Action, index int, contextName string) error {
	// Validate required fields
	if action.Name == "" {
		return fmt.Errorf("context (%s), action[%d]: name is required", contextName, index)
	}

	if action.Shortcut == "" {
		return fmt.Errorf("context (%s), action[%d] (%s): shortcut is required", contextName, index, action.Name)
	}

	if len(action.Shortcut) > 1 {
		return fmt.Errorf("context (%s), action[%d] (%s): shortcut must be single character, got '%s'",
			contextName, index, action.Name, action.Shortcut)
	}

	if action.Command == "" {
		return fmt.Errorf("context (%s), action[%d] (%s): command is required", contextName, index, action.Name)
	}

	// Validate command template syntax
	if err := validateCommandTemplate(action.Command); err != nil {
		return fmt.Errorf("context (%s), action[%d] (%s): invalid command template: %w",
			contextName, index, action.Name, err)
	}

	// Story 6.2: pod_pattern removed - no validation needed

	return nil
}

// validateCommandTemplate validates the Go template syntax in a command string
func validateCommandTemplate(command string) error {
	// Create a template with dummy data to validate syntax
	tmpl, err := template.New("command").Parse(command)
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	// Try to execute the template with dummy variables to catch undefined function errors
	// This validates that {{variable}} syntax works correctly
	data := map[string]string{
		"context":   "test-context",
		"namespace": "test-namespace",
		"pod":       "test-pod",
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	return nil
}

// validateFavorites validates the favorites format
func validateFavorites(favorites interface{}) error {
	if favorites == nil {
		return nil
	}

	// Try to parse favorites using the parser function
	_, err := parseFavorites(favorites)
	return err
}
