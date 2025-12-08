package pages

import (
	"fmt"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/api"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type projectDeleteKeyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
}

func (k projectDeleteKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Cancel}
}

func (k projectDeleteKeyMap) FullHelp() [][]key.Binding {
	return nil
}

func (p ProjectDeletePage) HelpKeys() help.KeyMap {
	return p.keys
}

type projectToDelete struct {
	ID    string
	Title string
}

type ProjectDeletePage struct {
	project  projectToDelete
	podCount int
	keys     projectDeleteKeyMap
	width    int
	height   int
}

func NewProjectDeletePage(s msg.Store, project *repo.Project) ProjectDeletePage {
	podCount := 0
	for _, p := range s.Pods() {
		if p.ProjectID == project.ID {
			podCount++
		}
	}

	return ProjectDeletePage{
		project:  projectToDelete{ID: project.ID, Title: project.Title},
		podCount: podCount,
		keys: projectDeleteKeyMap{
			Confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
			Cancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		},
	}
}

func (p ProjectDeletePage) Init() tea.Cmd {
	return nil
}

func (p ProjectDeletePage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
			}
		case tea.KeyEnter:
			if p.podCount > 0 {
				return p, nil
			}
			return p, tea.Batch(
				api.DeleteProject(p.project.ID),
				func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
				},
			)
		}
	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}
	return p, nil
}

func (p ProjectDeletePage) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Render("Delete Project")

	name := lipgloss.NewStyle().
		PaddingTop(1).
		Render(p.project.Title)

	var hint string
	if p.podCount > 0 {
		hint = styles.MutedStyle().
			PaddingTop(1).
			Render(fmt.Sprintf("Delete all %d pods first", p.podCount))
	} else {
		hint = styles.MutedStyle().
			PaddingTop(1).
			Render("Press enter to confirm")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, name, hint)

	card := components.Card(components.CardProps{
		Width:   40,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p ProjectDeletePage) Breadcrumbs() []string {
	return []string{"Projects", "Delete"}
}
