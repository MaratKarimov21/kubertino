package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewErrorModal(t *testing.T) {
	modal := NewErrorModal()

	assert.NotNil(t, modal)
	assert.False(t, modal.IsVisible)
	assert.Empty(t, modal.Message)
	assert.Empty(t, modal.Operation)
	assert.Empty(t, modal.Suggestion)
	assert.Nil(t, modal.RetryFunc)
}

func TestErrorModal_Show(t *testing.T) {
	modal := NewErrorModal()
	retryFunc := func() tea.Cmd { return nil }

	modal.Show("test error", "Test Operation", retryFunc)

	assert.True(t, modal.IsVisible)
	assert.Equal(t, "test error", modal.Message)
	assert.Equal(t, "Test Operation", modal.Operation)
	assert.Empty(t, modal.Suggestion)
	assert.NotNil(t, modal.RetryFunc)
}

func TestErrorModal_ShowWithSuggestion(t *testing.T) {
	modal := NewErrorModal()
	retryFunc := func() tea.Cmd { return nil }

	modal.ShowWithSuggestion("test error", "Test Operation", "Try this fix", retryFunc)

	assert.True(t, modal.IsVisible)
	assert.Equal(t, "test error", modal.Message)
	assert.Equal(t, "Test Operation", modal.Operation)
	assert.Equal(t, "Try this fix", modal.Suggestion)
	assert.NotNil(t, modal.RetryFunc)
}

func TestErrorModal_Hide(t *testing.T) {
	modal := NewErrorModal()
	retryFunc := func() tea.Cmd { return nil }

	// Show modal first
	modal.Show("test error", "Test Operation", retryFunc)
	assert.True(t, modal.IsVisible)

	// Hide modal
	modal.Hide()

	assert.False(t, modal.IsVisible)
	assert.Empty(t, modal.Message)
	assert.Empty(t, modal.Operation)
	assert.Empty(t, modal.Suggestion)
	assert.Nil(t, modal.RetryFunc)
}

func TestErrorModal_SetSize(t *testing.T) {
	modal := NewErrorModal()

	modal.SetSize(100, 50)

	assert.Equal(t, 100, modal.termWidth)
	assert.Equal(t, 50, modal.termHeight)
}

func TestErrorModal_View_NotVisible(t *testing.T) {
	modal := NewErrorModal()

	view := modal.View()

	assert.Empty(t, view)
}

func TestErrorModal_View_Visible(t *testing.T) {
	modal := NewErrorModal()
	modal.SetSize(80, 24)
	modal.Show("test error", "Test Operation", nil)

	view := modal.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Error: Test Operation")
	assert.Contains(t, view, "test error")
	assert.Contains(t, view, "[Press ESC to exit]")
}

func TestErrorModal_View_WithRetryFunc(t *testing.T) {
	modal := NewErrorModal()
	modal.SetSize(80, 24)
	retryFunc := func() tea.Cmd { return nil }

	modal.Show("test error", "Test Operation", retryFunc)

	view := modal.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "[Press Enter to retry]")
	assert.Contains(t, view, "[Press ESC to exit]")
}

func TestErrorModal_View_WithSuggestion(t *testing.T) {
	modal := NewErrorModal()
	modal.SetSize(80, 24)

	modal.ShowWithSuggestion("test error", "Test Operation", "Try this fix", nil)

	view := modal.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "test error")
	assert.Contains(t, view, "Try this fix")
}

func TestErrorModal_HandleKeyPress_NotVisible(t *testing.T) {
	modal := NewErrorModal()

	handled, cmd := modal.HandleKeyPress("enter")

	assert.False(t, handled)
	assert.Nil(t, cmd)
}

func TestErrorModal_HandleKeyPress_Enter_WithRetry(t *testing.T) {
	modal := NewErrorModal()
	retryCalled := false
	retryFunc := func() tea.Cmd {
		retryCalled = true
		return func() tea.Msg { return "test" }
	}

	modal.Show("test error", "Test Operation", retryFunc)

	handled, cmd := modal.HandleKeyPress("enter")

	assert.True(t, handled)
	assert.NotNil(t, cmd)            // Should return the command from retryFunc
	assert.False(t, modal.IsVisible) // Modal should be hidden after retry
	assert.True(t, retryCalled)      // RetryFunc should be called immediately
}

func TestErrorModal_HandleKeyPress_Enter_NoRetry(t *testing.T) {
	modal := NewErrorModal()
	modal.Show("test error", "Test Operation", nil)

	handled, cmd := modal.HandleKeyPress("enter")

	assert.True(t, handled)
	assert.Nil(t, cmd)
	assert.False(t, modal.IsVisible) // Modal should be hidden even without retry
}

func TestErrorModal_HandleKeyPress_ESC(t *testing.T) {
	modal := NewErrorModal()
	modal.Show("test error", "Test Operation", nil)
	assert.True(t, modal.IsVisible)

	handled, cmd := modal.HandleKeyPress("esc")

	assert.True(t, handled)
	assert.Nil(t, cmd)
	assert.False(t, modal.IsVisible) // Modal should be hidden
}

func TestErrorModal_HandleKeyPress_OtherKey(t *testing.T) {
	modal := NewErrorModal()
	modal.Show("test error", "Test Operation", nil)

	handled, cmd := modal.HandleKeyPress("q")

	assert.True(t, handled)         // Modal blocks all input
	assert.Nil(t, cmd)              // But doesn't execute any command
	assert.True(t, modal.IsVisible) // Modal stays visible
}
