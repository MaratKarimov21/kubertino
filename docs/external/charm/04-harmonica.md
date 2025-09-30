# Harmonica - Spring Animation Library Guide

## Overview

Harmonica is a spring physics-based animation library for terminal applications. It provides smooth, natural-feeling animations using damped harmonic oscillators (spring physics). Unlike time-based linear animations, spring animations feel more responsive and organic, making terminal UIs more engaging.

**Repository:** https://github.com/charmbracelet/harmonica
**Package:** `github.com/charmbracelet/harmonica`
**Import:** `"github.com/charmbracelet/harmonica"`

## Installation

```bash
go get github.com/charmbracelet/harmonica
```

## Core Concepts

### What is Spring Animation?

Spring animation simulates the motion of an object attached to a spring:
- Objects accelerate and decelerate naturally
- Can overshoot target (bounce) or smoothly approach it
- More realistic than linear interpolation
- Responsive to user input changes mid-animation

### Why Use Spring Physics in TUIs?

1. **Natural Motion**: Feels more organic than linear animations
2. **Interruptible**: Can change target mid-animation smoothly
3. **Responsive**: Adapts to input changes instantly
4. **Professional Feel**: Makes terminal apps feel polished
5. **Predictable**: Physics-based, not arbitrary timing

## Spring Physics Model

### The Damped Harmonic Oscillator

Harmonica implements a damped harmonic oscillator with the formula:

```
F = -kx - bv

Where:
- F: Force applied to the object
- k: Spring stiffness constant
- x: Displacement from target
- b: Damping coefficient
- v: Velocity
```

### Key Components

```go
type Spring struct {
    Value    float64  // Current position
    Velocity float64  // Current velocity
    Target   float64  // Target position
    Damping  float64  // Damping ratio (0-1+)
    Stiffness float64 // Spring stiffness
}
```

## Basic Usage

### Creating a Spring

```go
import "github.com/charmbracelet/harmonica"

// Create spring with default parameters
spring := harmonica.NewSpring(0.0, 15.0, 1.0)
// Start: 0.0, Stiffness: 15.0, Damping: 1.0

// Set target
spring.SetTarget(100.0)
```

### Updating the Spring

```go
// In your Bubble Tea Update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.FrameMsg:
        // Update spring with time delta
        // Typically 16.67ms (60 FPS) or frame time
        m.spring.Update(msg.Delta)

        // Check if animation complete
        if !m.spring.Done() {
            return m, tea.Tick(time.Millisecond*16, func(time.Time) tea.Msg {
                return tea.FrameMsg{}
            })
        }
    }
    return m, nil
}
```

### Getting Current Value

```go
currentPosition := spring.Value()
isAnimating := !spring.Done()
```

## Configuration Parameters

### Time Delta

The time elapsed since last update (in seconds or milliseconds):

```go
// Using seconds
spring.Update(0.016) // 60 FPS = 16ms = 0.016s

// Using milliseconds (convert to seconds)
delta := time.Since(lastUpdate).Seconds()
spring.Update(delta)
```

### Angular Velocity (Stiffness)

Controls how quickly the spring responds (10-300):

```go
// Low stiffness (10-50): Slow, gentle motion
spring := harmonica.NewSpring(0, 10, 1.0)

// Medium stiffness (50-150): Balanced motion
spring := harmonica.NewSpring(0, 100, 1.0)

// High stiffness (150-300): Snappy, quick motion
spring := harmonica.NewSpring(0, 250, 1.0)
```

**Common Values:**
- `10-30`: Smooth, floating animations
- `50-100`: General UI animations
- `150-250`: Snappy, responsive animations
- `300+`: Very fast, almost instant

### Damping Ratio

Controls oscillation behavior (0.0-2.0+):

```go
// Under-damped (< 1.0): Overshoots and bounces
spring := harmonica.NewSpring(0, 100, 0.5)  // Bouncy

// Critically damped (= 1.0): No overshoot, fastest
spring := harmonica.NewSpring(0, 100, 1.0)  // Perfect

// Over-damped (> 1.0): Slow, no overshoot
spring := harmonica.NewSpring(0, 100, 1.5)  // Sluggish
```

**Damping Values:**
- `0.0-0.5`: Very bouncy, multiple oscillations
- `0.5-0.8`: Slight overshoot, playful
- `1.0`: Critical damping, no overshoot (recommended)
- `1.0-2.0`: Over-damped, slower approach

## Damping Modes

### Under-Damped (Bouncy)

Springs with damping < 1.0 overshoot their target:

```go
bouncySpring := harmonica.NewSpring(0, 100, 0.6)
bouncySpring.SetTarget(100)

// Value will overshoot 100, then oscillate back
// Example progression: 0 → 80 → 110 → 95 → 102 → 99 → 100
```

**Use Cases:**
- Playful UI elements
- Attention-grabbing animations
- Game-like interactions
- Bouncy buttons or icons

### Critically-Damped (Smooth)

Damping = 1.0 reaches target quickly without overshoot:

```go
smoothSpring := harmonica.NewSpring(0, 100, 1.0)
smoothSpring.SetTarget(100)

// Value smoothly approaches 100 without overshooting
// Example progression: 0 → 60 → 85 → 95 → 99 → 100
```

**Use Cases:**
- Professional UIs
- Smooth scrolling
- Panel sliding
- General animations (recommended default)

### Over-Damped (Gentle)

Damping > 1.0 approaches target slowly:

```go
gentleSpring := harmonica.NewSpring(0, 100, 1.5)
gentleSpring.SetTarget(100)

// Value very gradually approaches 100
// Example progression: 0 → 40 → 70 → 85 → 93 → 97 → 99 → 100
```

**Use Cases:**
- Subtle background animations
- Gentle transitions
- Low-priority animations
- Accessibility (reduced motion)

## Integration with Bubble Tea

### Basic Animation Loop

```go
import (
    "time"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/harmonica"
)

type frameMsg struct {
    delta float64
}

type model struct {
    spring     harmonica.Spring
    lastUpdate time.Time
    animating  bool
}

func initialModel() model {
    return model{
        spring:     harmonica.NewSpring(0, 100, 1.0),
        lastUpdate: time.Now(),
        animating:  false,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "space":
            // Start animation
            m.spring.SetTarget(100)
            m.lastUpdate = time.Now()
            m.animating = true
            return m, m.tick()
        }

    case frameMsg:
        // Update spring
        m.spring.Update(msg.delta)

        // Continue animating or stop
        if !m.spring.Done() {
            m.animating = true
            return m, m.tick()
        }
        m.animating = false
    }

    return m, nil
}

func (m model) tick() tea.Cmd {
    return func() tea.Msg {
        now := time.Now()
        delta := now.Sub(m.lastUpdate).Seconds()
        m.lastUpdate = now
        return frameMsg{delta: delta}
    }
}

func (m model) View() string {
    return fmt.Sprintf("Position: %.2f\n", m.spring.Value())
}
```

### Smooth Scrolling

```go
type model struct {
    spring      harmonica.Spring
    scrollY     int
    targetY     int
    items       []string
    lastUpdate  time.Time
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "down", "j":
            m.targetY += 10
            m.spring.SetTarget(float64(m.targetY))
            m.lastUpdate = time.Now()
            return m, m.animate()

        case "up", "k":
            m.targetY -= 10
            if m.targetY < 0 {
                m.targetY = 0
            }
            m.spring.SetTarget(float64(m.targetY))
            m.lastUpdate = time.Now()
            return m, m.animate()
        }

    case frameMsg:
        m.spring.Update(msg.delta)
        m.scrollY = int(m.spring.Value())

        if !m.spring.Done() {
            return m, m.animate()
        }
    }

    return m, nil
}

func (m model) animate() tea.Cmd {
    return tea.Tick(time.Millisecond*16, func(time.Time) tea.Msg {
        now := time.Now()
        delta := now.Sub(m.lastUpdate).Seconds()
        m.lastUpdate = now
        return frameMsg{delta: delta}
    })
}
```

### Elastic Menu Selection

```go
type model struct {
    cursor     int
    spring     harmonica.Spring
    cursorY    float64
    items      []string
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            if m.cursor > 0 {
                m.cursor--
                // Bouncy spring for playful feel
                m.spring = harmonica.NewSpring(m.cursorY, 100, 0.7)
                m.spring.SetTarget(float64(m.cursor * 2))
                return m, m.animate()
            }

        case "down":
            if m.cursor < len(m.items)-1 {
                m.cursor++
                m.spring = harmonica.NewSpring(m.cursorY, 100, 0.7)
                m.spring.SetTarget(float64(m.cursor * 2))
                return m, m.animate()
            }
        }

    case frameMsg:
        m.spring.Update(msg.delta)
        m.cursorY = m.spring.Value()

        if !m.spring.Done() {
            return m, m.animate()
        }
    }

    return m, nil
}

func (m model) View() string {
    s := "Select an item:\n\n"

    for i, item := range m.items {
        cursor := "  "
        // Smooth cursor position based on spring
        if float64(i*2) <= m.cursorY && m.cursorY < float64((i+1)*2) {
            cursor = "→ "
        }
        s += fmt.Sprintf("%s%s\n", cursor, item)
    }

    return s
}
```

## Use Cases

### 1. Smooth Scrolling

```go
// Viewport with spring-based scrolling
type viewport struct {
    spring     harmonica.Spring
    offset     int
    targetOffset int
}

func (v *viewport) ScrollTo(offset int) tea.Cmd {
    v.targetOffset = offset
    v.spring.SetTarget(float64(offset))
    return v.animate()
}
```

### 2. Panel Sliding

```go
// Sliding sidebar
type sidebar struct {
    spring harmonica.Spring
    open   bool
    width  int
}

func (s *sidebar) Toggle() tea.Cmd {
    if s.open {
        s.spring.SetTarget(0)
        s.open = false
    } else {
        s.spring.SetTarget(float64(s.width))
        s.open = true
    }
    return s.animate()
}
```

### 3. Progress Indicators

```go
// Smooth progress bar
type progressBar struct {
    spring  harmonica.Spring
    current float64
}

func (p *progressBar) SetProgress(percent float64) tea.Cmd {
    p.spring.SetTarget(percent)
    return p.animate()
}
```

### 4. Natural Motion Transitions

```go
// Opacity fade with spring physics
type fadeEffect struct {
    spring harmonica.Spring
    alpha  float64
}

func (f *fadeEffect) FadeIn() tea.Cmd {
    f.spring.SetTarget(1.0)
    return f.animate()
}

func (f *fadeEffect) FadeOut() tea.Cmd {
    f.spring.SetTarget(0.0)
    return f.animate()
}
```

### 5. Elastic Effects

```go
// Button press effect
type button struct {
    spring harmonica.Spring
    scale  float64
    pressed bool
}

func (b *button) Press() tea.Cmd {
    // Scale down with bounce
    b.spring = harmonica.NewSpring(1.0, 200, 0.6)
    b.spring.SetTarget(0.9)
    b.pressed = true
    return b.animate()
}

func (b *button) Release() tea.Cmd {
    // Scale back with bounce
    b.spring.SetTarget(1.0)
    b.pressed = false
    return b.animate()
}
```

## Performance Considerations

### Frame Rate Optimization

```go
// Adaptive frame rate based on motion
func (m model) tick() tea.Cmd {
    var interval time.Duration

    // Faster updates when moving quickly
    if math.Abs(m.spring.Velocity) > 50 {
        interval = time.Millisecond * 8  // 120 FPS
    } else if math.Abs(m.spring.Velocity) > 10 {
        interval = time.Millisecond * 16 // 60 FPS
    } else {
        interval = time.Millisecond * 33 // 30 FPS
    }

    return tea.Tick(interval, func(time.Time) tea.Msg {
        return frameMsg{}
    })
}
```

### Early Termination

```go
// Stop animating when close enough to target
func (s *Spring) NearlyDone(threshold float64) bool {
    return math.Abs(s.Value() - s.Target) < threshold &&
           math.Abs(s.Velocity) < threshold
}

// Usage
if m.spring.NearlyDone(0.1) {
    m.spring.Value = m.spring.Target
    m.spring.Velocity = 0
    return m, nil // Stop animating
}
```

### Multiple Springs

```go
// Manage multiple independent animations
type model struct {
    springs map[string]*harmonica.Spring
    active  int // Count of active springs
}

func (m *model) UpdateSprings(delta float64) {
    m.active = 0
    for _, spring := range m.springs {
        spring.Update(delta)
        if !spring.Done() {
            m.active++
        }
    }
}

// Only animate if springs are active
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case frameMsg:
        m.UpdateSprings(msg.delta)
        if m.active > 0 {
            return m, m.animate()
        }
    }
    return m, nil
}
```

## Best Practices

### 1. Choose Appropriate Parameters

```go
// UI Element Type → Recommended Settings
var presets = map[string]SpringConfig{
    "scroll": {stiffness: 50, damping: 1.0},     // Smooth scrolling
    "menu":   {stiffness: 100, damping: 0.8},    // Slight bounce
    "panel":  {stiffness: 80, damping: 1.0},     // Smooth slide
    "button": {stiffness: 200, damping: 0.6},    // Bouncy press
    "fade":   {stiffness: 70, damping: 1.0},     // Smooth fade
}
```

### 2. Handle Interruptions

```go
// Allow interrupting animations smoothly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Change target mid-animation
        m.spring.SetTarget(newTarget)
        // Spring automatically adapts, using current velocity
        return m, m.animate()
    }
    return m, nil
}
```

### 3. Manage Time Delta Accurately

```go
type model struct {
    spring     harmonica.Spring
    lastUpdate time.Time
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    now := time.Now()

    switch msg := msg.(type) {
    case frameMsg:
        // Use actual time delta, not fixed value
        delta := now.Sub(m.lastUpdate).Seconds()
        m.spring.Update(delta)
        m.lastUpdate = now
    }

    return m, nil
}
```

### 4. Clamp Values When Necessary

```go
func (m *model) updatePosition() {
    m.spring.Update(delta)

    // Clamp to valid range
    pos := m.spring.Value()
    if pos < 0 {
        pos = 0
        m.spring.Value = 0
        m.spring.Velocity = 0
    } else if pos > float64(m.maxScroll) {
        pos = float64(m.maxScroll)
        m.spring.Value = pos
        m.spring.Velocity = 0
    }

    m.scrollY = int(pos)
}
```

### 5. Test Different Terminal Refresh Rates

```go
// Adapt to terminal capabilities
func (m model) getFrameInterval() time.Duration {
    // Some terminals refresh slower
    if m.terminalSlow {
        return time.Millisecond * 33 // 30 FPS
    }
    return time.Millisecond * 16 // 60 FPS
}
```

## Common Patterns

### Animation Manager

```go
type AnimationManager struct {
    springs    map[string]*harmonica.Spring
    lastUpdate time.Time
}

func (am *AnimationManager) Animate(name string, target float64) {
    if spring, exists := am.springs[name]; exists {
        spring.SetTarget(target)
    }
}

func (am *AnimationManager) Update() bool {
    now := time.Now()
    delta := now.Sub(am.lastUpdate).Seconds()
    am.lastUpdate = now

    anyActive := false
    for _, spring := range am.springs {
        spring.Update(delta)
        if !spring.Done() {
            anyActive = true
        }
    }

    return anyActive
}
```

### Easing Functions Wrapper

```go
// Use springs to implement easing functions
func EaseInOut(start, end, t float64) float64 {
    spring := harmonica.NewSpring(start, 100, 1.0)
    spring.SetTarget(end)

    // Simulate spring for duration
    totalTime := 0.0
    for totalTime < 1.0 {
        spring.Update(0.016)
        totalTime += 0.016
    }

    return spring.Value()
}
```

## Complete Example: Animated Dashboard

```go
package main

import (
    "fmt"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/harmonica"
    "github.com/charmbracelet/lipgloss"
)

type frameMsg struct {
    delta float64
}

type model struct {
    panels      []panel
    selected    int
    lastUpdate  time.Time
}

type panel struct {
    title  string
    spring harmonica.Spring
    height float64
}

func initialModel() model {
    return model{
        panels: []panel{
            {title: "CPU", spring: harmonica.NewSpring(0, 100, 1.0)},
            {title: "Memory", spring: harmonica.NewSpring(0, 100, 1.0)},
            {title: "Disk", spring: harmonica.NewSpring(0, 100, 1.0)},
        },
        lastUpdate: time.Now(),
    }
}

func (m model) Init() tea.Cmd {
    return m.tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit

        case "tab":
            m.selected = (m.selected + 1) % len(m.panels)
            // Animate selected panel expansion
            for i := range m.panels {
                if i == m.selected {
                    m.panels[i].spring.SetTarget(20)
                } else {
                    m.panels[i].spring.SetTarget(10)
                }
            }
            return m, m.tick()
        }

    case frameMsg:
        // Update all springs
        anyActive := false
        for i := range m.panels {
            m.panels[i].spring.Update(msg.delta)
            m.panels[i].height = m.panels[i].spring.Value()
            if !m.panels[i].spring.Done() {
                anyActive = true
            }
        }

        // Continue animating if needed
        if anyActive {
            return m, m.tick()
        }
    }

    return m, nil
}

func (m model) tick() tea.Cmd {
    return tea.Tick(time.Millisecond*16, func(time.Time) tea.Msg {
        now := time.Now()
        delta := now.Sub(m.lastUpdate).Seconds()
        m.lastUpdate = now
        return frameMsg{delta: delta}
    })
}

func (m model) View() string {
    var views []string

    for i, panel := range m.panels {
        style := lipgloss.NewStyle().
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("240")).
            Width(40).
            Height(int(panel.height))

        if i == m.selected {
            style = style.BorderForeground(lipgloss.Color("205"))
        }

        views = append(views, style.Render(panel.title))
    }

    return lipgloss.JoinVertical(lipgloss.Left, views...) +
        "\n\nPress Tab to switch panels, q to quit"
}

func main() {
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
    }
}
```

## Resources

- **Official Repository:** https://github.com/charmbracelet/harmonica
- **Package Documentation:** https://pkg.go.dev/github.com/charmbracelet/harmonica
- **Physics Background:** https://en.wikipedia.org/wiki/Harmonic_oscillator
- **Spring Animation Theory:** https://www.ryanjuckett.com/damped-springs/
- **Related:** [01-bubbletea.md](01-bubbletea.md), [03-lipgloss.md](03-lipgloss.md)