# BubbleZone - Mouse Tracking Library Guide

## Overview

BubbleZone is a library for handling mouse interactions in Bubble Tea applications. It provides a simple API for defining clickable regions (zones) and tracking mouse events, making it easy to add mouse support to terminal UIs without complex coordinate tracking.

**Repository:** https://github.com/charmbracelet/bubblezone
**Package:** `github.com/charmbracelet/bubblezone`
**Import:** `zone "github.com/charmbracelet/bubblezone"`

## Installation

```bash
go get github.com/charmbracelet/bubblezone
```

## Core Concepts

### What is BubbleZone?

BubbleZone wraps text content with invisible markers that identify clickable regions. When a mouse click occurs, you can determine which zone was clicked by checking the mouse coordinates against registered zones.

### How It Works

1. **Mark**: Wrap content with zone markers using a unique ID
2. **Render**: Display the marked content (markers are invisible)
3. **Handle**: Process mouse events and check which zone was clicked
4. **React**: Update your UI based on the clicked zone

### Key Components

- **Manager**: Tracks all zones and their positions
- **Zone ID**: Unique identifier for each clickable region
- **Bounds**: Coordinates defining the zone's clickable area
- **Mouse Event**: Click coordinates to check against zones

## Basic Usage

### 1. Create a Manager

```go
import zone "github.com/charmbracelet/bubblezone"

type model struct {
    zone *zone.Manager
}

func initialModel() model {
    return model{
        zone: zone.New(),
    }
}
```

### 2. Mark Clickable Content

```go
func (m model) View() string {
    // Mark content with zone ID
    button := m.zone.Mark("button-1", "[ Click Me ]")

    return fmt.Sprintf("Here's a button: %s", button)
}
```

### 3. Handle Mouse Events

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == tea.MouseLeft {
            // Check which zone was clicked
            if m.zone.Get("button-1").InBounds(msg) {
                // Button was clicked!
                return m, m.handleButtonClick()
            }
        }
    }
    return m, nil
}
```

### 4. Scan Content (Alternative)

Instead of checking specific zones, scan to find what was clicked:

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == tea.MouseLeft {
            // Find which zone (if any) was clicked
            if id := m.zone.Get("").Scan(msg); id != "" {
                return m, m.handleZoneClick(id)
            }
        }
    }
    return m, nil
}
```

## Core API

### Manager Methods

#### New() - Create Manager

```go
manager := zone.New()
```

#### Mark() - Mark Content as Zone

```go
// Mark content with unique ID
marked := manager.Mark("zone-id", "content")

// The marked string contains invisible markers
// When rendered, only "content" is visible
```

#### Get() - Retrieve Zone

```go
// Get zone by ID
z := manager.Get("zone-id")

// Check if zone exists
if z != nil {
    // Zone found
}
```

#### Scan() - Find Clicked Zone

```go
// Find which zone was clicked at mouse coordinates
clickedID := manager.Get("").Scan(mouseMsg)

if clickedID != "" {
    fmt.Printf("Clicked zone: %s\n", clickedID)
}
```

### Zone Methods

#### InBounds() - Check if Mouse is in Zone

```go
zone := manager.Get("zone-id")

if zone.InBounds(mouseMsg) {
    // Mouse is within this zone
}
```

#### Bounds() - Get Zone Coordinates

```go
zone := manager.Get("zone-id")

minX, maxX, minY, maxY := zone.Bounds()

fmt.Printf("Zone at (%d,%d) to (%d,%d)\n", minX, minY, maxX, maxY)
```

## Global vs Local Zone Managers

### Global Manager

Use the global manager for simple applications:

```go
import zone "github.com/charmbracelet/bubblezone"

func (m model) View() string {
    // Uses global manager
    return zone.Mark("button", "[ Click ]")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if zone.Get("button").InBounds(msg) {
            // Handle click
        }
    }
    return m, nil
}
```

### Local Manager (Recommended)

Use local managers for better encapsulation:

```go
type model struct {
    zones *zone.Manager
}

func initialModel() model {
    return model{
        zones: zone.New(),
    }
}

func (m model) View() string {
    return m.zones.Mark("button", "[ Click ]")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if m.zones.Get("button").InBounds(msg) {
            // Handle click
        }
    }
    return m, nil
}
```

## Integration with Bubble Tea Components

### Clickable Buttons

```go
type button struct {
    id      string
    label   string
    onClick func() tea.Cmd
}

func (b button) Render(zones *zone.Manager) string {
    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("63")).
        Padding(0, 2).
        Bold(true)

    content := style.Render(b.label)
    return zones.Mark(b.id, content)
}

func (b button) HandleClick(zones *zone.Manager, msg tea.MouseMsg) tea.Cmd {
    if zones.Get(b.id).InBounds(msg) {
        return b.onClick()
    }
    return nil
}

// Usage
type model struct {
    zones   *zone.Manager
    buttons []button
}

func (m model) View() string {
    var buttons []string
    for _, btn := range m.buttons {
        buttons = append(buttons, btn.Render(m.zones))
    }
    return lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == tea.MouseLeft {
            for _, btn := range m.buttons {
                if cmd := btn.HandleClick(m.zones, msg); cmd != nil {
                    return m, cmd
                }
            }
        }
    }
    return m, nil
}
```

### Clickable List Items

```go
type listItem struct {
    id   string
    text string
}

func (m model) View() string {
    var items []string

    for _, item := range m.items {
        style := lipgloss.NewStyle().Padding(0, 1)

        // Highlight if hovered
        if m.zones.Get(item.id).InBounds(m.lastMouse) {
            style = style.Background(lipgloss.Color("240"))
        }

        content := style.Render(item.text)
        marked := m.zones.Mark(item.id, content)
        items = append(items, marked)
    }

    return lipgloss.JoinVertical(lipgloss.Left, items...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        m.lastMouse = msg

        if msg.Type == tea.MouseLeft {
            // Find which item was clicked
            for _, item := range m.items {
                if m.zones.Get(item.id).InBounds(msg) {
                    return m, m.selectItem(item.id)
                }
            }
        }
    }
    return m, nil
}
```

### Interactive Menu

```go
type menuItem struct {
    id     string
    label  string
    action func() tea.Cmd
}

type menu struct {
    zones *zone.Manager
    items []menuItem
}

func (m menu) Render() string {
    var items []string

    for _, item := range m.items {
        style := lipgloss.NewStyle().
            Foreground(lipgloss.Color("252")).
            Padding(0, 2)

        content := style.Render(item.label)
        marked := m.zones.Mark(item.id, content)
        items = append(items, marked)
    }

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        items...,
    )
}

func (m menu) HandleClick(msg tea.MouseMsg) tea.Cmd {
    if msg.Type == tea.MouseLeft {
        for _, item := range m.items {
            if m.zones.Get(item.id).InBounds(msg) {
                return item.action()
            }
        }
    }
    return nil
}
```

## Mouse Event Handling Patterns

### Pattern 1: Direct Zone Checking

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == tea.MouseLeft {
            if m.zones.Get("button-1").InBounds(msg) {
                return m, m.action1()
            }
            if m.zones.Get("button-2").InBounds(msg) {
                return m, m.action2()
            }
        }
    }
    return m, nil
}
```

### Pattern 2: Scan and Switch

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == tea.MouseLeft {
            switch m.zones.Get("").Scan(msg) {
            case "button-1":
                return m, m.action1()
            case "button-2":
                return m, m.action2()
            case "button-3":
                return m, m.action3()
            }
        }
    }
    return m, nil
}
```

### Pattern 3: Hover Detection

```go
type model struct {
    zones      *zone.Manager
    hoveredID  string
    lastMouse  tea.MouseMsg
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        m.lastMouse = msg

        // Update hover state
        if msg.Type == tea.MouseMotion {
            newHovered := m.zones.Get("").Scan(msg)
            if newHovered != m.hoveredID {
                m.hoveredID = newHovered
                // Trigger re-render for hover effect
            }
        }

        // Handle clicks
        if msg.Type == tea.MouseLeft {
            if m.hoveredID != "" {
                return m, m.handleClick(m.hoveredID)
            }
        }
    }
    return m, nil
}

func (m model) View() string {
    // Render with hover effects
    for _, item := range m.items {
        style := lipgloss.NewStyle()
        if item.id == m.hoveredID {
            style = style.Background(lipgloss.Color("240"))
        }
        // Render item...
    }
}
```

### Pattern 4: Drag Detection

```go
type model struct {
    zones     *zone.Manager
    dragging  bool
    dragStart tea.MouseMsg
    dragID    string
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        switch msg.Type {
        case tea.MouseLeft:
            // Start drag
            if id := m.zones.Get("").Scan(msg); id != "" {
                m.dragging = true
                m.dragStart = msg
                m.dragID = id
            }

        case tea.MouseRelease:
            // End drag
            if m.dragging {
                m.dragging = false
                return m, m.handleDragEnd(m.dragID, msg)
            }

        case tea.MouseMotion:
            // Continue drag
            if m.dragging {
                return m, m.handleDrag(m.dragID, msg)
            }
        }
    }
    return m, nil
}
```

## Coordinate Tracking and Bounds Checking

### Understanding Coordinates

```go
// Mouse coordinates are relative to terminal window
// (0,0) is top-left corner

type tea.MouseMsg struct {
    X    int  // Column (horizontal position)
    Y    int  // Row (vertical position)
    Type MouseEventType
}

// Zone bounds define rectangular region
type Bounds struct {
    MinX int  // Left edge
    MaxX int  // Right edge
    MinY int  // Top edge
    MaxY int  // Bottom edge
}
```

### Manual Bounds Checking

```go
func isInBounds(msg tea.MouseMsg, minX, maxX, minY, maxY int) bool {
    return msg.X >= minX && msg.X <= maxX &&
           msg.Y >= minY && msg.Y <= maxY
}

// Usage without BubbleZone
if isInBounds(mouseMsg, 10, 30, 5, 7) {
    // Mouse is in region
}
```

### Relative Positioning

```go
type component struct {
    zones  *zone.Manager
    x, y   int  // Component position
    width  int
    height int
}

func (c *component) Render() string {
    content := lipgloss.NewStyle().
        Width(c.width).
        Height(c.height).
        Render("Content")

    return c.zones.Mark(c.id, content)
}

func (c *component) Contains(msg tea.MouseMsg) bool {
    // Check if mouse is within component bounds
    return msg.X >= c.x && msg.X < c.x+c.width &&
           msg.Y >= c.y && msg.Y < c.y+c.height
}
```

## Performance Optimization

### 1. Efficient Zone Management

```go
// Reuse zone IDs, don't create new ones each frame
type model struct {
    zones   *zone.Manager
    zoneIDs map[string]string  // Cache zone IDs
}

func (m *model) getZoneID(key string) string {
    if id, exists := m.zoneIDs[key]; exists {
        return id
    }
    id := fmt.Sprintf("zone-%s", key)
    m.zoneIDs[key] = id
    return id
}
```

### 2. Batch Zone Checks

```go
// Check multiple zones efficiently
func (m model) findClickedZone(msg tea.MouseMsg, ids []string) string {
    for _, id := range ids {
        if m.zones.Get(id).InBounds(msg) {
            return id
        }
    }
    return ""
}
```

### 3. Conditional Zone Creation

```go
// Only create zones for visible items
func (m model) View() string {
    var visible []string

    for i := m.scrollOffset; i < m.scrollOffset+m.visibleCount; i++ {
        if i >= len(m.items) {
            break
        }

        item := m.items[i]
        zoneID := fmt.Sprintf("item-%d", i)
        marked := m.zones.Mark(zoneID, item.Render())
        visible = append(visible, marked)
    }

    return lipgloss.JoinVertical(lipgloss.Left, visible...)
}
```

### 4. Early Exit on Mouse Events

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        // Ignore non-click events if not needed
        if msg.Type != tea.MouseLeft {
            return m, nil
        }

        // Check zones in priority order
        for _, id := range m.priorityZones {
            if m.zones.Get(id).InBounds(msg) {
                return m, m.handleZoneClick(id)
            }
        }
    }
    return m, nil
}
```

## Complete Working Examples

### Example 1: Button Grid

```go
package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    zone "github.com/charmbracelet/bubblezone"
)

type model struct {
    zones   *zone.Manager
    buttons []button
    clicked string
}

type button struct {
    id    string
    label string
}

func initialModel() model {
    return model{
        zones: zone.New(),
        buttons: []button{
            {"btn-1", "Button 1"},
            {"btn-2", "Button 2"},
            {"btn-3", "Button 3"},
            {"btn-4", "Button 4"},
        },
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }

    case tea.MouseMsg:
        if msg.Type == tea.MouseLeft {
            for _, btn := range m.buttons {
                if m.zones.Get(btn.id).InBounds(msg) {
                    m.clicked = btn.label
                }
            }
        }
    }

    return m, nil
}

func (m model) View() string {
    buttonStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("63")).
        Padding(1, 3).
        Margin(0, 1)

    var row1, row2 []string

    for i, btn := range m.buttons {
        content := buttonStyle.Render(btn.label)
        marked := m.zones.Mark(btn.id, content)

        if i < 2 {
            row1 = append(row1, marked)
        } else {
            row2 = append(row2, marked)
        }
    }

    grid := lipgloss.JoinVertical(lipgloss.Left,
        lipgloss.JoinHorizontal(lipgloss.Top, row1...),
        lipgloss.JoinHorizontal(lipgloss.Top, row2...),
    )

    status := fmt.Sprintf("\nLast clicked: %s", m.clicked)
    return grid + status + "\n\nPress 'q' to quit"
}

func main() {
    p := tea.NewProgram(initialModel(), tea.WithMouseAllMotion())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
    }
}
```

### Example 2: Interactive List

```go
package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    zone "github.com/charmbracelet/bubblezone"
)

type model struct {
    zones     *zone.Manager
    items     []string
    selected  int
    hovered   int
    lastMouse tea.MouseMsg
}

func initialModel() model {
    return model{
        zones: zone.New(),
        items: []string{
            "Item 1",
            "Item 2",
            "Item 3",
            "Item 4",
            "Item 5",
        },
        selected: -1,
        hovered:  -1,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }

    case tea.MouseMsg:
        m.lastMouse = msg

        // Update hover state
        m.hovered = -1
        for i := range m.items {
            zoneID := fmt.Sprintf("item-%d", i)
            if m.zones.Get(zoneID).InBounds(msg) {
                m.hovered = i
                break
            }
        }

        // Handle clicks
        if msg.Type == tea.MouseLeft && m.hovered != -1 {
            m.selected = m.hovered
        }
    }

    return m, nil
}

func (m model) View() string {
    var items []string

    for i, item := range m.items {
        style := lipgloss.NewStyle().
            Padding(0, 2).
            Width(30)

        // Apply styling based on state
        switch {
        case i == m.selected:
            style = style.
                Foreground(lipgloss.Color("229")).
                Background(lipgloss.Color("63")).
                Bold(true)
        case i == m.hovered:
            style = style.
                Background(lipgloss.Color("240"))
        }

        content := style.Render(item)
        zoneID := fmt.Sprintf("item-%d", i)
        marked := m.zones.Mark(zoneID, content)
        items = append(items, marked)
    }

    list := lipgloss.JoinVertical(lipgloss.Left, items...)

    info := fmt.Sprintf(
        "\nSelected: %d | Hovered: %d | Mouse: (%d, %d)",
        m.selected,
        m.hovered,
        m.lastMouse.X,
        m.lastMouse.Y,
    )

    return list + info + "\n\nPress 'q' to quit"
}

func main() {
    p := tea.NewProgram(
        initialModel(),
        tea.WithMouseAllMotion(),
    )
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
    }
}
```

### Example 3: Toolbar with Icons

```go
type toolbar struct {
    zones   *zone.Manager
    buttons []toolbarButton
}

type toolbarButton struct {
    id      string
    icon    string
    tooltip string
    action  func() tea.Cmd
}

func (t toolbar) Render() string {
    var buttons []string

    buttonStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("252")).
        Background(lipgloss.Color("236")).
        Padding(0, 1)

    for _, btn := range t.buttons {
        content := buttonStyle.Render(btn.icon)
        marked := t.zones.Mark(btn.id, content)
        buttons = append(buttons, marked)
    }

    return lipgloss.NewStyle().
        Background(lipgloss.Color("235")).
        Width(50).
        Render(lipgloss.JoinHorizontal(lipgloss.Top, buttons...))
}

func (t toolbar) HandleClick(msg tea.MouseMsg) tea.Cmd {
    if msg.Type == tea.MouseLeft {
        for _, btn := range t.buttons {
            if t.zones.Get(btn.id).InBounds(msg) {
                return btn.action()
            }
        }
    }
    return nil
}
```

## Best Practices

### 1. Enable Mouse Support in Program

```go
// Always enable mouse support when using BubbleZone
p := tea.NewProgram(
    initialModel(),
    tea.WithMouseAllMotion(),  // Track all mouse movement
    // or
    tea.WithMouseCellMotion(), // Track only when button pressed
)
```

### 2. Use Descriptive Zone IDs

```go
// Good: Descriptive and unique
m.zones.Mark("button-save", content)
m.zones.Mark("list-item-5", content)
m.zones.Mark("menu-file", content)

// Bad: Generic and confusing
m.zones.Mark("zone1", content)
m.zones.Mark("btn", content)
```

### 3. Handle Zone Cleanup

```go
// Clear zones when switching views
func (m model) switchView(newView int) model {
    m.currentView = newView
    m.zones = zone.New()  // Create fresh manager
    return m
}
```

### 4. Test Mouse Interactions

```go
func TestMouseClick(t *testing.T) {
    m := initialModel()

    // Simulate mouse click
    msg := tea.MouseMsg{
        X:    15,
        Y:    5,
        Type: tea.MouseLeft,
    }

    newModel, _ := m.Update(msg)
    // Assert expected behavior
}
```

### 5. Provide Keyboard Alternatives

```go
// Always support keyboard navigation as fallback
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Keyboard shortcuts
        switch msg.String() {
        case "enter":
            return m, m.selectCurrent()
        }

    case tea.MouseMsg:
        // Mouse interactions
        if msg.Type == tea.MouseLeft {
            // Handle click
        }
    }
    return m, nil
}
```

## Resources

- **Official Repository:** https://github.com/charmbracelet/bubblezone
- **Package Documentation:** https://pkg.go.dev/github.com/charmbracelet/bubblezone
- **Examples:** https://github.com/charmbracelet/bubblezone/tree/main/examples
- **Related:** [01-bubbletea.md](01-bubbletea.md), [03-lipgloss.md](03-lipgloss.md)