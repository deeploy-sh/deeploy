package page

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/tui/utils"
)

// domainItem wraps PodDomain to implement ScrollItem interface
type domainItem struct {
	domain model.PodDomain
}

func (d domainItem) Title() string       { return d.domain.Domain }
func (d domainItem) FilterValue() string { return d.domain.Domain }
func (d domainItem) Suffix() string {
	badges := fmt.Sprintf(":%d", d.domain.Port)
	if d.domain.Type == "auto" {
		badges += " auto"
	}
	return badges
}

var podDomainsCard = styles.CardProps{Width: styles.CardWidthMD, Padding: []int{1, 2}, Accent: true}

type podDomains struct {
	pod       *model.Pod
	project   *model.Project
	domains   components.ScrollList
	keyAdd    key.Binding
	keyAuto   key.Binding
	keyEdit   key.Binding
	keyDelete key.Binding
	keyOpen   key.Binding
	keyBack   key.Binding
	width     int
	height    int
}

func (m podDomains) HelpKeys() []key.Binding {
	return []key.Binding{m.keyAdd, m.keyAuto, m.keyEdit, m.keyDelete, m.keyOpen, m.keyBack}
}

func NewPodDomains(s msg.Store, pod *model.Pod, project *model.Project) podDomains {
	rawDomains := s.PodDomains(pod.ID)
	items := make([]components.ScrollItem, len(rawDomains))
	for i, d := range rawDomains {
		items[i] = domainItem{domain: d}
	}

	return podDomains{
		pod:       pod,
		project:   project,
		domains:   components.NewScrollList(items, components.ScrollListConfig{Width: podDomainsCard.InnerWidth(), Height: 8}),
		keyAdd:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new custom")),
		keyAuto:   key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "generate auto")),
		keyEdit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		keyDelete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		keyOpen:   key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open")),
		keyBack:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

func (m podDomains) Init() tea.Cmd {
	return nil
}

func (m podDomains) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(tmsg)

	case tea.MouseWheelMsg:
		m.domains, _ = m.domains.Update(tmsg)

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
	}

	return m, nil
}

func (m podDomains) handleKeyPress(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		podID := m.pod.ID
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDetail(s, podID)
				},
			}
		}

	case key.Matches(tmsg, m.keyAdd):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDomainsForm(pod, project, nil, false)
				},
			}
		}

	case key.Matches(tmsg, m.keyAuto):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDomainsForm(pod, project, nil, true)
				},
			}
		}

	case key.Matches(tmsg, m.keyEdit):
		if item := m.domains.SelectedItem(); item != nil {
			domain := item.(domainItem).domain
			pod := m.pod
			project := m.project
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewPodDomainsForm(pod, project, &domain, false)
					},
				}
			}
		}

	case key.Matches(tmsg, m.keyDelete):
		if item := m.domains.SelectedItem(); item != nil {
			domain := item.(domainItem).domain
			pod := m.pod
			project := m.project
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewPodDomainsDelete(domain, pod, project)
					},
				}
			}
		}

	case key.Matches(tmsg, m.keyOpen):
		if item := m.domains.SelectedItem(); item != nil {
			return m, utils.OpenBrowserCmd(item.(domainItem).domain.URL)
		}
	}

	// Let ScrollList handle navigation (up/down/j/k/mouse)
	m.domains, _ = m.domains.Update(tmsg)
	return m, nil
}

func (m podDomains) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Domains"))
	b.WriteString("\n")
	b.WriteString(styles.MutedStyle().Render("Configure domains for " + m.pod.Title))
	b.WriteString("\n\n")

	if len(m.domains.Items()) == 0 {
		b.WriteString(styles.MutedStyle().Render("No domains configured."))
		b.WriteString("\n\n")
		b.WriteString(styles.MutedStyle().Render("Press 'g' to generate an auto domain, or 'n' to add a custom one."))
		b.WriteString("\n")
		b.WriteString(styles.MutedStyle().Render("A domain is required before you can deploy."))
	} else {
		b.WriteString(m.domains.View())
	}

	card := styles.Card(podDomainsCard).Render(b.String())
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m podDomains) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title, "Domains"}
}
