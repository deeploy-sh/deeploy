package pages

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type podDetailKeyMap struct {
	EditPod   key.Binding
	DeletePod key.Binding
	Filter    key.Binding
	Back      key.Binding
}

func (k podDetailKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.EditPod, k.DeletePod, k.Filter, k.Back}
}

func (k podDetailKeyMap) FullHelp() [][]key.Binding {
	return nil
}

func newPodDetailKeyMap() podDetailKeyMap {
	return podDetailKeyMap{
		EditPod:   key.NewBinding(key.WithKeys("enter", "e"), key.WithHelp("enter", "edit pod")),
		DeletePod: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete pod")),
		Filter:    key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		Back:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

type PodDetailPage struct {
	pod     *repo.Pod
	keys    podDetailKeyMap
	help    help.Model
	loading bool
	width   int
	height  int
	err     error
}

type podDetailDataMsg struct {
	pod repo.Pod
}

type podDetailErrMsg struct{ err error }

func NewPodDetailPage(podID string, pod *repo.Pod) PodDetailPage {
	return PodDetailPage{
		keys:    newPodDetailKeyMap(),
		help:    styles.NewHelpModel(),
		loading: true,
		pod:     pod,
		// pod:     &repo.Pod{ID: podID},
	}
}

func (p PodDetailPage) Init() tea.Cmd {
	return nil
}

func (p PodDetailPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, p.keys.Back):
			// return p, func() tea.Msg {
			// 	return messages.ChangePageMsg{Page: NewDashboard()}
			// }
		case key.Matches(msg, p.keys.EditPod):
			// item := p.list.SelectedItem()
			// if item != nil {
			// 	pod := item.(components.PodItem).Pod
			// 	projectID := p.project.ID
			// 	return p, func() tea.Msg {
			// 		return messages.ChangePageMsg{Page: NewPodFormPage(projectID, &pod)}
			// 	}
			// }
		case key.Matches(msg, p.keys.DeletePod):
			// item := p.list.SelectedItem()
			// if item != nil {
			// 	pod := item.(components.PodItem).Pod
			// 	return p, func() tea.Msg {
			// 		return messages.ChangePageMsg{Page: NewPodDeletePage(&pod)}
			// 	}
			// }
		}

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		// List height for card content
		// listHeight := min((msg.Height-1)/2, 12)
		return p, nil

	case messages.PodUpdatedMsg:
		// pod := msg
		// items := p.list.Items()
		// for i, item := range items {
		// 	pi, ok := item.(components.PodItem)
		// 	if ok && pi.ID == pod.ID {
		// 		items[i] = components.PodItem{Pod: repo.Pod(pod)}
		// 		break
		// 	}
		// }
		// cmd := p.list.SetItems(items)
		return p, nil

	case messages.PodDeleteMsg:
		// pod := msg
		// items := p.list.Items()
		// for i, item := range items {
		// 	pi, ok := item.(components.PodItem)
		// 	if ok && pi.ID == pod.ID {
		// 		items = append(items[:i], items[i+1:]...)
		// 		break
		// 	}
		// }
		// cmd := p.list.SetItems(items)
		// return p, cmd
		return p, nil
	}

	// Pass other messages to the list
	var cmd tea.Cmd
	return p, cmd
}

func (p PodDetailPage) View() tea.View {
	helpView := p.help.View(p.keys)
	contentHeight := p.height - 1

	var content string

	// if p.loading {
	// 	content = styles.MutedStyle.Render("Loading...")
	// } else if p.err != nil {
	// 	content = styles.ErrorStyle.Render("Error: " + p.err.Error())
	// } else {
	// 	content = p.renderContent()
	// }

	content = p.pod.Title

	centered := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered, helpView))
}

func (p PodDetailPage) renderContent() string {
	// titleStyle := lipgloss.NewStyle().
	// 	Bold(true).
	// 	Foreground(styles.ColorForeground)

	// Project Header
	// header := titleStyle.Render(p.project.Title)
	// if p.project.Description != "" {
	// 	header += "\n" + styles.MutedStyle.Render(p.project.Description)
	// }
	//
	// // Pods list
	// var podsContent string
	// if len(p.list.Items()) == 0 {
	// 	podsContent = fmt.Sprintf("Pods (0)\n\n%s",
	// 		styles.MutedStyle.Render("No pods yet. Press 'n' to create one."))
	// } else {
	// 	podsContent = p.list.View()
	// }
	//
	// cardContent := lipgloss.JoinVertical(lipgloss.Left,
	// 	header,
	// 	"",
	// 	podsContent,
	// )
	//
	// card := components.Card(components.CardProps{
	// 	Width:   50,
	// 	Padding: []int{1, 2},
	// }).Render(cardContent)
	//
	// return card
	return ""
}

func (p PodDetailPage) Breadcrumbs() []string {
	// if p.project != nil && p.project.Title != "" {
	// 	return []string{"Pod", p.project.Title}
	// }
	return []string{"Pod", "Detail"}
}
