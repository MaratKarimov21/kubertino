# Charm Ecosystem Overview

## Introduction

The Charm ecosystem is a comprehensive suite of libraries for building modern Terminal User Interfaces (TUIs) in Go. Created by Charm (charmbracelet), these tools provide everything needed to create beautiful, interactive, and performant command-line applications.

## Core Libraries

### 1. Bubble Tea - TUI Framework
**Repository:** https://github.com/charmbracelet/bubbletea
**Package:** github.com/charmbracelet/bubbletea

The foundational framework for building terminal applications using the Elm Architecture (Model-View-Update pattern). Bubble Tea handles the program lifecycle, event loop, terminal management, and rendering.

**Key Features:**
- Functional, state-based architecture
- Event-driven message passing system
- Cross-platform support (Unix/Windows)
- Framerate-based rendering
- Keyboard and mouse event handling
- Terminal resize and focus management

### 2. Bubbles - Component Library
**Repository:** https://github.com/charmbracelet/bubbles
**Package:** github.com/charmbracelet/bubbles

A collection of pre-built, customizable UI components designed specifically for Bubble Tea applications.

**Available Components:**
- Text Input (single-line)
- Text Area (multi-line)
- Table
- List
- Viewport (scrollable content)
- Spinner
- Progress Bar
- Paginator
- Timer & Stopwatch
- File Picker
- Help & Key bindings

### 3. Lip Gloss - Styling Framework
**Repository:** https://github.com/charmbracelet/lipgloss
**Package:** github.com/charmbracelet/lipgloss

A declarative, CSS-like styling library for terminal applications. Provides comprehensive formatting, layout, and color support.

**Key Features:**
- Chainable style API
- Full color support (ANSI, 256-color, TrueColor)
- Adaptive colors for different terminal backgrounds
- Layout capabilities (padding, margins, alignment, borders)
- Advanced features (gradients, tables, lists)
- Automatic color profile detection

### 4. Harmonica - Animation Library
**Repository:** https://github.com/charmbracelet/harmonica
**Package:** github.com/charmbracelet/harmonica

A spring physics-based animation library for creating smooth, natural motion in terminal applications.

**Key Features:**
- Spring animation model with configurable damping
- Framework-agnostic (works in 2D/3D contexts)
- Three damping modes: under-damped, critically-damped, over-damped
- Smooth scrolling and elastic effects
- Optimized for terminal framerates

### 5. BubbleZone - Mouse Event Tracking
**Repository:** https://github.com/lrstanley/bubblezone
**Package:** github.com/lrstanley/bubblezone

A utility library that simplifies mouse event tracking for Bubble Tea components, enabling clickable regions in TUIs.

**Key Features:**
- Zero-width zone markers
- Global and local zone management
- Coordinate tracking and bounds checking
- Nested component support
- Optimized for complex interfaces

### 6. ntcharts - Terminal Charting
**Repository:** https://github.com/NimbleMarkets/ntcharts
**Package:** github.com/NimbleMarkets/ntcharts

A comprehensive charting library built specifically for Bubble Tea, providing various chart types for data visualization in terminals.

**Available Chart Types:**
- Canvas (foundation for custom plotting)
- Bar Charts (horizontal/vertical)
- Line Charts (standard, time-series, streaming, waveline)
- Heat Maps
- Scatter Charts
- Sparklines
- Candlestick/OHLC Charts

## Architecture Philosophy

The Charm ecosystem follows several key architectural principles:

### 1. Functional Design
All libraries embrace functional programming patterns, with immutable state and pure functions where possible.

### 2. Composability
Components are designed to work together seamlessly. You can combine Bubbles components, style them with Lip Gloss, animate with Harmonica, and make them interactive with BubbleZone.

### 3. Progressive Complexity
Start simple with basic components and styling, then add complexity as needed (animations, mouse interactions, charts).

### 4. Terminal-First
All libraries are optimized for terminal environments, respecting terminal capabilities and providing graceful degradation.

### 5. Developer Experience
Intuitive APIs, comprehensive documentation, and clear examples make it easy to build sophisticated TUIs.

## Integration Patterns

### Basic Stack
```
Bubble Tea (Framework)
    └── Bubbles (Components)
        └── Lip Gloss (Styling)
```

### Full Stack
```
Bubble Tea (Framework)
    ├── Bubbles (Components)
    │   └── Lip Gloss (Styling)
    ├── Harmonica (Animations)
    ├── BubbleZone (Mouse Tracking)
    └── ntcharts (Data Visualization)
```

## Use Cases

The Charm ecosystem is ideal for:

1. **Interactive CLI Tools** - Complex command-line applications with rich UIs
2. **Development Tools** - IDE-like tools, debuggers, deployment dashboards
3. **System Monitoring** - Real-time dashboards, log viewers, metrics displays
4. **Data Analysis** - Terminal-based data exploration and visualization
5. **DevOps Tools** - Kubernetes dashboards (like Kubertino), deployment tools
6. **Terminal Games** - Interactive games with smooth animations

## Industry Adoption

The Charm ecosystem is used by major organizations including:
- Microsoft
- Google
- GitHub
- AWS
- And many open-source projects

## Getting Started

### Installation

```bash
# Core framework
go get github.com/charmbracelet/bubbletea

# Component library
go get github.com/charmbracelet/bubbles

# Styling
go get github.com/charmbracelet/lipgloss

# Animation (optional)
go get github.com/charmbracelet/harmonica

# Mouse tracking (optional)
go get github.com/lrstanley/bubblezone

# Charting (optional)
go get github.com/NimbleMarkets/ntcharts
```

### Minimal Example

```go
package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct{}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    return "Hello, Bubble Tea! Press 'q' to quit.\n"
}

func main() {
    p := tea.NewProgram(model{})
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
    }
}
```

## Documentation Structure

This documentation is organized into the following sections:

1. **00-overview.md** (this file) - High-level overview of the ecosystem
2. **01-bubbletea.md** - Comprehensive Bubble Tea framework guide
3. **02-bubbles.md** - Component library reference
4. **03-lipgloss.md** - Styling and layout guide
5. **04-harmonica.md** - Animation library guide
6. **05-bubblezone.md** - Mouse event tracking guide
7. **06-ntcharts.md** - Terminal charting reference
8. **99-best-practices.md** - Patterns, tips, and best practices

## Resources

- **Official Tutorial:** https://github.com/charmbracelet/bubbletea/tree/master/tutorials
- **Examples:** https://github.com/charmbracelet/bubbletea/tree/master/examples
- **Community:** GitHub Discussions on each repository
- **Blog:** https://charm.sh/blog/

## License

All Charm libraries are open-source and MIT licensed, making them suitable for both commercial and personal projects.