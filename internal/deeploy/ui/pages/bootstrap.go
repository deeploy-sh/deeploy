package pages

import (
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
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

type checkTokenMsg struct {
	ok  bool
	err error
}

type checkingState = int

const (
	checkingStateInternet checkingState = iota
	checkingStateConfig
	checkingStateServer
	checkingStateToken
)

type bootstrap struct {
	internetOK    bool
	configOK      bool
	serverOK      bool
	tokenOK       bool
	checkingState checkingState
	err           error
}

func NewBootstrap() tea.Model {
	return &bootstrap{
		checkingState: checkingStateInternet,
	}
}

func (m bootstrap) Init() tea.Cmd {
	return checkInternet
}

func (m bootstrap) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			m.checkingState = checkingStateToken
			return m, checkToken
		}
		m.err = msg.err

	case checkTokenMsg:
		if msg.ok {
			m.tokenOK = true
		}
		m.err = msg.err
	}

	return m, nil
}

func (m bootstrap) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	switch m.checkingState {
	case checkingStateInternet:
		return "checking internet connection ..."
	case checkingStateConfig:
		return "checking config ..."
	case checkingStateServer:
		return "checking server connection ..."
	case checkingStateToken:
		return "checking token validation ..."
	}

	return ""
}

func checkInternet() tea.Msg {
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
	config, err := config.LoadConfig()
	needsSetup := err != nil || config == nil || config.Server == "" || config.Token == ""

	return checkConfigMsg{
		ok:         !needsSetup,
		needsSetup: needsSetup,
		err:        nil,
	}
}

func checkServer() tea.Msg {
	config, err := config.LoadConfig()
	if err != nil {
		return checkServerMsg{
			ok:  false,
			err: err,
		}
	}

	req, err := http.NewRequest(http.MethodHead, config.Server, nil)
	if err != nil {
		return checkServerMsg{
			ok:  false,
			err: err,
		}
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	_, err = client.Do(req)
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

func checkToken() tea.Msg {
	return ""
}
