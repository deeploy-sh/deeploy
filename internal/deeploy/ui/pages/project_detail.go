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
	help    help.Model
	loading bool
	width   int
	height  int
	err     error
}

type projectDetailErrMsg struct{ err error }

func NewProjectDetailPage(s Store, projectID string) ProjectDetailPage {
	delegate := components.NewPodDelegate(40)

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

	l := list.New(components.PodsToItems(pods), delegate, 0, 0)
	l.Title = "Pods"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorForeground)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return ProjectDetailPage{
		store:   s,
		pods:    l,
		keys:    newProjectDetailKeyMap(),
		help:    styles.NewHelpModel(),
		project: &project,
	}
}

func (p ProjectDetailPage) Init() tea.Cmd {
	return nil
}

func (p ProjectDetailPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Don't handle keys if filtering is active
		if p.pods.FilterState() == list.Filtering {
			// But allow esc to cancel filter
			if msg.Code == tea.KeyEscape {
				var cmd tea.Cmd
				p.pods, cmd = p.pods.Update(msg)
				return p, cmd
			}
			break
		}

		switch {
		case key.Matches(msg, p.keys.Back):
			return p, func() tea.Msg {
				return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewDashboard() }}
			}
		case key.Matches(msg, p.keys.NewPod):
			projectID := p.project.ID
			return p, func() tea.Msg {
				return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodFormPage(projectID, nil) }}
			}
		case key.Matches(msg, p.keys.EditPod):
			item := p.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				projectID := p.project.ID
				return p, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodFormPage(projectID, &pod) }}
				}
			}
		case key.Matches(msg, p.keys.SelectPod):
			item := p.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				projectID := p.project.ID
				return p, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodDetailPage(projectID, &pod) }}
				}
			}

		case key.Matches(msg, p.keys.DeletePod):
			item := p.pods.SelectedItem()
			if item != nil {
				pod := item.(components.PodItem).Pod
				return p, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewPodDeletePage(&pod) }}
				}
			}
		case key.Matches(msg, p.keys.EditProject):
			if p.project != nil {
				project := p.project
				return p, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewProjectFormPage(project) }}
				}
			}
		case key.Matches(msg, p.keys.DeleteProject):
			if p.project != nil {
				project := p.project
				return p, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewProjectDeletePage(project) }}
				}
			}
		}

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		// List height for card content
		listHeight := min((msg.Height-1)/2, 12)
		p.pods.SetSize(46, listHeight)
		return p, nil

	case projectDetailErrMsg:
		p.err = msg.err
		p.loading = false
		return p, nil

	case messages.PodCreatedMsg:
		newItem := components.PodItem{Pod: repo.Pod(msg)}
		cmd := p.pods.InsertItem(len(p.pods.Items()), newItem)
		return p, cmd

	case messages.PodUpdatedMsg:
		pod := msg
		items := p.pods.Items()
		for i, item := range items {
			pi, ok := item.(components.PodItem)
			if ok && pi.ID == pod.ID {
				items[i] = components.PodItem{Pod: repo.Pod(pod)}
				break
			}
		}
		cmd := p.pods.SetItems(items)
		return p, cmd

	case messages.PodDeleteMsg:
		pod := msg
		items := p.pods.Items()
		for i, item := range items {
			pi, ok := item.(components.PodItem)
			if ok && pi.ID == pod.ID {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
		cmd := p.pods.SetItems(items)
		return p, cmd

	case messages.ProjectUpdatedMsg:
		project := repo.Project(msg)
		p.project = &project
		return p, nil
	}

	// Pass other messages to the list
	var cmd tea.Cmd
	p.pods, cmd = p.pods.Update(msg)
	return p, cmd
}

func (p ProjectDetailPage) View() tea.View {
	helpView := p.help.View(p.keys)
	contentHeight := p.height - 1

	var content string

	if p.loading {
		content = styles.MutedStyle.Render("Loading...")
	} else if p.err != nil {
		content = styles.ErrorStyle.Render("Error: " + p.err.Error())
	} else {
		content = p.renderContent()
	}

	centered := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered, helpView))
}

func (p ProjectDetailPage) renderContent() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorForeground)

	// Project Header
	header := titleStyle.Render(p.project.Title)
	if p.project.Description != "" {
		header += "\n" + styles.MutedStyle.Render(p.project.Description)
	}

	// Pods list
	var podsContent string
	if len(p.pods.Items()) == 0 {
		podsContent = fmt.Sprintf("Pods (0)\n\n%s",
			styles.MutedStyle.Render("No pods yet. Press 'n' to create one."))
	} else {
		podsContent = p.pods.View()
	}

	cardContent := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		podsContent,
	)

	card := components.Card(components.CardProps{
		Width:   50,
		Padding: []int{1, 2},
	}).Render(cardContent)

	return card
}

func (p ProjectDetailPage) Breadcrumbs() []string {
	return []string{"Projects", p.project.Title}
}
