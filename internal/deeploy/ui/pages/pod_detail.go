package pages

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type PodDetailPage struct {
	pod          *repo.Pod
	project      *repo.Project
	keyEditPod   key.Binding
	keyDeletePod key.Binding
	keyFilter    key.Binding
	keyBack      key.Binding
	width        int
	height       int
}

func (m PodDetailPage) HelpKeys() []key.Binding {
	return []key.Binding{m.keyEditPod, m.keyDeletePod, m.keyFilter, m.keyBack}
}

func NewPodDetailPage(pod *repo.Pod, project *repo.Project) PodDetailPage {
	return PodDetailPage{
		pod:          pod,
		project:      project,
		keyEditPod:   key.NewBinding(key.WithKeys("enter", "e"), key.WithHelp("enter", "edit pod")),
		keyDeletePod: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete pod")),
		keyFilter:    key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		keyBack:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

func (m PodDetailPage) Init() tea.Cmd {
	return nil
}

func (m PodDetailPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(tmsg, m.keyBack):
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
