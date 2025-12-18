package page

import (
	"fmt"

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

type serverDomainDelete struct {
	domain     string
	input      textinput.Model
	keyConfirm key.Binding
	keyCancel  key.Binding
	width      int
	height     int
}

func (p serverDomainDelete) HelpKeys() []key.Binding {
	return []key.Binding{p.keyConfirm, p.keyCancel}
}

func NewServerDomainDelete(domain string) serverDomainDelete {
	card := styles.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = domain
	ti.Focus()
	ti.CharLimit = 100

	return serverDomainDelete{
		domain:     domain,
		input:      ti,
		keyConfirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (p serverDomainDelete) Init() tea.Cmd {
	return textinput.Blink
}

func (p serverDomainDelete) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.ServerDomainDeleted:
		// API deleted Traefik config - revert local config to original IP
		cfg, err := config.Load()
		if err != nil {
			return p, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewServerDomain() },
				}
			}
		}

		cfg.Server = cfg.ServerIP
		cfg.ServerIP = "" // Clear fallback
		config.Save(cfg)

		// Return to dashboard
		return p, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
			}
		}

	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			return p, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewServerDomain() },
				}
			}
		case tea.KeyEnter:
			// Only delete if input matches domain exactly
			if p.input.Value() != p.domain {
				return p, nil
			}
			return p, api.DeleteServerDomain()
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

func (p serverDomainDelete) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary()).
		Render("Delete Server Domain")

	domainName := lipgloss.NewStyle().
		Bold(true).
		Render(p.domain)

	hint := styles.MutedStyle().
		Render(fmt.Sprintf("Type '%s' to confirm", p.domain))

	warning := styles.WarningStyle().
		Render("This will disable HTTPS and revert to HTTP")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		domainName,
		"",
		warning,
		"",
		hint,
		"",
		p.input.View(),
	)

	card := styles.Card(styles.CardProps{
		Width:   50,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p serverDomainDelete) Breadcrumbs() []string {
	return []string{"server domain", "delete"}
}
