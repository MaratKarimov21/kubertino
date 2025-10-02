package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockKubeAdapter is a mock implementation of KubeAdapter for testing
type mockKubeAdapter struct {
	namespaces []string
	err        error
}

func (m *mockKubeAdapter) GetNamespaces(context string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.namespaces, nil
}

func newMockAdapter() *mockKubeAdapter {
	return &mockKubeAdapter{
		namespaces: []string{"default", "kube-system", "production", "staging"},
	}
}

func TestNewAppModel(t *testing.T) {
	t.Run("single context auto-selects", func(t *testing.T) {
		cfg := &config.Config{
			Version: "1.0",
			Contexts: []config.Context{
				{Name: "test-context"},
			},
		}

		model := NewAppModel(cfg, newMockAdapter())

		assert.NotNil(t, model.config, "config should be set")
		assert.Equal(t, cfg, model.config, "config should match")
		assert.NotNil(t, model.currentContext, "currentContext should be auto-selected")
		assert.Equal(t, "test-context", model.currentContext.Name, "should auto-select single context")
		assert.Equal(t, viewModeNamespaceView, model.viewMode, "should skip to namespace_view")
		assert.Nil(t, model.err, "err should be nil initially")
		assert.NotNil(t, model.keys, "keys should be initialized")
	})

	t.Run("multiple contexts show selection screen", func(t *testing.T) {
		cfg := &config.Config{
			Version: "1.0",
			Contexts: []config.Context{
				{Name: "context1"},
				{Name: "context2"},
				{Name: "context3"},
			},
		}

		model := NewAppModel(cfg, newMockAdapter())

		assert.NotNil(t, model.config, "config should be set")
		assert.Nil(t, model.currentContext, "currentContext should be nil")
		assert.Equal(t, viewModeContextSelection, model.viewMode, "should be in context_selection mode")
		assert.Equal(t, 0, model.selectedContextIndex, "should start at index 0")
		assert.Len(t, model.contexts, 3, "should have 3 contexts")
	})
}

func TestAppModel_Init(t *testing.T) {
	model := NewAppModel(&config.Config{}, newMockAdapter())
	cmd := model.Init()

	assert.Nil(t, cmd, "Init should return nil command when no context selected")
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
			model := NewAppModel(&config.Config{}, newMockAdapter())

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

func TestAppModel_ContextSelection_Navigation(t *testing.T) {
	tests := []struct {
		name          string
		initialIndex  int
		keyType       tea.KeyType
		expectedIndex int
	}{
		{
			name:          "down arrow increments index",
			initialIndex:  0,
			keyType:       tea.KeyDown,
			expectedIndex: 1,
		},
		{
			name:          "down arrow wraps around at end",
			initialIndex:  2,
			keyType:       tea.KeyDown,
			expectedIndex: 0,
		},
		{
			name:          "up arrow decrements index",
			initialIndex:  1,
			keyType:       tea.KeyUp,
			expectedIndex: 0,
		},
		{
			name:          "up arrow wraps around at start",
			initialIndex:  0,
			keyType:       tea.KeyUp,
			expectedIndex: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Contexts: []config.Context{
					{Name: "context1"},
					{Name: "context2"},
					{Name: "context3"},
				},
			}
			model := NewAppModel(cfg, newMockAdapter())
			model.selectedContextIndex = tt.initialIndex

			keyMsg := tea.KeyMsg{Type: tt.keyType}
			newModel, _ := model.Update(keyMsg)
			m := newModel.(AppModel)

			assert.Equal(t, tt.expectedIndex, m.selectedContextIndex, "index should be updated correctly")
			assert.Equal(t, viewModeContextSelection, m.viewMode, "should remain in context_selection mode")
		})
	}
}

func TestAppModel_ContextSelection_EnterKey(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{
			{Name: "context1"},
			{Name: "context2"},
			{Name: "context3"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.selectedContextIndex = 1

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, cmd := model.Update(keyMsg)
	m := newModel.(AppModel)

	assert.NotNil(t, m.currentContext, "currentContext should be set")
	assert.Equal(t, "context2", m.currentContext.Name, "should select context at index 1")
	assert.Equal(t, viewModeNamespaceView, m.viewMode, "should transition to namespace_view")
	assert.NotNil(t, cmd, "should return namespace fetch command")
	assert.True(t, m.namespacesLoading, "should set loading state")
}

func TestAppModel_ContextSelection_QuitKeys(t *testing.T) {
	tests := []struct {
		name    string
		keyType tea.KeyType
		keyStr  string
	}{
		{"quit with q", tea.KeyRunes, "q"},
		{"quit with esc", tea.KeyEsc, "esc"},
		{"quit with ctrl+c", tea.KeyCtrlC, "ctrl+c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Contexts: []config.Context{
					{Name: "context1"},
					{Name: "context2"},
				},
			}
			model := NewAppModel(cfg, newMockAdapter())

			var keyMsg tea.KeyMsg
			if tt.keyType == tea.KeyRunes {
				keyMsg = tea.KeyMsg{Type: tt.keyType, Runes: []rune(tt.keyStr)}
			} else {
				keyMsg = tea.KeyMsg{Type: tt.keyType}
			}

			_, cmd := model.Update(keyMsg)

			assert.NotNil(t, cmd, "should return quit command")
			msg := cmd()
			assert.IsType(t, tea.QuitMsg{}, msg, "should return QuitMsg")
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
			model := NewAppModel(&config.Config{}, newMockAdapter())

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
			model := NewAppModel(&config.Config{}, newMockAdapter())
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
	model := NewAppModel(&config.Config{}, newMockAdapter())
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
			model := NewAppModel(&config.Config{}, newMockAdapter())
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

func TestAppModel_View_ContextSelection(t *testing.T) {
	t.Run("shows context selection screen", func(t *testing.T) {
		cfg := &config.Config{
			Contexts: []config.Context{
				{Name: "production", FavoriteNamespaces: []string{"default", "api"}},
				{Name: "staging"},
				{Name: "development", FavoriteNamespaces: []string{"default"}},
			},
		}
		model := NewAppModel(cfg, newMockAdapter())
		model.selectedContextIndex = 0

		view := model.View()

		assert.Contains(t, view, "Select Kubernetes Context", "should show header")
		assert.Contains(t, view, "production", "should show context name")
		assert.Contains(t, view, "staging", "should show second context")
		assert.Contains(t, view, "development", "should show third context")
		assert.Contains(t, view, "2 namespaces", "should show namespace count for production")
		assert.Contains(t, view, "1 namespaces", "should show namespace count for development")
		assert.Contains(t, view, "↑/↓ Navigate", "should show navigation hint")
		assert.Contains(t, view, "Enter: Select", "should show selection hint")
	})

	t.Run("switches to namespace view after context selection", func(t *testing.T) {
		cfg := &config.Config{
			Contexts: []config.Context{
				{Name: "test-context"},
			},
		}
		model := NewAppModel(cfg, newMockAdapter())

		view := model.View()

		assert.NotContains(t, view, "Select Kubernetes Context", "should not show context selection")
		// Should show the basic layout instead (namespace_view mode)
	})
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

	assert.NotNil(t, keyMap.Up, "Up keys should be defined")
	assert.Contains(t, keyMap.Up, "up", "Should contain 'up' key")
	assert.Contains(t, keyMap.Up, "k", "Should contain 'k' key")

	assert.NotNil(t, keyMap.Down, "Down keys should be defined")
	assert.Contains(t, keyMap.Down, "down", "Should contain 'down' key")
	assert.Contains(t, keyMap.Down, "j", "Should contain 'j' key")

	assert.NotNil(t, keyMap.Enter, "Enter keys should be defined")
	assert.Contains(t, keyMap.Enter, "enter", "Should contain 'enter' key")
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestAppModel_NamespaceNavigation(t *testing.T) {
	tests := []struct {
		name          string
		initialIndex  int
		keyType       tea.KeyType
		expectedIndex int
		numNamespaces int
	}{
		{
			name:          "down arrow increments index",
			initialIndex:  0,
			keyType:       tea.KeyDown,
			expectedIndex: 1,
			numNamespaces: 4,
		},
		{
			name:          "down arrow wraps around at end",
			initialIndex:  3,
			keyType:       tea.KeyDown,
			expectedIndex: 0,
			numNamespaces: 4,
		},
		{
			name:          "up arrow decrements index",
			initialIndex:  2,
			keyType:       tea.KeyUp,
			expectedIndex: 1,
			numNamespaces: 4,
		},
		{
			name:          "up arrow wraps around at start",
			initialIndex:  0,
			keyType:       tea.KeyUp,
			expectedIndex: 3,
			numNamespaces: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Contexts: []config.Context{
					{Name: "test-context"},
				},
			}
			model := NewAppModel(cfg, newMockAdapter())
			model.viewMode = viewModeNamespaceView
			model.namespaces = []string{"default", "kube-system", "production", "staging"}
			model.selectedNamespaceIndex = tt.initialIndex

			keyMsg := tea.KeyMsg{Type: tt.keyType}
			newModel, _ := model.Update(keyMsg)
			m := newModel.(AppModel)

			assert.Equal(t, tt.expectedIndex, m.selectedNamespaceIndex, "index should be updated correctly")
			assert.Equal(t, viewModeNamespaceView, m.viewMode, "should remain in namespace_view mode")
		})
	}
}

func TestAppModel_NamespaceFetched(t *testing.T) {
	t.Run("successful namespace fetch", func(t *testing.T) {
		cfg := &config.Config{
			Contexts: []config.Context{
				{Name: "test-context"},
			},
		}
		model := NewAppModel(cfg, newMockAdapter())
		model.namespacesLoading = true

		msg := namespaceFetchedMsg{
			namespaces: []string{"default", "kube-system", "production"},
			err:        nil,
		}
		newModel, _ := model.Update(msg)
		m := newModel.(AppModel)

		assert.False(t, m.namespacesLoading, "loading should be false")
		assert.Nil(t, m.namespacesError, "error should be nil")
		assert.Len(t, m.namespaces, 3, "should have 3 namespaces")
		assert.Contains(t, m.namespaces, "default")
		assert.Contains(t, m.namespaces, "kube-system")
		assert.Contains(t, m.namespaces, "production")
	})

	t.Run("failed namespace fetch", func(t *testing.T) {
		cfg := &config.Config{
			Contexts: []config.Context{
				{Name: "test-context"},
			},
		}
		model := NewAppModel(cfg, newMockAdapter())
		model.namespacesLoading = true

		testErr := &testError{msg: "connection timeout"}
		msg := namespaceFetchedMsg{
			namespaces: nil,
			err:        testErr,
		}
		newModel, _ := model.Update(msg)
		m := newModel.(AppModel)

		assert.False(t, m.namespacesLoading, "loading should be false")
		assert.NotNil(t, m.namespacesError, "error should be set")
		assert.Equal(t, testErr, m.namespacesError, "error should match")
	})
}

func TestAppModel_SortNamespacesWithFavorites(t *testing.T) {
	model := AppModel{}

	t.Run("favorites sorted first", func(t *testing.T) {
		namespaces := []string{"default", "kube-system", "production", "staging", "dev"}
		favorites := []string{"production", "staging"}

		sorted := model.sortNamespacesWithFavorites(namespaces, favorites)

		assert.Len(t, sorted, 5)
		assert.Equal(t, "production", sorted[0], "first favorite should be first")
		assert.Equal(t, "staging", sorted[1], "second favorite should be second")
		assert.Contains(t, sorted[2:], "default", "non-favorites should follow")
		assert.Contains(t, sorted[2:], "kube-system", "non-favorites should follow")
		assert.Contains(t, sorted[2:], "dev", "non-favorites should follow")
	})

	t.Run("no favorites", func(t *testing.T) {
		namespaces := []string{"default", "kube-system"}
		favorites := []string{}

		sorted := model.sortNamespacesWithFavorites(namespaces, favorites)

		assert.Equal(t, namespaces, sorted, "order should remain unchanged")
	})

	t.Run("all favorites", func(t *testing.T) {
		namespaces := []string{"default", "kube-system"}
		favorites := []string{"default", "kube-system"}

		sorted := model.sortNamespacesWithFavorites(namespaces, favorites)

		assert.Len(t, sorted, 2)
		assert.Contains(t, sorted, "default")
		assert.Contains(t, sorted, "kube-system")
	})
}

func TestAppModel_RenderNamespaceList(t *testing.T) {
	t.Run("shows loading state", func(t *testing.T) {
		cfg := &config.Config{
			Contexts: []config.Context{
				{Name: "test-context"},
			},
		}
		model := NewAppModel(cfg, newMockAdapter())
		model.currentContext = &cfg.Contexts[0]
		model.viewMode = viewModeNamespaceView
		model.namespacesLoading = true

		view := model.View()

		assert.Contains(t, view, "Loading namespaces", "should show loading message")
		assert.Contains(t, view, "Namespaces", "should show header")
	})

	t.Run("shows error state", func(t *testing.T) {
		cfg := &config.Config{
			Contexts: []config.Context{
				{Name: "test-context"},
			},
		}
		model := NewAppModel(cfg, newMockAdapter())
		model.currentContext = &cfg.Contexts[0]
		model.viewMode = viewModeNamespaceView
		model.namespacesError = &testError{msg: "connection failed"}

		view := model.View()

		assert.Contains(t, view, "Error fetching namespaces", "should show error message")
		assert.Contains(t, view, "connection failed", "should show error details")
	})

	t.Run("shows namespace list", func(t *testing.T) {
		cfg := &config.Config{
			Contexts: []config.Context{
				{Name: "test-context", FavoriteNamespaces: []string{"production"}},
			},
		}
		model := NewAppModel(cfg, newMockAdapter())
		model.currentContext = &cfg.Contexts[0]
		model.viewMode = viewModeNamespaceView
		model.termWidth = 80
		model.termHeight = 24
		model.namespaces = []string{"production", "default", "kube-system"}

		view := model.View()

		assert.Contains(t, view, "Namespaces (3)", "should show namespace count")
		assert.Contains(t, view, "production", "should show namespace")
		assert.Contains(t, view, "default", "should show namespace")
		assert.Contains(t, view, "kube-system", "should show namespace")
		assert.Contains(t, view, "★", "should show favorite indicator")
		assert.Contains(t, view, "Navigate", "should show navigation hint")
	})
}
