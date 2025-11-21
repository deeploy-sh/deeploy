package pages

import (
	"strings"

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
	err         string
}

func NewConnectPage() connectPage {
	ti := textinput.New()
	ti.Placeholder = "e.g. 123.45.67.89:8090"
	ti.Width = 80 // HACK: only because of: https://github.com/charmbracelet/bubbles/issues/779
	ti.Focus()

	return connectPage{
		serverInput: ti,
	}
}

func (p connectPage) Init() tea.Cmd {
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
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := m.serverInput.Value()
			err := utils.ValidateServer(input)
			if err != nil {
				m.err = err.Error()
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

func (p connectPage) View() string {
	var b strings.Builder

	b.WriteString("Connect to deeploy.sh server\n\n")
	b.WriteString(p.serverInput.View())
	if p.err != "" {
		b.WriteString(styles.ErrorStyle.Render("\n* " + p.err))
	}
	if p.status != "" {
		b.WriteString(p.status)
	}

	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	card := components.Card(components.CardProps{Width: 50}).Render(b.String())

	view := lipgloss.JoinVertical(lipgloss.Center, logo, card)
	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}

func (p *connectPage) resetErr() {
	p.err = ""
}
