package page

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type info struct {
	tuiVersion    string
	serverVersion string
	latestVersion string
	width         int
	height        int
	keyBack       key.Binding
}

func NewInfo(tuiVersion, serverVersion, latestVersion string) info {
	return info{
		tuiVersion:    tuiVersion,
		serverVersion: serverVersion,
		latestVersion: latestVersion,
		keyBack:       key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc", "back")),
	}
}

func (m info) Breadcrumbs() []string {
	return []string{"about"}
}

func (m info) HelpKeys() []key.Binding {
	return []key.Binding{m.keyBack}
}

func (m info) Init() tea.Cmd {
	return nil
}

func (m info) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height

	case tea.KeyPressMsg:
		if key.Matches(tmsg, m.keyBack) {
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			}
		}
	}

	return m, nil
}

func (m info) View() tea.View {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary())

	labelStyle := lipgloss.NewStyle().
		Width(10)

	mutedStyle := styles.MutedStyle()

	// Version info
	tuiStatus := m.tuiVersion
	serverStatus := m.serverVersion
	if serverStatus == "" {
		serverStatus = "..."
	}

	// Check for updates
	tuiNeedsUpdate := false
	serverNeedsUpdate := false
	if m.latestVersion != "" {
		if m.tuiVersion != "dev" && m.tuiVersion != m.latestVersion {
			tuiNeedsUpdate = true
		}
		if m.serverVersion != "" && m.serverVersion != "dev" && m.serverVersion != m.latestVersion {
			serverNeedsUpdate = true
		}
	}

	if tuiNeedsUpdate {
		tuiStatus += fmt.Sprintf("  →  %s  ⬆", m.latestVersion)
	} else if m.latestVersion != "" {
		tuiStatus += "  ✓ latest"
	}

	if serverNeedsUpdate {
		serverStatus += fmt.Sprintf("  →  %s  ⬆", m.latestVersion)
	} else if m.latestVersion != "" && m.serverVersion != "" {
		serverStatus += "  ✓ latest"
	}

	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Version Info"))
	sb.WriteString("\n\n")

	sb.WriteString(labelStyle.Render("TUI"))
	sb.WriteString(tuiStatus)
	sb.WriteString("\n")

	sb.WriteString(labelStyle.Render("Server"))
	sb.WriteString(serverStatus)
	sb.WriteString("\n")

	// Show update commands if needed
	if tuiNeedsUpdate || serverNeedsUpdate {
		sb.WriteString("\n")
		sb.WriteString(mutedStyle.Render(strings.Repeat("─", 40)))
		sb.WriteString("\n\n")

		if tuiNeedsUpdate {
			sb.WriteString(titleStyle.Render("Update TUI"))
			sb.WriteString("\n")
			sb.WriteString(mutedStyle.Render("curl -fsSL https://deeploy.sh/tui | bash"))
			sb.WriteString("\n\n")
		}

		if serverNeedsUpdate {
			sb.WriteString(titleStyle.Render("Update Server (on VPS)"))
			sb.WriteString("\n")
			sb.WriteString(mutedStyle.Render("curl -fsSL https://deeploy.sh/server | bash"))
			sb.WriteString("\n")
		}
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
