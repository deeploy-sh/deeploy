package components

import (
	"fmt"

	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

// ProjectItem wraps a Project for use in ScrollList
type ProjectItem struct {
	repo.Project
	PodCount int
}

func (i ProjectItem) Title() string       { return i.Project.Title }
func (i ProjectItem) Suffix() string      { return fmt.Sprintf("(%d)", i.PodCount) }
func (i ProjectItem) FilterValue() string { return i.Project.Title }

// ProjectsToItems converts a slice of Projects to ScrollItems
func ProjectsToItems(projects []repo.Project, pods []repo.Pod) []ScrollItem {
	// Count pods per project
	podCounts := make(map[string]int)
	for _, p := range pods {
		podCounts[p.ProjectID]++
	}

	items := make([]ScrollItem, len(projects))
	for i, p := range projects {
		items[i] = ProjectItem{Project: p, PodCount: podCounts[p.ID]}
	}
	return items
}

// PodItem wraps a Pod for use in ScrollList
type PodItem struct {
	repo.Pod
}

func (i PodItem) Title() string       { return i.Pod.Title }
func (i PodItem) FilterValue() string { return i.Pod.Title }

// PodsToItems converts a slice of Pods to ScrollItems
func PodsToItems(pods []repo.Pod) []ScrollItem {
	items := make([]ScrollItem, len(pods))
	for i, p := range pods {
		items[i] = PodItem{Pod: p}
	}
	return items
}
