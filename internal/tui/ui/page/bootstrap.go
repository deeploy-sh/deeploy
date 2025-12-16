package page

import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type bootstrap struct {
	width, height int
	offline       bool
	spinner       spinner.Model
}

func NewBootstrap() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.ColorPrimary())

	return &bootstrap{
		spinner: s,
	}
}

func (m *bootstrap) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *bootstrap) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

func (m *bootstrap) View() tea.View {
	spinner := m.spinner.View()
	height := m.height + headerHeight + footerHeight
	if m.offline {
		return tea.NewView(components.Centered(m.width, height, spinner+" ⚡ deeploy.sh\n\ncan't connect. retrying..."))
	}
	return tea.NewView(components.Centered(m.width, height, spinner+" ⚡ deeploy.sh"))
}
