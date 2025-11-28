package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
)

type bootstrap struct {
	width, height int
	offline       bool
}

func NewBootstrap() tea.Model {
	return &bootstrap{}
}

func (m *bootstrap) Init() tea.Cmd {
	return nil
}

func (m *bootstrap) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *bootstrap) View() string {
	if m.offline {
		return components.Centered(m.width, m.height, "◐ offline")
	}
	return components.Centered(m.width, m.height, "◐ deeploy.sh")
}
