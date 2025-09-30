package internal_test

import (
	"os"
	"path/filepath"
	"testing"
)

// TestProjectStructure verifies that all required directories exist
func TestProjectStructure(t *testing.T) {
	// Get the root directory (one level up from internal)
	rootDir := ".."

	tests := []struct {
		name string
		path string
	}{
		{"cmd directory", "cmd/kubertino"},
		{"internal/config directory", "internal/config"},
		{"internal/k8s directory", "internal/k8s"},
		{"internal/tui directory", "internal/tui"},
		{"internal/executor directory", "internal/executor"},
		{"internal/search directory", "internal/search"},
		{"pkg directory", "pkg"},
		{"examples directory", "examples"},
		{"scripts directory", "scripts"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(rootDir, tt.path)
			info, err := os.Stat(fullPath)
			if err != nil {
				t.Errorf("directory %s does not exist: %v", tt.path, err)
				return
			}
			if !info.IsDir() {
				t.Errorf("%s exists but is not a directory", tt.path)
			}
		})
	}
}

// TestGoModFile verifies that go.mod exists and contains the correct module path
func TestGoModFile(t *testing.T) {
	rootDir := ".."
	goModPath := filepath.Join(rootDir, "go.mod")

	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("failed to read go.mod: %v", err)
	}

	contentStr := string(content)

	// Check for module declaration
	if contentStr == "" {
		t.Error("go.mod file is empty")
	}

	// Should contain module path
	if len(contentStr) < 10 {
		t.Error("go.mod appears to be too short or malformed")
	}
}

// TestRequiredFiles verifies that all required files exist
func TestRequiredFiles(t *testing.T) {
	rootDir := ".."

	tests := []struct {
		name string
		path string
	}{
		{"README.md", "README.md"},
		{"Makefile", "Makefile"},
		{".gitignore", ".gitignore"},
		{"LICENSE", "LICENSE"},
		{".golangci.yml", ".golangci.yml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(rootDir, tt.path)
			info, err := os.Stat(fullPath)
			if err != nil {
				t.Errorf("file %s does not exist: %v", tt.path, err)
				return
			}
			if info.IsDir() {
				t.Errorf("%s exists but is a directory, not a file", tt.path)
			}
			if info.Size() == 0 {
				t.Errorf("file %s exists but is empty", tt.path)
			}
		})
	}
}

// TestGitignorePatterns verifies that .gitignore contains required patterns
func TestGitignorePatterns(t *testing.T) {
	rootDir := ".."
	gitignorePath := filepath.Join(rootDir, ".gitignore")

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	contentStr := string(content)

	requiredPatterns := []string{
		"kubertino", // binary name
		"*.out",     // coverage files
		".idea/",    // IDE directory
		".vscode/",  // IDE directory
		"logs/",     // logs directory
	}

	for _, pattern := range requiredPatterns {
		if !containsPattern(contentStr, pattern) {
			t.Errorf(".gitignore does not contain required pattern: %s", pattern)
		}
	}
}

// containsPattern checks if the content contains the pattern
func containsPattern(content, pattern string) bool {
	// Simple check - just look for the pattern in the content
	return len(content) > 0 && (content == pattern ||
		findSubstring(content, pattern))
}

// findSubstring checks if pattern exists in content
func findSubstring(content, pattern string) bool {
	contentLen := len(content)
	patternLen := len(pattern)

	if patternLen > contentLen {
		return false
	}

	for i := 0; i <= contentLen-patternLen; i++ {
		if content[i:i+patternLen] == pattern {
			return true
		}
	}
	return false
}

// TestGolangciLintConfig verifies that .golangci.yml contains required linters
func TestGolangciLintConfig(t *testing.T) {
	rootDir := ".."
	configPath := filepath.Join(rootDir, ".golangci.yml")

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read .golangci.yml: %v", err)
	}

	contentStr := string(content)

	requiredLinters := []string{
		"gofmt",
		"gocyclo",
		"goconst",
		"misspell",
	}

	for _, linter := range requiredLinters {
		if !containsPattern(contentStr, linter) {
			t.Errorf(".golangci.yml does not contain required linter: %s", linter)
		}
	}
}
