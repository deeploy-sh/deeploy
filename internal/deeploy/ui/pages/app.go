package pages

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
)

const footerHeight = 1

type state int

const (
	StateConnecting state = iota
	StateOffline
	StateNeedsSetup
	StateNeedsAuth
	StateReady
)

type app struct {
	currentPage      tea.Model
	width            int
	height           int
	heartbeatStarted bool
	state            state
	//
	serverInput textinput.Model
	tokenInput  textinput.Model
}

func NewApp() tea.Model {
	si := textinput.New()
	si.Placeholder = "https://deeploy.example.com"

	ti := textinput.New()
	ti.Placeholder = "your-token"
	ti.EchoMode = textinput.EchoPassword

	return &app{
		state:       StateConnecting,
		serverInput: si,
		tokenInput:  ti,
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(utils.CheckConnection)
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case utils.ConnectionResultMsg:
		switch {
		case msg.NeedsSetup:
			m.state = StateNeedsSetup
		case msg.NeedsAuth:
			m.state = StateNeedsAuth
		case msg.Offline:
			m.state = StateOffline
		default:
			m.state = StateReady
		}

		if m.heartbeatStarted {
			return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return utils.CheckConnection()
			})
		}

		m.heartbeatStarted = true
		return m, tea.Batch(
			func() tea.Msg {
				return messages.ChangePageMsg{Page: NewDashboard()}
			},
			tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return utils.CheckConnection()
			}),
		)

	case tea.KeyMsg:
		if m.state == StateOffline && msg.Type != tea.KeyCtrlC {
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

		if !m.heartbeatStarted {
			m.heartbeatStarted = true
			return m, tea.Batch(
				m.currentPage.Init(),
				utils.CheckConnection,
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

func (m app) centered(content string) string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m app) View() string {
	if m.currentPage == nil {
		return ""
	}
	switch m.state {
	// case StateOffline:
	// 	return m.centered("◐ Offline")

	case StateConnecting:
		return m.centered("◐ deeploy.sh")

	case StateNeedsSetup:
		return m.centered(fmt.Sprintf(
			"⚡ deeploy.sh\n\n"+
				"Server URL:\n%s\n\n"+
				"[Enter] connect",
			m.serverInput.View(),
		))

	case StateNeedsAuth:
		return m.centered(fmt.Sprintf(
			"⚡ deeploy.sh\n\n"+
				"Token:\n%s\n\n"+
				"[Enter] login",
			m.tokenInput.View(),
		))
	}

	var status string
	var statusStyle lipgloss.Style

	logo := "⚡ deeploy.sh"
	if m.state == StateOffline {
		status = "● reconnecting"
		statusStyle = styles.OfflineStyle
	} else {
		status = "● online"
		statusStyle = styles.OnlineStyle
	}

	gap := m.width - lipgloss.Width(logo) - lipgloss.Width(status) - 2
	if gap < 1 {
		gap = 1
	}

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
