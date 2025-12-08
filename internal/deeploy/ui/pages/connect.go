package pages

import (
	"log"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
)

type connectPage struct {
	serverInput textinput.Model
	keyConnect  key.Binding
	status      string
	width       int
	height      int
	err         error
}

func (p connectPage) HelpKeys() []key.Binding {
	return []key.Binding{p.keyConnect}
}

func NewConnectPage(err error) connectPage {
	log.Println(err)
	card := components.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = "http://123.45.67.89:8090"
	ti.CharLimit = 50
	ti.Focus()

	return connectPage{
		serverInput: ti,
		keyConnect:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "connect")),
		err:         err,
	}
}

func (m connectPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m connectPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
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
					PageFactory: func(s msg.Store) tea.Model { return NewAuthPage(input) },
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

func (m connectPage) View() tea.View {
	content := "Connect to deeploy.sh server\n\n" + m.serverInput.View()
	if m.err != nil {
		content += styles.ErrorStyle().Render("\n* " + m.err.Error())
	}

	card := components.Card(components.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}).Render(content)

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m connectPage) Breadcrumbs() []string {
	return []string{"Connect"}
}

func (m *connectPage) resetErr() {
	m.err = nil
}
