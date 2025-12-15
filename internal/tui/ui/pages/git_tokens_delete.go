package pages

import (
	"fmt"

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

type GitTokenDeletePage struct {
	token      model.GitToken
	input      textinput.Model
	keyConfirm key.Binding
	keyCancel  key.Binding
	width      int
	height     int
}

func (p GitTokenDeletePage) HelpKeys() []key.Binding {
	return []key.Binding{p.keyConfirm, p.keyCancel}
}

func NewGitTokenDeletePage(token model.GitToken) GitTokenDeletePage {
	card := styles.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = token.Name
	ti.Focus()
	ti.CharLimit = 100

	return GitTokenDeletePage{
		token:      token,
		input:      ti,
		keyConfirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (p GitTokenDeletePage) Init() tea.Cmd {
	return textinput.Blink
}

func (p GitTokenDeletePage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.GitTokenDeleted:
		return p, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewGitTokensPage() },
			}
		}

	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			return p, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewGitTokensPage() },
				}
			}
		case tea.KeyEnter:
			// Only delete if input matches token name exactly
			if p.input.Value() != p.token.Name {
				return p, nil
			}
			return p, api.DeleteGitToken(p.token.ID)
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

func (p GitTokenDeletePage) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary()).
		Render(fmt.Sprintf("Delete Token"))

	tokenName := lipgloss.NewStyle().
		Bold(true).
		Render(p.token.Name)

	hint := styles.MutedStyle().
		Render("Type '" + p.token.Name + "' to confirm")

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		tokenName,
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

func (p GitTokenDeletePage) Breadcrumbs() []string {
	return []string{"Settings", "Git Tokens", "Delete"}
}
