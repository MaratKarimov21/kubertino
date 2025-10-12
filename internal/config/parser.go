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

	// Story 6.2: First parse into a generic structure to detect deprecated fields
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Check for deprecated fields at root level and in contexts/actions
	if err := checkDeprecatedFields(rawConfig); err != nil {
		return nil, err
	}

	// Parse YAML into Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Parse favorites dual-format
	if config.Favorites != nil {
		if _, err := parseFavorites(config.Favorites); err != nil {
			return nil, fmt.Errorf("invalid favorites format: %w", err)
		}
	}

	// Story 6.2: Pattern matching removed - no compilation needed

	return &config, nil
}

// FavoritesFormat represents the format of favorites configuration
type FavoritesFormat string

const (
	FavoritesFormatPerContext FavoritesFormat = "per-context" // map[string][]string
	FavoritesFormatGlobal     FavoritesFormat = "global"      // []string
)

// ParsedFavorites represents parsed favorites with format information
type ParsedFavorites struct {
	Format     FavoritesFormat
	PerContext map[string][]string // For per-context format
	GlobalList []string            // For global format
}

// parseFavorites detects and parses the favorites format
func parseFavorites(favorites interface{}) (*ParsedFavorites, error) {
	if favorites == nil {
		return &ParsedFavorites{Format: FavoritesFormatGlobal, GlobalList: []string{}}, nil
	}

	switch v := favorites.(type) {
	case map[string]interface{}:
		// Per-context format: map[string][]string
		perContext := make(map[string][]string)
		for contextName, value := range v {
			// Value should be a list of strings
			list, ok := value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("per-context favorites: context '%s' value must be a list", contextName)
			}

			namespaces := make([]string, 0, len(list))
			for i, item := range list {
				ns, ok := item.(string)
				if !ok {
					return nil, fmt.Errorf("per-context favorites: context '%s' item %d is not a string", contextName, i)
				}
				namespaces = append(namespaces, ns)
			}
			perContext[contextName] = namespaces
		}
		return &ParsedFavorites{Format: FavoritesFormatPerContext, PerContext: perContext}, nil

	case []interface{}:
		// Global format: []string
		globalList := make([]string, 0, len(v))
		for i, item := range v {
			ns, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("global favorites: item %d is not a string", i)
			}
			globalList = append(globalList, ns)
		}
		return &ParsedFavorites{Format: FavoritesFormatGlobal, GlobalList: globalList}, nil

	default:
		return nil, fmt.Errorf("favorites must be either a map (per-context) or a list (global), got %T", favorites)
	}
}

// GetFavorites returns the favorite namespaces for a given context
func GetFavorites(config *Config, contextName string) ([]string, error) {
	if config.Favorites == nil {
		return []string{}, nil
	}

	parsed, err := parseFavorites(config.Favorites)
	if err != nil {
		return nil, err
	}

	switch parsed.Format {
	case FavoritesFormatPerContext:
		if namespaces, ok := parsed.PerContext[contextName]; ok {
			return namespaces, nil
		}
		return []string{}, nil
	case FavoritesFormatGlobal:
		return parsed.GlobalList, nil
	default:
		return nil, fmt.Errorf("unknown favorites format: %s", parsed.Format)
	}
}

// MergeActions combines global actions with per-context actions
// Per-context actions override global actions with the same shortcut
func MergeActions(globalActions, contextActions []Action) []Action {
	// Start with all global actions
	result := make([]Action, len(globalActions))
	copy(result, globalActions)

	// Then process per-context actions
	for _, contextAction := range contextActions {
		// Check if this shortcut exists in global actions
		overrideIndex := -1
		for i, action := range result {
			if action.Shortcut == contextAction.Shortcut {
				overrideIndex = i
				break
			}
		}

		if overrideIndex >= 0 {
			// Override the global action
			result[overrideIndex] = contextAction
		} else {
			// Add new action
			result = append(result, contextAction)
		}
	}

	return result
}

// checkDeprecatedFields detects and rejects deprecated fields from Story 6.2
func checkDeprecatedFields(rawConfig map[string]interface{}) error {
	deprecatedFields := map[string]string{
		"pod_pattern":         "Use manual pod selection workflow instead. All pods must be selected manually via Tab → Enter.",
		"default_pod_pattern": "Use manual pod selection workflow instead. Remove this field from context configuration.",
		"container":           "Specify container in command template instead: kubectl exec ... -c container-name",
		"type":                "All actions now execute as local commands. Remove 'type' field and use full kubectl commands in 'command' field.",
	}

	// Check contexts
	if contexts, ok := rawConfig["contexts"].([]interface{}); ok {
		for i, ctx := range contexts {
			contextMap, ok := ctx.(map[string]interface{})
			if !ok {
				continue
			}

			// Check deprecated fields in context
			for field, msg := range deprecatedFields {
				if _, exists := contextMap[field]; exists {
					return fmt.Errorf("context[%d]: deprecated field '%s' is no longer supported. %s", i, field, msg)
				}
			}

			// Check actions within context
			if actions, ok := contextMap["actions"].([]interface{}); ok {
				for j, act := range actions {
					actionMap, ok := act.(map[string]interface{})
					if !ok {
						continue
					}

					for field, msg := range deprecatedFields {
						if _, exists := actionMap[field]; exists {
							return fmt.Errorf("context[%d], action[%d]: deprecated field '%s' is no longer supported. %s", i, j, field, msg)
						}
					}
				}
			}
		}
	}

	// Check global actions
	if actions, ok := rawConfig["actions"].([]interface{}); ok {
		for i, act := range actions {
			actionMap, ok := act.(map[string]interface{})
			if !ok {
				continue
			}

			for field, msg := range deprecatedFields {
				if _, exists := actionMap[field]; exists {
					return fmt.Errorf("global action[%d]: deprecated field '%s' is no longer supported. %s", i, field, msg)
				}
			}
		}
	}

	return nil
}
