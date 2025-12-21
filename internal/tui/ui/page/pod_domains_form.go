package page

import (
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

type podDomainsForm struct {
	domain       *model.PodDomain // nil = create, otherwise edit
	pod          *model.Pod
	project      *model.Project
	isAuto       bool // only relevant for create
	domainInput  textinput.Model
	portInput    textinput.Model
	focusedField int
	keySave      key.Binding
	keyBack      key.Binding
	keyTab       key.Binding
	keyShiftTab  key.Binding
	width        int
	height       int
}

const (
	fieldDomain = iota
	fieldPort
)

func (m podDomainsForm) HelpKeys() []key.Binding {
	return []key.Binding{m.keySave, m.keyTab, m.keyBack}
}

func NewPodDomainsForm(pod *model.Pod, project *model.Project, domain *model.PodDomain, isAuto bool) podDomainsForm {
	card := styles.CardProps{Width: styles.CardWidthMD, Padding: []int{1, 2}, Accent: true}
	inputWidth := card.InnerWidth()

	domainInput := components.NewTextInput(inputWidth)
	domainInput.Placeholder = "app.example.com"
	domainInput.CharLimit = 100

	portInput := components.NewTextInput(inputWidth)
	portInput.Placeholder = "8080"
	portInput.CharLimit = 5
	portInput.SetValue("8080")

	// Set values if editing
	if domain != nil {
		domainInput.SetValue(domain.Domain)
		portInput.SetValue(strconv.Itoa(domain.Port))
	}

	// Set initial focus
	if isAuto {
		// Auto domain: focus on port (domain is auto-generated)
		portInput.Focus()
	} else {
		domainInput.Focus()
	}

	focusedField := fieldDomain
	if isAuto {
		focusedField = fieldPort
	}

	return podDomainsForm{
		domain:       domain,
		pod:          pod,
		project:      project,
		isAuto:       isAuto,
		domainInput:  domainInput,
		portInput:    portInput,
		focusedField: focusedField,
		keySave:      key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyBack:      key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		keyTab:       key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
		keyShiftTab:  key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
	}
}

func (m podDomainsForm) Init() tea.Cmd {
	return textinput.Blink
}

func (m podDomainsForm) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(tmsg)

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	// Update focused input for blink messages
	var cmd tea.Cmd
	switch m.focusedField {
	case fieldDomain:
		m.domainInput, cmd = m.domainInput.Update(tmsg)
	case fieldPort:
		m.portInput, cmd = m.portInput.Update(tmsg)
	}
	return m, cmd
}

func (m *podDomainsForm) handleKeyPress(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDomains(s, pod, project)
				},
			}
		}

	case key.Matches(tmsg, m.keySave):
		return m.save()

	case key.Matches(tmsg, m.keyTab):
		if m.isAuto && m.domain == nil {
			// Auto create: only port field, no tab
			return m, nil
		}
		m.focusedField = (m.focusedField + 1) % 2
		return m, m.updateFocus()

	case key.Matches(tmsg, m.keyShiftTab):
		if m.isAuto && m.domain == nil {
			// Auto create: only port field, no tab
			return m, nil
		}
		m.focusedField = (m.focusedField + 1) % 2
		return m, m.updateFocus()
	}

	// Update focused input
	var cmd tea.Cmd
	switch m.focusedField {
	case fieldDomain:
		m.domainInput, cmd = m.domainInput.Update(tmsg)
	case fieldPort:
		m.portInput, cmd = m.portInput.Update(tmsg)
	}
	return m, cmd
}

func (m *podDomainsForm) blurAll() {
	m.domainInput.Blur()
	m.portInput.Blur()
}

func (m *podDomainsForm) updateFocus() tea.Cmd {
	m.blurAll()
	switch m.focusedField {
	case fieldDomain:
		return m.domainInput.Focus()
	case fieldPort:
		return m.portInput.Focus()
	}
	return nil
}

func (m *podDomainsForm) save() (tea.Model, tea.Cmd) {
	port := 8080
	pVal, err := strconv.Atoi(m.portInput.Value())
	if err == nil && pVal > 0 {
		port = pVal
	}

	// Auto domain create
	if m.domain == nil && m.isAuto {
		return m, tea.Batch(
			func() tea.Msg { return msg.StartLoading{Text: "Generating domain"} },
			api.GenerateAutoDomain(m.pod.ID, port, true),
		)
	}

	// Custom domain - validate
	domain := strings.TrimSpace(m.domainInput.Value())
	if domain == "" {
		return m, nil
	}

	// Create
	if m.domain == nil {
		return m, tea.Batch(
			func() tea.Msg { return msg.StartLoading{Text: "Creating domain"} },
			api.CreatePodDomain(m.pod.ID, domain, port, true),
		)
	}

	// Update
	return m, tea.Batch(
		func() tea.Msg { return msg.StartLoading{Text: "Updating domain"} },
		api.UpdatePodDomain(m.pod.ID, m.domain.ID, domain, port, true),
	)
}

func (m podDomainsForm) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())

	// Title based on mode
	if m.domain == nil {
		if m.isAuto {
			b.WriteString(titleStyle.Render("Generate Auto Domain"))
			b.WriteString("\n")
			b.WriteString(styles.MutedStyle().Render("Domain will be auto-generated based on pod name"))
		} else {
			b.WriteString(titleStyle.Render("New Domain"))
		}
	} else {
		b.WriteString(titleStyle.Render("Edit Domain"))
		b.WriteString("\n")
		typeLabel := "custom"
		if m.domain.Type == "auto" {
			typeLabel = "auto"
		}
		b.WriteString(styles.MutedStyle().Render("Type: " + typeLabel))
	}
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorMuted())
	activeLabel := lipgloss.NewStyle().Foreground(styles.ColorPrimary())

	// Domain field (not for auto create)
	if !(m.domain == nil && m.isAuto) {
		if m.focusedField == fieldDomain {
			b.WriteString(activeLabel.Render("Domain"))
		} else {
			b.WriteString(labelStyle.Render("Domain"))
		}
		b.WriteString("\n")
		b.WriteString(m.domainInput.View())
		b.WriteString("\n\n")
	}

	// Port field
	if m.focusedField == fieldPort {
		b.WriteString(activeLabel.Render("Port"))
	} else {
		b.WriteString(labelStyle.Render("Port"))
	}
	b.WriteString("\n")
	b.WriteString(m.portInput.View())
	b.WriteString("\n\n")

	// SSL info
	b.WriteString(labelStyle.Render("SSL"))
	b.WriteString("\n")
	b.WriteString(styles.MutedStyle().Render("Automatic (Let's Encrypt)"))

	card := styles.Card(styles.CardProps{
		Width:   styles.CardWidthMD,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(b.String())

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m podDomainsForm) Breadcrumbs() []string {
	if m.domain == nil {
		return []string{"Projects", m.project.Title, "Pods", m.pod.Title, "Domains", "New"}
	}
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title, "Domains", "Edit"}
}
