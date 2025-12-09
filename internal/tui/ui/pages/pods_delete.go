package pages

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
)

type podToDelete struct {
	ID        string
	Title     string
	ProjectID string
}

type PodDeletePage struct {
	pod        podToDelete
	input      textinput.Model
	keyConfirm key.Binding
	keyCancel  key.Binding
	width      int
	height     int
}

func (p PodDeletePage) HelpKeys() []key.Binding {
	return []key.Binding{p.keyConfirm, p.keyCancel}
}

func NewPodDeletePage(pod *repo.Pod) PodDeletePage {
	card := components.CardProps{Width: 40, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = pod.Title
	ti.Focus()
	ti.CharLimit = 100

	return PodDeletePage{
		pod:        podToDelete{ID: pod.ID, Title: pod.Title, ProjectID: pod.ProjectID},
		input:      ti,
		keyConfirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
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
