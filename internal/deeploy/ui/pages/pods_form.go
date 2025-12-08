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

type PodFormPage struct {
	titleInput textinput.Model
	keys       formKeyMap
	projectID  string
	pod        *repo.Pod
	width      int
	height     int
}

func (p PodFormPage) HelpKeys() help.KeyMap {
	return p.keys
}

func NewPodFormPage(projectID string, pod *repo.Pod) PodFormPage {
	card := components.CardProps{Width: 40, Padding: []int{1, 2}, Accent: true}
	titleInput := components.NewTextInput(card.InnerWidth())
	titleInput.Focus()
	titleInput.Placeholder = "Title"
	if pod != nil {
		titleInput.SetValue(pod.Title)
	}

	podFormPage := PodFormPage{
		titleInput: titleInput,
		keys:       newFormKeyMap(),
		projectID:  projectID,
	}
	if pod != nil {
		podFormPage.pod = pod
	}
	return podFormPage
}

func (p PodFormPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p PodFormPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			projectID := p.projectID
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
			}
		}
		switch {
		case key.Matches(tmsg, p.keys.Save):
			if len(p.titleInput.Value()) > 0 {
				projectID := p.projectID
				var apiCmd tea.Cmd
				if p.pod != nil {
					p.pod.Title = p.titleInput.Value()
					apiCmd = api.UpdatePod(p.pod)
				} else {
					apiCmd = api.CreatePod(p.titleInput.Value(), p.projectID)
				}
				return p, tea.Batch(
					apiCmd,
					func() tea.Msg {
						return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
					},
				)
			}
		}

	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}

	p.titleInput, cmd = p.titleInput.Update(tmsg)
	return p, cmd
}

func (p PodFormPage) View() tea.View {
	var titleText string
	if p.pod != nil {
		titleText = "Edit Pod"
	} else {
		titleText = "New Pod"
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Background(styles.ColorBackgroundPanel()).
		Render(titleText)

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		p.titleInput.View(),
	)

	card := components.Card(components.CardProps{
		Width:   40,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p PodFormPage) Breadcrumbs() []string {
	if p.pod != nil {
		return []string{"Projects", "Pods", "Edit"}
	}
	return []string{"Projects", "Pods", "New"}
}

func (p PodFormPage) HasFocusedInput() bool {
	return p.titleInput.Focused()
}
