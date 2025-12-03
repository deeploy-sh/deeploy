package pages

import (
	"encoding/json"
	"fmt"
	"log"
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

type dashboardProjectsMsg []repo.Project
type dashboardErrMsg struct{ err error }

type dashboardKeyMap struct {
	Search key.Binding
	New    key.Binding
	Select key.Binding
	Filter key.Binding
}

func (k dashboardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.New, k.Select, k.Filter}
}

func (k dashboardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Search, k.New, k.Select, k.Filter}}
}

func newDashboardKeyMap() dashboardKeyMap {
	return dashboardKeyMap{
		Search: key.NewBinding(key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "search")),
		New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new project")),
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Filter: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
	}
}

type DashboardPage struct {
	list    list.Model
	keys    dashboardKeyMap
	help    help.Model
	loading bool
	width   int
	height  int
	err     error
}

func NewDashboard() DashboardPage {
	delegate := components.NewProjectDelegate(40)
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Projects"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorForeground)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return DashboardPage{
		list:    l,
		keys:    newDashboardKeyMap(),
		help:    styles.NewHelpModel(),
		loading: true,
	}
}

func (p DashboardPage) Init() tea.Cmd {
	return loadDashboardProjects
}

func (p DashboardPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Don't handle keys if filtering is active
		if p.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, p.keys.New):
			return p, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewProjectFormPage(nil)}
			}
		case key.Matches(msg, p.keys.Select):
			item := p.list.SelectedItem()
			if item != nil {
				i := item.(components.ProjectItem)
				return p, func() tea.Msg {
					return messages.ChangePageMsg{Page: NewProjectDetailPage(i.ID)}
				}
			}
		}

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		// List size for card content
		listHeight := min((msg.Height-1)/2, 15)
		p.list.SetSize(56, listHeight)
		log.Println(p.list.Height(), p.list.Width())
		return p, nil

	case dashboardProjectsMsg:
		items := make([]list.Item, len(msg))
		for i, p := range msg {
			items[i] = components.ProjectItem{Project: p}
		}
		cmd := p.list.SetItems(items)
		p.loading = false
		return p, cmd

	case dashboardErrMsg:
		p.err = msg.err
		p.loading = false
		return p, nil

	case messages.ProjectCreatedMsg:
		newItem := components.ProjectItem{Project: repo.Project(msg)}
		cmd := p.list.InsertItem(len(p.list.Items()), newItem)
		return p, cmd

	case messages.ProjectUpdatedMsg:
		project := msg
		items := p.list.Items()
		for i, item := range items {
			if pi, ok := item.(components.ProjectItem); ok && pi.Project.ID == project.ID {
				items[i] = components.ProjectItem{Project: repo.Project(project)}
				break
			}
		}
		cmd := p.list.SetItems(items)
		return p, cmd

	case messages.ProjectDeleteMsg:
		project := msg
		items := p.list.Items()
		for i, item := range items {
			if pi, ok := item.(components.ProjectItem); ok && pi.Project.ID == project.ID {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
		cmd := p.list.SetItems(items)
		return p, cmd
	}

	// Pass other messages to the list
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p DashboardPage) View() tea.View {
	helpView := p.help.View(p.keys)
	contentHeight := p.height - 1

	var content string

	if p.loading {
		content = styles.MutedStyle.Render("Loading projects...")
	} else if p.err != nil {
		content = styles.ErrorStyle.Render("Error: " + p.err.Error())
	} else if len(p.list.Items()) == 0 {
		content = p.renderEmptyState()
	} else {
		content = p.renderList()
	}

	centered := lipgloss.Place(p.width, contentHeight, lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered, helpView))
}

func (p DashboardPage) renderEmptyState() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorForeground).
		MarginBottom(1)

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("No projects yet"),
		styles.MutedStyle.Render("Press 'n' to create your first project"),
		"",
		styles.MutedStyle.Render("or use Ctrl+K to search"),
	)
}

func (p DashboardPage) renderList() string {
	card := components.Card(components.CardProps{
		Width:   60,
		Padding: []int{1, 2},
	}).Render(p.list.View())

	return card
}

func (p DashboardPage) Breadcrumbs() []string {
	return []string{"Dashboard"}
}

func loadDashboardProjects() tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return dashboardErrMsg{err: err}
	}

	r, err := http.NewRequest("GET", cfg.Server+"/api/projects", nil)
	if err != nil {
		return dashboardErrMsg{err: err}
	}
	r.Header.Set("Authorization", "Bearer "+cfg.Token)

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return dashboardErrMsg{err: err}
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return dashboardErrMsg{err: fmt.Errorf("unauthorized")}
	}

	var projects []repo.Project
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		return dashboardErrMsg{err: err}
	}

	return dashboardProjectsMsg(projects)
}
