package tui

import tea "github.com/charmbracelet/bubbletea"

// KeyMap defines keyboard bindings for the TUI
type KeyMap struct {
	Quit     []string // Keys that trigger quit (q, esc, ctrl+c)
	Up       []string // Keys for navigating up (up arrow, k)
	Down     []string // Keys for navigating down (down arrow, j)
	Enter    []string // Keys for selection (enter)
	Tab      []string // Keys for switching focus forward (tab) - Story 3.3, 4.1: Namespaces → Pods → Actions
	ShiftTab []string // Keys for switching focus backward (shift+tab) - Story 3.3, 4.1: Actions → Pods → Namespaces
}

// DefaultKeyMap returns the default keyboard bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit:     []string{"q", "esc", "ctrl+c"},
		Up:       []string{"up", "k"},
		Down:     []string{"down", "j"},
		Enter:    []string{"enter"},
		Tab:      []string{"tab"},       // Story 3.3, 4.1: Three-panel focus switching
		ShiftTab: []string{"shift+tab"}, // Story 3.3, 4.1: Three-panel backward focus switching
	}
}

// KeyMatches checks if a key message matches any of the provided key bindings
func KeyMatches(msg tea.KeyMsg, keys []string) bool {
	keyStr := msg.String()
	for _, k := range keys {
		if keyStr == k {
			return true
		}
	}
	return false
}
