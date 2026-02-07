# Building Effective Bubble Tea TUI Programs

## 📝 Introduction
Bubble Tea is a TUI (Terminal User Interface) framework for Go based on the Elm Architecture. While powerful, it has a learning curve. These notes cover essential tips for building robust, maintainable TUIs.

## 1. Keep the Event Loop Fast ⚡

### Intuition
Bubble Tea processes messages in a single-threaded event loop. Slow `Update()` or `View()` methods cause lag and unresponsive UI.

### Key Points
- **Never block** in `Update()` or `View()`
- **Offload expensive operations** to `tea.Cmd`
- Messages can **back up** if processing is slow

### Code Example
```go
// ❌ DON'T: Block in Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        time.Sleep(time.Minute) // Blocks event loop!
        return m, nil
    }
    return m, nil
}

// ✅ DO: Offload to command
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        return m, func() tea.Msg {
            time.Sleep(time.Minute) // Runs in goroutine
            return operationDoneMsg{}
        }
    }
    return m, nil
}
```

## 2. Debug by Dumping Messages to a File 🔍

### Intuition
Understanding the message flow is crucial for debugging. Messages arrive from multiple sources and seeing them helps understand program state.

### Implementation
```go
import "github.com/davecgh/go-spew/spew"

type model struct {
    dump io.Writer
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if m.dump != nil {
        spew.Fdump(m.dump, msg) // Pretty-print to file
    }
    // ... rest of Update logic
}

func main() {
    var dump *os.File
    if _, ok := os.LookupEnv("DEBUG"); ok {
        dump, _ = os.OpenFile("messages.log", 
            os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
    }
    p := tea.NewProgram(model{dump: dump})
    // ...
}
```
**Usage:** `tail -f messages.log` in another terminal

## 3. Live Reload Code Changes 🔄

### Intuition
Quick feedback loop speeds up development. Recompile and restart on file changes.

### Simple Implementation
```bash
# build_and_watch.sh
while true; do
    go build -o app && pkill -f 'app'
    inotifywait -e attrib $(find . -name '*.go')
done

# run_forever.sh
while true; do
    ./app
done
```

### Better Tools
- **watchexec**: Works with TTY programs
- **air**: Live reload for Go (may need TTY workarounds)

## 4. Receiver Methods: Value vs Pointer 🎯

### Intuition
Bubble Tea follows functional patterns, but Go allows both value and pointer receivers.

### Key Insights
- **Documentation examples use value receivers** (functional style)
- **Pointer receivers can modify state** but risk race conditions
- **Stick to message flow**: Changes in `Update()`, not in goroutines

### Code Examples
```go
// ✅ Acceptable pointer receiver usage
func (m *model) helperMethod() {
    m.someField = "updated" // Called from within Update()
}

// ❌ Risky: Modifying outside event loop
func (m *model) Init() tea.Cmd {
    go func() {
        m.content = "changed" // RACE CONDITION!
    }()
    return nil
}
```

### Recommendation
- Use **value receivers** for main model methods
- Use **pointer receivers** for helper methods called from `Update()`
- **Never modify model** in goroutines without synchronization

## 5. Message Ordering is Not Guaranteed 🎲

### Intuition
Concurrent message sources mean messages may arrive out of order.

### Sources of Messages
1. **User input** (ordered - single goroutine)
2. **Commands** (unordered - concurrent goroutines)
3. **Programmatic sends** (unordered)
4. **System signals** (ordered)

### Demonstration
```go
type nMsg int

func main() {
    p := tea.NewProgram(model{})
    
    // Concurrent sends - order is unpredictable
    for i := range 10 {
        go p.Send(nMsg(i))
    }
    // Output might be: [0 1 9 8 5 6 4 2 3 7]
}
```

### Solutions for Ordered Processing
```go
// Option 1: Process in Update() directly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        m.items = append(m.items, getNextItem()) // Direct update
        return m, nil // No command
    }
    return m, nil
}

// Option 2: Use tea.Sequence for ordered commands
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, tea.Sequence(
        firstCommand,
        secondCommand,
        thirdCommand,
    )
}

// Option 3: Embrace concurrency - redesign if order isn't critical
```

## 6. Build a Tree of Models 🌳

### Intuition
Complex TUIs need hierarchical organization. Parent models route messages and compose views.

### Architecture Pattern
```
Root Model
├── Header Model
├── Content Area Model
│   ├── List Model
│   ├── Detail Model
│   └── Form Model
└── Footer Model
```

### Implementation Strategy
```go
type rootModel struct {
    children map[string]tea.Model
    current  string // ID of active child
    history  []string // Navigation stack
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Route messages appropriately
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // Send to all children
        for id, child := range m.children {
            updated, _ := child.Update(msg)
            m.children[id] = updated
        }
    case tea.KeyMsg:
        if isGlobalKey(msg) {
            // Handle in root
            return m.handleGlobalKey(msg)
        } else {
            // Route to current child
            child := m.children[m.current]
            updated, cmd := child.Update(msg)
            m.children[m.current] = updated
            return m, cmd
        }
    }
    return m, nil
}

func (m rootModel) View() string {
    // Compose views from children
    var views []string
    for _, child := range m.children {
        views = append(views, child.View())
    }
    return lipgloss.JoinVertical(lipgloss.Top, views...)
}
```

### Best Practices
- **Cache child models** for performance
- **Clear routing strategy** (global vs local messages)
- **Separate concerns** (each model handles its domain)

## 7. Layout Arithmetic is Error-Prone 📐

### Intuition
Manual height/width calculations break easily with changes. Use Lipgloss helpers.

### Problematic Code
```go
func (m model) View() string {
    // Hard-coded heights - fragile!
    contentHeight := m.height - 1 - 1 // Header(1) + Footer(1)
    // Breaks if borders added later
}
```

### Solution: Use Lipgloss Dimensions
```go
import "github.com/charmbracelet/lipgloss"

func (m model) View() string {
    header := lipgloss.NewStyle().
        Width(m.width).
        BorderBottom(true).
        Render("Header")
    
    footer := lipgloss.NewStyle().
        Width(m.width).
        Render("Footer")
    
    // Dynamic calculation
    contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)
    
    content := lipgloss.NewStyle().
        Width(m.width).
        Height(contentHeight). // Adapts to changes
        Render("Content")
    
    return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
```

### Key Functions
- `lipgloss.Height(string)`: Get rendered height
- `lipgloss.Width(string)`: Get rendered width
- **Always recalculate** on window resize

## 8. Recovering Your Terminal 🚑

### Problem
Panics in commands (not event loop) leave terminal in raw mode.

### Symptoms
- No cursor
- Input not echoed
- Malformed output

### Quick Recovery
```bash
# Reset terminal
reset
# Or
stty sane
# Or
tput reset
```

### Prevention
```go
defer func() {
    if r := recover(); r != nil {
        // Log panic
        // Optionally restart program
    }
}()

// Run command that might panic
```

## 9. End-to-End Testing with Teatest 🧪

### Intuition
Automate user interaction testing. Teatest simulates keypresses and validates output.

### Example Test
```go
import "github.com/charmbracelet/x/teatest"

func TestQuitConfirmation(t *testing.T) {
    m := model{}
    tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))
    
    // Wait for initial state
    teatest.WaitFor(t, tm.Output(),
        func(b []byte) bool {
            return strings.Contains(string(b), "Running.")
        },
        teatest.WithDuration(3*time.Second),
    )
    
    // Send Ctrl+C
    tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
    
    // Verify confirmation appears
    teatest.WaitFor(t, tm.Output(),
        func(b []byte) bool {
            return strings.Contains(string(b), "Quit?")
        },
    )
    
    // Send 'y' to quit
    tm.Type("y")
    
    // Wait for program to finish
    tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
```

### Golden File Testing
```go
tm := teatest.NewTestModel(t, m, 
    teatest.WithInitialTermSize(80, 24),
)

// First run creates golden file
// Subsequent runs compare output
teatest.RequireEqualOutput(t, tm)
```

## 10. Record Demos with VHS 🎥

### Intuition
Create professional demos and screenshots declaratively.

### Example VHS Tape (.tape file)
```vhs
Output demo.gif

Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"

Type "go run main.go"
Enter
Sleep 1s

# Interact with program
Ctrl+a
Type "help"
Enter
Sleep 1s

# Take screenshot
Screenshot screenshot.png

# Show some feature
Type "list"
Enter
Sleep 2s
```

### Usage
```bash
# Record demo
vhs demo.tape

# Creates:
# - demo.gif (animated)
# - screenshot.png (static)
```

### Benefits
- **Declarative scripting**
- **Consistent recordings**
- **Version-controlled demos**
- **Automated documentation**

## 11. Additional Best Practices 🌟

### Code Organization
1. **Separate models** by concern
2. **Use interfaces** for reusable components
3. **Centralize styling** with Lipgloss
4. **Implement navigation stack** for complex flows

### Performance Tips
1. **Virtualize long lists** (render only visible items)
2. **Debounce rapid updates**
3. **Cache expensive View() calculations**
4. **Use `tea.Batch`** for multiple commands

### Error Handling
```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    defer func() {
        if err := recover(); err != nil {
            m.error = fmt.Sprintf("Error: %v", err)
        }
    }()
    // ... update logic
}
```

### Testing Strategy
1. **Unit test models** in isolation
2. **Integration tests** with teatest
3. **Visual regression tests** with VHS
4. **Fuzz test** message handling

## 📚 Resources & References

- **Official Docs**: https://github.com/charmbracelet/bubbletea
- **Bubbles Library**: https://github.com/charmbracelet/bubbles
- **Lipgloss (Styling)**: https://github.com/charmbracelet/lipgloss
- **Teatest**: https://github.com/charmbracelet/x/tree/main/teatest
- **VHS**: https://github.com/charmbracelet/vhs
- **Example Project (PUG)**: https://github.com/leg100/pug

## 🎯 Key Takeaways

1. **Respect the event loop** - keep it fast
2. **Debug systematically** - message dumps are invaluable
3. **Architect carefully** - model trees scale better
4. **Test thoroughly** - teatest + VHS = confidence
5. **Layout dynamically** - avoid hard-coded dimensions
6. **Embrace concurrency** but manage message ordering
7. **Recover gracefully** - don't leave terminal broken

Building TUIs with Bubble Tea requires understanding both the Elm architecture patterns and Go's concurrency model. Start simple, add complexity gradually, and leverage the excellent tooling from Charm.