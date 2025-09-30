# Bubbles - Component Library Reference

## Overview

Bubbles is a collection of pre-built, customizable TUI components for Bubble Tea applications. Each component follows the Bubble Tea architecture pattern (Init, Update, View) and can be easily integrated into your applications.

**Repository:** https://github.com/charmbracelet/bubbles
**Package:** `github.com/charmbracelet/bubbles`

## Installation

```bash
go get github.com/charmbracelet/bubbles
```

## Component Catalog

### 1. Text Input - Single-Line Input Field

**Import:** `github.com/charmbracelet/bubbles/textinput`

A single-line text input with cursor, unicode support, and scrolling for long text.

#### Features
- Cursor positioning and navigation
- Character masking (for passwords)
- Unicode support
- Input validation
- Placeholder text
- Character limit
- Focus management

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/textinput"

type model struct {
    textInput textinput.Model
}

func initialModel() model {
    ti := textinput.New()
    ti.Placeholder = "Enter your name"
    ti.Focus()
    ti.CharLimit = 50
    ti.Width = 20

    return model{
        textInput: ti,
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.textInput, cmd = m.textInput.Update(msg)
    return m, cmd
}

func (m model) View() string {
    return fmt.Sprintf(
        "What's your name?\n\n%s\n\n(Press Enter to submit)",
        m.textInput.View(),
    )
}
```

#### Advanced Configuration

```go
ti := textinput.New()

// Appearance
ti.Placeholder = "Email address"
ti.Width = 30
ti.Prompt = "> "
ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

// Behavior
ti.CharLimit = 100
ti.Focus()  // or ti.Blur()
ti.EchoMode = textinput.EchoPassword  // Mask input
ti.EchoCharacter = '•'

// Validation
ti.Validate = func(s string) error {
    if !strings.Contains(s, "@") {
        return errors.New("invalid email")
    }
    return nil
}

// Get value
value := ti.Value()
```

#### Styling

```go
import "github.com/charmbracelet/lipgloss"

ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
```

---

### 2. Text Area - Multi-Line Text Editor

**Import:** `github.com/charmbracelet/bubbles/textarea`

Multi-line text input with vertical scrolling, line numbers, and editor-like functionality.

#### Features
- Multi-line editing
- Vertical scrolling
- Line numbers (optional)
- Cursor navigation (arrow keys, page up/down)
- Character limit
- Placeholder text
- Auto-height or fixed height

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/textarea"

type model struct {
    textarea textarea.Model
}

func initialModel() model {
    ta := textarea.New()
    ta.Placeholder = "Enter your description..."
    ta.Focus()

    return model{textarea: ta}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.textarea, cmd = m.textarea.Update(msg)
    return m, cmd
}

func (m model) View() string {
    return fmt.Sprintf(
        "Enter description:\n\n%s\n\n%d chars",
        m.textarea.View(),
        len(m.textarea.Value()),
    )
}
```

#### Configuration

```go
ta := textarea.New()
ta.SetWidth(60)
ta.SetHeight(10)
ta.CharLimit = 500
ta.ShowLineNumbers = true
ta.Placeholder = "Type your message here..."

// Get content
content := ta.Value()
lines := ta.Line() // Current line number
```

---

### 3. Viewport - Scrollable Content Container

**Import:** `github.com/charmbracelet/bubbles/viewport`

A scrollable viewport for displaying large content with vertical scrolling.

#### Features
- Vertical scrolling
- Mouse wheel support
- High-performance rendering
- Percentage-based scrolling
- Auto-fit content

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/viewport"

type model struct {
    viewport viewport.Model
    content  string
    ready    bool
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        if !m.ready {
            m.viewport = viewport.New(msg.Width, msg.Height-2)
            m.viewport.SetContent(m.content)
            m.ready = true
        } else {
            m.viewport.Width = msg.Width
            m.viewport.Height = msg.Height - 2
        }
    }

    var cmd tea.Cmd
    m.viewport, cmd = m.viewport.Update(msg)
    return m, cmd
}

func (m model) View() string {
    if !m.ready {
        return "Initializing..."
    }
    return m.viewport.View() + "\n" + m.viewportInfo()
}

func (m model) viewportInfo() string {
    return fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
}
```

#### Methods

```go
vp := viewport.New(width, height)

// Content management
vp.SetContent(content)
vp.SetYOffset(offset)

// Scrolling
vp.LineUp(n)
vp.LineDown(n)
vp.HalfViewUp()
vp.HalfViewDown()
vp.GotoTop()
vp.GotoBottom()

// Information
vp.ScrollPercent()  // 0.0 - 1.0
vp.AtTop()          // bool
vp.AtBottom()       // bool
vp.PastBottom()     // bool
```

---

### 4. List - Selectable Item List

**Import:** `github.com/charmbracelet/bubbles/list`

A feature-rich list component with filtering, pagination, and help text.

#### Features
- Item selection
- Filtering/search
- Pagination
- Customizable rendering
- Built-in help
- Status bar

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/list"

type item struct {
    title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
    list list.Model
}

func initialModel() model {
    items := []list.Item{
        item{title: "Raspberry Pi", desc: "A small computer"},
        item{title: "Arduino", desc: "A microcontroller"},
        item{title: "ESP32", desc: "WiFi + Bluetooth chip"},
    }

    l := list.New(items, list.NewDefaultDelegate(), 0, 0)
    l.Title = "Hardware Projects"

    return model{list: l}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.list.SetSize(msg.Width, msg.Height)

    case tea.KeyMsg:
        if msg.String() == "enter" {
            selected := m.list.SelectedItem().(item)
            // Handle selection
        }
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}

func (m model) View() string {
    return m.list.View()
}
```

#### Custom Delegate (Advanced)

```go
type itemDelegate struct{}

func (d itemDelegate) Height() int                           { return 2 }
func (d itemDelegate) Spacing() int                          { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
    i, ok := item.(item)
    if !ok {
        return
    }

    str := fmt.Sprintf("%d. %s - %s", index+1, i.title, i.desc)

    // Highlight selected item
    if index == m.Index() {
        str = selectedStyle.Render("> " + str)
    }

    fmt.Fprint(w, str)
}
```

---

### 5. Table - Tabular Data Display

**Import:** `github.com/charmbracelet/bubbles/table`

Display and navigate tabular data with columns and rows.

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/table"

type model struct {
    table table.Model
}

func initialModel() model {
    columns := []table.Column{
        {Title: "Name", Width: 20},
        {Title: "Age", Width: 10},
        {Title: "Email", Width: 30},
    }

    rows := []table.Row{
        {"Alice", "25", "alice@example.com"},
        {"Bob", "30", "bob@example.com"},
        {"Charlie", "35", "charlie@example.com"},
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(7),
    )

    s := table.DefaultStyles()
    s.Header = s.Header.
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        BorderBottom(true).
        Bold(false)
    s.Selected = s.Selected.
        Foreground(lipgloss.Color("229")).
        Background(lipgloss.Color("57")).
        Bold(false)
    t.SetStyles(s)

    return model{table: t}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

func (m model) View() string {
    return m.table.View()
}
```

---

### 6. Spinner - Loading Indicator

**Import:** `github.com/charmbracelet/bubbles/spinner`

Animated spinner for indicating ongoing operations.

#### Built-in Spinner Styles

```go
spinner.Line
spinner.Dot
spinner.MiniDot
spinner.Jump
spinner.Pulse
spinner.Points
spinner.Globe
spinner.Moon
spinner.Monkey
// ... and many more
```

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/spinner"

type model struct {
    spinner spinner.Model
}

func initialModel() model {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
    return model{spinner: s}
}

func (m model) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.spinner, cmd = m.spinner.Update(msg)
    return m, cmd
}

func (m model) View() string {
    return fmt.Sprintf("%s Loading...", m.spinner.View())
}
```

---

### 7. Progress - Progress Bar

**Import:** `github.com/charmbracelet/bubbles/progress`

Visual progress indicator with customizable appearance.

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/progress"

type model struct {
    progress progress.Model
    percent  float64
}

func initialModel() model {
    return model{
        progress: progress.New(progress.WithDefaultGradient()),
        percent:  0.0,
    }
}

func (m model) Init() tea.Cmd {
    return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tickMsg:
        if m.percent >= 1.0 {
            return m, tea.Quit
        }

        m.percent += 0.05
        return m, tea.Batch(
            m.progress.SetPercent(m.percent),
            tickCmd(),
        )
    }

    return m, nil
}

func (m model) View() string {
    return "\n" +
        m.progress.View() + "\n\n" +
        fmt.Sprintf("%.0f%%", m.percent*100)
}

func tickCmd() tea.Cmd {
    return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}
```

#### Styling

```go
prog := progress.New(
    progress.WithDefaultGradient(),
    progress.WithWidth(50),
    progress.WithoutPercentage(),
)

// Custom colors
prog := progress.New(
    progress.WithGradient("#FF0000", "#00FF00"),
)
```

---

### 8. Paginator - Pagination Control

**Import:** `github.com/charmbracelet/bubbles/paginator`

Handle pagination logic and UI for large datasets.

#### Basic Usage

```go
import "github.com/charmbracelet/bubbles/paginator"

type model struct {
    paginator paginator.Model
    items     []string
}

func initialModel() model {
    items := make([]string, 100)
    for i := range items {
        items[i] = fmt.Sprintf("Item %d", i+1)
    }

    p := paginator.New()
    p.Type = paginator.Dots
    p.PerPage = 10
    p.SetTotalPages(len(items))

    return model{
        paginator: p,
        items:     items,
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.paginator, cmd = m.paginator.Update(msg)
    return m, cmd
}

func (m model) View() string {
    start, end := m.paginator.GetSliceBounds(len(m.items))
    items := m.items[start:end]

    var s string
    for _, item := range items {
        s += item + "\n"
    }

    s += "\n" + m.paginator.View()
    return s
}
```

#### Paginator Types

```go
paginator.Arabic    // 1 2 3 4 5
paginator.Dots      // • • • • •
```

---

### 9. Timer & Stopwatch

**Import:**
- `github.com/charmbracelet/bubbles/timer`
- `github.com/charmbracelet/bubbles/stopwatch`

#### Timer (Countdown)

```go
import "github.com/charmbracelet/bubbles/timer"

type model struct {
    timer timer.Model
}

func initialModel() model {
    return model{
        timer: timer.NewWithInterval(5*time.Minute, time.Second),
    }
}

func (m model) Init() tea.Cmd {
    return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case timer.TickMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd

    case timer.TimeoutMsg:
        return m, tea.Quit
    }

    return m, nil
}

func (m model) View() string {
    return m.timer.View()
}
```

#### Stopwatch (Count Up)

```go
import "github.com/charmbracelet/bubbles/stopwatch"

type model struct {
    stopwatch stopwatch.Model
}

func initialModel() model {
    return model{
        stopwatch: stopwatch.NewWithInterval(time.Second),
    }
}

func (m model) Init() tea.Cmd {
    return m.stopwatch.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.stopwatch, cmd = m.stopwatch.Update(msg)
    return m, cmd
}

func (m model) View() string {
    return fmt.Sprintf("Elapsed: %s", m.stopwatch.View())
}
```

---

### 10. Help - Keyboard Shortcut Display

**Import:** `github.com/charmbracelet/bubbles/help`

Automatic help text generation from key bindings.

#### Basic Usage

```go
import (
    "github.com/charmbracelet/bubbles/help"
    "github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
    Up    key.Binding
    Down  key.Binding
    Enter key.Binding
    Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down},
        {k.Enter, k.Quit},
    }
}

var keys = keyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "move down"),
    ),
    Enter: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "select"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
}

type model struct {
    help help.Model
    keys keyMap
}

func (m model) View() string {
    return "\n" + m.help.View(m.keys)
}
```

---

### 11. File Picker

**Import:** `github.com/charmbracelet/bubbles/filepicker`

Navigate filesystem and select files.

```go
import "github.com/charmbracelet/bubbles/filepicker"

type model struct {
    filepicker filepicker.Model
}

func initialModel() model {
    fp := filepicker.New()
    fp.AllowedTypes = []string{".go", ".txt", ".md"}
    fp.CurrentDirectory = "."

    return model{filepicker: fp}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.filepicker, cmd = m.filepicker.Update(msg)

    if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
        // Handle file selection
        return m, tea.Quit
    }

    return m, cmd
}
```

## Component Composition Best Practices

### 1. Delegate Message Handling

```go
type model struct {
    input1 textinput.Model
    input2 textinput.Model
    focused int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    // Route to focused component
    if m.focused == 0 {
        var cmd tea.Cmd
        m.input1, cmd = m.input1.Update(msg)
        cmds = append(cmds, cmd)
    } else {
        var cmd tea.Cmd
        m.input2, cmd = m.input2.Update(msg)
        cmds = append(cmds, cmd)
    }

    return m, tea.Batch(cmds...)
}
```

### 2. Synchronize Component Sizes

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.viewport.Width = msg.Width
        m.viewport.Height = msg.Height - 5
        m.list.SetSize(msg.Width, msg.Height-5)
    }
    return m, nil
}
```

## Resources

- **Repository:** https://github.com/charmbracelet/bubbles
- **Package Docs:** https://pkg.go.dev/github.com/charmbracelet/bubbles
- **Examples:** https://github.com/charmbracelet/bubbles/tree/master/examples