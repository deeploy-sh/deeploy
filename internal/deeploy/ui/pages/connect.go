package pages

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/deeploy-sh/deeploy/internal/deeploy/viewtypes"
)

type connectPage struct {
	serverInput textinput.Model
	status      string
	waiting     bool
	width       int
	height      int
	err         string
}

type authCallback struct {
	token string
	err   error
}

func NewConnectPage() connectPage {
	ti := textinput.New()
	ti.Placeholder = "e.g. 123.45.67.89:8090"
	ti.Focus()

	return connectPage{
		serverInput: ti,
	}
}

func (p connectPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m connectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		m.resetErr()
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			err := utils.ValidateServer(m.serverInput.Value())
			if err != nil {
				m.err = err.Error()
				return m, nil
			}
			m.waiting = true
			return m, m.startBrowserAuth()
		}
	case messages.AuthSuccessMsg:
		return m, func() tea.Msg {
			return viewtypes.Dashboard
		}
	}
	m.serverInput, cmd = m.serverInput.Update(msg)
	return m, cmd
}

func (p connectPage) View() string {
	var b strings.Builder

	if p.waiting {
		b.WriteString("✨ Browser opened for authentication. Waiting for completion.")
	} else {
		b.WriteString("CONNECT TO SERVER\n\n")
		b.WriteString(styles.FocusedStyle.Render("Server "))
		b.WriteString(p.serverInput.View())
		if p.err != "" {
			b.WriteString(styles.ErrorStyle.Render("\n* " + p.err))
		}
		if p.status != "" {
			b.WriteString(p.status)
		}
	}

	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	card := components.Card(components.CardProps{Width: 50}).Render(b.String())

	view := lipgloss.JoinVertical(0.5, logo, card)
	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func (p *connectPage) resetErr() {
	p.err = ""
}

// Starts a local server for ayth callback
func startLocalAuthServer() (int, chan authCallback) {
	callback := make(chan authCallback)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /callback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		token, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		callback <- authCallback{token: string(token)}
		w.Write([]byte("OK"))
		return
	})

	// Get a free random port
	listener, _ := net.Listen("tcp", "localhost:0")
	port := listener.Addr().(*net.TCPAddr).Port

	go http.Serve(listener, mux)

	return port, callback
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "bsd", etc.
		cmd = "xdg-open"
	}

	return exec.Command(cmd, append(args, url)...).Start()
}

func (m connectPage) startBrowserAuth() tea.Cmd {
	return func() tea.Msg {
		port, callback := startLocalAuthServer()

		// Open browser
		authURL := fmt.Sprintf(
			"%s?cli=true&port=%d",
			m.serverInput.Value(),
			port,
		)
		openBrowser(authURL)

		// Waiting for token
		result := <-callback
		if result.err != nil {
			return messages.AuthErrorMsg{Err: result.err}
		}

		// Save config
		cfg := config.Config{
			Server: m.serverInput.Value(),
			Token:  result.token,
		}
		if err := config.SaveConfig(&cfg); err != nil {
			return messages.AuthErrorMsg{Err: err}
		}

		return messages.ChangePageMsg{Page: NewDashboard()}
	}
}
