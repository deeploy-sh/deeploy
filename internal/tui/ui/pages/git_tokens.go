package pages

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type gitTokensMode int

const (
	modeList gitTokensMode = iota
	modeAdd
)

type GitTokensPage struct {
	tokens       []api.GitToken
	selected     int
	mode         gitTokensMode
	nameInput    textinput.Model
	tokenInput   textinput.Model
	focusedInput int
	loading      bool
	keyAdd       key.Binding
	keyDelete    key.Binding
	keyBack      key.Binding
	keySave      key.Binding
	keyTab       key.Binding
	width        int
	height       int
}

func (m GitTokensPage) HelpKeys() []key.Binding {
	if m.mode == modeAdd {
		return []key.Binding{m.keySave, m.keyBack}
	}
	return []key.Binding{m.keyAdd, m.keyDelete, m.keyBack}
}

func NewGitTokensPage() GitTokensPage {
	nameInput := components.NewTextInput(40)
	nameInput.Placeholder = "Token name (e.g. GitHub Personal)"
	nameInput.CharLimit = 50

	tokenInput := components.NewTextInput(40)
	tokenInput.Placeholder = "ghp_xxxx or glpat-xxxx"
	tokenInput.CharLimit = 200
	tokenInput.EchoMode = textinput.EchoPassword

	return GitTokensPage{
		nameInput:  nameInput,
		tokenInput: tokenInput,
		keyAdd:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add token")),
		keyDelete:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		keyBack:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		keySave:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
		keyTab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
		loading:    true,
	}
}

func (m GitTokensPage) Init() tea.Cmd {
	return tea.Batch(api.FetchGitTokens(), textinput.Blink)
}

func (m GitTokensPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.GitTokensLoaded:
		m.loading = false
		if tokens, ok := tmsg.Tokens.([]api.GitToken); ok {
			m.tokens = tokens
		}
		return m, nil

	case msg.GitTokenCreated, msg.GitTokenDeleted:
		m.mode = modeList
		m.nameInput.SetValue("")
		m.tokenInput.SetValue("")
		return m, api.FetchGitTokens()

	case tea.KeyPressMsg:
		if m.mode == modeAdd {
			switch {
			case key.Matches(tmsg, m.keyBack):
				m.mode = modeList
				m.nameInput.SetValue("")
				m.tokenInput.SetValue("")
				return m, nil
			case key.Matches(tmsg, m.keyTab):
				m.focusedInput = (m.focusedInput + 1) % 2
				if m.focusedInput == 0 {
					m.nameInput.Focus()
					m.tokenInput.Blur()
				} else {
					m.nameInput.Blur()
					m.tokenInput.Focus()
				}
				return m, nil
			case key.Matches(tmsg, m.keySave):
				name := m.nameInput.Value()
				token := m.tokenInput.Value()
				if name != "" && token != "" {
					provider := detectProvider(token)
					m.loading = true
					return m, api.CreateGitToken(name, provider, token)
				}
				return m, nil
			}

			var cmd tea.Cmd
			if m.focusedInput == 0 {
				m.nameInput, cmd = m.nameInput.Update(tmsg)
			} else {
				m.tokenInput, cmd = m.tokenInput.Update(tmsg)
			}
			return m, cmd
		}

		// List mode
		switch {
		case key.Matches(tmsg, m.keyBack):
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			}
		case key.Matches(tmsg, m.keyAdd):
			m.mode = modeAdd
			m.focusedInput = 0
			m.nameInput.Focus()
			return m, textinput.Blink
		case key.Matches(tmsg, m.keyDelete):
			if len(m.tokens) > 0 && m.selected < len(m.tokens) {
				tokenID := m.tokens[m.selected].ID
				m.loading = true
				return m, api.DeleteGitToken(tokenID)
			}
		case tmsg.Code == tea.KeyUp:
			if m.selected > 0 {
				m.selected--
			}
		case tmsg.Code == tea.KeyDown:
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

	if m.loading {
		b.WriteString("Loading...")
		return tea.NewView(lipgloss.NewStyle().Padding(2, 4).Render(b.String()))
	}

	if m.mode == modeAdd {
		b.WriteString("Add New Token\n\n")

		labelStyle := lipgloss.NewStyle().Width(10)
		b.WriteString(labelStyle.Render("Name:"))
		b.WriteString(m.nameInput.View())
		b.WriteString("\n\n")

		b.WriteString(labelStyle.Render("Token:"))
		b.WriteString(m.tokenInput.View())
		b.WriteString("\n\n")

		b.WriteString(styles.MutedStyle().Render("Provider will be auto-detected from token format"))
	} else {
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
	}

	content := lipgloss.NewStyle().Padding(2, 4).Render(b.String())
	return tea.NewView(content)
}

func (m GitTokensPage) Breadcrumbs() []string {
	return []string{"Settings", "Git Tokens"}
}

func detectProvider(token string) string {
	switch {
	case strings.HasPrefix(token, "ghp_"), strings.HasPrefix(token, "github_pat_"):
		return "github"
	case strings.HasPrefix(token, "glpat-"):
		return "gitlab"
	case strings.HasPrefix(token, "ATBB"):
		return "bitbucket"
	default:
		return "github"
	}
}
