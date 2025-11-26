package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
)

const footerHeight = 1

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
		var cmd tea.Cmd

		m.currentPage, cmd = m.currentPage.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - footerHeight

		pageMsg := tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height, // Reduzierte Höhe!
		}
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(pageMsg)
		return m, cmd

	case messages.ChangePageMsg:
		newPage := msg.Page
		m.currentPage = newPage

		pageMsg := tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		}
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(pageMsg)
		return m, tea.Batch(
			newPage.Init(),
			cmd,
		)

	default:
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(msg)
		return m, cmd
	}
}

type FooterMenuItem struct {
	Key  string
	Desc string
}

func (m app) View() string {
	footerMenuItems := []FooterMenuItem{
		// {Key: "esc", Desc: "back"},
		{Key: "ctrl+c", Desc: "quit"},
	}

	var footerText strings.Builder

	for i, v := range footerMenuItems {
		footerText.WriteString(styles.FocusedStyle.Render(v.Key))
		footerText.WriteString(" ")
		footerText.WriteString(v.Desc)
		if len(footerMenuItems)-1 != i {
			footerText.WriteString(" • ")
		}
	}
	footer := lipgloss.
		NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		Render(footerText.String())

	view := lipgloss.JoinVertical(lipgloss.Left, m.currentPage.View(), footer)

	return lipgloss.Place(m.width, m.height+footerHeight, lipgloss.Left, lipgloss.Bottom, view)
}
