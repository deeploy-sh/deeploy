package pages

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type podDetailKeyMap struct {
	EditPod   key.Binding
	DeletePod key.Binding
	Filter    key.Binding
	Back      key.Binding
}

func (k podDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.EditPod, k.DeletePod, k.Filter, k.Back}
}

func (k podDetailKeyMap) FullHelp() [][]key.Binding {
	return nil
}

func newPodDetailKeyMap() podDetailKeyMap {
	return podDetailKeyMap{
		EditPod:   key.NewBinding(key.WithKeys("enter", "e"), key.WithHelp("enter", "edit pod")),
		DeletePod: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete pod")),
		Filter:    key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

type PodDetailPage struct {
	pod     *repo.Pod
	project *repo.Project
	keys    podDetailKeyMap
	help    help.Model
	loading bool
	width   int
	height  int
}

func NewPodDetailPage(pod *repo.Pod, project *repo.Project) PodDetailPage {
	return PodDetailPage{
		pod:     pod,
		project: project,
		keys:    newPodDetailKeyMap(),
		help:    styles.NewHelpModel(),
		loading: true,
	}
}

func (m PodDetailPage) Init() tea.Cmd {
	return nil
}

func (m PodDetailPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return ChangePageMsg{
					// Page: NewDashboard()
					PageFactory: func(s Store) tea.Model {
						return NewProjectDetailPage(s, m.project.ID)
					},
				}
			}
		case key.Matches(msg, m.keys.EditPod):
			// item := p.list.SelectedItem()
			// if item != nil {
			// 	pod := item.(components.PodItem).Pod
			// 	projectID := p.project.ID
			// 	return p, func() tea.Msg {
			// 		return messages.ChangePageMsg{Page: NewPodFormPage(projectID, &pod)}
			// 	}
			// }
		case key.Matches(msg, m.keys.DeletePod):
			// item := p.list.SelectedItem()
			// if item != nil {
			// 	pod := item.(components.PodItem).Pod
			// 	return p, func() tea.Msg {
			// 		return messages.ChangePageMsg{Page: NewPodDeletePage(&pod)}
			// 	}
			// }
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// List height for card content
		// listHeight := min((msg.Height-1)/2, 12)
		return m, nil

	case messages.PodUpdatedMsg:
		// pod := msg
		// items := p.list.Items()
		// for i, item := range items {
		// 	pi, ok := item.(components.PodItem)
		// 	if ok && pi.ID == pod.ID {
		// 		items[i] = components.PodItem{Pod: repo.Pod(pod)}
		// 		break
		// 	}
		// }
		// cmd := p.list.SetItems(items)
		return m, nil

	case messages.PodDeleteMsg:
		// pod := msg
		// items := p.list.Items()
		// for i, item := range items {
		// 	pi, ok := item.(components.PodItem)
		// 	if ok && pi.ID == pod.ID {
		// 		items = append(items[:i], items[i+1:]...)
		// 		break
		// 	}
		// }
		// cmd := p.list.SetItems(items)
		// return p, cmd
		return m, nil
	}

	// Pass other messages to the list
	var cmd tea.Cmd
	return m, cmd
}

func (m PodDetailPage) View() tea.View {
	helpView := m.help.View(m.keys)
	contentHeight := m.height - 1

	var content string

	content = m.pod.Title

	centered := lipgloss.Place(m.width, contentHeight,
		lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered, helpView))
}

func (m PodDetailPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title}
}
