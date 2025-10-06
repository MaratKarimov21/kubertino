package config

import "regexp"

// Config represents the complete user configuration from ~/.kubertino.yml
type Config struct {
	Version    string      `yaml:"version"`
	Kubeconfig string      `yaml:"kubeconfig,omitempty"` // Optional kubeconfig path override
	Actions    []Action    `yaml:"actions,omitempty"`    // Global actions for all contexts
	Favorites  interface{} `yaml:"favorites,omitempty"`  // map[string][]string OR []string
	Contexts   []Context   `yaml:"contexts"`
}

// Context represents a Kubernetes context with its settings
type Context struct {
	Name              string   `yaml:"name"`
	DefaultPodPattern string   `yaml:"default_pod_pattern,omitempty"`
	Actions           []Action `yaml:"actions,omitempty"` // Per-context actions (extend/override global)

	// Runtime fields (not in YAML)
	CompiledPattern *regexp.Regexp `yaml:"-"` // Compiled default_pod_pattern for performance
}

// Action represents a configurable action with a shortcut
type Action struct {
	Name        string `yaml:"name"`
	Shortcut    string `yaml:"shortcut"`
	Command     string `yaml:"command"`               // Template with {{context}}, {{namespace}}, {{pod}}
	PodPattern  string `yaml:"pod_pattern,omitempty"` // Optional pod regex override
	Container   string `yaml:"container,omitempty"`   // Optional container name for multi-container pods
	Destructive bool   `yaml:"destructive,omitempty"` // Requires confirmation (optional)
}
