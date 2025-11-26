package pages

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/deeploy-sh/deeploy/internal/deeploy/viewtypes"
)

type connectPage struct {
	serverInput textinput.Model
	status      string
	width       int
	height      int
	err         error
}

func NewConnectPage(err error) connectPage {
	log.Println(err)
	ti := textinput.New()
	ti.Placeholder = "e.g. 123.45.67.89:8090"
	ti.Width = 30 // HACK: only because of: https://github.com/charmbracelet/bubbles/issues/779
	ti.Focus()

	return connectPage{
		serverInput: ti,
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
	case tea.KeyMsg:
		m.resetErr()
		switch msg.Type {
		case tea.KeyEnter:
			input := m.serverInput.Value()
			err := utils.ValidateServer(input)
			if err != nil {
				m.err = err
				return m, nil
			}
			return m, func() tea.Msg {
				return messages.ChangePageMsg{
					Page: NewAuthPage(input),
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

func (m connectPage) View() string {
	logo := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		MarginBottom(1).
		Render("🔥deeploy.sh")

	content := "Connect to deeploy.sh server\n\n" + m.serverInput.View()
	if m.err != nil {
		content += styles.ErrorStyle.Render("\n* " + m.err.Error())
	}

	card := components.Card(components.CardProps{Width: 50}).Render(content)

	view := lipgloss.JoinVertical(lipgloss.Center, logo, card)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
}

func (m *connectPage) resetErr() {
	m.err = nil
}
