package page

import (
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

const (
	fieldName = iota
	fieldToken
)

type gitTokenForm struct {
	nameInput    textinput.Model
	tokenInput   textinput.Model
	focusedField int
	keySave      key.Binding
	keyCancel    key.Binding
	keyTab       key.Binding
	keyShiftTab  key.Binding
	width        int
	height       int
}

func (m gitTokenForm) HelpKeys() []key.Binding {
	return []key.Binding{m.keySave, m.keyTab, m.keyCancel}
}

func NewGitTokenForm() gitTokenForm {
	card := styles.CardProps{Width: styles.CardWidthMD, Padding: []int{1, 2}, Accent: true}
	inputWidth := card.InnerWidth()

	nameInput := components.NewTextInput(inputWidth)
	nameInput.Placeholder = "e.g. GitHub Personal"
	nameInput.CharLimit = 50
	nameInput.Focus()

	tokenInput := components.NewTextInput(inputWidth)
	tokenInput.Placeholder = "ghp_xxxx or glpat-xxxx"
	tokenInput.CharLimit = 200
	tokenInput.EchoMode = textinput.EchoPassword

	return gitTokenForm{
		nameInput:   nameInput,
		tokenInput:  tokenInput,
		keySave:     key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyCancel:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		keyTab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next")),
		keyShiftTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev")),
	}
}

func (m gitTokenForm) Init() tea.Cmd {
	return textinput.Blink
}

func (m gitTokenForm) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.GitTokenCreated:
		return m, tea.Batch(
			api.LoadData(),
			func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewGitTokens(s.GitTokens()) },
				}
			},
		)

	case tea.KeyPressMsg:
		return m.handleKeyPress(tmsg)

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	// Blink passthrough
	var cmd tea.Cmd
	switch m.focusedField {
	case fieldName:
		m.nameInput, cmd = m.nameInput.Update(tmsg)
	case fieldToken:
		m.tokenInput, cmd = m.tokenInput.Update(tmsg)
	}
	return m, cmd
}

func (m *gitTokenForm) handleKeyPress(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyCancel):
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewGitTokens(s.GitTokens()) },
			}
		}

	case key.Matches(tmsg, m.keySave):
		name := strings.TrimSpace(m.nameInput.Value())
		token := strings.TrimSpace(m.tokenInput.Value())
		if name != "" && token != "" {
			provider := detectProvider(token)
			return m, tea.Batch(
				func() tea.Msg { return msg.StartLoading{Text: "Creating token"} },
				api.CreateGitToken(name, provider, token),
			)
		}
		return m, nil

	case key.Matches(tmsg, m.keyTab):
		m.focusedField = (m.focusedField + 1) % 2
		return m, m.updateFocus()

	case key.Matches(tmsg, m.keyShiftTab):
		m.focusedField = (m.focusedField + 1) % 2 // Only 2 fields, so +1 works for both directions
		return m, m.updateFocus()
	}

	// Update focused input
	var cmd tea.Cmd
	switch m.focusedField {
	case fieldName:
		m.nameInput, cmd = m.nameInput.Update(tmsg)
	case fieldToken:
		m.tokenInput, cmd = m.tokenInput.Update(tmsg)
	}
	return m, cmd
}

func (m *gitTokenForm) updateFocus() tea.Cmd {
	m.nameInput.Blur()
	m.tokenInput.Blur()

	switch m.focusedField {
	case fieldName:
		return m.nameInput.Focus()
	case fieldToken:
		return m.tokenInput.Focus()
	}
	return nil
}

func (m gitTokenForm) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Add Git Token"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorMuted())
	activeLabel := lipgloss.NewStyle().Foreground(styles.ColorPrimary())

	// Name
	if m.focusedField == fieldName {
		b.WriteString(activeLabel.Render("Name"))
	} else {
		b.WriteString(labelStyle.Render("Name"))
	}
	b.WriteString("\n")
	b.WriteString(m.nameInput.View())
	b.WriteString("\n\n")

	// Token
	if m.focusedField == fieldToken {
		b.WriteString(activeLabel.Render("Token"))
	} else {
		b.WriteString(labelStyle.Render("Token"))
	}
	b.WriteString("\n")
	b.WriteString(m.tokenInput.View())
	b.WriteString("\n\n")

	b.WriteString(styles.MutedStyle().Render("Provider will be auto-detected from token format"))

	card := styles.Card(styles.CardProps{
		Width:   styles.CardWidthMD,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(b.String())

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m gitTokenForm) Breadcrumbs() []string {
	return []string{"Settings", "Git Tokens", "Add"}
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
