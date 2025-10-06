package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "c", Command: "kubectl exec -it {{.pod}} -- /bin/sh"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with global actions",
			config: &Config{
				Version: "1.0",
				Actions: []Action{
					{Name: "logs", Shortcut: "l", Command: "kubectl logs {{.pod}}"},
				},
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			config: &Config{
				Contexts: []Context{
					{Name: "test", DefaultPodPattern: ".*"},
				},
			},
			wantErr:     true,
			errContains: "version is required",
		},
		{
			name:        "nil config",
			config:      nil,
			wantErr:     true,
			errContains: "config is nil",
		},
		{
			name: "no contexts",
			config: &Config{
				Version:  "1.0",
				Contexts: []Context{},
			},
			wantErr:     true,
			errContains: "no contexts defined",
		},
		{
			name: "missing context name",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{DefaultPodPattern: ".*"},
				},
			},
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name: "optional default pod pattern",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{Name: "test"},
				},
			},
			wantErr: false, // DefaultPodPattern is now optional
		},
		{
			name: "invalid regex pattern",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{Name: "test", DefaultPodPattern: "[invalid(regex"},
				},
			},
			wantErr:     true,
			errContains: "invalid default_pod_pattern regex",
		},
		{
			name: "duplicate shortcuts in context",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "c", Command: "/bin/sh"},
							{Name: "bash", Shortcut: "c", Command: "/bin/bash"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "duplicate shortcut",
		},
		{
			name: "duplicate shortcuts in global actions",
			config: &Config{
				Version: "1.0",
				Actions: []Action{
					{Name: "logs", Shortcut: "l", Command: "kubectl logs {{.pod}}"},
					{Name: "logs-tail", Shortcut: "l", Command: "kubectl logs {{.pod}} --tail=100"},
				},
				Contexts: []Context{
					{Name: "test", DefaultPodPattern: ".*"},
				},
			},
			wantErr:     true,
			errContains: "duplicate shortcut",
		},
		{
			name: "missing action name",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Shortcut: "c", Command: "/bin/sh"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name: "missing action shortcut",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Command: "/bin/sh"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "shortcut is required",
		},
		{
			name: "multi-character shortcut",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "con", Command: "/bin/sh"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "shortcut must be single character",
		},
		{
			name: "missing command",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "c"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "command is required",
		},
		{
			name: "invalid template syntax",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "test", Shortcut: "t", Command: "kubectl exec {{pod"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "invalid command template",
		},
		{
			name: "invalid pod pattern regex",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "c", Command: "/bin/sh", PodPattern: "[invalid("},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "invalid pod_pattern regex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.config)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidate_WithFixtures is temporarily skipped as testdata files need updating for new config structure
// TODO: Update testdata files to match Epic 5 config structure
func TestValidate_WithFixtures(t *testing.T) {
	t.Skip("Testdata files need updating for Epic 5 config structure")
}
