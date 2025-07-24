package pages

import (
	"github.com/axadrn/deeploy/internal/tui/messages"
	"github.com/axadrn/deeploy/internal/tui/ui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type DashboardPage struct {
	width  int
	height int
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewDashboard() DashboardPage {
	return DashboardPage{}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p DashboardPage) Init() tea.Cmd {
	return nil
}

func (p DashboardPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "p" {
			return p, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewProjectPage()}
			}
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	}
	return p, nil

}

func (p DashboardPage) View() string {
	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	menu := components.Card(components.CardProps{
		Padding: []int{0, 2},
	}).Render("[P]rojects [S]ettings")

	view := lipgloss.JoinVertical(0.5, logo, menu)

	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)
	return layout

}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////
