package tui

import tea "github.com/charmbracelet/bubbletea"

// KeyMap defines keyboard bindings for the TUI
type KeyMap struct {
	Quit  []string // Keys that trigger quit (q, esc, ctrl+c)
	Up    []string // Keys for navigating up (up arrow, k)
	Down  []string // Keys for navigating down (down arrow, j)
	Enter []string // Keys for selection (enter)
}

// DefaultKeyMap returns the default keyboard bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit:  []string{"q", "esc", "ctrl+c"},
		Up:    []string{"up", "k"},
		Down:  []string{"down", "j"},
		Enter: []string{"enter"},
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
