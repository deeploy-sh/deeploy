package pages

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
)

type checkInternetMsg struct {
	ok  bool
	err error
}

type checkConfigMsg struct {
	ok         bool
	needsSetup bool
	err        error
}

type checkServerMsg struct {
	ok  bool
	err error
}

type checkAuthMsg struct {
	ok  bool
	err error
}

type checkingState = int

const (
	checkingStateInternet checkingState = iota
	checkingStateConfig
	checkingStateServer
	checkingStateAuth
)

type bootstrap struct {
	internetOK    bool
	configOK      bool
	serverOK      bool
	authOK        bool
	checkingState checkingState
	width, height int
	spinner       spinner.Model
	err           error
}

func NewBootstrap() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &bootstrap{
		spinner:       s,
		checkingState: checkingStateInternet,
	}
}

func (m bootstrap) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		checkInternet,
	)
}

func (m bootstrap) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlA:
			return m, nil
		}

	case checkInternetMsg:
		if msg.ok {
			m.internetOK = true
			m.checkingState = checkingStateConfig
			return m, checkConfig
		}
		// todo: implement retry modal
		m.err = msg.err

	case checkConfigMsg:
		if msg.ok {
			m.configOK = true
			m.checkingState = checkingStateServer
			return m, checkServer
		}
		m.err = msg.err
		return m, func() tea.Msg {
			return messages.ChangePageMsg{
				Page: NewConnectPage(),
			}
		}

	case checkServerMsg:
		if msg.ok {
			m.serverOK = true
			m.checkingState = checkingStateAuth
			return m, checkAuth
		}
		m.err = msg.err
		return m, func() tea.Msg {
			return messages.ChangePageMsg{
				Page: NewConnectPage(),
			}
		}

	case checkAuthMsg:
		if msg.ok {
			m.authOK = true
			return m, func() tea.Msg {
				return messages.ChangePageMsg{
					Page: NewDashboard(),
				}
			}
		}
		m.err = msg.err
		return m, func() tea.Msg {
			return messages.ChangePageMsg{
				Page: NewAuthPage(""),
			}
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m bootstrap) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var b strings.Builder

	spinner := m.spinner.View()

	switch m.checkingState {
	case checkingStateInternet:
		b.WriteString(spinner + " checking internet connection...")
	case checkingStateConfig:
		b.WriteString(spinner + " checking config...")
	case checkingStateServer:
		b.WriteString(spinner + " checking server...")
	case checkingStateAuth:
		b.WriteString(spinner + " checking auth...")

	}

	logo := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	card := components.Card(components.CardProps{Width: 50}).Render(b.String())

	view := lipgloss.JoinVertical(0.5, logo, card)
	layout := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}

func checkInternet() tea.Msg {
	time.Sleep(1 * time.Second)

	endpoints := []string{
		"https://www.google.com",
		"https://1.1.1.1", // Cloudflare
		"https://8.8.8.8", // Google DNS
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, endpoint := range endpoints {
		req, _ := http.NewRequest(http.MethodHead, endpoint, nil)
		_, err := client.Do(req)
		if err == nil {
			return checkInternetMsg{
				ok:  true,
				err: nil,
			}
		}
	}

	return checkInternetMsg{
		ok:  false,
		err: fmt.Errorf("no internet connection"),
	}
}

func checkConfig() tea.Msg {
	time.Sleep(1 * time.Second)

	config, err := config.Load()
	needsSetup := err != nil || config == nil || config.Server == "" || config.Token == ""

	return checkConfigMsg{
		ok:         !needsSetup,
		needsSetup: needsSetup,
		err:        nil,
	}
}

func checkServer() tea.Msg {
	time.Sleep(1 * time.Second)

	config, err := config.Load()
	if err != nil {
		return checkServerMsg{
			ok:  false,
			err: err,
		}

	}

	err = utils.ValidateServer(config.Server)
	if err != nil {
		return checkServerMsg{
			ok:  false,
			err: err,
		}
	}

	return checkServerMsg{
		ok:  true,
		err: nil,
	}
}

func checkAuth() tea.Msg {
	time.Sleep(1 * time.Second)

	_, err := utils.Request(utils.RequestProps{
		Method: "GET",
		URL:    "/dashboard",
	})
	if err != nil {
		return checkAuthMsg{
			ok:  false,
			err: err,
		}
	}

	return checkAuthMsg{
		ok:  true,
		err: nil,
	}
}
