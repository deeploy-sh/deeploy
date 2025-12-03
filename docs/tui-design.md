# Deeploy TUI Design Plan v0.1

## The Concept: Navigation Palette + Direct Action Keys

```
PALETTE (ctrl+k) = Navigation (where do I want to go?)
DIRECT KEYS      = Actions (what do I want to do?)
VIM KEYS         = List Navigation (j/k/g/G)
```

---

## Core UX

### Fullscreen View with Header + Footer

```
+------------------------------------------------------------------+
|  deeploy > Projects > my-api                    server [online]   |  ← Header (Breadcrumb)
+------------------------------------------------------------------+
|                                                                   |
|  PODS                                                             |
|                                                                   |
|  > web-pod           running     3h      nginx:latest             |
|    api-pod           running     3h      node:18                  |
|    worker-pod        stopped     1d      python:3.11              |
|                                                                   |
|                                                                   |
+------------------------------------------------------------------+
|  j/k nav  enter open  l logs  x shell  s stop  n new   ctrl+k go  |  ← Footer (Context Keys)
+------------------------------------------------------------------+
```

### Navigation Palette (ctrl+k)

```
+------------------------------------------------------------------+
|  deeploy > Projects                                               |
+------------------------------------------------------------------+
|                                                                   |
|     +--------------------------------------------------------+   |
|     |  > logs web_                                           |   |
|     +--------------------------------------------------------+   |
|     |  > logs my-api/web-pod       View pod logs             |   |
|     |    logs my-api/api-pod       View pod logs             |   |
|     |    logs staging/worker       View pod logs             |   |
|     |    project my-api            Open project              |   |
|     |    projects                  All projects              |   |
|     +--------------------------------------------------------+   |
|                                                                   |
+------------------------------------------------------------------+
|  j/k nav   enter select   esc close                               |
+------------------------------------------------------------------+
```

---

## Keyboard Model

### Global Keys (everywhere)

| Key | Action |
|-----|--------|
| `ctrl+k` | Open Navigation Palette |
| `?` | Help overlay |
| `esc` | Back / Close overlay |
| `q` | Quit (with confirmation if needed) |
| `ctrl+c` | Force quit |

### List Navigation (in all lists)

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `ctrl+d` | Half page down |
| `ctrl+u` | Half page up |
| `enter` | Select / Open |
| `/` | Filter list (inline search) |

### Context Actions (in footer, per view)

**Projects List:**
| Key | Action |
|-----|--------|
| `enter` | Open project |
| `n` | New project |
| `e` | Edit project |
| `d` | Delete project |
| `r` | Refresh |

**Project Detail / Pods List:**
| Key | Action |
|-----|--------|
| `enter` | Open pod |
| `n` | New pod |
| `l` | View logs |
| `x` | Shell/exec |
| `s` | Start/stop |
| `r` | Restart |
| `d` | Delete pod |

**Pod Logs View:**
| Key | Action |
|-----|--------|
| `f` | Follow (tail) |
| `p` | Pause |
| `/` | Search in logs |
| `g` | Jump to top |
| `G` | Jump to bottom |

---

## Navigation Palette Commands

### Static Commands (always available)

```
projects              → Projects List
settings              → Global Settings
help                  → Help View
```

### Dynamic Commands (generated from data)

```
project <name>        → Project Detail
  - project my-api
  - project staging

pods <project>        → Pods List
  - pods my-api
  - pods staging

logs <project>/<pod>  → Pod Logs (DIRECT!)
  - logs my-api/web-pod
  - logs my-api/api-pod
  - logs staging/worker

settings <project>    → Project Settings
  - settings my-api
  - settings staging
```

### Fuzzy Search

User types: `logs web`
Matches:
1. `logs my-api/web-pod`
2. `logs staging/web-frontend`

**8 characters → directly at the logs!**

---

## View Hierarchy

```
Dashboard (optional, or directly to Projects)
├── Projects List
│   ├── Project Detail (shows Pods)
│   │   ├── Pod Detail
│   │   ├── Pod Logs
│   │   └── Pod Shell (external terminal)
│   └── Project Settings
├── Global Settings
└── Help
```

**Navigation:**
- `enter` = drill down
- `esc` = back
- `ctrl+k` = jump to ANY view

---

## Visual Design

### Colors (Minimal)

```go
Primary     = "204"   // Pink/Magenta - Focus, Selection
Success     = "42"    // Green - Running, OK
Warning     = "214"   // Orange - Building, Warning
Error       = "9"     // Red - Stopped, Error
TextDim     = "244"   // Gray - Secondary text
Surface     = "236"   // Dark - Palette background
```

### Typography

```
Header:     Bold, Primary color for current location
List Items: Normal, Status color for status indicator
Footer:     Dim, Key highlighted
```

### List Item Format

```
  NAME               STATUS      AGE     IMAGE
> web-pod            running     3h      nginx:latest
  api-pod            running     3h      node:18
  worker-pod         stopped     1d      python:3.11
```

- `>` cursor for selected
- Status colored (green/red/yellow)
- Aligned columns

---

## Components

### 1. Header (Breadcrumb)

```go
// components/header.go
type Header struct {
    Breadcrumb []string  // ["deeploy", "Projects", "my-api"]
    ServerStatus string  // "online" / "offline"
}

func (h Header) View() string {
    path := strings.Join(h.Breadcrumb, " > ")
    status := fmt.Sprintf("[%s]", h.ServerStatus)
    // left-align path, right-align status
}
```

### 2. Footer (Context Keys)

```go
// components/footer.go
type FooterKey struct {
    Key  string  // "l"
    Desc string  // "logs"
}

type Footer struct {
    Keys []FooterKey
}

func (f Footer) View() string {
    // "j/k nav  l logs  x shell  ctrl+k go"
}
```

### 3. List (Navigable)

```go
// components/list.go
type List struct {
    Items    []ListItem
    Selected int
    Height   int
}

// Handles j/k/g/G/ctrl+d/ctrl+u
func (l *List) Update(msg tea.Msg) tea.Cmd
```

### 4. Palette (Navigation)

```go
// components/palette.go
type Palette struct {
    Input    textinput.Model
    Commands []Command
    Filtered []Command
    Selected int
    Visible  bool
}

// Fuzzy search, j/k navigation, enter to execute
func (p *Palette) Update(msg tea.Msg) tea.Cmd
```

---

## File Structure

```
internal/deeploy/ui/
├── app.go                    # Root, handles palette overlay
├── components/
│   ├── header.go             # Breadcrumb header
│   ├── footer.go             # Context keys footer
│   ├── list.go               # Navigable list
│   └── palette.go            # Navigation palette
├── views/
│   ├── projects_list.go      # Projects list
│   ├── project_detail.go     # Project with pods
│   ├── pod_logs.go           # Log viewer
│   ├── settings.go           # Settings
│   └── help.go               # Help
├── commands/
│   └── registry.go           # Palette commands
└── styles/
    └── styles.go             # Colors, spacing
```

---

## Implementation Order

### Phase 1: Foundation
1. `components/header.go` - Breadcrumb
2. `components/footer.go` - Context keys
3. `components/list.go` - j/k/g/G navigation
4. `styles/styles.go` - Colors

### Phase 2: Views
5. `views/projects_list.go` - Main list
6. `views/project_detail.go` - Pods list
7. `views/pod_logs.go` - Log viewer
8. `views/settings.go` - Settings
9. `views/help.go` - Help

### Phase 3: Palette
10. `components/palette.go` - Overlay
11. `commands/registry.go` - Static + dynamic commands
12. Fuzzy search integration (`sahilm/fuzzy`)

### Phase 4: Polish
13. Inline `/` filter for lists
14. Delete confirmation (inline `[y/n]`)
15. Loading states
16. Error handling

---

## Data Loading & State

### App State (Global)

```go
// app.go
type App struct {
    // UI State
    currentView   tea.Model
    palette       Palette
    showPalette   bool

    // Data (loaded at start, refreshed on actions)
    projects      []Project
    pods          map[string][]Pod  // projectID -> pods

    // Connection
    serverURL     string
    token         string
    isConnected   bool
}
```

### Loading Strategy: Everything at App Start

```go
func (app *App) Init() tea.Cmd {
    return tea.Batch(
        app.checkConnection(),   // Ping server
        app.loadProjects(),      // GET /projects
    )
}

func (app *App) loadProjects() tea.Cmd {
    return func() tea.Msg {
        projects, err := api.GetProjects()
        if err != nil {
            return ErrorMsg{err}
        }

        // Parallel: Load pods for each project
        var allPods = make(map[string][]Pod)
        for _, p := range projects {
            pods, _ := api.GetPods(p.ID)
            allPods[p.ID] = pods
        }

        return DataLoadedMsg{
            Projects: projects,
            Pods:     allPods,
        }
    }
}
```

### Refresh Triggers

| Trigger | When |
|---------|------|
| App Start | Initial load |
| `r` key | Manual refresh |
| After Create | New project/pod created |
| After Delete | Project/pod deleted |
| After Deploy | Deployment finished |

```go
func (app *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case ProjectCreatedMsg:
        return app, app.refreshData()

    case PodDeletedMsg:
        return app, app.refreshData()

    case tea.KeyMsg:
        if msg.String() == "r" {
            return app, app.refreshData()
        }
    }
}
```

### Palette Command Generation (from cached data)

```go
func (app *App) generatePaletteCommands() []Command {
    var cmds []Command

    // Static commands
    cmds = append(cmds,
        Command{Name: "projects", Desc: "All projects"},
        Command{Name: "settings", Desc: "Global settings"},
        Command{Name: "help", Desc: "Show help"},
    )

    // Dynamic from cached data
    for _, p := range app.projects {
        cmds = append(cmds, Command{
            Name: fmt.Sprintf("project %s", p.Name),
            Desc: "Open project",
        })

        for _, pod := range app.pods[p.ID] {
            cmds = append(cmds, Command{
                Name: fmt.Sprintf("logs %s/%s", p.Name, pod.Name),
                Desc: "View pod logs",
            })
        }
    }

    return cmds
}
```

### Why This Works for deeploy

| Metric | Typical Value | Impact |
|--------|---------------|--------|
| Projects | 5-20 | Negligible |
| Pods per Project | 1-10 | Negligible |
| Total Items | 50-200 | <10KB memory |
| API Calls at Start | 1 + N | <500ms |
| Fuzzy Search | Instant | No delay |

**No lazy loading needed. No complex caching. Simple and fast.**

---

## Tech Details

### Dependencies

```go
// go.mod additions
github.com/sahilm/fuzzy       // Fuzzy search for palette
// Already have:
github.com/charmbracelet/bubbletea
github.com/charmbracelet/bubbles
github.com/charmbracelet/lipgloss
```

---

## Decisions Made

| Topic | Decision | Rationale |
|-------|----------|-----------|
| **Palette** | Deep navigation (logs, settings per pod) | "wow" factor, VS Code-like |
| **Actions** | Direct keys in footer | Speed, muscle memory |
| **Mouse** | Keyboard only (v0.1) | Simpler, SSH-friendly |
| **Confirmations** | Inline `[y/n]` | Fast, minimal |
| **Layout** | Fullscreen views | Focus, simplicity |
| **Data Loading** | Everything at app start | Simple, fast for small datasets |

---

## Success Criteria v0.1

- [ ] `ctrl+k` opens palette with fuzzy search
- [ ] Can navigate to any view via palette ("logs web" → pod logs)
- [ ] j/k/g/G works in all lists
- [ ] Footer shows context-appropriate keys
- [ ] Header shows breadcrumb path
- [ ] Direct keys work (n/d/l/x/s)
- [ ] Looks clean and professional
- [ ] "wow" reaction from users
