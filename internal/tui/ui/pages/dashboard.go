package pages

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type DashboardPage struct {
	store     msg.Store
	list      components.ScrollList
	keyNew    key.Binding
	keyEdit   key.Binding
	keyDelete key.Binding
	keySelect key.Binding
	width     int
	height    int
}

func (m DashboardPage) HelpKeys() []key.Binding {
	return []key.Binding{m.keyNew, m.keyEdit, m.keyDelete, m.keySelect}
}

func NewDashboard(s msg.Store) DashboardPage {
	card := styles.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	l := components.NewScrollList(components.ProjectsToItems(s.Projects(), s.Pods()), components.ScrollListConfig{
		Width:  card.InnerWidth(),
		Height: 15,
	})

	return DashboardPage{
		store:     s,
		list:      l,
		keyNew:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new project")),
		keyEdit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit project")),
		keyDelete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete project")),
		keySelect: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}
}

func (m DashboardPage) Init() tea.Cmd {
	return nil
}

func (m DashboardPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(tmsg)

	switch tmsg := tmsg.(type) {
	case msg.DataLoaded:
		m.list.SetItems(components.ProjectsToItems(tmsg.Projects, tmsg.Pods))
		return m, cmd

	case tea.KeyPressMsg:
		switch {
		case key.Matches(tmsg, m.keyNew):
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewProjectFormPage(nil) },
				}
			}
		case key.Matches(tmsg, m.keyEdit):
			item := m.list.SelectedItem()
			if item != nil {
				project := item.(components.ProjectItem).Project
				return m, func() tea.Msg {
					return msg.ChangePage{
						PageFactory: func(s msg.Store) tea.Model { return NewProjectFormPage(&project) },
					}
				}
			}
		case key.Matches(tmsg, m.keyDelete):
			item := m.list.SelectedItem()
			if item != nil {
				project := item.(components.ProjectItem).Project
				return m, func() tea.Msg {
					return msg.ChangePage{
						PageFactory: func(s msg.Store) tea.Model { return NewProjectDeletePage(s, &project) },
					}
				}
			}
		case key.Matches(tmsg, m.keySelect):
			item := m.list.SelectedItem()
			if item != nil {
				i := item.(components.ProjectItem)
				return m, func() tea.Msg {
					return msg.ChangePage{
						PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, i.ID) },
					}
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

func (m DashboardPage) View() tea.View {
	contentHeight := m.height

	var content string

	if len(m.list.Items()) == 0 {
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
	)
}

func (m DashboardPage) renderList() string {
	card := styles.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	w := card.InnerWidth()

	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render("Projects")

	list := lipgloss.NewStyle().
		Width(w).
		Height(m.list.Height()).
		Background(styles.ColorBackgroundPanel()).
		Render(m.list.View())

	content := lipgloss.JoinVertical(lipgloss.Left, title, list)

	return styles.Card(card).Render(content)
}

func (m DashboardPage) Breadcrumbs() []string {
	return []string{"Dashboard"}
}
