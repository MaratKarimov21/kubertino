package config

import (
	"fmt"
	"regexp"
	"strings"
)

// Validate validates the configuration and returns an error if invalid
func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if len(cfg.Contexts) == 0 {
		return fmt.Errorf("no contexts defined")
	}

	for i, ctx := range cfg.Contexts {
		if err := validateContext(&ctx, i); err != nil {
			return err
		}
	}

	return nil
}

// validateContext validates a single context
func validateContext(ctx *Context, index int) error {
	// Validate required fields
	if ctx.Name == "" {
		return fmt.Errorf("context[%d]: name is required", index)
	}

	if ctx.DefaultPodPattern == "" {
		return fmt.Errorf("context[%d] (%s): default_pod_pattern is required", index, ctx.Name)
	}

	// Validate regex patterns
	if _, err := regexp.Compile(ctx.DefaultPodPattern); err != nil {
		return fmt.Errorf("context[%d] (%s): invalid default_pod_pattern regex: %w", index, ctx.Name, err)
	}

	// Validate actions
	shortcuts := make(map[string]string)
	for j, action := range ctx.Actions {
		if err := validateAction(&action, j, ctx.Name); err != nil {
			return err
		}

		// Check for duplicate shortcuts
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

	if action.Type == "" {
		return fmt.Errorf("context (%s), action[%d] (%s): type is required", contextName, index, action.Name)
	}

	// Validate action type
	validTypes := map[string]bool{
		"pod_exec": true,
		"url":      true,
		"local":    true,
	}
	if !validTypes[action.Type] {
		return fmt.Errorf("context (%s), action[%d] (%s): invalid type '%s', must be one of: pod_exec, url, local",
			contextName, index, action.Name, action.Type)
	}

	// Validate type-specific fields
	switch action.Type {
	case "pod_exec":
		if action.Command == "" {
			return fmt.Errorf("context (%s), action[%d] (%s): pod_exec action requires command field",
				contextName, index, action.Name)
		}
		// Check for command injection characters
		if strings.ContainsAny(action.Command, ";|&") {
			return fmt.Errorf("context (%s), action[%d] (%s): command contains unsafe characters (;|&)",
				contextName, index, action.Name)
		}
	case "url":
		if action.URL == "" {
			return fmt.Errorf("context (%s), action[%d] (%s): url action requires url field",
				contextName, index, action.Name)
		}
	}

	// Validate pod pattern if provided
	if action.PodPattern != "" {
		if _, err := regexp.Compile(action.PodPattern); err != nil {
			return fmt.Errorf("context (%s), action[%d] (%s): invalid pod_pattern regex: %w",
				contextName, index, action.Name, err)
		}
	}

	return nil
}
