package page

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type serverDomain struct {
	domainInput   textinput.Model
	domain        string // domain from DB (source of truth)
	serverIP      string // original IP:port for fallback
	pendingDomain string // domain being saved
	loading       bool   // waiting for initial load
	saving        bool   // waiting for API response
	deleting      bool   // waiting for delete API response
	width         int
	height        int
	keyBack       key.Binding
	keySave       key.Binding
	keyDelete     key.Binding
	err           error
}

func NewServerDomain() serverDomain {
	// Load server IP from local config (for fallback)
	cfg, _ := config.Load()
	serverIP := ""
	if cfg != nil {
		serverIP = cfg.ServerIP
	}

	// Text input for domain
	card := styles.CardProps{Width: 50, Padding: []int{1, 2}, Accent: false}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = "deeploy.yourdomain.com"
	ti.CharLimit = 100
	ti.Focus()

	return serverDomain{
		domainInput: ti,
		serverIP:    serverIP,
		loading:     true,
		keyBack:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		keySave:     key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyDelete:   key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "delete")),
	}
}

func (m serverDomain) Breadcrumbs() []string {
	return []string{"server domain"}
}

func (m serverDomain) HelpKeys() []key.Binding {
	keys := []key.Binding{m.keyBack, m.keySave}
	// Only show delete if a domain is configured and we have fallback IP
	if m.domain != "" && m.serverIP != "" {
		keys = append(keys, m.keyDelete)
	}
	return keys
}

func (m serverDomain) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, api.GetServerDomain())
}

func (m serverDomain) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height

	case msg.ServerDomainLoaded:
		m.loading = false
		m.domain = tmsg.Domain
		return m, nil

	case msg.ServerDomainSet:
		// API success - save to local config and go back
		cfg, err := config.Load()
		if err != nil {
			m.err = fmt.Errorf("failed to load config")
			m.saving = false
			return m, nil
		}

		// Save original IP for fallback (only if not already set)
		if cfg.ServerIP == "" {
			cfg.ServerIP = cfg.Server
		}

		cfg.Server = "https://" + m.pendingDomain
		err = config.Save(cfg)
		if err != nil {
			m.err = fmt.Errorf("failed to save local config")
			m.saving = false
			return m, nil
		}

		// Return to dashboard
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
			}
		}

	case msg.ServerDomainDeleted:
		// API success - revert to original IP
		cfg, err := config.Load()
		if err != nil {
			m.err = fmt.Errorf("failed to load config")
			m.deleting = false
			return m, nil
		}

		cfg.Server = cfg.ServerIP
		cfg.ServerIP = "" // Clear fallback
		err = config.Save(cfg)
		if err != nil {
			m.err = fmt.Errorf("failed to save local config")
			m.deleting = false
			return m, nil
		}

		// Return to dashboard
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
			}
		}

	case msg.Error:
		m.err = tmsg.Err
		m.loading = false
		m.saving = false
		m.deleting = false
		return m, nil

	case tea.KeyPressMsg:
		if m.loading || m.saving || m.deleting {
			return m, nil // ignore input while loading/saving/deleting
		}

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

			// Call API to set domain (this writes Traefik config)
			m.pendingDomain = domain
			m.saving = true
			return m, api.SetServerDomain(domain)
		}

		if key.Matches(tmsg, m.keyDelete) {
			// Only allow delete if we have a fallback IP
			if m.serverIP == "" {
				m.err = fmt.Errorf("no fallback IP available")
				return m, nil
			}

			m.deleting = true
			return m, api.DeleteServerDomain()
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

	// Current domain (from DB)
	sb.WriteString(labelStyle.Render("Current:"))
	if m.loading {
		sb.WriteString(mutedStyle.Render("loading..."))
	} else if m.domain != "" {
		sb.WriteString("https://" + m.domain)
	} else {
		sb.WriteString(mutedStyle.Render("not set"))
	}
	sb.WriteString("\n")

	// Status
	sb.WriteString(labelStyle.Render("Status:"))
	if m.loading {
		sb.WriteString(mutedStyle.Render("..."))
	} else if m.domain != "" {
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

	// Use serverIP for DNS example, fallback to placeholder
	ip := extractIP(m.serverIP)

	sb.WriteString("1. Create DNS A-Record:\n")
	sb.WriteString(mutedStyle.Render(fmt.Sprintf("   yourdomain.com → %s", ip)))
	sb.WriteString("\n\n")

	sb.WriteString("2. Enter domain:\n")
	sb.WriteString("   ")
	sb.WriteString(m.domainInput.View())
	sb.WriteString("\n\n")

	sb.WriteString(mutedStyle.Render("Traefik will automatically get an SSL cert."))

	// Status messages
	if m.loading {
		sb.WriteString("\n\n")
		sb.WriteString(styles.MutedStyle().Render("Loading..."))
	} else if m.saving {
		sb.WriteString("\n\n")
		sb.WriteString(styles.MutedStyle().Render("Saving..."))
	} else if m.deleting {
		sb.WriteString("\n\n")
		sb.WriteString(styles.MutedStyle().Render("Deleting..."))
	} else if m.err != nil {
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
