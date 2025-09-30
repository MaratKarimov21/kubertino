package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppModel(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}

	model := NewAppModel(cfg)

	assert.NotNil(t, model.config, "config should be set")
	assert.Equal(t, cfg, model.config, "config should match")
	assert.Nil(t, model.currentContext, "currentContext should be nil initially")
	assert.Nil(t, model.err, "err should be nil initially")
	assert.Equal(t, 0, model.width, "width should be 0 initially")
	assert.Equal(t, 0, model.height, "height should be 0 initially")
	assert.NotNil(t, model.keys, "keys should be initialized")
}

func TestAppModel_Init(t *testing.T) {
	model := NewAppModel(&config.Config{})
	cmd := model.Init()

	assert.Nil(t, cmd, "Init should return nil command for Story 1.4")
}

func TestAppModel_Update_KeyHandling(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		keyType    tea.KeyType
		shouldQuit bool
	}{
		{
			name:       "quit with q",
			key:        "q",
			keyType:    tea.KeyRunes,
			shouldQuit: true,
		},
		{
			name:       "quit with esc",
			key:        "esc",
			keyType:    tea.KeyEsc,
			shouldQuit: true,
		},
		{
			name:       "quit with ctrl+c",
			key:        "ctrl+c",
			keyType:    tea.KeyCtrlC,
			shouldQuit: true,
		},
		{
			name:       "no quit with enter",
			key:        "enter",
			keyType:    tea.KeyEnter,
			shouldQuit: false,
		},
		{
			name:       "no quit with tab",
			key:        "tab",
			keyType:    tea.KeyTab,
			shouldQuit: false,
		},
		{
			name:       "no quit with up arrow",
			key:        "up",
			keyType:    tea.KeyUp,
			shouldQuit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewAppModel(&config.Config{})

			var keyMsg tea.KeyMsg
			switch tt.keyType {
			case tea.KeyRunes:
				keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			case tea.KeyEsc:
				keyMsg = tea.KeyMsg{Type: tea.KeyEsc}
			case tea.KeyCtrlC:
				keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
			case tea.KeyEnter:
				keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
			case tea.KeyTab:
				keyMsg = tea.KeyMsg{Type: tea.KeyTab}
			case tea.KeyUp:
				keyMsg = tea.KeyMsg{Type: tea.KeyUp}
			}

			newModel, cmd := model.Update(keyMsg)

			if tt.shouldQuit {
				assert.NotNil(t, cmd, "Expected quit command for key: %s", tt.key)
				// Verify it's actually a quit command by checking if it returns nil msg
				if cmd != nil {
					msg := cmd()
					assert.IsType(t, tea.QuitMsg{}, msg, "Command should return QuitMsg")
				}
			} else {
				// For non-quit keys, we expect nil command or the model unchanged
				assert.NotNil(t, newModel, "Model should not be nil")
			}
		})
	}
}

func TestAppModel_Update_WindowSize(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small terminal", 80, 24},
		{"medium terminal", 120, 40},
		{"large terminal", 200, 60},
		{"zero dimensions", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewAppModel(&config.Config{})

			msg := tea.WindowSizeMsg{
				Width:  tt.width,
				Height: tt.height,
			}

			newModel, cmd := model.Update(msg)

			require.NotNil(t, newModel)
			m := newModel.(AppModel)

			assert.Equal(t, tt.width, m.width, "width should be updated")
			assert.Equal(t, tt.height, m.height, "height should be updated")
			assert.Nil(t, cmd, "WindowSizeMsg should not return a command")
		})
	}
}

func TestAppModel_View_WithError(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected []string // Substrings that should be present
	}{
		{
			name:   "config error",
			errMsg: "config file not found",
			expected: []string{
				"Error:",
				"config file not found",
				"Press 'q' or ESC to quit",
			},
		},
		{
			name:   "validation error",
			errMsg: "context validation failed",
			expected: []string{
				"Error:",
				"context validation failed",
				"quit",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewAppModel(&config.Config{})
			model.err = assert.AnError // Use a generic error

			// Override with specific error message for testing
			model.err = &testError{msg: tt.errMsg}

			view := model.View()

			assert.NotEmpty(t, view, "View should not be empty")

			for _, expected := range tt.expected {
				assert.Contains(t, view, expected, "View should contain: %s", expected)
			}
		})
	}
}

func TestAppModel_View_NoError(t *testing.T) {
	model := NewAppModel(&config.Config{})
	model.width = 100
	model.height = 40

	view := model.View()

	assert.NotEmpty(t, view, "View should not be empty when no error")
	// Basic layout should be rendered
	assert.NotContains(t, view, "Error:", "View should not show error when err is nil")
}

func TestAppModel_View_ResponsiveLayout(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small", 80, 24},
		{"medium", 120, 40},
		{"large", 200, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewAppModel(&config.Config{})
			model.width = tt.width
			model.height = tt.height

			view := model.View()

			assert.NotEmpty(t, view, "View should not be empty")
			// Verify quit instructions are present
			assert.True(t,
				strings.Contains(view, "quit") ||
					strings.Contains(view, "Quit") ||
					strings.Contains(view, "q"),
				"View should contain quit instructions",
			)
		})
	}
}

func TestKeyMatches(t *testing.T) {
	tests := []struct {
		name     string
		keyMsg   tea.KeyMsg
		keys     []string
		expected bool
	}{
		{
			name:     "matches q",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")},
			keys:     []string{"q", "esc"},
			expected: true,
		},
		{
			name:     "matches esc",
			keyMsg:   tea.KeyMsg{Type: tea.KeyEsc},
			keys:     []string{"q", "esc"},
			expected: true,
		},
		{
			name:     "matches ctrl+c",
			keyMsg:   tea.KeyMsg{Type: tea.KeyCtrlC},
			keys:     []string{"ctrl+c"},
			expected: true,
		},
		{
			name:     "no match",
			keyMsg:   tea.KeyMsg{Type: tea.KeyEnter},
			keys:     []string{"q", "esc"},
			expected: false,
		},
		{
			name:     "empty keys list",
			keyMsg:   tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")},
			keys:     []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KeyMatches(tt.keyMsg, tt.keys)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultKeyMap(t *testing.T) {
	keyMap := DefaultKeyMap()

	assert.NotNil(t, keyMap.Quit, "Quit keys should be defined")
	assert.Contains(t, keyMap.Quit, "q", "Should contain 'q' key")
	assert.Contains(t, keyMap.Quit, "esc", "Should contain 'esc' key")
	assert.Contains(t, keyMap.Quit, "ctrl+c", "Should contain 'ctrl+c' key")
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
