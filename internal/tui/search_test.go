package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/stretchr/testify/assert"
)

const testSearchQuery = "kube"

func TestSearchMode_Activation(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production"}

	// Press '/' to activate search
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	updatedModel, _ := model.Update(msg)
	model = updatedModel.(AppModel)

	assert.True(t, model.searchMode, "search mode should be activated")
	assert.Equal(t, "", model.searchQuery, "search query should be empty initially")
	assert.Equal(t, model.namespaces, model.filteredNamespaces, "filtered namespaces should be initialized")
	assert.Equal(t, 0, model.selectedNamespaceIndex, "selected index should be reset")
}

func TestSearchMode_Deactivation(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production"}

	// Activate search
	model.searchMode = true
	model.searchQuery = testSearchQuery
	model.filteredNamespaces = []string{"kube-system"}

	// Press ESC to deactivate search
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)
	model = updatedModel.(AppModel)

	assert.False(t, model.searchMode, "search mode should be deactivated")
	assert.Equal(t, "", model.searchQuery, "search query should be cleared")
	assert.Nil(t, model.filteredNamespaces, "filtered namespaces should be cleared")
}

func TestSearchMode_CharacterInput(t *testing.T) {
	tests := []struct {
		name          string
		initialQuery  string
		inputChar     rune
		expectedQuery string
		shouldFilter  bool
	}{
		{
			name:          "single character",
			initialQuery:  "",
			inputChar:     'k',
			expectedQuery: "k",
			shouldFilter:  true,
		},
		{
			name:          "multiple characters",
			initialQuery:  "kub",
			inputChar:     'e',
			expectedQuery: testSearchQuery,
			shouldFilter:  true,
		},
		{
			name:          "dash character",
			initialQuery:  testSearchQuery,
			inputChar:     '-',
			expectedQuery: "kube-",
			shouldFilter:  true,
		},
		{
			name:          "dot character",
			initialQuery:  "app",
			inputChar:     '.',
			expectedQuery: "app.",
			shouldFilter:  true,
		},
		{
			name:          "underscore character",
			initialQuery:  "my",
			inputChar:     '_',
			expectedQuery: "my_",
			shouldFilter:  true,
		},
		{
			name:          "number character",
			initialQuery:  "v",
			inputChar:     '1',
			expectedQuery: "v1",
			shouldFilter:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Version:  "1.0",
				Contexts: []config.Context{},
			}
			model := NewAppModel(cfg, newMockAdapter())
			model.viewMode = viewModeNamespaceView
			model.namespaces = []string{"kube-system", "default", "production"}
			model.searchMode = true
			model.searchQuery = tt.initialQuery
			model.filteredNamespaces = model.namespaces

			// Type character
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.inputChar}}
			updatedModel, _ := model.Update(msg)
			model = updatedModel.(AppModel)

			assert.Equal(t, tt.expectedQuery, model.searchQuery, "search query should be updated")
			assert.NotNil(t, model.filteredNamespaces, "filtered namespaces should be updated")
		})
	}
}

func TestSearchMode_Backspace(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production"}
	model.searchMode = true
	model.searchQuery = testSearchQuery
	model.filteredNamespaces = []string{"kube-system"}

	// Press backspace
	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	updatedModel, _ := model.Update(msg)
	model = updatedModel.(AppModel)

	assert.Equal(t, "kub", model.searchQuery, "last character should be removed")
	assert.NotNil(t, model.filteredNamespaces, "filtered namespaces should be updated")
}

func TestSearchMode_BackspaceOnEmpty(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production"}
	model.searchMode = true
	model.searchQuery = ""
	model.filteredNamespaces = model.namespaces

	// Press backspace on empty query
	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	updatedModel, _ := model.Update(msg)
	model = updatedModel.(AppModel)

	assert.Equal(t, "", model.searchQuery, "query should remain empty")
}

func TestSearchMode_EnterSelectsNamespace(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production"}
	model.searchMode = true
	model.searchQuery = testSearchQuery
	model.filteredNamespaces = []string{"kube-system"}
	model.selectedNamespaceIndex = 0

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := model.Update(msg)
	model = updatedModel.(AppModel)

	// Should deactivate search (actual namespace selection is Story 2.4+)
	assert.False(t, model.searchMode, "search mode should be deactivated after Enter")
}

func TestSearchMode_Navigation(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "kube-public", "default"}
	model.searchMode = true
	model.searchQuery = testSearchQuery
	model.filteredNamespaces = []string{"kube-system", "kube-public"}
	model.selectedNamespaceIndex = 0

	// Navigate down in filtered list
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, _ := model.Update(msg)
	model = updatedModel.(AppModel)

	assert.Equal(t, 1, model.selectedNamespaceIndex, "should navigate in filtered list")

	// Navigate up in filtered list
	msg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ = model.Update(msg)
	model = updatedModel.(AppModel)

	assert.Equal(t, 0, model.selectedNamespaceIndex, "should navigate back to first item")
}

func TestSearchMode_EmptyQuery(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production"}
	model.searchMode = true
	model.searchQuery = "test"
	model.filteredNamespaces = []string{}

	// Clear query by updating to empty
	model.updateSearchQuery("")

	assert.Equal(t, "", model.searchQuery, "query should be empty")
	assert.Equal(t, model.namespaces, model.filteredNamespaces, "should show all namespaces")
}

func TestSearchMode_FuzzyMatching(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production", "kubertino-app"}

	// Test fuzzy matching with testSearchQuery
	model.updateSearchQuery(testSearchQuery)

	assert.Greater(t, len(model.filteredNamespaces), 0, "should have matches")
	assert.Contains(t, model.filteredNamespaces, "kube-system", "should match kube-system")
	assert.Contains(t, model.filteredNamespaces, "kubertino-app", "should match kubertino-app")
}

func TestSearchMode_NoMatches(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.namespaces = []string{"kube-system", "default", "production"}

	// Search for something that doesn't match
	model.updateSearchQuery("xyz")

	assert.Equal(t, 0, len(model.filteredNamespaces), "should have no matches")
}

func TestRenderNamespaceList_SearchMode(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.currentContext = &cfg.Contexts[0]
	model.namespaces = []string{"kube-system", "default"}
	model.height = 20
	model.searchMode = true
	model.searchQuery = testSearchQuery
	model.filteredNamespaces = []string{"kube-system"}

	output := model.renderNamespaceList(0)

	// Story 7.2: Search box now has yellow border and separate label
	assert.Contains(t, output, "Search", "should show Search label")
	assert.Contains(t, output, "kube_", "should show search query with cursor")
	assert.Contains(t, output, "kube-system", "should show filtered namespace")
	assert.NotContains(t, output, "default", "should not show unmatched namespace")
	assert.Contains(t, output, "ESC: Cancel", "should show ESC hint")
}

func TestRenderNamespaceList_SearchModeNoMatches(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.currentContext = &cfg.Contexts[0]
	model.namespaces = []string{"kube-system", "default"}
	model.height = 20
	model.searchMode = true
	model.searchQuery = "xyz"
	model.filteredNamespaces = []string{}

	output := model.renderNamespaceList(0)

	assert.Contains(t, output, "No matches found", "should show no matches message")
	// Story 7.2: Search box now has yellow border and separate label
	assert.Contains(t, output, "Search", "should show Search label")
	assert.Contains(t, output, "xyz_", "should show search query with cursor")
}

func TestRenderNamespaceList_NormalMode(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.viewMode = viewModeNamespaceView
	model.currentContext = &cfg.Contexts[0]
	model.namespaces = []string{"kube-system", "default"}
	model.height = 20
	model.searchMode = false

	output := model.renderNamespaceList(0)

	assert.NotContains(t, output, "Search:", "should not show search input box")
	assert.Contains(t, output, "/: Search", "should show / key hint")
	assert.Contains(t, output, "kube-system", "should show all namespaces")
	assert.Contains(t, output, "default", "should show all namespaces")
}

func TestGetMatchIndices(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.namespaces = []string{"kube-system"}
	model.searchMode = true
	model.searchQuery = "ks"

	indices := model.getMatchIndices("kube-system")

	assert.NotNil(t, indices, "should return match indices")
	assert.Greater(t, len(indices), 0, "should have at least one match index")
}

func TestGetMatchIndices_NotInSearchMode(t *testing.T) {
	cfg := &config.Config{
		Version: "1.0",
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}
	model := NewAppModel(cfg, newMockAdapter())
	model.namespaces = []string{"kube-system"}
	model.searchMode = false

	indices := model.getMatchIndices("kube-system")

	assert.Nil(t, indices, "should return nil when not in search mode")
}

// Story 7.6: Test cursor position preservation when exiting search mode
func TestSearchMode_DeactivationPreservesCursor(t *testing.T) {
	tests := []struct {
		name                    string
		namespaces              []string
		filteredNamespaces      []string
		selectedIndexInFiltered int
		expectedIndexAfterExit  int
	}{
		{
			name:                    "Cursor preserved when exiting search with matched namespace",
			namespaces:              []string{"default", "kube-system", "production", "kubertino-app"},
			filteredNamespaces:      []string{"kube-system", "kubertino-app"},
			selectedIndexInFiltered: 1, // "kubertino-app" selected in filtered list
			expectedIndexAfterExit:  3, // Should map to index 3 in full list
		},
		{
			name:                    "Cursor preserved at first match",
			namespaces:              []string{"default", "kube-system", "production"},
			filteredNamespaces:      []string{"kube-system"},
			selectedIndexInFiltered: 0, // "kube-system" selected
			expectedIndexAfterExit:  1, // Should map to index 1 in full list
		},
		{
			name:                    "Cursor clamped when filtered index out of bounds",
			namespaces:              []string{"default", "production"},
			filteredNamespaces:      []string{"default"},
			selectedIndexInFiltered: 10, // Invalid index
			expectedIndexAfterExit:  0,  // Should clamp to 0
		},
		{
			name:                    "Empty filtered list resets to 0",
			namespaces:              []string{"default", "production"},
			filteredNamespaces:      []string{},
			selectedIndexInFiltered: 0,
			expectedIndexAfterExit:  0, // Should reset to 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Version: "1.0",
				Contexts: []config.Context{
					{Name: "test-context"},
				},
			}
			model := NewAppModel(cfg, newMockAdapter())
			model.viewMode = viewModeNamespaceView
			model.namespaces = tt.namespaces
			model.searchMode = true
			model.searchQuery = "kube"
			model.filteredNamespaces = tt.filteredNamespaces
			model.selectedNamespaceIndex = tt.selectedIndexInFiltered
			model.termHeight = 40 // Set height for viewport adjustment

			// Press ESC to exit search mode
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			updatedModel, _ := model.Update(msg)
			model = updatedModel.(AppModel)

			assert.False(t, model.searchMode, "search mode should be deactivated")
			assert.Equal(t, tt.expectedIndexAfterExit, model.selectedNamespaceIndex,
				"cursor should be at correct position after exiting search")
		})
	}
}
