package pages

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type GitTokensPage struct {
	tokens    []api.GitToken
	selected  int
	keyAdd    key.Binding
	keyDelete key.Binding
	keyBack   key.Binding
	width     int
	height    int
}

func (m GitTokensPage) HelpKeys() []key.Binding {
	return []key.Binding{m.keyAdd, m.keyDelete, m.keyBack}
}

func NewGitTokensPage() GitTokensPage {
	return GitTokensPage{
		keyAdd:    key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
		keyDelete: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		keyBack:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

func (m GitTokensPage) Init() tea.Cmd {
	return api.FetchGitTokens()
}

func (m GitTokensPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.GitTokensLoaded:
		if tokens, ok := tmsg.Tokens.([]api.GitToken); ok {
			m.tokens = tokens
		}
		return m, nil

	case msg.GitTokenDeleted:
		return m, api.FetchGitTokens()

	case tea.KeyPressMsg:
		switch {
		case key.Matches(tmsg, m.keyBack):
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			}

		case key.Matches(tmsg, m.keyAdd):
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewGitTokenFormPage() },
				}
			}

		case key.Matches(tmsg, m.keyDelete):
			if len(m.tokens) > 0 && m.selected < len(m.tokens) {
				token := m.tokens[m.selected]
				return m, func() tea.Msg {
					return msg.ChangePage{
						PageFactory: func(s msg.Store) tea.Model { return NewGitTokenDeletePage(token) },
					}
				}
			}

		case tmsg.Code == tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			}
		case tmsg.Code == tea.KeyDown:
			if m.selected < len(m.tokens)-1 {
				m.selected++
			}
		// Ctrl+P = previous (up)
		case tmsg.String() == "ctrl+p":
			if m.selected > 0 {
				m.selected--
			}
		// Ctrl+N = next (down)
		case tmsg.String() == "ctrl+n":
			if m.selected < len(m.tokens)-1 {
				m.selected++
			}
		}

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
	}

	return m, nil
}

func (m GitTokensPage) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Git Tokens"))
	b.WriteString("\n")
	b.WriteString(styles.MutedStyle().Render("Personal Access Tokens for private repositories"))
	b.WriteString("\n\n")

	if len(m.tokens) == 0 {
		b.WriteString(styles.MutedStyle().Render("No tokens configured. Press 'a' to add one."))
	} else {
		for i, t := range m.tokens {
			cursor := "  "
			style := lipgloss.NewStyle()
			if i == m.selected {
				cursor = "> "
				style = style.Foreground(styles.ColorPrimary())
			}

			providerBadge := fmt.Sprintf("[%s]", t.Provider)
			line := fmt.Sprintf("%s%s %s", cursor, style.Render(t.Name), styles.MutedStyle().Render(providerBadge))
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	card := styles.Card(styles.CardProps{
		Width:   60,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(b.String())

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m GitTokensPage) Breadcrumbs() []string {
	return []string{"Settings", "Git Tokens"}
}
