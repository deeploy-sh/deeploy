package page

import (
	"fmt"
	"strconv"
	"strings"

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

type podDomainsEdit struct {
	domain       model.PodDomain
	pod          *model.Pod
	project      *model.Project
	domainInput  textinput.Model
	portInput    textinput.Model
	focusedInput int
	keySave      key.Binding
	keyTab       key.Binding
	keyCancel    key.Binding
	width        int
	height       int
	// Note: SSL toggle removed - SSL is automatic in production via Let's Encrypt
}

func (p podDomainsEdit) HelpKeys() []key.Binding {
	return []key.Binding{p.keySave, p.keyTab, p.keyCancel}
}

func NewPodDomainsEdit(domain model.PodDomain, pod *model.Pod, project *model.Project) podDomainsEdit {
	domainInput := components.NewTextInput(40)
	domainInput.Placeholder = "app.example.com"
	domainInput.CharLimit = 100
	domainInput.SetValue(domain.Domain)
	domainInput.Focus()

	portInput := components.NewTextInput(10)
	portInput.Placeholder = "8080"
	portInput.CharLimit = 5
	portInput.SetValue(strconv.Itoa(domain.Port))

	return podDomainsEdit{
		domain:      domain,
		pod:         pod,
		project:     project,
		domainInput: domainInput,
		portInput:   portInput,
		keySave:     key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyTab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
		keyCancel:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		// Note: SSL toggle removed - SSL is automatic in production
	}
}

func (p podDomainsEdit) Init() tea.Cmd {
	return textinput.Blink
}

func (p podDomainsEdit) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.PodDomainUpdated:
		pod := p.pod
		project := p.project
		return p, tea.Batch(
			func() tea.Msg { return msg.ShowStatus{Text: "Saved. Restart or deploy to apply.", Type: msg.StatusSuccess} },
			func() tea.Msg { return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodDomains(pod, project) }} },
		)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(tmsg, p.keyCancel):
			pod := p.pod
			project := p.project
			return p, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewPodDomains(pod, project) },
				}
			}

		case key.Matches(tmsg, p.keyTab):
			// Toggle between domain and port (no SSL toggle anymore)
			p.focusedInput = (p.focusedInput + 1) % 2
			p.domainInput.Blur()
			p.portInput.Blur()
			switch p.focusedInput {
			case 0:
				p.domainInput.Focus()
			case 1:
				p.portInput.Focus()
			}
			return p, nil

		case key.Matches(tmsg, p.keySave):
			domain := strings.TrimSpace(p.domainInput.Value())
			if domain == "" {
				return p, nil
			}

			port := 8080
			pVal, err := strconv.Atoi(p.portInput.Value())
			if err == nil && pVal > 0 {
				port = pVal
			}

			// SSL is always enabled - it's automatic in production via Let's Encrypt
			return p, tea.Batch(
				func() tea.Msg { return msg.StartLoading{Text: "Updating domain"} },
				api.UpdatePodDomain(p.pod.ID, p.domain.ID, domain, port, true),
			)
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

func (p podDomainsEdit) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Edit Domain"))
	b.WriteString("\n\n")

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

	// SSL info (no toggle - SSL is automatic in production)
	b.WriteString(labelStyle.Render("SSL:"))
	b.WriteString(styles.MutedStyle().Render("Automatic (Let's Encrypt)"))
	b.WriteString("\n")

	return tea.NewView(p.centeredCard(b.String()))
}

func (p podDomainsEdit) centeredCard(content string) string {
	card := styles.Card(styles.CardProps{
		Width:   60,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	return lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)
}

func (p podDomainsEdit) Breadcrumbs() []string {
	return []string{"Projects", p.project.Title, "Pods", p.pod.Title, "Domains", "Edit"}
}
