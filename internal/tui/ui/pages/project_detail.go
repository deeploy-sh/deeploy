package pages

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type ProjectDetailPage struct {
	store          msg.Store
	project        *model.Project
	pods           components.ScrollList
	keyNewPod      key.Binding
	keySelectPod   key.Binding
	keyDeletePod   key.Binding
	keyEditProject key.Binding
	keyBack        key.Binding
	width          int
	height         int
}

func (m ProjectDetailPage) HelpKeys() []key.Binding {
	return []key.Binding{m.keyNewPod, m.keySelectPod, m.keyDeletePod, m.keyEditProject, m.keyBack}
}

func NewProjectDetailPage(s msg.Store, projectID string) ProjectDetailPage {
	var project model.Project
	for _, p := range s.Projects() {
		if p.ID == projectID {
			project = p
			break
		}
	}

	var pods []model.Pod
	for _, p := range s.Pods() {
		if p.ProjectID == projectID {
			pods = append(pods, p)
		}
	}

	card := styles.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	l := components.NewScrollList(components.PodsToItems(pods), components.ScrollListConfig{
		Width:  card.InnerWidth(),
		Height: 15,
	})

	return ProjectDetailPage{
		store:          s,
		pods:           l,
		project:        &project,
		keyNewPod:      key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new pod")),
		keyDeletePod:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete pod")),
		keySelectPod:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select pod")),
		keyEditProject: key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit project")),
		keyBack:        key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
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
		var pods []model.Pod
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
		case key.Matches(tmsg, m.keyBack):
			return m, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
			}
		case key.Matches(tmsg, m.keyNewPod):
			projectID := m.project.ID
			return m, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodFormPage(projectID, nil) }}
			}
		case key.Matches(tmsg, m.keySelectPod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodDetailPage(s, pod.ID) }}
				}
			}
		case key.Matches(tmsg, m.keyDeletePod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodDeletePage(&pod) }}
				}
			}
		case key.Matches(tmsg, m.keyEditProject):
			if m.project != nil {
				project := m.project
				return m, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectFormPage(project) }}
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
	if len(m.pods.Items()) == 0 {
		return m.renderEmptyState()
	}

	card := styles.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	w := card.InnerWidth()

	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render(m.project.Title + " > Pods")

	podsContent := lipgloss.NewStyle().
		Width(w).
		Height(m.pods.Height()).
		Background(styles.ColorBackgroundPanel()).
		Render(m.pods.View())

	content := lipgloss.JoinVertical(lipgloss.Left, title, podsContent)

	return styles.Card(card).Render(content)
}

func (m ProjectDetailPage) renderEmptyState() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorForeground()).
		MarginBottom(1)

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("No pods yet"),
		styles.MutedStyle().Render("Press 'n' to create your first pod"),
	)
}

func (m ProjectDetailPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title}
}
