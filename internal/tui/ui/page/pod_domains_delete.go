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

type podDomainsDelete struct {
	domain     model.PodDomain
	pod        *model.Pod
	project    *model.Project
	input      textinput.Model
	keyConfirm key.Binding
	keyCancel  key.Binding
	width      int
	height     int
}

func (p podDomainsDelete) HelpKeys() []key.Binding {
	return []key.Binding{p.keyConfirm, p.keyCancel}
}

func NewPodDomainsDelete(domain model.PodDomain, pod *model.Pod, project *model.Project) podDomainsDelete {
	card := styles.CardProps{Width: styles.CardWidthMD, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = domain.Domain
	ti.Focus()
	ti.CharLimit = 100

	return podDomainsDelete{
		domain:     domain,
		pod:        pod,
		project:    project,
		input:      ti,
		keyConfirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (p podDomainsDelete) Init() tea.Cmd {
	return textinput.Blink
}

func (p podDomainsDelete) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			pod := p.pod
			project := p.project
			return p, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewPodDomains(s, pod, project) },
				}
			}
		case tea.KeyEnter:
			// Only delete if input matches domain exactly
			if p.input.Value() != p.domain.Domain {
				return p, nil
			}
			return p, tea.Batch(
				func() tea.Msg { return msg.StartLoading{Text: "Deleting domain"} },
				api.DeletePodDomain(p.pod.ID, p.domain.ID),
			)
		}

	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}

	var cmd tea.Cmd
	p.input, cmd = p.input.Update(tmsg)
	return p, cmd
}

func (p podDomainsDelete) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary()).
		Render(fmt.Sprintf("Delete Domain"))

	domainName := lipgloss.NewStyle().
		Bold(true).
		Render(p.domain.Domain)

	hint := styles.MutedStyle().
		Render("Type '" + p.domain.Domain + "' to confirm")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		domainName,
		"",
		hint,
		"",
		p.input.View(),
	)

	card := styles.Card(styles.CardProps{
		Width:   styles.CardWidthMD,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p podDomainsDelete) Breadcrumbs() []string {
	return []string{"Projects", p.project.Title, "Pods", p.pod.Title, "Domains", "Delete"}
}
