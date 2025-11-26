package pages

import (
	"errors"
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

var (
	ErrNoInternetConnection = errors.New("no internet connection")
	ErrMissingConfig        = errors.New("missing config")
	ErrMissingServer        = errors.New("missing server")
	ErrMissingToken         = errors.New("missing token")
)

type checkInternetMsg struct {
	ok bool
}

type checkConfigMsg struct {
	ok  bool
	err error
}

type checkServerMsg struct {
	ok  bool
	err error
}

type checkAuthMsg struct {
	ok  bool
	err error
}

type state = int

const (
	stateCheckingInternet state = iota
	stateNoInternet
	stateCheckingConfig
	stateCheckingServer
	stateCheckingAuth
)

type bootstrap struct {
	internetOK    bool
	configOK      bool
	serverOK      bool
	authOK        bool
	state         state
	width, height int
	spinner       spinner.Model
	err           error
}

func NewBootstrap() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &bootstrap{
		spinner: s,
		state:   stateCheckingInternet,
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
			m.state = stateCheckingConfig
			return m, checkConfig
		}

		m.state = stateNoInternet
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return checkInternet()
		})

	case checkConfigMsg:
		if msg.ok {
			m.configOK = true
			m.state = stateCheckingServer
			return m, checkServer
		}
		m.err = msg.err
		return m, func() tea.Msg {
			return messages.ChangePageMsg{
				Page: NewConnectPage(m.err),
			}
		}

	case checkServerMsg:
		if msg.ok {
			m.serverOK = true
			m.state = stateCheckingAuth
			return m, checkAuth
		}
		m.err = msg.err
		return m, func() tea.Msg {
			return messages.ChangePageMsg{
				Page: NewConnectPage(m.err),
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
	var b strings.Builder

	spinner := m.spinner.View()

	switch m.state {
	case stateCheckingInternet:
		b.WriteString(spinner + " Checking internet...")
	case stateNoInternet:
		b.WriteString(spinner + " No internet. Retrying...")
	case stateCheckingConfig:
		b.WriteString(spinner + " Checking config...")
	case stateCheckingServer:
		b.WriteString(spinner + " Checking server...")
	case stateCheckingAuth:
		b.WriteString(spinner + " Checking auth...")
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
	isOnline := utils.IsOnline()
	return checkInternetMsg{
		ok: isOnline,
	}
}

func checkConfig() tea.Msg {
	time.Sleep(1 * time.Second)

	config, err := config.Load()

	if err != nil || config == nil {
		err = ErrMissingConfig
	} else if config.Server == "" {
		err = ErrMissingServer
	} else if config.Token == "" {
		err = ErrMissingToken
	}

	return checkConfigMsg{
		ok:  err == nil,
		err: err,
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
		utils.DeleteCfgToken()
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
