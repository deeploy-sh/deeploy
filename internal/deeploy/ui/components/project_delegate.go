package components

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type ProjectItem struct {
	repo.Project
}

func (i ProjectItem) FilterValue() string { return i.Title }

func ProjectsToItems(projects []repo.Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = ProjectItem{Project: p}
	}
	return items
}

type ProjectDelegate struct {
	width int
}

func NewProjectDelegate(width int) ProjectDelegate {
	return ProjectDelegate{width: width}
}

func (d ProjectDelegate) Height() int                             { return 1 }
func (d ProjectDelegate) Spacing() int                            { return 0 }
func (d ProjectDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d ProjectDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	project, ok := item.(ProjectItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	// Base style
	lineStyle := lipgloss.NewStyle().
		Width(d.width).
		Background(styles.ColorBackgroundPanel())

	var line string
	if isSelected {
		line = lineStyle.
			Foreground(styles.ColorPrimary()).
			Bold(true).
			Render(fmt.Sprintf("%s", project.Title))
	} else {
		line = lineStyle.
			Foreground(styles.ColorForeground()).
			Render(fmt.Sprintf("%s", project.Title))
	}

	fmt.Fprint(w, line)
}
