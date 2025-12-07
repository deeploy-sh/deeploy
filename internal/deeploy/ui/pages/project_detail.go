package pages

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
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
	store   Store
	project *repo.Project
	pods    components.ScrollList
	keys    projectDetailKeyMap
	loading bool
	width   int
	height  int
	err     error
}

type projectDetailErrMsg struct{ err error }

func NewProjectDetailPage(s Store, projectID string) ProjectDetailPage {
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
	l := components.NewScrollList(components.PodsToItems(pods), card.InnerWidth(), 15)

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

func (m ProjectDetailPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Back):
			return m, func() tea.Msg {
				return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewDashboard(s) }}
			}
		case key.Matches(msg, m.keys.NewPod):
			projectID := m.project.ID
			return m, func() tea.Msg {
				return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodFormPage(projectID, nil) }}
			}
		case key.Matches(msg, m.keys.EditPod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				projectID := m.project.ID
				return m, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodFormPage(projectID, &pod) }}
				}
			}
		case key.Matches(msg, m.keys.SelectPod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				return m, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodDetailPage(&pod, m.project) }}
				}
			}

		case key.Matches(msg, m.keys.DeletePod):
			item := m.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				return m, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodDeletePage(&pod) }}
				}
			}
		case key.Matches(msg, m.keys.EditProject):
			if m.project != nil {
				project := m.project
				return m, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewProjectFormPage(project) }}
				}
			}
		case key.Matches(msg, m.keys.DeleteProject):
			if m.project != nil {
				project := m.project
				return m, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewProjectDeletePage(project) }}
				}
			}
		case msg.Code == tea.KeyUp:
			m.pods.CursorUp()
			return m, nil
		case msg.Code == tea.KeyDown:
			m.pods.CursorDown()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case projectDetailErrMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case messages.PodCreatedMsg:
		items := m.pods.Items()
		items = append(items, components.PodItem{Pod: repo.Pod(msg)})
		m.pods.SetItems(items)
		return m, nil

	case messages.PodUpdatedMsg:
		pod := msg
		items := m.pods.Items()
		for i, item := range items {
			pi, ok := item.(components.PodItem)
			if ok && pi.ID == pod.ID {
				items[i] = components.PodItem{Pod: repo.Pod(pod)}
				break
			}
		}
		m.pods.SetItems(items)
		return m, nil

	case messages.PodDeleteMsg:
		pod := msg
		items := m.pods.Items()
		for i, item := range items {
			pi, ok := item.(components.PodItem)
			if ok && pi.ID == pod.ID {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
		m.pods.SetItems(items)
		return m, nil

	case messages.ProjectUpdatedMsg:
		project := repo.Project(msg)
		m.project = &project
		return m, nil
	}

	return m, nil
}

func (m ProjectDetailPage) View() tea.View {
	contentHeight := m.height

	var content string

	if m.loading {
		content = styles.MutedStyle().Render("Loading...")
	} else if m.err != nil {
		content = styles.ErrorStyle().Render("Error: " + m.err.Error())
	} else {
		content = m.renderContent()
	}

	centered := lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(centered)
}

func (m ProjectDetailPage) renderContent() string {
	card := components.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	w := card.InnerWidth()

	// Custom title (like dashboard)
	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render(m.project.Title + " > Pods")

	// Pods list with background (like dashboard)
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
