package config

// Config represents the complete user configuration from ~/.kubertino.yml
type Config struct {
	Version  string    `yaml:"version"`
	Contexts []Context `yaml:"contexts"`
}

// Context represents a Kubernetes context with its settings
type Context struct {
	Name               string   `yaml:"name"`
	Kubeconfig         string   `yaml:"kubeconfig"`
	ClusterURL         string   `yaml:"cluster_url,omitempty"`
	DefaultPodPattern  string   `yaml:"default_pod_pattern"`
	FavoriteNamespaces []string `yaml:"favorite_namespaces"`
	Actions            []Action `yaml:"actions"`
}

// Action represents a configurable action with a shortcut
type Action struct {
	Name       string `yaml:"name"`
	Shortcut   string `yaml:"shortcut"`
	Type       string `yaml:"type"` // pod_exec, url, local
	Command    string `yaml:"command,omitempty"`
	URL        string `yaml:"url,omitempty"`
	PodPattern string `yaml:"pod_pattern,omitempty"`
}
