package pages

import (
	"log"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type PodDeletePage struct {
	pod      *repo.Pod
	decision int
	keys     deleteKeyMap
	help     help.Model
	width    int
	height   int
}

func NewPodDeletePage(pod *repo.Pod) PodDeletePage {
	return PodDeletePage{
		pod:  pod,
		keys: newDeleteKeyMap(),
		help: styles.NewHelpModel(),
	}
}

func (p PodDeletePage) Init() tea.Cmd {
	return nil
}

func (p PodDeletePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
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
				return messages.ChangePageMsg{Page: NewProjectDetailPage(projectID)}
			}
		case tea.KeyEnter:
			projectID := p.pod.ProjectID
			if p.decision == confirmNo {
				return p, func() tea.Msg {
					return messages.ChangePageMsg{Page: NewProjectDetailPage(projectID)}
				}
			}
			return p, tea.Batch(
				p.DeletePod,
				func() tea.Msg {
					return messages.ChangePageMsg{Page: NewProjectDetailPage(projectID)}
				},
			)
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	}
	return p, nil
}

func (p PodDeletePage) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 0, 1, 0).
		Render("Delete " + p.pod.Title + "?")

	// Button Styles
	baseButton := lipgloss.NewStyle().
		Padding(0, 3).
		Width(1).
		MarginRight(1)

	activeButton := baseButton.
		Background(styles.ColorPrimary).
		Foreground(lipgloss.Color("0"))

	inactiveButton := baseButton.
		Background(lipgloss.Color("237"))

	// Render Buttons based on decision
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

	helpView := p.help.View(p.keys)
	contentHeight := p.height - 1 // 1 für help

	centered := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered, helpView))
}

func (p PodDeletePage) Breadcrumbs() []string {
	return []string{"Projects", "Pods", "Delete"}
}

func (p PodDeletePage) DeletePod() tea.Msg {
	res, err := utils.Request("DELETE", "/pods/"+p.pod.ID, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	return messages.PodDeleteMsg(p.pod)
}
