package pages

import (
	"log"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/deeploy-sh/deeploy/internal/deeploy/viewtypes"
)

type connectKeyMap struct {
	Connect key.Binding
}

func (k connectKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Connect}
}

func (k connectKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Connect}}
}

func newConnectKeyMap() connectKeyMap {
	return connectKeyMap{
		Connect: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "connect")),
	}
}

type connectPage struct {
	serverInput textinput.Model
	keys        connectKeyMap
	help        help.Model
	status      string
	width       int
	height      int
	err         error
}

func NewConnectPage(err error) connectPage {
	log.Println(err)
	ti := textinput.New()
	ti.Placeholder = "e.g. 123.45.67.89:8090"
	ti.CharLimit = 50
	ti.Focus()

	return connectPage{
		serverInput: ti,
		keys:        newConnectKeyMap(),
		help:        styles.NewHelpModel(),
		err:         err,
	}
}

func (m connectPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m connectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyPressMsg:
		m.resetErr()
		switch msg.Code {
		case tea.KeyEnter:
			input := m.serverInput.Value()
			err := utils.ValidateServer(input)
			if err != nil {
				m.err = err
				return m, nil
			}
			return m, func() tea.Msg {
				return ChangePageMsg{
					PageFactory: func(s Store) tea.Model { return NewAuthPage(input) },
				}
			}
		}
	case messages.AuthSuccessMsg:
		return m, func() tea.Msg {
			return viewtypes.Dashboard
		}
	}
	m.serverInput, cmd = m.serverInput.Update(msg)
	return m, cmd
}

func (m connectPage) View() tea.View {
	content := "Connect to deeploy.sh server\n\n" + m.serverInput.View()
	if m.err != nil {
		content += styles.ErrorStyle().Render("\n* " + m.err.Error())
	}

	card := components.Card(components.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}).Render(content)
	helpView := m.help.View(m.keys)
	contentHeight := m.height - 1 // 1 für help

	// Card vertikal zentrieren
	centered := lipgloss.Place(m.width, contentHeight,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered, helpView))
}

func (m connectPage) Breadcrumbs() []string {
	return []string{"Connect"}
}

func (m *connectPage) resetErr() {
	m.err = nil
}
