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
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
)

var (
	ErrNoInternet    = errors.New("no internet")
	ErrMissingConfig = errors.New("missing config")
	ErrMissingServer = errors.New("missing server")
	ErrMissingToken  = errors.New("missing token")
)

type checkInternetMsg struct {
	err error
}

type checkConfigMsg struct {
	err error
}

type checkServerMsg struct {
	err error
}

type checkAuthMsg struct {
	err error
}

type bootstrap struct {
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
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
				return checkInternet()
			})
		}
		return m, checkConfig

	case checkConfigMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, func() tea.Msg {
				return messages.ChangePageMsg{
					Page: NewConnectPage(m.err),
				}
			}

		}
		return m, checkServer

	case checkServerMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, func() tea.Msg {
				return messages.ChangePageMsg{
					Page: NewConnectPage(m.err),
				}
			}
		}
		return m, checkAuth

	case checkAuthMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, func() tea.Msg {
				return messages.ChangePageMsg{
					Page: NewAuthPage(""),
				}
			}
		}
		return m, func() tea.Msg {
			return messages.ChangePageMsg{
				Page: NewDashboard(),
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

	if m.err != nil {
		if errors.Is(m.err, ErrNoInternet) {
			b.WriteString(spinner + " No internet. Retrying...")
		}
	} else {
		b.WriteString(spinner + " deeploy.sh")
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, b.String())
}

func checkInternet() tea.Msg {
	isOnline := utils.IsOnline()
	var err error
	if !isOnline {
		err = ErrNoInternet
	}
	return checkInternetMsg{
		err: err,
	}
}

func checkConfig() tea.Msg {
	config, err := config.Load()

	if err != nil || config == nil {
		err = ErrMissingConfig
	} else if config.Server == "" {
		err = ErrMissingServer
	} else if config.Token == "" {
		err = ErrMissingToken
	}

	return checkConfigMsg{
		err: err,
	}
}

func checkServer() tea.Msg {
	config, err := config.Load()
	if err != nil {
		return checkServerMsg{
			err: err,
		}

	}

	err = utils.ValidateServer(config.Server)
	if err != nil {
		return checkServerMsg{
			err: err,
		}
	}

	return checkServerMsg{
		err: nil,
	}
}

func checkAuth() tea.Msg {
	_, err := utils.Request(utils.RequestProps{
		Method: "GET",
		URL:    "/dashboard",
	})
	if err != nil {
		utils.DeleteCfgToken()
		return checkAuthMsg{
			err: err,
		}
	}

	return checkAuthMsg{
		err: nil,
	}
}
