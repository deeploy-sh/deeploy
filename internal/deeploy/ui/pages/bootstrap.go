package pages

import tea "github.com/charmbracelet/bubbletea"

// check internet connectivity -> yes = next check, no = show view (no internet connection, try again)
// check if there's a server set/connected/authenticated -> yes = show dashboard, no = show connect page

type bootstrap struct{}

func NewBootstrap() tea.Model {
	return &bootstrap{}
}

func (m bootstrap) Init() tea.Cmd
func (m bootstrap) Update(tea.Msg) (tea.Model, tea.Cmd)
func (m bootstrap) View() string
