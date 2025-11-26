package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
)

type app struct {
	currentPage tea.Model
	width       int
	height      int
}

func NewApp() app {
	return app{
		currentPage: NewBootstrap(),
	}
}

func (m app) Init() tea.Cmd {
	return m.currentPage.Init()
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// return a, cmd
		// return a, nil

	case messages.ChangePageMsg:
		newPage := msg.Page
		m.currentPage = newPage

		// Batch window size and init commands together
		// This prevents double rendering by ensuring both happen in sequence
		return m, tea.Batch(
			func() tea.Msg {
				return tea.WindowSizeMsg{
					Width:  m.width,
					Height: m.height,
				}
			},
			// INFO: do we really need this here?
			// newPage.Init(),
		)
	}

	// All other messages go to current page
	currentPage := m.currentPage
	updatedPage, cmd := currentPage.Update(msg)
	m.currentPage = updatedPage
	return m, cmd
}

type FooterMenuItem struct {
	Key  string
	Desc string
}

func (m app) View() string {
	main := m.currentPage.View()

	footerMenuItems := []FooterMenuItem{
		{Key: ":", Desc: "menu"},
		{Key: "esc", Desc: "back"},
		{Key: "q", Desc: "quit"},
	}

	var footer strings.Builder

	for i, v := range footerMenuItems {
		footer.WriteString(styles.FocusedStyle.Render(v.Key))
		footer.WriteString(" ")
		footer.WriteString(v.Desc)
		if len(footerMenuItems)-1 != i {
			footer.WriteString(" • ")
		}
	}

	footerCard := components.Card(components.CardProps{
		Width:   m.width,
		Padding: []int{0, 1},
	}).Render(footer.String())

	horizontal := lipgloss.JoinHorizontal(0.5, main)
	view := lipgloss.JoinVertical(lipgloss.Bottom, horizontal, footerCard)

	return view
}
