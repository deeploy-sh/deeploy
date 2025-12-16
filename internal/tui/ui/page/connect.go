package page

import (
	"log"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/tui/utils"
)

type connect struct {
	serverInput textinput.Model
	keyconnect  key.Binding
	status      string
	width       int
	height      int
	err         error
}

func (p connect) HelpKeys() []key.Binding {
	return []key.Binding{p.keyconnect}
}

func NewConnect(err error) connect {
	log.Println(err)
	card := styles.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = "http://123.45.67.89:8090"
	ti.CharLimit = 50
	ti.Focus()

	return connect{
		serverInput: ti,
		keyconnect:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "connect")),
		err:         err,
	}
}

func (m connect) Init() tea.Cmd {
	return textinput.Blink
}

func (m connect) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
	case tea.KeyPressMsg:
		m.resetErr()
		switch tmsg.Code {
		case tea.KeyEnter:
			input := m.serverInput.Value()
			err := utils.ValidateServer(input)
			if err != nil {
				m.err = err
				return m, nil
			}
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewAuth(input) },
				}
			}
		}
	case msg.AuthSuccess:
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
			}
		}
	}
	m.serverInput, cmd = m.serverInput.Update(tmsg)
	return m, cmd
}

func (m connect) View() tea.View {
	content := "connect to deeploy.sh server\n\n" + m.serverInput.View()
	if m.err != nil {
		content += styles.ErrorStyle().Render("\n* " + m.err.Error())
	}

	card := styles.Card(styles.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}).Render(content)

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m connect) Breadcrumbs() []string {
	return []string{"connect"}
}

func (m *connect) resetErr() {
	m.err = nil
}
