package page

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type auth struct {
	keyauthenticate key.Binding
	keyQuit         key.Binding
	isReauth        bool
	waiting         bool
	width           int
	height          int
	serverURL       string
	err             string
}

func (p auth) HelpKeys() []key.Binding {
	return []key.Binding{p.keyauthenticate, p.keyQuit}
}

type authCallback struct {
	token string
	err   error
}

func NewAuth(server string) auth {
	isReauth := server == ""
	if isReauth {
		cfg, _ := config.Load()
		server = cfg.Server
	}
	return auth{
		keyauthenticate: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "authenticate")),
		keyQuit:         key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
		isReauth:        isReauth,
		serverURL:       server,
	}
}

func (p auth) Init() tea.Cmd {
	return nil
}

func (m auth) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
	case tea.KeyPressMsg:
		m.resetErr()
		switch {
		case tmsg.String() == "ctrl+c" || tmsg.Code == tea.KeyEscape:
			return m, tea.Quit
		case tmsg.Code == tea.KeyEnter:
			m.waiting = true
			return m, m.startBrowserauth()
		}
	}
	return m, cmd
}

func (p auth) View() tea.View {
	var b strings.Builder

	if p.waiting {
		b.WriteString("Waiting for browser authentication...")
	} else {
		title := "authenticate"
		if p.isReauth {
			title = "Re-authenticate"
		}
		b.WriteString(title + "\n\n")
		b.WriteString("Server: " + p.serverURL + "\n\n")
		b.WriteString("Press enter to open browser")
	}

	card := styles.Card(styles.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}).Render(b.String())

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p auth) Breadcrumbs() []string {
	return []string{"auth"}
}

func (p *auth) resetErr() {
	p.err = ""
}

func startLocalauthServer() (int, chan authCallback) {
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
	})

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
	default:
		cmd = "xdg-open"
	}

	return exec.Command(cmd, append(args, url)...).Start()
}

func (m auth) startBrowserauth() tea.Cmd {
	return func() tea.Msg {
		port, callback := startLocalauthServer()

		authURL := fmt.Sprintf(
			"%s/auth?cli=true&port=%d",
			m.serverURL,
			port,
		)
		openBrowser(authURL)

		result := <-callback
		if result.err != nil {
			return msg.AuthError{Err: result.err}
		}

		cfg := config.Config{
			Server: m.serverURL,
			Token:  result.token,
		}
		err := config.Save(&cfg)
		if err != nil {
			return msg.AuthError{Err: err}
		}

		return msg.AuthSuccess{}
	}
}
