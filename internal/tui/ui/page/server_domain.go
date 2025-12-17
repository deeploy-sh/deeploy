package page

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type serverDomain struct {
	domainInput textinput.Model
	currentURL  string
	width       int
	height      int
	keyBack     key.Binding
	keySave     key.Binding
	err         error
}

func NewServerDomain() serverDomain {
	// Load current server URL from config
	cfg, _ := config.Load()
	currentURL := ""
	if cfg != nil {
		currentURL = cfg.Server
	}

	// Text input for domain
	card := styles.CardProps{Width: 50, Padding: []int{1, 2}, Accent: false}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = "deeploy.yourdomain.com"
	ti.CharLimit = 100
	ti.Focus()

	return serverDomain{
		domainInput: ti,
		currentURL:  currentURL,
		keyBack:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		keySave:     key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
	}
}

func (m serverDomain) Breadcrumbs() []string {
	return []string{"server domain"}
}

func (m serverDomain) HelpKeys() []key.Binding {
	return []key.Binding{m.keyBack, m.keySave}
}

func (m serverDomain) Init() tea.Cmd {
	return textinput.Blink
}

func (m serverDomain) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height

	case tea.KeyPressMsg:
		m.err = nil

		if key.Matches(tmsg, m.keyBack) {
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			}
		}

		if key.Matches(tmsg, m.keySave) {
			domain := strings.TrimSpace(m.domainInput.Value())

			// Validate domain
			if domain == "" {
				m.err = fmt.Errorf("domain is required")
				return m, nil
			}

			// Remove protocol if user accidentally included it
			domain = strings.TrimPrefix(domain, "https://")
			domain = strings.TrimPrefix(domain, "http://")

			// Save to config with https://
			cfg, err := config.Load()
			if err != nil {
				m.err = fmt.Errorf("failed to load config")
				return m, nil
			}

			cfg.Server = "https://" + domain
			err = config.Save(cfg)
			if err != nil {
				m.err = fmt.Errorf("failed to save config")
				return m, nil
			}

			// Return to dashboard - heartbeat will reconnect automatically
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			}
		}
	}

	m.domainInput, cmd = m.domainInput.Update(tmsg)
	return m, cmd
}

func (m serverDomain) View() tea.View {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary())

	labelStyle := lipgloss.NewStyle().
		Width(10)

	mutedStyle := styles.MutedStyle()

	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Server Domain"))
	sb.WriteString("\n\n")

	// Current URL
	sb.WriteString(labelStyle.Render("Current:"))
	if m.currentURL != "" {
		sb.WriteString(m.currentURL)
	} else {
		sb.WriteString(mutedStyle.Render("not set"))
	}
	sb.WriteString("\n")

	// Status
	sb.WriteString(labelStyle.Render("Status:"))
	if strings.HasPrefix(m.currentURL, "https://") {
		sb.WriteString(styles.SuccessStyle().Render("Encrypted (HTTPS)"))
	} else {
		sb.WriteString(styles.WarningStyle().Render("Not encrypted"))
	}
	sb.WriteString("\n\n")

	// Separator
	sb.WriteString(mutedStyle.Render(strings.Repeat("─", 44)))
	sb.WriteString("\n\n")

	// Instructions
	sb.WriteString(titleStyle.Render("To enable HTTPS:"))
	sb.WriteString("\n\n")

	// Extract IP from current URL for the DNS example
	ip := extractIP(m.currentURL)

	sb.WriteString("1. Create DNS A-Record:\n")
	sb.WriteString(mutedStyle.Render(fmt.Sprintf("   yourdomain.com → %s", ip)))
	sb.WriteString("\n\n")

	sb.WriteString("2. Enter domain:\n")
	sb.WriteString("   ")
	sb.WriteString(m.domainInput.View())
	sb.WriteString("\n\n")

	sb.WriteString(mutedStyle.Render("Traefik will automatically get an SSL cert."))

	// Error message
	if m.err != nil {
		sb.WriteString("\n\n")
		sb.WriteString(styles.ErrorStyle().Render("* " + m.err.Error()))
	}

	cardStyle := styles.Card(styles.CardProps{
		Width:   50,
		Padding: []int{1, 2},
	})
	card := cardStyle.Render(sb.String())

	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		card,
	)

	return tea.NewView(centered)
}

// extractIP extracts the IP/host from a URL like "http://123.45.67.89:8090"
func extractIP(url string) string {
	if url == "" {
		return "YOUR_VPS_IP"
	}

	// Remove protocol
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	// Remove port
	if idx := strings.Index(url, ":"); idx != -1 {
		url = url[:idx]
	}

	// Remove path
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}

	return url
}
