package pages

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
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

func (m PodDetailPage) HelpKeys() help.KeyMap {
	return m.keys
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
	width   int
	height  int
}

func NewPodDetailPage(pod *repo.Pod, project *repo.Project) PodDetailPage {
	return PodDetailPage{
		pod:     pod,
		project: project,
		keys:    newPodDetailKeyMap(),
	}
}

func (m PodDetailPage) Init() tea.Cmd {
	return nil
}

func (m PodDetailPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(tmsg, m.keys.Back):
			projectID := m.project.ID
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewProjectDetailPage(s, projectID)
					},
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	return m, nil
}

func (m PodDetailPage) View() tea.View {
	contentHeight := m.height
	content := m.pod.Title

	centered := lipgloss.Place(m.width, contentHeight,
		lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(centered)
}

func (m PodDetailPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title}
}
