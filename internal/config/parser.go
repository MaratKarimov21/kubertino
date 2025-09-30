package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Parse loads and parses a YAML configuration file
func Parse(filename string) (*Config, error) {
	// Expand tilde in path
	if len(filename) >= 2 && filename[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		filename = filepath.Join(homeDir, filename[2:])
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}
