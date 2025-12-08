package pages

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type projectDetailKeyMap struct {
	NewPod        key.Binding
	EditPod       key.Binding
	SelectPod     key.Binding
	DeletePod     key.Binding
	EditProject   key.Binding
	DeleteProject key.Binding
	Back          key.Binding
}

func (k projectDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NewPod, k.EditPod, k.SelectPod, k.DeletePod, k.EditProject, k.DeleteProject, k.Back}
}

func (k projectDetailKeyMap) FullHelp() [][]key.Binding {
	return nil
}

func (m ProjectDetailPage) HelpKeys() help.KeyMap {
	return m.keys
}

func newProjectDetailKeyMap() projectDetailKeyMap {
	return projectDetailKeyMap{
		NewPod:        key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new pod")),
		EditPod:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit pod")),
		DeletePod:     key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete pod")),
		SelectPod:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select pod")),
		EditProject:   key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "edit project")),
		DeleteProject: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete project")),
		Back:          key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

type ProjectDetailPage struct {
	store   msg.Store
	project *repo.Project
	pods    components.ScrollList
	keys    projectDetailKeyMap
	width   int
	height  int
}

func NewProjectDetailPage(s msg.Store, projectID string) ProjectDetailPage {
	var project repo.Project
	for _, p := range s.Projects() {
		if p.ID == projectID {
			project = p
			break
		}
	}

	var pods []repo.Pod
	for _, p := range s.Pods() {
		if p.ProjectID == projectID {
			pods = append(pods, p)
		}
	}

	card := components.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	l := components.NewScrollList(components.PodsToItems(pods), components.ScrollListConfig{
		Width:  card.InnerWidth(),
		Height: 15,
	})

	return ProjectDetailPage{
		store:   s,
		pods:    l,
		keys:    newProjectDetailKeyMap(),
		project: &project,
	}
}

func (m ProjectDetailPage) Init() tea.Cmd {
	return nil
}

func (m ProjectDetailPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.pods, cmd = m.pods.Update(tmsg)

	switch tmsg := tmsg.(type) {
	case msg.DataLoaded:
		// Filter pods for this project
		var pods []repo.Pod
		for _, p := range tmsg.Pods {
			if p.ProjectID == m.project.ID {
				pods = append(pods, p)
			}
		}
		m.pods.SetItems(components.PodsToItems(pods))

		// Update project data too (for title changes)
		for _, p := range tmsg.Projects {
			if p.ID == m.project.ID {
				m.project = &p
				break
			}
		}
		return m, cmd

	case tea.KeyPressMsg:
		switch {
		case key.Matches(tmsg, m.keys.Back):
			return m, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
			}
		case key.Matches(tmsg, m.keys.NewPod):
			projectID := m.project.ID
			return m, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodFormPage(projectID, nil) }}
			}
		case key.Matches(tmsg, m.keys.EditPod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				projectID := m.project.ID
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodFormPage(projectID, &pod) }}
				}
			}
		case key.Matches(tmsg, m.keys.SelectPod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodDetailPage(&pod, m.project) }}
				}
			}
		case key.Matches(tmsg, m.keys.DeletePod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodDeletePage(&pod) }}
				}
			}
		case key.Matches(tmsg, m.keys.EditProject):
			if m.project != nil {
				project := m.project
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectFormPage(project) }}
				}
			}
		case key.Matches(tmsg, m.keys.DeleteProject):
			if m.project != nil {
				project := m.project
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDeletePage(s, project) }}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, cmd
	}

	return m, cmd
}

func (m ProjectDetailPage) View() tea.View {
	contentHeight := m.height
	content := m.renderContent()
	centered := lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, content)
	return tea.NewView(centered)
}

func (m ProjectDetailPage) renderContent() string {
	card := components.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	w := card.InnerWidth()

	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render(m.project.Title + " > Pods")

	var podsContent string
	if len(m.pods.Items()) == 0 {
		podsContent = styles.MutedStyle().Render("No pods yet. Press 'n' to create one.")
	} else {
		podsContent = lipgloss.NewStyle().
			Width(w).
			Height(m.pods.Height()).
			Background(styles.ColorBackgroundPanel()).
			Render(m.pods.View())
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, podsContent)

	return components.Card(card).Render(content)
}

func (m ProjectDetailPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title}
}
