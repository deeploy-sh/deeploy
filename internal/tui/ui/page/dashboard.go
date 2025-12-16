package page

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type dashboard struct {
	store     msg.Store
	list      components.ScrollList
	keyNew    key.Binding
	keyDelete key.Binding
	keySelect key.Binding
	width     int
	height    int
}

func (m dashboard) HelpKeys() []key.Binding {
	return []key.Binding{m.keyNew, m.keyDelete, m.keySelect}
}

func NewDashboard(s msg.Store) dashboard {
	card := styles.CardProps{Width: 50, Padding: []int{1, 1}, Accent: true}
	l := components.NewScrollList(components.ProjectsToItems(s.Projects(), s.Pods()), components.ScrollListConfig{
		Width:  card.InnerWidth(),
		Height: 15,
	})

	return dashboard{
		store:     s,
		list:      l,
		keyNew:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new project")),
		keyDelete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete project")),
		keySelect: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}
}

func (m dashboard) Init() tea.Cmd {
	return nil
}

func (m dashboard) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
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
					PageFactory: func(s msg.Store) tea.Model { return NewProjectForm(nil) },
				}
			}
		case key.Matches(tmsg, m.keyDelete):
			item := m.list.SelectedItem()
			if item != nil {
				project := item.(components.ProjectItem).Project
				return m, func() tea.Msg {
					return msg.ChangePage{
						PageFactory: func(s msg.Store) tea.Model { return NewProjectDelete(s, &project) },
					}
				}
			}
		case key.Matches(tmsg, m.keySelect):
			item := m.list.SelectedItem()
			if item != nil {
				i := item.(components.ProjectItem)
				return m, func() tea.Msg {
					return msg.ChangePage{
						PageFactory: func(s msg.Store) tea.Model { return NewProjectDetail(s, i.ID) },
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

func (m dashboard) View() tea.View {
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

func (m dashboard) renderEmptyState() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorForeground()).
		MarginBottom(1)

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("No projects yet"),
		styles.MutedStyle().Render("Press 'n' to create your first project"),
	)
}

func (m dashboard) renderList() string {
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

func (m dashboard) Breadcrumbs() []string {
	return []string{"dashboard"}
}
