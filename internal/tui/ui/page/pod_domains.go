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
	"github.com/deeploy-sh/deeploy/internal/tui/utils"
)

type podDomainsMode int

const (
	modeDomainList podDomainsMode = iota
	modeDomainAdd
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
	pod          *model.Pod
	project      *model.Project
	domains      []model.PodDomain
	list         components.ScrollList
	mode         podDomainsMode
	domainInput  textinput.Model
	portInput    textinput.Model
	isAuto       bool
	focusedInput int
	keyAdd       key.Binding
	keyAuto      key.Binding
	keyEdit      key.Binding
	keyDelete    key.Binding
	keyOpen      key.Binding
	keyBack      key.Binding
	keySave      key.Binding
	keyTab       key.Binding
	width        int
	height       int
	// Note: SSL toggle removed - SSL is now always enabled automatically
	// via Let's Encrypt in production (see docker.go RunContainer)
}

func (m podDomains) HelpKeys() []key.Binding {
	if m.mode == modeDomainAdd {
		return []key.Binding{m.keySave, m.keyTab, m.keyBack}
	}
	return []key.Binding{m.keyAdd, m.keyAuto, m.keyEdit, m.keyDelete, m.keyOpen, m.keyBack}
}

func NewPodDomains(s msg.Store, pod *model.Pod, project *model.Project) podDomains {
	domainInput := components.NewTextInput(40)
	domainInput.Placeholder = "app.example.com"
	domainInput.CharLimit = 100

	portInput := components.NewTextInput(10)
	portInput.Placeholder = "8080"
	portInput.CharLimit = 5

	// Get domains from store
	domains := s.PodDomains(pod.ID)
	items := make([]components.ScrollItem, len(domains))
	for i, d := range domains {
		items[i] = domainItem{domain: d}
	}

	return podDomains{
		pod:         pod,
		project:     project,
		domains:     domains,
		list:        components.NewScrollList(items, components.ScrollListConfig{Width: podDomainsCard.InnerWidth(), Height: 8}),
		domainInput: domainInput,
		portInput:   portInput,
		keyAdd:      key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new custom")),
		keyAuto:     key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "generate auto")),
		keyEdit:     key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		keyDelete:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		keyOpen:     key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open")),
		keyBack:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		keySave:     key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyTab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	}
}

func (m podDomains) Init() tea.Cmd {
	return textinput.Blink
}

func (m podDomains) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		if m.mode == modeDomainAdd {
			return m.handleAddMode(tmsg)
		}
		return m.handleListMode(tmsg)

	case tea.MouseWheelMsg:
		if m.mode == modeDomainList {
			m.list, _ = m.list.Update(tmsg)
		}

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
	}

	return m, nil
}

func (m podDomains) handleListMode(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
		m.mode = modeDomainAdd
		m.isAuto = false
		m.focusedInput = 0
		m.domainInput.Focus()
		m.portInput.SetValue("8080")
		return m, textinput.Blink

	case key.Matches(tmsg, m.keyAuto):
		m.mode = modeDomainAdd
		m.isAuto = true
		m.focusedInput = 1
		m.domainInput.Blur()
		m.portInput.Focus()
		m.portInput.SetValue("8080")
		return m, textinput.Blink

	case key.Matches(tmsg, m.keyEdit):
		if item := m.list.SelectedItem(); item != nil {
			domain := item.(domainItem).domain
			pod := m.pod
			project := m.project
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewPodDomainsEdit(domain, pod, project)
					},
				}
			}
		}

	case key.Matches(tmsg, m.keyDelete):
		if item := m.list.SelectedItem(); item != nil {
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
		if item := m.list.SelectedItem(); item != nil {
			return m, utils.OpenBrowserCmd(item.(domainItem).domain.URL)
		}
	}

	// Let ScrollList handle navigation (up/down/j/k/mouse)
	m.list, _ = m.list.Update(tmsg)
	return m, nil
}

func (m podDomains) handleAddMode(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		m.mode = modeDomainList
		m.domainInput.SetValue("")
		m.portInput.SetValue("")
		m.isAuto = false
		return m, nil

	case key.Matches(tmsg, m.keyTab):
		if m.isAuto {
			// Only port for auto domains (SSL is automatic)
			// Keep focus on port field
			return m, nil
		}
		// For custom domains: toggle between domain and port
		m.focusedInput = (m.focusedInput + 1) % 2
		m.domainInput.Blur()
		m.portInput.Blur()
		switch m.focusedInput {
		case 0:
			m.domainInput.Focus()
		case 1:
			m.portInput.Focus()
		}
		return m, nil

	case key.Matches(tmsg, m.keySave):
		port := 8080
		pVal, err := strconv.Atoi(m.portInput.Value())
		if err == nil && pVal > 0 {
			port = pVal
		}

		// SSL is always enabled - it's automatic in production via Let's Encrypt
		if m.isAuto {
			return m, tea.Batch(
				func() tea.Msg { return msg.StartLoading{Text: "Generating domain"} },
				api.GenerateAutoDomain(m.pod.ID, port, true),
			)
		}

		domain := strings.TrimSpace(m.domainInput.Value())
		if domain != "" {
			return m, tea.Batch(
				func() tea.Msg { return msg.StartLoading{Text: "Creating domain"} },
				api.CreatePodDomain(m.pod.ID, domain, port, true),
			)
		}
		return m, nil
	}

	var cmd tea.Cmd
	switch m.focusedInput {
	case 0:
		m.domainInput, cmd = m.domainInput.Update(tmsg)
	case 1:
		m.portInput, cmd = m.portInput.Update(tmsg)
	}
	return m, cmd
}

func (m podDomains) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Domains"))
	b.WriteString("\n")
	b.WriteString(styles.MutedStyle().Render("Configure domains for " + m.pod.Title))
	b.WriteString("\n\n")

	if m.mode == modeDomainAdd {
		b.WriteString(m.renderAddMode())
	} else {
		b.WriteString(m.renderListMode())
	}

	return tea.NewView(m.centeredCard(b.String()))
}

func (m podDomains) renderListMode() string {
	if len(m.domains) == 0 {
		var b strings.Builder
		b.WriteString(styles.MutedStyle().Render("No domains configured."))
		b.WriteString("\n\n")
		b.WriteString(styles.MutedStyle().Render("Press 'g' to generate an auto domain, or 'n' to add a custom one."))
		b.WriteString("\n")
		b.WriteString(styles.MutedStyle().Render("A domain is required before you can deploy."))
		return b.String()
	}
	return m.list.View()
}

func (m podDomains) renderAddMode() string {
	var b strings.Builder

	if m.isAuto {
		b.WriteString("Generate Auto Domain\n\n")
		b.WriteString(styles.MutedStyle().Render("Domain will be auto-generated based on pod name"))
		b.WriteString("\n\n")
	} else {
		b.WriteString("Add Custom Domain\n\n")
	}

	labelStyle := lipgloss.NewStyle().Width(12)
	activeLabel := lipgloss.NewStyle().Width(12).Foreground(styles.ColorPrimary())

	// Domain field (only for custom domains)
	if !m.isAuto {
		if m.focusedInput == 0 {
			b.WriteString(activeLabel.Render("Domain:"))
		} else {
			b.WriteString(labelStyle.Render("Domain:"))
		}
		b.WriteString(m.domainInput.View())
		b.WriteString("\n\n")
	}

	// Port field
	if m.focusedInput == 1 {
		b.WriteString(activeLabel.Render("Port:"))
	} else {
		b.WriteString(labelStyle.Render("Port:"))
	}
	b.WriteString(m.portInput.View())
	b.WriteString("\n\n")

	// SSL info (no toggle - SSL is automatic in production)
	b.WriteString(labelStyle.Render("SSL:"))
	b.WriteString(styles.MutedStyle().Render("Automatic (Let's Encrypt)"))
	b.WriteString("\n")

	return b.String()
}

func (m podDomains) centeredCard(content string) string {
	card := styles.Card(podDomainsCard).Render(content)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
}

func (m podDomains) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title, "Domains"}
}
