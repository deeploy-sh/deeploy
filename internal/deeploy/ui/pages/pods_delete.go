package pages

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/api"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type podDeleteKeyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
}

func (k podDeleteKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Cancel}
}

func (k podDeleteKeyMap) FullHelp() [][]key.Binding {
	return nil
}

type podToDelete struct {
	ID        string
	Title     string
	ProjectID string
}

type PodDeletePage struct {
	pod    podToDelete
	input  textinput.Model
	keys   podDeleteKeyMap
	width  int
	height int
}

func (p PodDeletePage) HelpKeys() help.KeyMap {
	return p.keys
}

func NewPodDeletePage(pod *repo.Pod) PodDeletePage {
	ti := textinput.New()
	ti.Placeholder = pod.Title
	ti.Focus()
	ti.CharLimit = 100

	return PodDeletePage{
		pod:   podToDelete{ID: pod.ID, Title: pod.Title, ProjectID: pod.ProjectID},
		input: ti,
		keys: podDeleteKeyMap{
			Confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
			Cancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		},
	}
}

func (p PodDeletePage) Init() tea.Cmd {
	return textinput.Blink
}

func (p PodDeletePage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			projectID := p.pod.ProjectID
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
			}
		case tea.KeyEnter:
			projectID := p.pod.ProjectID
			// Only delete if input matches pod title exactly
			if p.input.Value() != p.pod.Title {
				return p, nil
			}
			return p, tea.Batch(
				api.DeletePod(p.pod.ID),
				func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
				},
			)
		}
	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}

	p.input, cmd = p.input.Update(tmsg)
	return p, cmd
}

func (p PodDeletePage) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Render("Delete Pod")

	name := lipgloss.NewStyle().
		PaddingTop(1).
		Render(p.pod.Title)

	hint := styles.MutedStyle().
		PaddingTop(1).
		PaddingBottom(1).
		Render("Type '" + p.pod.Title + "' to confirm")

	content := lipgloss.JoinVertical(lipgloss.Left, title, name, hint, p.input.View())

	card := components.Card(components.CardProps{
		Width:   40,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p PodDeletePage) Breadcrumbs() []string {
	return []string{"Projects", "Pods", "Delete"}
}
