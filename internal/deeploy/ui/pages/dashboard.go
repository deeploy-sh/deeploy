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

type dashboardProjectsMsg []repo.Project
type dashboardErrMsg struct{ err error }

type dashboardKeyMap struct {
	Search key.Binding
	New    key.Binding
	Select key.Binding
}

func (m DashboardPage) HelpKeys() help.KeyMap {
	return m.keys
}
func (k dashboardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.New, k.Select}
}

func (k dashboardKeyMap) FullHelp() [][]key.Binding {
	return nil
}

func newDashboardKeyMap() dashboardKeyMap {
	return dashboardKeyMap{
		Search: key.NewBinding(key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "search")),
		New:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new project")),
		Select: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}
}

type DashboardPage struct {
	store  Store
	list   components.ScrollList
	keys   dashboardKeyMap
	width  int
	height int
	err    error
}

func NewDashboard(s Store) DashboardPage {
	card := components.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	l := components.NewScrollList(components.ProjectsToItems(s.Projects()), components.ScrollListConfig{
		Width:  card.InnerWidth(),
		Height: 15,
	})

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
	// ScrollList handles navigation (Up/Down/Ctrl+N/P)
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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
		return m, cmd

	case dashboardProjectsMsg:
		items := make([]components.ScrollItem, len(msg))
		for i, p := range msg {
			items[i] = components.ProjectItem{Project: p}
		}
		m.list.SetItems(items)
		return m, cmd

	case dashboardErrMsg:
		m.err = msg.err
		return m, cmd

	case messages.ProjectCreatedMsg:
		items := m.list.Items()
		items = append(items, components.ProjectItem{Project: repo.Project(msg)})
		m.list.SetItems(items)
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
		m.list.SetItems(items)
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
		m.list.SetItems(items)
		return m, cmd
	}

	return m, cmd
}

func (m DashboardPage) View() tea.View {
	contentHeight := m.height

	var content string

	if m.err != nil {
		content = styles.ErrorStyle().Render("Error: " + m.err.Error())
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
		Foreground(styles.ColorForeground()).
		MarginBottom(1)

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("No projects yet"),
		styles.MutedStyle().Render("Press 'n' to create your first project"),
		"",
		styles.MutedStyle().Render("or use Ctrl+K to search"),
	)
}

func (m DashboardPage) renderList() string {
	card := components.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	w := card.InnerWidth()

	// Custom title (Bubbles built-in title has layout bugs)
	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render("Projects")

	// List with background (like Crush does)
	list := lipgloss.NewStyle().
		Width(w).
		Height(m.list.Height()).
		Background(styles.ColorBackgroundPanel()).
		Render(m.list.View())

	content := lipgloss.JoinVertical(lipgloss.Left, title, list)

	return components.Card(card).Render(content)
}

func (m DashboardPage) Breadcrumbs() []string {
	return []string{"Dashboard"}
}
