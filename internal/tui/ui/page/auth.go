package page

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/tui/utils"
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
			return m, m.startBrowserAuth()
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

	card := styles.Card(styles.CardProps{Width: styles.CardWidthMD, Padding: []int{1, 2}, Accent: true}).Render(b.String())

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

// generateSessionID creates a random session ID for CLI auth
func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// startBrowserAuth opens browser and polls for auth completion
func (m auth) startBrowserAuth() tea.Cmd {
	return func() tea.Msg {
		sessionID := generateSessionID()

		// Open browser with session ID
		authURL := fmt.Sprintf("%s/auth?cli=true&session=%s", m.serverURL, sessionID)
		utils.OpenBrowser(authURL)

		// Poll for auth completion
		pollURL := fmt.Sprintf("%s/api/auth/poll?session=%s", m.serverURL, sessionID)
		client := &http.Client{Timeout: 10 * time.Second}

		for i := 0; i < 150; i++ { // 5 minutes max (150 * 2s)
			time.Sleep(2 * time.Second)

			resp, err := client.Get(pollURL)
			if err != nil {
				continue // Network error, retry
			}

			if resp.StatusCode == http.StatusNotFound {
				resp.Body.Close()
				return msg.AuthError{Err: errors.New("session expired")}
			}

			if resp.StatusCode == http.StatusAccepted {
				resp.Body.Close()
				continue // Still pending
			}

			if resp.StatusCode == http.StatusOK {
				var result struct {
					Token string `json:"token"`
				}
				json.NewDecoder(resp.Body).Decode(&result)
				resp.Body.Close()

				if result.Token != "" {
					cfg := config.Config{
						Server: m.serverURL,
						Token:  result.Token,
					}
					if err := config.Save(&cfg); err != nil {
						return msg.AuthError{Err: err}
					}
					return msg.AuthSuccess{}
				}
			}

			resp.Body.Close()
		}

		return msg.AuthError{Err: errors.New("authentication timeout")}
	}
}
