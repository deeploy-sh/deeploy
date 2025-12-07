package pages

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
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
	Filter        key.Binding
	Back          key.Binding
}

func (k projectDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NewPod, k.EditPod, k.SelectPod, k.DeletePod, k.EditProject, k.DeleteProject, k.Filter, k.Back}
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
		EditPod:       key.NewBinding(key.WithKeys("e"), key.WithHelp("ej", "edit pod")),
		DeletePod:     key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete pod")),
		SelectPod:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select pod")),
		EditProject:   key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "edit project")),
		DeleteProject: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete project")),
		Filter:        key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Back:          key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

type ProjectDetailPage struct {
	store   Store
	project *repo.Project
	pods    list.Model
	keys    projectDetailKeyMap
	loading bool
	width   int
	height  int
	err     error
}

type projectDetailErrMsg struct{ err error }

func NewProjectDetailPage(s Store, projectID string) ProjectDetailPage {
	delegate := components.NewPodDelegate()
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

	l := list.New(components.PodsToItems(pods), delegate, 46, 15)
	l.Title = project.Title + " > Pods"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorForeground())
	l.Styles.TitleBar = lipgloss.NewStyle().Padding(0, 0, 1, 0)
	l.SetShowTitle(true)
	l.InfiniteScrolling = true
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

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
		// Don't handle keys if filtering is active
		if m.pods.FilterState() == list.Filtering {
			// But allow esc to cancel filter
			if msg.Code == tea.KeyEscape {
				var cmd tea.Cmd
				m.pods, cmd = m.pods.Update(msg)
				return m, cmd
			}
			break
		}

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
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// List height for card content
		listHeight := min((msg.Height-1)/2, 15)
		m.pods.SetSize(56, listHeight)
		return m, nil

	case projectDetailErrMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case messages.PodCreatedMsg:
		newItem := components.PodItem{Pod: repo.Pod(msg)}
		cmd := m.pods.InsertItem(len(m.pods.Items()), newItem)
		return m, cmd

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
		cmd := m.pods.SetItems(items)
		return m, cmd

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
		cmd := m.pods.SetItems(items)
		return m, cmd

	case messages.ProjectUpdatedMsg:
		project := repo.Project(msg)
		m.project = &project
		return m, nil
	}

	// Pass other messages to the list
	var cmd tea.Cmd
	m.pods, cmd = m.pods.Update(msg)
	return m, cmd
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
	// Pods list
	var podsContent string
	if len(m.pods.Items()) == 0 {
		podsContent = fmt.Sprintf("Pods (0)\n\n%s",
			styles.MutedStyle().Render("No pods yet. Press 'n' to create one."))
	} else {
		podsContent = m.pods.View()
	}

	cardContent := lipgloss.JoinVertical(lipgloss.Left, podsContent)

	card := components.Card(components.CardProps{
		Width:   50,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(cardContent)

	return card
}

func (m ProjectDetailPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title}
}
