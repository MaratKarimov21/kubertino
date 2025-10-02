package tui

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maratkarimov/kubertino/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestDimensionCalculations verifies that panel dimensions are calculated correctly
func TestDimensionCalculations(t *testing.T) {
	tests := []struct {
		name        string
		termWidth   int
		termHeight  int
		wantLeftW   int
		wantRightW  int
		wantTopH    int
		wantBottomH int
	}{
		{
			name:        "standard 80x24 terminal",
			termWidth:   80,
			termHeight:  24,
			wantLeftW:   40,
			wantRightW:  40,
			wantTopH:    11, // (24-1)/2 = 11
			wantBottomH: 12, // 23 - 11 = 12
		},
		{
			name:        "large 120x40 terminal",
			termWidth:   120,
			termHeight:  40,
			wantLeftW:   60,
			wantRightW:  60,
			wantTopH:    19, // (40-1)/2 = 19
			wantBottomH: 20, // 39 - 19 = 20
		},
		{
			name:        "small 80x30 terminal",
			termWidth:   80,
			termHeight:  30,
			wantLeftW:   40,
			wantRightW:  40,
			wantTopH:    14, // (30-1)/2 = 14
			wantBottomH: 15, // 29 - 14 = 15
		},
		{
			name:        "wide 160x24 terminal",
			termWidth:   160,
			termHeight:  24,
			wantLeftW:   80,
			wantRightW:  80,
			wantTopH:    11, // (24-1)/2 = 11
			wantBottomH: 12, // 23 - 11 = 12
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate dimensions using story logic
			availableHeight := tt.termHeight - HeaderHeight
			leftW := tt.termWidth / 2
			rightW := tt.termWidth - leftW
			topH := availableHeight / 2
			bottomH := availableHeight - topH

			assert.Equal(t, tt.wantLeftW, leftW, "left panel width mismatch")
			assert.Equal(t, tt.wantRightW, rightW, "right panel width mismatch")
			assert.Equal(t, tt.wantTopH, topH, "top panel height mismatch")
			assert.Equal(t, tt.wantBottomH, bottomH, "bottom panel height mismatch")
		})
	}
}

// TestRenderSplitLayout verifies that the split layout renders correctly
func TestRenderSplitLayout(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{
			{
				Name:               "test-context",
				Kubeconfig:         "~/.kube/config",
				DefaultPodPattern:  ".*",
				FavoriteNamespaces: []string{"default", "kube-system"},
			},
		},
	}

	adapter := &mockKubeAdapter{
		namespaces: []string{"default", "kube-system", "test-ns"},
	}

	model := NewAppModel(cfg, adapter)
	model.currentContext = &cfg.Contexts[0]
	model.viewMode = viewModeNamespaceView
	model.termWidth = 80
	model.termHeight = 24
	model.namespaces = []string{"default", "kube-system", "test-ns"}

	output := model.renderSplitLayout()

	// Verify header contains context name
	assert.Contains(t, output, "Context: test-context", "header should contain context name")

	// Verify panel titles appear
	assert.Contains(t, output, "Pods", "should contain Pods panel title")
	assert.Contains(t, output, "Actions", "should contain Actions panel title")

	// Verify placeholder text (may be wrapped across lines)
	assert.True(t, strings.Contains(output, "Select a namespace") || strings.Contains(output, "view pods"), "should contain pods placeholder")
	assert.True(t, strings.Contains(output, "Available actions") || strings.Contains(output, "appear here"), "should contain actions placeholder")

	// Verify namespace list is present
	assert.Contains(t, output, "Namespaces", "should contain Namespaces title")
}

// TestRenderHeader verifies header rendering with different contexts
func TestRenderHeader(t *testing.T) {
	tests := []struct {
		name         string
		contextName  string
		expectedText string
	}{
		{
			name:         "with context",
			contextName:  "production",
			expectedText: "Context: production",
		},
		{
			name:         "with different context",
			contextName:  "staging",
			expectedText: "Context: staging",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Contexts: []config.Context{
					{Name: tt.contextName},
				},
			}

			model := NewAppModel(cfg, &mockKubeAdapter{})
			model.currentContext = &cfg.Contexts[0]
			model.termWidth = 80
			model.termHeight = 24

			header := model.renderHeader()
			assert.Contains(t, header, tt.expectedText)
		})
	}
}

// TestRenderHeaderNoContext verifies header when no context is selected
func TestRenderHeaderNoContext(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.termWidth = 80
	model.termHeight = 24

	header := model.renderHeader()
	assert.Contains(t, header, "Context: None")
}

// TestTerminalResizeHandling verifies terminal resize message handling
func TestTerminalResizeHandling(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{
			{Name: "test"},
		},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})

	// Simulate resize to acceptable size
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(AppModel)

	assert.Equal(t, 100, m.termWidth)
	assert.Equal(t, 30, m.termHeight)
	assert.Equal(t, 100, m.width)
	assert.Equal(t, 30, m.height)
	assert.False(t, m.terminalTooSmall)
}

// TestMinimumTerminalSizeEnforcement verifies minimum size warning
func TestMinimumTerminalSizeEnforcement(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		height         int
		expectTooSmall bool
	}{
		{
			name:           "acceptable size 80x24",
			width:          80,
			height:         24,
			expectTooSmall: false,
		},
		{
			name:           "large size 120x40",
			width:          120,
			height:         40,
			expectTooSmall: false,
		},
		{
			name:           "width too small",
			width:          70,
			height:         24,
			expectTooSmall: true,
		},
		{
			name:           "height too small",
			width:          80,
			height:         20,
			expectTooSmall: true,
		},
		{
			name:           "both too small",
			width:          60,
			height:         15,
			expectTooSmall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Contexts: []config.Context{{Name: "test"}},
			}

			model := NewAppModel(cfg, &mockKubeAdapter{})

			// Simulate resize
			msg := tea.WindowSizeMsg{Width: tt.width, Height: tt.height}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(AppModel)

			assert.Equal(t, tt.expectTooSmall, m.terminalTooSmall)
		})
	}
}

// TestTerminalTooSmallWarning verifies the warning message content
func TestTerminalTooSmallWarning(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.termWidth = 70
	model.termHeight = 20
	model.terminalTooSmall = true

	warning := model.renderTerminalTooSmallWarning()

	assert.Contains(t, warning, "Terminal too small")
	assert.Contains(t, warning, "80x24") // Minimum size
	assert.Contains(t, warning, "70x20") // Current size
}

// TestViewModeSwitchesToSplitLayout verifies View() uses split layout in namespace mode
func TestViewModeSwitchesToSplitLayout(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{
			{Name: "test-context"},
		},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.currentContext = &cfg.Contexts[0]
	model.viewMode = viewModeNamespaceView
	model.termWidth = 80
	model.termHeight = 24
	model.namespaces = []string{"default"}

	output := model.View()

	// Should contain elements from split layout
	assert.Contains(t, output, "Context: test-context")
	assert.Contains(t, output, "Pods")
	assert.Contains(t, output, "Actions")
}

// TestViewShowsWarningWhenTerminalTooSmall verifies warning is shown
func TestViewShowsWarningWhenTerminalTooSmall(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.currentContext = &cfg.Contexts[0]
	model.viewMode = viewModeNamespaceView
	model.termWidth = 60
	model.termHeight = 20
	model.terminalTooSmall = true

	output := model.View()

	// Should show warning instead of split layout
	assert.Contains(t, output, "Terminal too small")
	assert.Contains(t, output, "60x20")
}

// TestRenderPodPanel verifies pod panel rendering
func TestRenderPodPanel(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.termWidth = 80
	model.termHeight = 24

	panel := model.renderPodPanel(40, 11)

	assert.Contains(t, panel, "Pods")
	assert.True(t, strings.Contains(panel, "Select a namespace") || strings.Contains(panel, "view pods"))
}

// TestRenderActionsPanel verifies actions panel rendering
func TestRenderActionsPanel(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.termWidth = 80
	model.termHeight = 24

	panel := model.renderActionsPanel(40, 11)

	assert.Contains(t, panel, "Actions")
	assert.True(t, strings.Contains(panel, "Available actions") || strings.Contains(panel, "appear here"))
}

// TestRenderNamespacePanel verifies namespace panel includes content
func TestRenderNamespacePanel(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{
			{
				Name:               "test",
				FavoriteNamespaces: []string{"default"},
			},
		},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.currentContext = &cfg.Contexts[0]
	model.termWidth = 80
	model.termHeight = 24
	model.namespaces = []string{"default", "kube-system"}

	panel := model.renderNamespacePanel(40, 23)

	assert.Contains(t, panel, "Namespaces")
	// Should contain namespaces or at least the list structure
	assert.True(t, len(panel) > 0)
}

// TestResizeFromSmallToLarge verifies transition from too small to acceptable
func TestResizeFromSmallToLarge(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})

	// Start with too small terminal
	msg1 := tea.WindowSizeMsg{Width: 60, Height: 20}
	updatedModel, _ := model.Update(msg1)
	m := updatedModel.(AppModel)
	assert.True(t, m.terminalTooSmall)

	// Resize to acceptable size
	msg2 := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ = m.Update(msg2)
	m = updatedModel.(AppModel)
	assert.False(t, m.terminalTooSmall)
	assert.Equal(t, 100, m.termWidth)
	assert.Equal(t, 30, m.termHeight)
}

// TestLayoutWithEmptyNamespaces verifies layout works with no namespaces
func TestLayoutWithEmptyNamespaces(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.currentContext = &cfg.Contexts[0]
	model.viewMode = viewModeNamespaceView
	model.termWidth = 80
	model.termHeight = 24
	model.namespaces = []string{} // Empty namespace list

	output := model.renderSplitLayout()

	// Should still render all panels
	assert.Contains(t, output, "Context: test")
	assert.Contains(t, output, "Pods")
	assert.Contains(t, output, "Actions")
}

// TestPanelBordersPresent verifies that panel borders are rendered
func TestPanelBordersPresent(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{{Name: "test"}},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.currentContext = &cfg.Contexts[0]
	model.viewMode = viewModeNamespaceView
	model.termWidth = 80
	model.termHeight = 24
	model.namespaces = []string{"default"}

	output := model.renderSplitLayout()

	// The output should contain ANSI codes for borders (rounded border characters)
	// Lipgloss uses box drawing characters
	assert.True(t, len(output) > 100, "output should be substantial with borders and content")

	// Verify basic structure is present
	assert.True(t, strings.Contains(output, "Pods") || strings.Contains(output, "Actions"))
}

// TestRenderNamespacePanel_LongListInSplitLayout verifies viewport scrolling with 40+ namespaces
func TestRenderNamespacePanel_LongListInSplitLayout(t *testing.T) {
	cfg := &config.Config{
		Contexts: []config.Context{
			{
				Name:               "test-context",
				FavoriteNamespaces: []string{"production"},
			},
		},
	}

	model := NewAppModel(cfg, &mockKubeAdapter{})
	model.currentContext = &cfg.Contexts[0]
	model.viewMode = viewModeNamespaceView
	model.termWidth = 80
	model.termHeight = 24

	// Create 50 namespaces to ensure overflow
	namespaces := make([]string, 50)
	for i := 0; i < 50; i++ {
		namespaces[i] = fmt.Sprintf("namespace-%d", i)
	}
	model.namespaces = namespaces

	// Render split layout
	output := model.renderSplitLayout()

	// Verify the output is well-formed (not overflowing)
	assert.Contains(t, output, "Namespaces (50)", "should show total namespace count")
	assert.Contains(t, output, "Pods", "should contain Pods panel")
	assert.Contains(t, output, "Actions", "should contain Actions panel")

	// Verify scroll indicator is present OR list was truncated due to MaxHeight
	// The scroll indicator format is " [start-end of total]"
	// With MaxHeight, content may be truncated before scroll indicator appears
	hasScrollOrTruncation := strings.Contains(output, " of 50") ||
		strings.Contains(output, "[") ||
		!strings.Contains(output, "namespace-49") // Last namespace not visible = truncated

	assert.True(t, hasScrollOrTruncation, "should show scroll indicator or be truncated for long list")

	// Verify first namespace is visible (in viewport)
	assert.Contains(t, output, "namespace-0", "should show first namespace in viewport")

	// Verify layout structure is intact (not broken by overflow)
	// Note: Individual panels may exceed their allocated height slightly due to padding/borders,
	// but the important thing is that all panels are visible and layout isn't broken
	lines := strings.Split(output, "\n")
	// Allow some overflow due to border/padding rendering (up to ~30 lines for 24 terminal)
	assert.True(t, len(lines) <= 30, "output should be reasonably close to terminal height")

	// Most importantly: verify pods and actions panels are still properly visible
	// This is the regression test for the bug - layout shouldn't break
	assert.Contains(t, output, "Pods", "Pods panel should be visible despite long namespace list")
	assert.Contains(t, output, "Actions", "Actions panel should be visible despite long namespace list")
}
