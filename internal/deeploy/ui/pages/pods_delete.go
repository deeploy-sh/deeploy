package pages

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/api"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type PodDeletePage struct {
	pod      *repo.Pod
	decision int
	keys     deleteKeyMap
	width    int
	height   int
}

func (p PodDeletePage) HelpKeys() help.KeyMap {
	return p.keys
}

func NewPodDeletePage(pod *repo.Pod) PodDeletePage {
	return PodDeletePage{
		pod:  pod,
		keys: newDeleteKeyMap(),
	}
}

func (p PodDeletePage) Init() tea.Cmd {
	return nil
}

func (p PodDeletePage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyLeft, 'h':
			p.decision = confirmNo
			return p, nil
		case tea.KeyRight, 'l':
			p.decision = confirmYes
			return p, nil
		case tea.KeyTab:
			if p.decision == confirmNo {
				p.decision = confirmYes
			} else {
				p.decision = confirmNo
			}
		case tea.KeyEscape:
			projectID := p.pod.ProjectID
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
			}
		case tea.KeyEnter:
			projectID := p.pod.ProjectID
			if p.decision == confirmNo {
				return p, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
				}
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
	return p, nil
}

func (p PodDeletePage) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 0, 1, 0).
		Render("Delete " + p.pod.Title + "?")

	baseButton := lipgloss.NewStyle().
		Padding(0, 3).
		Width(1).
		MarginRight(1)

	activeButton := baseButton.
		Background(styles.ColorPrimary()).
		Foreground(lipgloss.Color("0"))

	inactiveButton := baseButton.
		Background(lipgloss.Color("237"))

	var noButton, yesButton string
	if p.decision == confirmNo {
		noButton = activeButton.Render("NO")
		yesButton = inactiveButton.Render("YES")
	} else {
		noButton = inactiveButton.Render("NO")
		yesButton = activeButton.Render("YES")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, noButton, yesButton)
	content := lipgloss.JoinVertical(0.5, title, buttons)

	contentHeight := p.height

	centered := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(centered)
}

func (p PodDeletePage) Breadcrumbs() []string {
	return []string{"Projects", "Pods", "Delete"}
}
