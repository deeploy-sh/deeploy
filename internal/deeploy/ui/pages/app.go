package pages

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
)

const footerHeight = 1

type heartbeatMsg struct {
	ok bool
}

type app struct {
	currentPage      tea.Model
	width            int
	height           int
	isOffline        bool
	heartbeatStarted bool
}

func NewApp() app {
	return app{
		currentPage: NewBootstrap(),
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(m.currentPage.Init())
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case heartbeatMsg:
		m.isOffline = !msg.ok
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return checkHearbeat()
		})

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
		m.currentPage = msg.Page

		pageMsg := tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		}
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(pageMsg)

		_, initial := m.currentPage.(DashboardPage)
		if initial && !m.heartbeatStarted {
			m.heartbeatStarted = true
			return m, tea.Batch(
				m.currentPage.Init(),
				checkHearbeat,
				cmd,
			)
		}

		return m, tea.Batch(
			m.currentPage.Init(),
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
	if m.isOffline && m.heartbeatStarted {
		logo := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render("🔥deeploy.sh\n")

		view := lipgloss.JoinVertical(0.5, logo, "No internet. Retrying...")
		layout := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
		return layout

	}

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

func checkHearbeat() tea.Msg {
	isOnline := utils.IsOnline()
	return heartbeatMsg{
		ok: isOnline,
	}
}
