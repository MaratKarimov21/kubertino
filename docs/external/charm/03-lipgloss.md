# Lip Gloss - Styling Library Comprehensive Guide

## Overview

Lip Gloss is a powerful styling library for terminal applications that brings CSS-like styling to the command line. It provides a declarative, chainable API for creating beautiful, consistent terminal UIs with support for colors, borders, alignment, spacing, and advanced layout features.

**Repository:** https://github.com/charmbracelet/lipgloss
**Package:** `github.com/charmbracelet/lipgloss`
**Import:** `lipgloss "github.com/charmbracelet/lipgloss"`

## Installation

```bash
go get github.com/charmbracelet/lipgloss
```

## Core Philosophy

Lip Gloss follows CSS-inspired principles:
- **Declarative Styling**: Define what you want, not how to achieve it
- **Style Composition**: Chain methods to build complex styles
- **Immutability**: Styles don't mutate; operations return new styles
- **Terminal-Aware**: Adapts to terminal capabilities automatically

## Getting Started

### Basic Usage

```go
import "github.com/charmbracelet/lipgloss"

func main() {
    style := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")).
        Padding(1, 4)

    fmt.Println(style.Render("Hello, Lip Gloss!"))
}
```

### Creating Styles

```go
// Method 1: Chain methods
var titleStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("205")).
    MarginTop(1)

// Method 2: Copy and modify existing styles
var subtitleStyle = titleStyle.Copy().
    Bold(false).
    Foreground(lipgloss.Color("241"))

// Method 3: Inline rendering
text := lipgloss.NewStyle().Bold(true).Render("Important")
```

## Color Support

### Color Types

Lip Gloss supports multiple color formats with automatic terminal capability detection.

#### ANSI Colors (16 colors)

```go
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("1"))    // Red
    Background(lipgloss.Color("4"))    // Blue

// Named ANSI colors (0-15)
// 0: Black    4: Blue      8: Bright Black   12: Bright Blue
// 1: Red      5: Magenta   9: Bright Red     13: Bright Magenta
// 2: Green    6: Cyan     10: Bright Green   14: Bright Cyan
// 3: Yellow   7: White    11: Bright Yellow  15: Bright White
```

#### 256-Color Palette

```go
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("205")).  // Pink
    Background(lipgloss.Color("235"))   // Dark gray

// 0-15:   Standard ANSI colors
// 16-231: 6×6×6 RGB cube
// 232-255: Grayscale
```

#### True Color (24-bit RGB)

```go
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#FF00FF")).     // Hex format
    Background(lipgloss.Color("rgb(255,0,0)")) // RGB format
```

### Adaptive Colors

Colors that adapt to light/dark terminal themes:

```go
var adaptiveColor = lipgloss.AdaptiveColor{
    Light: "#000000",  // Color for light backgrounds
    Dark:  "#FFFFFF",  // Color for dark backgrounds
}

style := lipgloss.NewStyle().
    Foreground(adaptiveColor)
```

### Complete Color API

```go
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("205")).          // Text color
    Background(lipgloss.Color("235")).          // Background color
    BorderForeground(lipgloss.Color("240")).    // Border color
    UnderlineColor(lipgloss.Color("212")).      // Underline color
    StrikethroughColor(lipgloss.Color("160"))   // Strikethrough color
```

## Text Formatting

### Basic Formatting

```go
style := lipgloss.NewStyle().
    Bold(true).                // Bold text
    Italic(true).              // Italic text
    Underline(true).           // Underlined text
    Strikethrough(true).       // Strikethrough text
    Blink(true).               // Blinking text (rarely supported)
    Faint(true).               // Dim/faint text
    Reverse(true)              // Reverse video (swap fg/bg)
```

### Text Transformation

```go
style := lipgloss.NewStyle().
    Transform(strings.ToUpper)  // Transform text

fmt.Println(style.Render("hello"))  // Output: HELLO
```

### Inline Formatting

```go
// Quick inline formatting without creating styles
bold := lipgloss.NewStyle().Bold(true).Render("Important")
colored := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Pink")
```

## Layout and Spacing

### Width and Height

```go
style := lipgloss.NewStyle().
    Width(50).              // Set fixed width
    MaxWidth(80).           // Set maximum width
    Height(10).             // Set fixed height
    MaxHeight(20)           // Set maximum height

// Get dimensions
w := lipgloss.Width("text")   // Calculate display width
h := lipgloss.Height("text")  // Calculate height (line count)
```

### Padding

Space inside the border/background:

```go
// Uniform padding
style := lipgloss.NewStyle().Padding(1)          // All sides: 1

// Vertical and horizontal
style := lipgloss.NewStyle().Padding(1, 2)       // V: 1, H: 2

// Individual sides
style := lipgloss.NewStyle().Padding(1, 2, 3, 4) // T: 1, R: 2, B: 3, L: 4

// Individual side methods
style := lipgloss.NewStyle().
    PaddingTop(1).
    PaddingRight(2).
    PaddingBottom(1).
    PaddingLeft(2)
```

### Margin

Space outside the border:

```go
// Same syntax as padding
style := lipgloss.NewStyle().Margin(1)           // All sides
style := lipgloss.NewStyle().Margin(1, 2)        // V: 1, H: 2
style := lipgloss.NewStyle().Margin(1, 2, 3, 4)  // T, R, B, L

// Individual methods
style := lipgloss.NewStyle().
    MarginTop(1).
    MarginRight(2).
    MarginBottom(1).
    MarginLeft(2)
```

### Alignment

#### Horizontal Alignment

```go
style := lipgloss.NewStyle().
    Width(50).
    Align(lipgloss.Left)      // Left align (default)

style := lipgloss.NewStyle().
    Width(50).
    Align(lipgloss.Center)    // Center align

style := lipgloss.NewStyle().
    Width(50).
    Align(lipgloss.Right)     // Right align
```

#### Vertical Alignment

```go
style := lipgloss.NewStyle().
    Height(10).
    AlignVertical(lipgloss.Top)     // Top align (default)

style := lipgloss.NewStyle().
    Height(10).
    AlignVertical(lipgloss.Middle)  // Middle align

style := lipgloss.NewStyle().
    Height(10).
    AlignVertical(lipgloss.Bottom)  // Bottom align
```

## Borders

### Border Styles

```go
// Predefined border styles
lipgloss.NormalBorder()       // ┌─┐│└─┘
lipgloss.RoundedBorder()      // ╭─╮│╰─╯
lipgloss.BlockBorder()        // ▛▀▜▌▙▄▟
lipgloss.OuterHalfBlockBorder() // ▛▀▜▐▙▄▟▌
lipgloss.InnerHalfBlockBorder() // ▗▄▖▐▝▀▘▌
lipgloss.ThickBorder()        // ┏━┓┃┗━┛
lipgloss.DoubleBorder()       // ╔═╗║╚═╝
lipgloss.HiddenBorder()       // Invisible border (for spacing)

// Usage
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder())
```

### Custom Borders

```go
// Define custom border characters
customBorder := lipgloss.Border{
    Top:         "─",
    Bottom:      "─",
    Left:        "│",
    Right:       "│",
    TopLeft:     "╔",
    TopRight:    "╗",
    BottomLeft:  "╚",
    BottomRight: "╝",
}

style := lipgloss.NewStyle().
    Border(customBorder)
```

### Selective Borders

```go
// Individual border sides
style := lipgloss.NewStyle().
    Border(lipgloss.NormalBorder()).
    BorderTop(true).
    BorderRight(true).
    BorderBottom(true).
    BorderLeft(false)  // No left border

// Or all at once
style := lipgloss.NewStyle().
    Border(lipgloss.NormalBorder(), true, false, true, false)
    // Top, Right, Bottom, Left
```

### Border Colors

```go
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("240"))

// Individual border colors
style := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderTopForeground(lipgloss.Color("205")).
    BorderRightForeground(lipgloss.Color("205")).
    BorderBottomForeground(lipgloss.Color("240")).
    BorderLeftForeground(lipgloss.Color("240"))
```

## Advanced Features

### Gradients

Create color gradients across text:

```go
// Color gradient
gradient := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#FF00FF")).
    Background(lipgloss.Color("#00FF00"))

// Not directly supported - implement custom gradient
func GradientText(text string, startColor, endColor lipgloss.Color) string {
    // Manual implementation needed
    // Interpolate colors character by character
}
```

### Tables

Create formatted tables with Lip Gloss:

```go
import "github.com/charmbracelet/lipgloss/table"

func renderTable() string {
    t := table.New().
        Border(lipgloss.NormalBorder()).
        BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
        StyleFunc(func(row, col int) lipgloss.Style {
            switch {
            case row == 0:
                return lipgloss.NewStyle().
                    Foreground(lipgloss.Color("229")).
                    Bold(true).
                    Align(lipgloss.Center)
            case row%2 == 0:
                return lipgloss.NewStyle().
                    Foreground(lipgloss.Color("246"))
            default:
                return lipgloss.NewStyle().
                    Foreground(lipgloss.Color("252"))
            }
        }).
        Headers("NAME", "AGE", "LOCATION").
        Row("Alice", "25", "New York").
        Row("Bob", "30", "San Francisco").
        Row("Charlie", "35", "London")

    return t.Render()
}
```

### Lists

Format lists with consistent styling:

```go
import "github.com/charmbracelet/lipgloss/list"

func renderList() string {
    l := list.New().
        Enumerator(list.Bullet).
        EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))).
        ItemStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("252"))).
        Item("First item").
        Item("Second item").
        Item("Third item")

    return l.String()
}

// Available enumerators
list.Bullet     // •
list.Alphabet   // a, b, c...
list.Roman      // i, ii, iii...
list.Arabic     // 1, 2, 3...
```

### Trees

Create tree structures:

```go
import "github.com/charmbracelet/lipgloss/tree"

func renderTree() string {
    t := tree.New().
        Root("Root").
        Item("Child 1").
        Item(
            tree.New().
                Root("Child 2").
                Item("Grandchild 1").
                Item("Grandchild 2"),
        ).
        Item("Child 3")

    return t.String()
}
```

## Position and Place

### Absolute Positioning

Place content at specific coordinates:

```go
// Place string at position
positioned := lipgloss.Place(
    80, 24,        // Canvas width and height
    lipgloss.Center, lipgloss.Middle,  // Horizontal, vertical position
    "Centered Text",
    lipgloss.WithWhitespaceChars("·"),
)
```

### Position Constants

```go
// Horizontal positions
lipgloss.Left
lipgloss.Center
lipgloss.Right

// Vertical positions
lipgloss.Top
lipgloss.Middle
lipgloss.Bottom
```

### PlaceHorizontal and PlaceVertical

```go
// Place horizontally only
text := lipgloss.PlaceHorizontal(50, lipgloss.Center, "Centered")

// Place vertically only
text := lipgloss.PlaceVertical(20, lipgloss.Middle, "Middle")
```

## Joining Content

### JoinHorizontal

Combine strings horizontally:

```go
left := lipgloss.NewStyle().
    Width(30).
    Background(lipgloss.Color("240")).
    Render("Left Panel")

right := lipgloss.NewStyle().
    Width(30).
    Background(lipgloss.Color("235")).
    Render("Right Panel")

// Join at top
combined := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

// Join at center
combined := lipgloss.JoinHorizontal(lipgloss.Middle, left, right)

// Join at bottom
combined := lipgloss.JoinHorizontal(lipgloss.Bottom, left, right)
```

### JoinVertical

Combine strings vertically:

```go
header := lipgloss.NewStyle().
    Width(60).
    Background(lipgloss.Color("62")).
    Render("Header")

content := lipgloss.NewStyle().
    Width(60).
    Height(10).
    Background(lipgloss.Color("240")).
    Render("Content")

// Join at left
combined := lipgloss.JoinVertical(lipgloss.Left, header, content)

// Join at center
combined := lipgloss.JoinVertical(lipgloss.Center, header, content)

// Join at right
combined := lipgloss.JoinVertical(lipgloss.Right, header, content)
```

## Style Composition

### Inherit and Override

```go
// Base style
baseStyle := lipgloss.NewStyle().
    Padding(1).
    Margin(1).
    Border(lipgloss.RoundedBorder())

// Inherit and modify
successStyle := baseStyle.Copy().
    Foreground(lipgloss.Color("42")).
    BorderForeground(lipgloss.Color("42"))

errorStyle := baseStyle.Copy().
    Foreground(lipgloss.Color("196")).
    BorderForeground(lipgloss.Color("196"))

warningStyle := baseStyle.Copy().
    Foreground(lipgloss.Color("214")).
    BorderForeground(lipgloss.Color("214"))
```

### Style Inheritance

```go
// Create style hierarchy
var (
    // Base theme
    primaryColor   = lipgloss.Color("205")
    secondaryColor = lipgloss.Color("240")

    // Base components
    boxStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1)

    // Specific components
    titleBox = boxStyle.Copy().
        BorderForeground(primaryColor).
        Bold(true)

    contentBox = boxStyle.Copy().
        BorderForeground(secondaryColor)
)
```

## Whitespace and Special Characters

### Custom Whitespace

```go
style := lipgloss.NewStyle().
    Width(50).
    WhitespaceChars("·")  // Use · instead of space

// Useful for debugging layout
debugStyle := lipgloss.NewStyle().
    Width(40).
    Height(10).
    Background(lipgloss.Color("240")).
    WhitespaceChars("·")
```

### Whitespace Options

```go
// Apply background color to whitespace
style := lipgloss.NewStyle().
    Width(50).
    Background(lipgloss.Color("240")).
    WhitespaceForeground(lipgloss.Color("235"))
```

## Inline Content

### Inline Mode

Render without width constraints:

```go
style := lipgloss.NewStyle().
    Inline(true).
    Bold(true).
    Foreground(lipgloss.Color("205"))

// Will not pad to width
text := style.Render("Short")
```

## Performance Optimization

### Style Reuse

```go
// Good: Reuse styles
var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

func render() string {
    return titleStyle.Render("Title 1") + "\n" +
           titleStyle.Render("Title 2")
}

// Bad: Create styles repeatedly
func renderBad() string {
    return lipgloss.NewStyle().Bold(true).Render("Title 1") + "\n" +
           lipgloss.NewStyle().Bold(true).Render("Title 2")
}
```

### Caching Rendered Content

```go
type Component struct {
    content      string
    cachedRender string
    dirty        bool
    style        lipgloss.Style
}

func (c *Component) Render() string {
    if c.dirty {
        c.cachedRender = c.style.Render(c.content)
        c.dirty = false
    }
    return c.cachedRender
}
```

### Width and Height Calculation

```go
// Efficient width calculation
width := lipgloss.Width(str)

// Pre-calculate for layout
maxWidth := 0
for _, line := range lines {
    w := lipgloss.Width(line)
    if w > maxWidth {
        maxWidth = w
    }
}
```

## Complete Examples

### Example 1: Styled Panel

```go
func StyledPanel(title, content string) string {
    titleStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("63")).
        Padding(0, 1).
        Width(60)

    contentStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("252")).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("63")).
        Padding(1).
        Width(60)

    return lipgloss.JoinVertical(
        lipgloss.Left,
        titleStyle.Render(title),
        contentStyle.Render(content),
    )
}
```

### Example 2: Two-Column Layout

```go
func TwoColumnLayout(left, right string, width int) string {
    colWidth := (width - 4) / 2

    leftStyle := lipgloss.NewStyle().
        Width(colWidth).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("63")).
        Padding(1)

    rightStyle := lipgloss.NewStyle().
        Width(colWidth).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("212")).
        Padding(1)

    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        leftStyle.Render(left),
        " ",
        rightStyle.Render(right),
    )
}
```

### Example 3: Status Messages

```go
var (
    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42")).
        Bold(true).
        Padding(0, 1).
        Render

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("196")).
        Bold(true).
        Padding(0, 1).
        Render

    infoStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("33")).
        Bold(true).
        Padding(0, 1).
        Render
)

func StatusMessage(msgType, text string) string {
    switch msgType {
    case "success":
        return successStyle("✓ " + text)
    case "error":
        return errorStyle("✗ " + text)
    case "info":
        return infoStyle("ℹ " + text)
    default:
        return text
    }
}
```

### Example 4: Dashboard Layout

```go
func Dashboard(width, height int) string {
    // Header
    headerStyle := lipgloss.NewStyle().
        Width(width).
        Bold(true).
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("63")).
        Padding(0, 2).
        Align(lipgloss.Center)

    header := headerStyle.Render("Application Dashboard")

    // Left sidebar
    sidebarWidth := width / 4
    sidebarStyle := lipgloss.NewStyle().
        Width(sidebarWidth).
        Height(height - 5).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("240")).
        Padding(1)

    sidebar := sidebarStyle.Render(
        "Menu\n\n" +
        "• Dashboard\n" +
        "• Settings\n" +
        "• Help\n" +
        "• Quit",
    )

    // Main content
    contentWidth := width - sidebarWidth - 6
    contentStyle := lipgloss.NewStyle().
        Width(contentWidth).
        Height(height - 5).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("240")).
        Padding(1)

    content := contentStyle.Render("Main Content Area")

    // Combine
    body := lipgloss.JoinHorizontal(
        lipgloss.Top,
        sidebar,
        " ",
        content,
    )

    // Footer
    footerStyle := lipgloss.NewStyle().
        Width(width).
        Foreground(lipgloss.Color("240")).
        Padding(0, 2)

    footer := footerStyle.Render("Press 'q' to quit")

    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        "\n",
        body,
        "\n",
        footer,
    )
}
```

### Example 5: Progress Indicator

```go
func ProgressBar(percent float64, width int) string {
    filled := int(percent * float64(width))
    empty := width - filled

    bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)

    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("205"))

    label := fmt.Sprintf("%.0f%%", percent*100)

    return fmt.Sprintf("%s %s",
        style.Render(bar),
        label,
    )
}
```

## Common Patterns

### Theme System

```go
type Theme struct {
    Primary   lipgloss.Color
    Secondary lipgloss.Color
    Success   lipgloss.Color
    Error     lipgloss.Color
    Warning   lipgloss.Color
    Info      lipgloss.Color
}

var DefaultTheme = Theme{
    Primary:   lipgloss.Color("205"),
    Secondary: lipgloss.Color("240"),
    Success:   lipgloss.Color("42"),
    Error:     lipgloss.Color("196"),
    Warning:   lipgloss.Color("214"),
    Info:      lipgloss.Color("33"),
}

func (t Theme) ButtonStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("229")).
        Background(t.Primary).
        Padding(0, 2).
        Bold(true)
}
```

### Responsive Layouts

```go
func ResponsiveLayout(content string, width int) string {
    var columns int
    switch {
    case width < 80:
        columns = 1
    case width < 120:
        columns = 2
    default:
        columns = 3
    }

    // Adjust layout based on columns
    colWidth := (width - (columns + 1)) / columns

    // Build layout...
}
```

## Best Practices

1. **Reuse Styles**: Define styles once, reuse throughout application
2. **Composition Over Duplication**: Use `Copy()` to create variants
3. **Responsive Design**: Handle different terminal sizes gracefully
4. **Color Accessibility**: Test colors in both light and dark terminals
5. **Performance**: Cache rendered content when possible
6. **Consistent Spacing**: Use consistent padding/margin values
7. **Border Consistency**: Use same border style for related components
8. **Width Management**: Always handle width constraints properly

## Integration with Bubble Tea

```go
type model struct {
    width  int
    height int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

func (m model) View() string {
    style := lipgloss.NewStyle().
        Width(m.width).
        Height(m.height).
        Align(lipgloss.Center).
        AlignVertical(lipgloss.Middle).
        Border(lipgloss.RoundedBorder())

    return style.Render("Responsive Content")
}
```

## Resources

- **Official Repository:** https://github.com/charmbracelet/lipgloss
- **Package Documentation:** https://pkg.go.dev/github.com/charmbracelet/lipgloss
- **Examples:** https://github.com/charmbracelet/lipgloss/tree/master/examples
- **Gallery:** https://github.com/charmbracelet/lipgloss/blob/master/GALLERY.md
- **Color Reference:** https://www.ditig.com/256-colors-cheat-sheet
- **Related:** [01-bubbletea.md](01-bubbletea.md), [02-bubbles.md](02-bubbles.md)