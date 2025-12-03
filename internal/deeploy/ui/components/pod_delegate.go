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

type PodItem struct {
	repo.Pod
}

func (i PodItem) FilterValue() string { return i.Title }

func PodsToItems(pods []repo.Pod) []list.Item {
	items := make([]list.Item, len(pods))
	for i, p := range pods {
		items[i] = PodItem{Pod: p}
	}
	return items
}

type PodDelegate struct {
	width int
}

func NewPodDelegate(width int) PodDelegate {
	return PodDelegate{width: width}
}

func (d PodDelegate) Height() int                             { return 1 }
func (d PodDelegate) Spacing() int                            { return 0 }
func (d PodDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d PodDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	pod, ok := item.(PodItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	var line string
	if isSelected {
		line = lipgloss.NewStyle().
			Foreground(styles.ColorPrimary).
			Bold(true).
			Render(fmt.Sprintf("> %s", pod.Title))
	} else {
		line = fmt.Sprintf("  %s", pod.Title)
	}

	fmt.Fprint(w, line)
}
