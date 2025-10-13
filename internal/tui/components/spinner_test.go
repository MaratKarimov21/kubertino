package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner()

	assert.NotNil(t, spinner)
	assert.False(t, spinner.IsActive)
	assert.Empty(t, spinner.Message)
	assert.Equal(t, 0, spinner.FrameIndex)
	assert.NotEmpty(t, spinner.Frames) // Should have default frames
	assert.Len(t, spinner.Frames, 10)  // Default has 10 unicode frames
}

func TestSpinner_Start(t *testing.T) {
	spinner := NewSpinner()

	spinner.Start("Loading...")

	assert.True(t, spinner.IsActive)
	assert.Equal(t, "Loading...", spinner.Message)
	assert.Equal(t, 0, spinner.FrameIndex) // Reset to 0 on start
}

func TestSpinner_Stop(t *testing.T) {
	spinner := NewSpinner()
	spinner.Start("Loading...")
	assert.True(t, spinner.IsActive)

	spinner.Stop()

	assert.False(t, spinner.IsActive)
	assert.Empty(t, spinner.Message)
}

func TestSpinner_Tick(t *testing.T) {
	spinner := NewSpinner()
	spinner.Start("Loading...")

	// Initial frame
	assert.Equal(t, 0, spinner.FrameIndex)

	// Tick once
	spinner.Tick()
	assert.Equal(t, 1, spinner.FrameIndex)

	// Tick again
	spinner.Tick()
	assert.Equal(t, 2, spinner.FrameIndex)
}

func TestSpinner_Tick_Wraps(t *testing.T) {
	spinner := NewSpinner()
	spinner.Frames = []string{"1", "2", "3"}
	spinner.Start("Loading...")

	// Advance to last frame
	spinner.FrameIndex = 2

	// Tick should wrap to 0
	spinner.Tick()
	assert.Equal(t, 0, spinner.FrameIndex)
}

func TestSpinner_Tick_WhenInactive(t *testing.T) {
	spinner := NewSpinner()
	spinner.FrameIndex = 5

	// Tick when inactive should not change frame
	spinner.Tick()
	assert.Equal(t, 5, spinner.FrameIndex)
}

func TestSpinner_View_NotActive(t *testing.T) {
	spinner := NewSpinner()

	view := spinner.View()

	assert.Empty(t, view)
}

func TestSpinner_View_Active(t *testing.T) {
	spinner := NewSpinner()
	spinner.Start("Loading data...")

	view := spinner.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Loading data...")
	// Should contain the first frame
	assert.Contains(t, view, spinner.Frames[0])
}

func TestSpinner_View_NoFrames(t *testing.T) {
	spinner := NewSpinner()
	spinner.Frames = []string{} // Empty frames
	spinner.Start("Loading...")

	view := spinner.View()

	// Should show message even without frames
	assert.Equal(t, "Loading...", view)
}

func TestSpinner_View_AdvancedFrame(t *testing.T) {
	spinner := NewSpinner()
	spinner.Frames = []string{"A", "B", "C"}
	spinner.Start("Loading...")

	// Advance to frame 1
	spinner.Tick()

	view := spinner.View()

	assert.Contains(t, view, "B")
	assert.Contains(t, view, "Loading...")
}

func TestTickCmd(t *testing.T) {
	cmd := TickCmd()

	assert.NotNil(t, cmd)

	// Execute the command
	msg := cmd()

	// Should return SpinnerTickMsg
	_, ok := msg.(SpinnerTickMsg)
	assert.True(t, ok)
}
