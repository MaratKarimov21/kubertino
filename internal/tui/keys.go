package tui

import tea "github.com/charmbracelet/bubbletea"

// KeyMap defines keyboard bindings for the TUI
type KeyMap struct {
	Quit []string // Keys that trigger quit (q, esc, ctrl+c)
}

// DefaultKeyMap returns the default keyboard bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: []string{"q", "esc", "ctrl+c"},
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
