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

type app struct {
	currentPage      tea.Model
	width            int
	height           int
	heartbeatStarted bool
	offline          bool
	bootstrapped     bool
}

func NewApp() tea.Model {
	return &app{
		currentPage: NewBootstrap(),
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(
		m.currentPage.Init(),
		// INFO: use tick here to show bootstrap(logo) min. 1 second
		tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
			return utils.CheckConnection()
		}),
	)
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case utils.ConnectionResultMsg:
		switch {
		case msg.NeedsSetup:
			return m, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewConnectPage(nil)}
			}
		case msg.NeedsAuth:
			return m, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewAuthPage("")}
			}
		case msg.Offline:
			m.offline = true
			if !m.bootstrapped {
				bp, ok := m.currentPage.(*bootstrap)
				if ok {
					bp.offline = true
				}
				return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
					return utils.CheckConnection()
				})
			}
		case msg.Online:
			m.offline = false
		}

		if m.heartbeatStarted {
			return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return utils.CheckConnection()
			})
		}

		m.heartbeatStarted = true
		m.bootstrapped = true

		return m, tea.Batch(
			func() tea.Msg {
				return messages.ChangePageMsg{Page: NewDashboard()}
			},
			tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return utils.CheckConnection()
			}),
		)

	case tea.KeyMsg:
		if m.offline && msg.Type != tea.KeyCtrlC {
			return m, nil // block app
		}
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		var cmd tea.Cmd

		m.currentPage, cmd = m.currentPage.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - footerHeight - 2

		pageMsg := tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height, // Reduced height!
		}
		var cmd tea.Cmd
		if m.currentPage == nil {
			return m, nil
		}
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
	_, ok := m.currentPage.(*bootstrap)
	if ok {
		return m.currentPage.View()
	}

	var status string
	var statusStyle lipgloss.Style

	logo := "⚡ deeploy.sh"
	if m.offline {
		status = "● reconnecting"
		statusStyle = styles.OfflineStyle
	} else {
		status = "● online"
		statusStyle = styles.OnlineStyle
	}

	gap := max(m.width-lipgloss.Width(logo)-lipgloss.Width(status)-2, 1)

	headerContent := logo + strings.Repeat(" ", gap) + statusStyle.Render(status)

	header := lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Render(headerContent)

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

	view := lipgloss.JoinVertical(lipgloss.Left, header, m.currentPage.View(), footer)

	return lipgloss.Place(m.width, m.height+footerHeight, lipgloss.Left, lipgloss.Bottom, view)
}
