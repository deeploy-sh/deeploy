# Plan: Fix UX Bugs

## 1. Pod count mit space-between

**File:** `internal/deeploy/ui/components/items.go`

```go
func (i ProjectItem) Title() string {
    // Needs width context - pass width to Title or use fixed width
    // Option: Just add spaces manually for now
    return fmt.Sprintf("%-30s (%d)", i.Project.Title, i.PodCount)
}
```

Or better: Let the list component handle the layout with space-between.

## 2. Remove ctrl+k from Dashboard

**File:** `internal/deeploy/ui/pages/dashboard.go`

Remove `Search` keybinding since palette is now triggered by `:` globally in app.go.

## 3. Pod palette action

**File:** `internal/deeploy/ui/pages/app.go`

```go
for _, p := range m.pods {
    pod := p
    items = append(items, components.PaletteItem{
        ItemTitle:   pod.Title,
        Description: pod.Description,
        Category:    "pod",
        Action: func() tea.Msg {
            return msg.ChangePage{
                PageFactory: func(s msg.Store) tea.Model {
                    // Need project reference for PodDetailPage
                    // Find project for this pod
                    var project *repo.Project
                    for _, pr := range s.Projects() {
                        if pr.ID == pod.ProjectID {
                            project = &pr
                            break
                        }
                    }
                    return NewPodDetailPage(&pod, project)
                },
            }
        },
    })
}
```

## Files to modify

1. `internal/deeploy/ui/components/items.go` - Space-between for pod count
2. `internal/deeploy/ui/pages/dashboard.go` - Remove Search keybinding
3. `internal/deeploy/ui/pages/app.go` - Fix pod palette action
