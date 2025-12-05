package pages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type projectDetailKeyMap struct {
	NewPod        key.Binding
	EditPod       key.Binding
	DeletePod     key.Binding
	EditProject   key.Binding
	DeleteProject key.Binding
	Filter        key.Binding
	Back          key.Binding
}

func (k projectDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NewPod, k.EditPod, k.DeletePod, k.Filter, k.Back}
}

func (k projectDetailKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.NewPod, k.EditPod, k.DeletePod, k.EditProject, k.DeleteProject, k.Filter, k.Back}}
}

func newProjectDetailKeyMap() projectDetailKeyMap {
	return projectDetailKeyMap{
		NewPod:        key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new pod")),
		EditPod:       key.NewBinding(key.WithKeys("enter", "e"), key.WithHelp("enter", "edit pod")),
		DeletePod:     key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete pod")),
		EditProject:   key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "edit project")),
		DeleteProject: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "delete project")),
		Filter:        key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Back:          key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

type ProjectDetailPage struct {
	project *repo.Project
	list    list.Model
	keys    projectDetailKeyMap
	help    help.Model
	loading bool
	width   int
	height  int
	err     error
}

type projectDetailDataMsg struct {
	project repo.Project
	pods    []repo.Pod
}

type projectDetailErrMsg struct{ err error }

func NewProjectDetailPage(projectID string) ProjectDetailPage {
	delegate := components.NewPodDelegate(40)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Pods"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorForeground)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return ProjectDetailPage{
		list:    l,
		keys:    newProjectDetailKeyMap(),
		help:    styles.NewHelpModel(),
		loading: true,
		project: &repo.Project{ID: projectID},
	}
}

func (p ProjectDetailPage) Init() tea.Cmd {
	projectID := p.project.ID
	return func() tea.Msg {
		return loadProjectDetail(projectID)
	}
}

func (p ProjectDetailPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Don't handle keys if filtering is active
		if p.list.FilterState() == list.Filtering {
			// But allow esc to cancel filter
			if msg.Code == tea.KeyEscape {
				var cmd tea.Cmd
				p.list, cmd = p.list.Update(msg)
				return p, cmd
			}
			break
		}

		switch {
		case key.Matches(msg, p.keys.Back):
			return p, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewDashboard()}
			}
		case key.Matches(msg, p.keys.NewPod):
			projectID := p.project.ID
			return p, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewPodFormPage(projectID, nil)}
			}
		case key.Matches(msg, p.keys.EditPod):
			if item := p.list.SelectedItem(); item != nil {
				pod := item.(components.PodItem).Pod
				projectID := p.project.ID
				return p, func() tea.Msg {
					return messages.ChangePageMsg{Page: NewPodFormPage(projectID, &pod)}
				}
			}
		case key.Matches(msg, p.keys.DeletePod):
			if item := p.list.SelectedItem(); item != nil {
				pod := item.(components.PodItem).Pod
				return p, func() tea.Msg {
					return messages.ChangePageMsg{Page: NewPodDeletePage(&pod)}
				}
			}
		case key.Matches(msg, p.keys.EditProject):
			if p.project != nil {
				project := p.project
				return p, func() tea.Msg {
					return messages.ChangePageMsg{Page: NewProjectFormPage(project)}
				}
			}
		case key.Matches(msg, p.keys.DeleteProject):
			if p.project != nil {
				project := p.project
				return p, func() tea.Msg {
					return messages.ChangePageMsg{Page: NewProjectDeletePage(project)}
				}
			}
		}

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		// List height for card content
		listHeight := min((msg.Height-1)/2, 12)
		p.list.SetSize(46, listHeight)
		return p, nil

	case projectDetailDataMsg:
		p.project = &msg.project
		items := components.PodsToItems(msg.pods)
		cmd := p.list.SetItems(items)
		p.loading = false
		return p, cmd

	case projectDetailErrMsg:
		p.err = msg.err
		p.loading = false
		return p, nil

	case messages.PodCreatedMsg:
		newItem := components.PodItem{Pod: repo.Pod(msg)}
		cmd := p.list.InsertItem(len(p.list.Items()), newItem)
		return p, cmd

	case messages.PodUpdatedMsg:
		pod := msg
		items := p.list.Items()
		for i, item := range items {
			if pi, ok := item.(components.PodItem); ok && pi.ID == pod.ID {
				items[i] = components.PodItem{Pod: repo.Pod(pod)}
				break
			}
		}
		cmd := p.list.SetItems(items)
		return p, cmd

	case messages.PodDeleteMsg:
		pod := msg
		items := p.list.Items()
		for i, item := range items {
			if pi, ok := item.(components.PodItem); ok && pi.ID == pod.ID {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
		cmd := p.list.SetItems(items)
		return p, cmd

	case messages.ProjectUpdatedMsg:
		project := repo.Project(msg)
		p.project = &project
		return p, nil
	}

	// Pass other messages to the list
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
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
	if len(p.list.Items()) == 0 {
		podsContent = fmt.Sprintf("Pods (0)\n\n%s",
			styles.MutedStyle.Render("No pods yet. Press 'n' to create one."))
	} else {
		podsContent = p.list.View()
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
	if p.project != nil && p.project.Title != "" {
		return []string{"Projects", p.project.Title}
	}
	return []string{"Projects", "Detail"}
}

func loadProjectDetail(projectID string) tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return projectDetailErrMsg{err: err}
	}

	// Load project
	projectReq, err := http.NewRequest("GET", cfg.Server+"/api/projects/"+projectID, nil)
	if err != nil {
		return projectDetailErrMsg{err: err}
	}
	projectReq.Header.Set("Authorization", "Bearer "+cfg.Token)

	client := http.Client{}
	projectRes, err := client.Do(projectReq)
	if err != nil {
		return projectDetailErrMsg{err: err}
	}
	defer projectRes.Body.Close()

	var project repo.Project
	if err := json.NewDecoder(projectRes.Body).Decode(&project); err != nil {
		return projectDetailErrMsg{err: err}
	}

	// Load pods
	podsReq, err := http.NewRequest("GET", cfg.Server+"/api/pods/project/"+projectID, nil)
	if err != nil {
		return projectDetailErrMsg{err: err}
	}
	podsReq.Header.Set("Authorization", "Bearer "+cfg.Token)

	podsRes, err := client.Do(podsReq)
	if err != nil {
		return projectDetailErrMsg{err: err}
	}
	defer podsRes.Body.Close()

	var pods []repo.Pod
	if err := json.NewDecoder(podsRes.Body).Decode(&pods); err != nil {
		return projectDetailErrMsg{err: err}
	}

	return projectDetailDataMsg{
		project: project,
		pods:    pods,
	}
}
