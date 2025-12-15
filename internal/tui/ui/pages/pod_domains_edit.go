package pages

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type PodDomainsEditPage struct {
	domain       api.PodDomain
	pod          *repo.Pod
	project      *repo.Project
	domainInput  textinput.Model
	portInput    textinput.Model
	sslEnabled   bool
	focusedInput int
	loading      bool
	keySave      key.Binding
	keyTab       key.Binding
	keyToggle    key.Binding
	keyCancel    key.Binding
	width        int
	height       int
}

func (p PodDomainsEditPage) HelpKeys() []key.Binding {
	return []key.Binding{p.keySave, p.keyTab, p.keyToggle, p.keyCancel}
}

func NewPodDomainsEditPage(domain api.PodDomain, pod *repo.Pod, project *repo.Project) PodDomainsEditPage {
	domainInput := components.NewTextInput(40)
	domainInput.Placeholder = "app.example.com"
	domainInput.CharLimit = 100
	domainInput.SetValue(domain.Domain)
	domainInput.Focus()

	portInput := components.NewTextInput(10)
	portInput.Placeholder = "8080"
	portInput.CharLimit = 5
	portInput.SetValue(strconv.Itoa(domain.Port))

	return PodDomainsEditPage{
		domain:      domain,
		pod:         pod,
		project:     project,
		domainInput: domainInput,
		portInput:   portInput,
		sslEnabled:  domain.SSLEnabled,
		keySave:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
		keyTab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
		keyToggle:   key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle SSL")),
		keyCancel:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (p PodDomainsEditPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p PodDomainsEditPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.PodDomainUpdated:
		pod := p.pod
		project := p.project
		return p, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewPodDomainsPage(pod, project) },
			}
		}

	case tea.KeyPressMsg:
		switch {
		case key.Matches(tmsg, p.keyCancel):
			pod := p.pod
			project := p.project
			return p, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewPodDomainsPage(pod, project) },
				}
			}

		case key.Matches(tmsg, p.keyTab):
			p.focusedInput = (p.focusedInput + 1) % 3
			p.domainInput.Blur()
			p.portInput.Blur()
			switch p.focusedInput {
			case 0:
				p.domainInput.Focus()
			case 1:
				p.portInput.Focus()
			}
			return p, nil

		case key.Matches(tmsg, p.keyToggle):
			if p.focusedInput == 2 {
				p.sslEnabled = !p.sslEnabled
			}
			return p, nil

		case key.Matches(tmsg, p.keySave):
			domain := strings.TrimSpace(p.domainInput.Value())
			if domain == "" {
				return p, nil
			}

			port := 8080
			if pVal, err := strconv.Atoi(p.portInput.Value()); err == nil && pVal > 0 {
				port = pVal
			}

			p.loading = true
			return p, api.UpdatePodDomain(p.pod.ID, p.domain.ID, domain, port, p.sslEnabled)
		}

	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}

	var cmd tea.Cmd
	switch p.focusedInput {
	case 0:
		p.domainInput, cmd = p.domainInput.Update(tmsg)
	case 1:
		p.portInput, cmd = p.portInput.Update(tmsg)
	}
	return p, cmd
}

func (p PodDomainsEditPage) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Edit Domain"))
	b.WriteString("\n\n")

	if p.loading {
		b.WriteString("Saving...")
		return tea.NewView(p.centeredCard(b.String()))
	}

	// Type badge (read-only)
	typeLabel := "custom"
	if p.domain.Type == "auto" {
		typeLabel = "auto"
	}
	b.WriteString(styles.MutedStyle().Render(fmt.Sprintf("Type: %s", typeLabel)))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Width(12)
	activeLabel := lipgloss.NewStyle().Width(12).Foreground(styles.ColorPrimary())

	// Domain field
	if p.focusedInput == 0 {
		b.WriteString(activeLabel.Render("Domain:"))
	} else {
		b.WriteString(labelStyle.Render("Domain:"))
	}
	b.WriteString(p.domainInput.View())
	b.WriteString("\n\n")

	// Port field
	if p.focusedInput == 1 {
		b.WriteString(activeLabel.Render("Port:"))
	} else {
		b.WriteString(labelStyle.Render("Port:"))
	}
	b.WriteString(p.portInput.View())
	b.WriteString("\n\n")

	// SSL toggle
	if p.focusedInput == 2 {
		b.WriteString(activeLabel.Render("SSL:"))
	} else {
		b.WriteString(labelStyle.Render("SSL:"))
	}
	if p.sslEnabled {
		b.WriteString("[x] Enabled")
	} else {
		b.WriteString("[ ] Disabled")
	}
	b.WriteString("\n")

	return tea.NewView(p.centeredCard(b.String()))
}

func (p PodDomainsEditPage) centeredCard(content string) string {
	card := styles.Card(styles.CardProps{
		Width:   60,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)
}

func (p PodDomainsEditPage) Breadcrumbs() []string {
	return []string{"Projects", p.project.Title, "Pods", p.pod.Title, "Domains", "Edit"}
}
