package pages

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
)

type authKeyMap struct {
	Authenticate key.Binding
	Quit         key.Binding
}

func (k authKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Authenticate, k.Quit}
}

func (k authKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Authenticate, k.Quit}}
}

func newAuthKeyMap() authKeyMap {
	return authKeyMap{
		Authenticate: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "authenticate")),
		Quit:         key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

type authPage struct {
	keys      authKeyMap
	help      help.Model
	isReauth  bool
	waiting   bool
	width     int
	height    int
	serverURL string
	err       string
}

type authCallback struct {
	token string
	err   error
}

func NewAuthPage(server string) authPage {
	isReauth := server == ""
	if isReauth {
		cfg, _ := config.Load()
		server = cfg.Server
	}
	return authPage{
		keys:      newAuthKeyMap(),
		help:      styles.NewHelpModel(),
		isReauth:  isReauth,
		serverURL: server,
	}
}

func (p authPage) Init() tea.Cmd {
	return nil
}

func (m authPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, m.startBrowserAuth()
		}
	case msg.AuthSuccess:
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
			}
		}
	}
	return m, cmd
}

func (p authPage) View() tea.View {
	var b strings.Builder

	if p.waiting {
		b.WriteString("Browser opened for authentication.\nWaiting for completion...")
	} else {
		if p.isReauth {
			b.WriteString(fmt.Sprintf("> Re-authenticating %v \n", p.serverURL))
			b.WriteString("> Press enter to authenticate")
		} else {
			b.WriteString(fmt.Sprintf("> Authenticating %v \n", p.serverURL))
			b.WriteString("> Press enter to authenticate")
		}
	}

	card := components.Card(components.CardProps{Width: 50, Padding: []int{1, 2}, Accent: true}).Render(b.String())
	helpView := p.help.View(p.keys)
	contentHeight := p.height - 1

	centered := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centered, helpView))
}

func (p authPage) Breadcrumbs() []string {
	return []string{"Auth"}
}

func (p *authPage) resetErr() {
	p.err = ""
}

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

func (m authPage) startBrowserAuth() tea.Cmd {
	return func() tea.Msg {
		port, callback := startLocalAuthServer()

		authURL := fmt.Sprintf(
			"%s?cli=true&port=%d",
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
		if err := config.Save(&cfg); err != nil {
			return msg.AuthError{Err: err}
		}

		return msg.ChangePage{
			PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
		}
	}
}
