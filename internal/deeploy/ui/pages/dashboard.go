package pages

import (
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

type dashboardProjectsMsg []repo.Project
type dashboardErrMsg struct{ err error }

type dashboardKeyMap struct {
	Search key.Binding
	New    key.Binding
	Select key.Binding
	Filter key.Binding
}

func (m DashboardPage) HelpKeys() help.KeyMap {
	return m.keys // gibt einfach den existierenden KeyMap zurück
}
func (k dashboardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.New, k.Select, k.Filter}
}

func (k dashboardKeyMap) FullHelp() [][]key.Binding {
	return nil
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
	store  Store
	list   list.Model
	keys   dashboardKeyMap
	width  int
	height int
	err    error
}

func NewDashboard(s Store) DashboardPage {
	delegate := components.NewProjectDelegate(40)
	l := list.New(components.ProjectsToItems(s.Projects()), delegate, 0, 0)
	l.Title = "Projects"
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorForeground)
	l.SetShowTitle(true)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	return DashboardPage{
		store: s,
		list:  l,
		keys:  newDashboardKeyMap(),
	}
}

func (m DashboardPage) Init() tea.Cmd {
	return nil
}

func (m DashboardPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Don't handle keys if filtering is active
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.New):
			return m, func() tea.Msg {
				return ChangePageMsg{
					PageFactory: func(s Store) tea.Model { return NewProjectFormPage(nil) },
				}
			}
		case key.Matches(msg, m.keys.Select):
			item := m.list.SelectedItem()
			if item != nil {
				i := item.(components.ProjectItem)
				return m, func() tea.Msg {
					return ChangePageMsg{
						PageFactory: func(s Store) tea.Model { return NewProjectDetailPage(s, i.ID) },
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// List size for card content
		listHeight := min((msg.Height-1)/2, 15)
		m.list.SetSize(56, listHeight)
		return m, nil

	case dashboardProjectsMsg:
		items := make([]list.Item, len(msg))
		for i, p := range msg {
			items[i] = components.ProjectItem{Project: p}
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case dashboardErrMsg:
		m.err = msg.err
		return m, nil

	case messages.ProjectCreatedMsg:
		newItem := components.ProjectItem{Project: repo.Project(msg)}
		cmd := m.list.InsertItem(len(m.list.Items()), newItem)
		return m, cmd

	case messages.ProjectUpdatedMsg:
		project := msg
		items := m.list.Items()
		for i, item := range items {
			if pi, ok := item.(components.ProjectItem); ok && pi.Project.ID == project.ID {
				items[i] = components.ProjectItem{Project: repo.Project(project)}
				break
			}
		}
		cmd := m.list.SetItems(items)
		return m, cmd

	case messages.ProjectDeleteMsg:
		project := msg
		items := m.list.Items()
		for i, item := range items {
			if pi, ok := item.(components.ProjectItem); ok && pi.Project.ID == project.ID {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
		cmd := m.list.SetItems(items)
		return m, cmd
	}

	// Pass other messages to the list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m DashboardPage) View() tea.View {
	contentHeight := m.height

	var content string

	if m.err != nil {
		content = styles.ErrorStyle.Render("Error: " + m.err.Error())
	} else if len(m.list.Items()) == 0 {
		content = m.renderEmptyState()
	} else {
		content = m.renderList()
	}

	centered := lipgloss.Place(m.width, contentHeight, lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered))
}

func (m DashboardPage) renderEmptyState() string {
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

func (m DashboardPage) renderList() string {
	card := components.Card(components.CardProps{
		Width:   50,
		Padding: []int{1, 2},
	}).Render(m.list.View())

	return card
}

func (m DashboardPage) Breadcrumbs() []string {
	return []string{"Dashboard"}
}
