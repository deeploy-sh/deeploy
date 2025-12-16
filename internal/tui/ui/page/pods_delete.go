package page

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type podToDelete struct {
	ID        string
	Title     string
	ProjectID string
}

type podDelete struct {
	pod        podToDelete
	input      textinput.Model
	keyConfirm key.Binding
	keyCancel  key.Binding
	width      int
	height     int
}

func (p podDelete) HelpKeys() []key.Binding {
	return []key.Binding{p.keyConfirm, p.keyCancel}
}

func NewPodDelete(pod *model.Pod) podDelete {
	card := styles.CardProps{Width: 40, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = pod.Title
	ti.Focus()
	ti.CharLimit = 100

	return podDelete{
		pod:        podToDelete{ID: pod.ID, Title: pod.Title, ProjectID: pod.ProjectID},
		input:      ti,
		keyConfirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (p podDelete) Init() tea.Cmd {
	return textinput.Blink
}

func (p podDelete) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case msg.PodDeleted:
		projectID := p.pod.ProjectID
		return p, tea.Batch(
			api.LoadData(),
			func() tea.Msg { return msg.ShowStatus{Text: "Pod deleted", Type: msg.StatusSuccess} },
			func() tea.Msg { return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetail(s, projectID) }} },
		)

	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			projectID := p.pod.ProjectID
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetail(s, projectID) }}
			}
		case tea.KeyEnter:
			// Only delete if input matches pod title exactly
			if p.input.Value() != p.pod.Title {
				return p, nil
			}
			return p, api.DeletePod(p.pod.ID)
		}
	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}

	p.input, cmd = p.input.Update(tmsg)
	return p, cmd
}

func (p podDelete) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Render(fmt.Sprintf("Delete Pod (%v)", p.pod.Title))

	hint := styles.MutedStyle().
		PaddingTop(1).
		PaddingBottom(1).
		Render("Type '" + p.pod.Title + "' to confirm")

	content := lipgloss.JoinVertical(lipgloss.Left, title, hint, p.input.View())

	card := styles.Card(styles.CardProps{
		Width:   40,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p podDelete) Breadcrumbs() []string {
	return []string{"Projects", "Pods", "Delete"}
}
