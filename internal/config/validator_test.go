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
							{Name: "console", Shortcut: "c", Type: "pod_exec", Command: "/bin/sh"},
						},
					},
				},
			},
			wantErr: false,
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
			name: "missing default pod pattern",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{Name: "test"},
				},
			},
			wantErr:     true,
			errContains: "default_pod_pattern is required",
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
			name: "duplicate shortcuts",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "c", Type: "pod_exec", Command: "/bin/sh"},
							{Name: "bash", Shortcut: "c", Type: "pod_exec", Command: "/bin/bash"},
						},
					},
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
							{Shortcut: "c", Type: "pod_exec", Command: "/bin/sh"},
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
							{Name: "console", Type: "pod_exec", Command: "/bin/sh"},
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
							{Name: "console", Shortcut: "con", Type: "pod_exec", Command: "/bin/sh"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "shortcut must be single character",
		},
		{
			name: "missing action type",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "c", Command: "/bin/sh"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "type is required",
		},
		{
			name: "invalid action type",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "test", Shortcut: "t", Type: "invalid_type", Command: "test"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "invalid type",
		},
		{
			name: "pod_exec missing command",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "console", Shortcut: "c", Type: "pod_exec"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "pod_exec action requires command field",
		},
		{
			name: "url action missing url",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "grafana", Shortcut: "g", Type: "url"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "url action requires url field",
		},
		{
			name: "command injection - semicolon",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "malicious", Shortcut: "m", Type: "pod_exec", Command: "/bin/sh; rm -rf /"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "command contains unsafe characters",
		},
		{
			name: "command injection - pipe",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "malicious", Shortcut: "m", Type: "pod_exec", Command: "cat /etc/passwd | grep root"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "command contains unsafe characters",
		},
		{
			name: "command injection - ampersand",
			config: &Config{
				Version: "1.0",
				Contexts: []Context{
					{
						Name:              "test",
						DefaultPodPattern: ".*",
						Actions: []Action{
							{Name: "malicious", Shortcut: "m", Type: "pod_exec", Command: "sleep 10 & rm -rf /"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "command contains unsafe characters",
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
							{Name: "console", Shortcut: "c", Type: "pod_exec", Command: "/bin/sh", PodPattern: "[invalid("},
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

func TestValidate_WithFixtures(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid config",
			filename: "../testdata/valid-config.yml",
			wantErr:  false,
		},
		{
			name:        "missing name",
			filename:    "../testdata/invalid-missing-name.yml",
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name:        "duplicate shortcuts",
			filename:    "../testdata/invalid-duplicate-shortcuts.yml",
			wantErr:     true,
			errContains: "duplicate shortcut",
		},
		{
			name:        "invalid regex",
			filename:    "../testdata/invalid-regex.yml",
			wantErr:     true,
			errContains: "invalid default_pod_pattern regex",
		},
		{
			name:        "invalid action type",
			filename:    "../testdata/invalid-action-type.yml",
			wantErr:     true,
			errContains: "invalid type",
		},
		{
			name:        "command injection",
			filename:    "../testdata/invalid-command-injection.yml",
			wantErr:     true,
			errContains: "command contains unsafe characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Parse(tt.filename)
			require.NoError(t, err)

			err = Validate(cfg)

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
