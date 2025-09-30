# Bubble Tea Framework - Comprehensive Guide

## Overview

Bubble Tea is a powerful framework for building terminal user interfaces (TUIs) in Go. It's based on The Elm Architecture, a functional design pattern that emphasizes simplicity, testability, and maintainability.

**Repository:** https://github.com/charmbracelet/bubbletea
**Package:** `github.com/charmbracelet/bubbletea`
**Import:** `tea "github.com/charmbracelet/bubbletea"`

## Core Architecture

### The Elm Architecture (MVU Pattern)

Bubble Tea implements the Model-View-Update (MVU) pattern:

```
┌─────────────┐
│    Init     │  Initialize application state
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Update    │◄─── Messages (Events)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│    View     │  Render UI
└─────────────┘
```

### Three Core Concepts

#### 1. Model
The **Model** represents your application's state. It can be any Go type (struct, int, string, etc.).

```go
type model struct {
    cursor   int
    choices  []string
    selected map[int]struct{}
}
```

#### 2. Update
The **Update** function handles incoming messages and updates the model accordingly. It returns an updated model and optionally a command to run.

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "up":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        }
    }
    return m, nil
}
```

#### 3. View
The **View** function renders your UI as a string based on the current model state.

```go
func (m model) View() string {
    s := "What should we buy?\n\n"

    for i, choice := range m.choices {
        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }

        checked := " "
        if _, ok := m.selected[i]; ok {
            checked = "x"
        }

        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    s += "\nPress q to quit.\n"
    return s
}
```

## Program Lifecycle

### 1. Initialization

```go
func (m model) Init() tea.Cmd {
    // Return a command to run on startup
    // Return nil if no initial command needed
    return nil
}
```

### 2. Creating and Running a Program

```go
func main() {
    p := tea.NewProgram(
        initialModel(),
        tea.WithAltScreen(),        // Use full terminal window
        tea.WithMouseCellMotion(),  // Enable mouse support
    )

    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
```

### 3. Program Options

```go
// Terminal modes
tea.WithAltScreen()              // Use alternate screen buffer (full window)
tea.WithoutRenderer()            // Disable rendering (for testing)

// Input handling
tea.WithMouseCellMotion()        // Track mouse motion with cells
tea.WithMouseAllMotion()         // Track all mouse motion
tea.WithoutSignalHandler()       // Disable default signal handling

// Context and lifecycle
tea.WithContext(ctx)             // Provide context for cancellation
tea.WithInput(reader)            // Custom input source
tea.WithOutput(writer)           // Custom output destination

// Focus and visibility
tea.WithReportFocus()            // Report terminal focus/blur events
tea.WithoutBracketedPaste()      // Disable bracketed paste mode
```

## Messages and Events

### Built-in Message Types

#### KeyMsg - Keyboard Input
```go
type KeyMsg struct {
    Type  KeyType
    Runes []rune
    Alt   bool
}

// Usage in Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "enter":
            return m, m.submitForm()
        case "up", "k":
            m.cursor--
        case "down", "j":
            m.cursor++
        }
    }
    return m, nil
}
```

#### MouseMsg - Mouse Events
```go
type MouseMsg struct {
    X, Y   int
    Type   MouseEventType
    Button MouseButton
}

// Mouse event types
MouseLeft     // Left button
MouseRight    // Right button
MouseMiddle   // Middle button
MouseRelease  // Button release
MouseWheelUp  // Scroll up
MouseWheelDown // Scroll down
MouseMotion   // Mouse moved

// Usage
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseMsg:
        if msg.Type == tea.MouseLeft {
            m.clickX, m.clickY = msg.X, msg.Y
        }
    }
    return m, nil
}
```

#### WindowSizeMsg - Terminal Resize
```go
type WindowSizeMsg struct {
    Width  int
    Height int
}

// Usage
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Recalculate layouts
    }
    return m, nil
}
```

#### Focus Messages
```go
type FocusMsg struct{}  // Terminal gained focus
type BlurMsg struct{}   // Terminal lost focus

// Usage (requires tea.WithReportFocus() option)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.FocusMsg:
        m.focused = true
        return m, m.fetchLatestData()
    case tea.BlurMsg:
        m.focused = false
    }
    return m, nil
}
```

### Custom Messages

Define your own message types for application-specific events:

```go
type tickMsg time.Time
type apiResponseMsg struct {
    data string
    err  error
}
type progressMsg float64
```

## Commands (Cmd)

Commands are functions that perform I/O and return messages. They enable side effects in the otherwise pure functional architecture.

### Command Type

```go
type Cmd func() Msg
```

### Built-in Commands

#### Quit - Exit the program
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if shouldQuit {
        return m, tea.Quit
    }
    return m, nil
}
```

#### Batch - Run multiple commands concurrently
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, tea.Batch(
        fetchUserData,
        fetchSettings,
        startTimer,
    )
}
```

#### Sequence - Run commands sequentially
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, tea.Sequence(
        authenticate,
        fetchProfile,
        loadDashboard,
    )
}
```

#### Tick - Time-based messages
```go
// Send a message after a duration
func (m model) Init() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
```

#### Every - Clock-synced periodic messages
```go
// Send messages every N duration, synced to system clock
func (m model) Init() tea.Cmd {
    return tea.Every(time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
```

### Custom Commands

Create commands that perform I/O operations:

```go
// HTTP request command
func fetchDataCmd(url string) tea.Cmd {
    return func() tea.Msg {
        resp, err := http.Get(url)
        if err != nil {
            return apiResponseMsg{err: err}
        }
        defer resp.Body.Close()

        data, err := io.ReadAll(resp.Body)
        if err != nil {
            return apiResponseMsg{err: err}
        }

        return apiResponseMsg{data: string(data)}
    }
}

// File reading command
func readFileCmd(path string) tea.Cmd {
    return func() tea.Msg {
        data, err := os.ReadFile(path)
        return fileReadMsg{data: data, err: err}
    }
}
```

## Advanced Patterns

### 1. Subscriptions

Long-running processes that continuously send messages:

```go
type model struct {
    ticker *time.Ticker
    done   chan bool
}

func (m model) Init() tea.Cmd {
    m.ticker = time.NewTicker(time.Second)
    return listenForTicks(m.ticker.C, m.done)
}

func listenForTicks(tickChan <-chan time.Time, done <-chan bool) tea.Cmd {
    return func() tea.Msg {
        select {
        case t := <-tickChan:
            return tickMsg(t)
        case <-done:
            return nil
        }
    }
}
```

### 2. Child Components (Composition)

Delegate messages to child components:

```go
type model struct {
    textInput textinput.Model
    spinner   spinner.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    // Update child components
    var cmd tea.Cmd
    m.textInput, cmd = m.textInput.Update(msg)
    cmds = append(cmds, cmd)

    m.spinner, cmd = m.spinner.Update(msg)
    cmds = append(cmds, cmd)

    // Handle own messages
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            return m, m.submit()
        }
    }

    return m, tea.Batch(cmds...)
}

func (m model) View() string {
    return fmt.Sprintf(
        "%s\n\n%s\n",
        m.textInput.View(),
        m.spinner.View(),
    )
}
```

### 3. External I/O with Send()

Send messages from outside the Update loop:

```go
func main() {
    p := tea.NewProgram(initialModel())

    // Send messages from another goroutine
    go func() {
        for event := range eventChannel {
            p.Send(externalEventMsg(event))
        }
    }()

    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### 4. Graceful Shutdown

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "ctrl+c" {
            return m, tea.Sequence(
                m.saveState(),
                m.cleanup(),
                tea.Quit,
            )
        }
    }
    return m, nil
}
```

## Program Control

### Send Messages Programmatically

```go
p := tea.NewProgram(model{})

// From another goroutine
go func() {
    p.Send(customMsg{})
}()

p.Run()
```

### Kill Program

```go
p.Kill()  // Forcefully terminate
```

### Suspend and Resume

```go
// Suspend (e.g., for running external editor)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if shouldSuspend {
        return m, tea.Suspend
    }
    return m, nil
}
```

### Execute External Commands

```go
func openEditorCmd(file string) tea.Cmd {
    c := exec.Command("vim", file)
    return tea.ExecProcess(c, func(err error) tea.Msg {
        return editorFinishedMsg{err: err}
    })
}
```

## Printing Outside TUI

Sometimes you need to print before TUI starts or after it exits:

```go
func main() {
    fmt.Println("Starting application...")

    p := tea.NewProgram(model{})
    finalModel, err := p.Run()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Final state: %+v\n", finalModel)
}

// Or use tea.Printf during execution (thread-safe)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, func() tea.Msg {
        tea.Printf("Log message: %s\n", msg)
        return nil
    }
}
```

## Performance Optimization

### 1. Efficient Rendering

```go
// Only render what changed
func (m model) View() string {
    if !m.dirty {
        return m.cachedView
    }

    m.cachedView = m.render()
    m.dirty = false
    return m.cachedView
}
```

### 2. Debouncing

```go
type model struct {
    lastKeyTime time.Time
    debounce    time.Duration
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        now := time.Now()
        if now.Sub(m.lastKeyTime) < m.debounce {
            return m, nil
        }
        m.lastKeyTime = now
        // Process key
    }
    return m, nil
}
```

### 3. Lazy Loading

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case loadMoreMsg:
        if !m.loading && m.hasMore {
            m.loading = true
            return m, fetchNextPage(m.page)
        }
    }
    return m, nil
}
```

## Testing

### Unit Testing Models

```go
func TestUpdate(t *testing.T) {
    m := model{cursor: 0}

    // Simulate key press
    newModel, _ := m.Update(tea.KeyMsg{
        Type:  tea.KeyRunes,
        Runes: []rune{'j'},
    })

    m = newModel.(model)
    if m.cursor != 1 {
        t.Errorf("Expected cursor at 1, got %d", m.cursor)
    }
}
```

### Testing with WithoutRenderer

```go
func TestProgram(t *testing.T) {
    m := model{}
    p := tea.NewProgram(m, tea.WithoutRenderer())

    go func() {
        time.Sleep(100 * time.Millisecond)
        p.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
    }()

    finalModel, err := p.Run()
    if err != nil {
        t.Fatal(err)
    }

    // Assert on finalModel
}
```

## Best Practices

1. **Keep Models Immutable** - Always return new model instances
2. **Separate Concerns** - Use child components for complex UIs
3. **Handle All Message Types** - Use type switches with default cases
4. **Commands for I/O** - Never do I/O in Update(), always use commands
5. **Responsive to Resize** - Always handle WindowSizeMsg
6. **Graceful Errors** - Show user-friendly error messages, don't panic
7. **Fast View()** - Keep rendering logic simple and fast
8. **Test Without Renderer** - Use tea.WithoutRenderer() for automated tests

## Common Patterns

### Loading States
```go
type model struct {
    loading bool
    data    string
    err     error
}

func (m model) View() string {
    if m.loading {
        return "Loading..."
    }
    if m.err != nil {
        return fmt.Sprintf("Error: %v", m.err)
    }
    return m.data
}
```

### Multi-Screen Navigation
```go
type screen int

const (
    menuScreen screen = iota
    settingsScreen
    helpScreen
)

type model struct {
    currentScreen screen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok {
        switch key.String() {
        case "1":
            m.currentScreen = menuScreen
        case "2":
            m.currentScreen = settingsScreen
        }
    }
    return m, nil
}

func (m model) View() string {
    switch m.currentScreen {
    case menuScreen:
        return m.renderMenu()
    case settingsScreen:
        return m.renderSettings()
    default:
        return ""
    }
}
```

## Resources

- **Official Repository:** https://github.com/charmbracelet/bubbletea
- **Package Documentation:** https://pkg.go.dev/github.com/charmbracelet/bubbletea
- **Tutorials:** https://github.com/charmbracelet/bubbletea/tree/master/tutorials
- **Examples:** https://github.com/charmbracelet/bubbletea/tree/master/examples
- **Community:** GitHub Discussions