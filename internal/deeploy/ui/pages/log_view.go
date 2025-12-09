package pages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
)

// logsResponse from the API
type logsResponse struct {
	Logs   []string `json:"logs"`
	Status string   `json:"status"`
}

// pollLogsMsg triggers a poll
type pollLogsMsg struct{}

// logsUpdated contains fetched logs
type logsUpdated struct {
	logs   []string
	status string
}

type LogViewPage struct {
	podID     string
	podTitle  string
	viewport  viewport.Model
	logs      []string
	status    string
	keyBack   key.Binding
	keyDeploy key.Binding
	width     int
	height    int
}

func NewLogViewPage(podID, podTitle string) LogViewPage {
	vp := viewport.New()
	return LogViewPage{
		podID:     podID,
		podTitle:  podTitle,
		viewport:  vp,
		status:    "building",
		keyBack:   key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc/q", "back")),
		keyDeploy: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "redeploy")),
	}
}

func (m LogViewPage) Init() tea.Cmd {
	return tea.Batch(
		m.fetchLogs(),
		m.schedulePoll(),
	)
}

func (m LogViewPage) schedulePoll() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return pollLogsMsg{}
	})
}

func (m LogViewPage) fetchLogs() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.Load()
		if err != nil {
			return logsUpdated{logs: []string{"Error: " + err.Error()}, status: "error"}
		}

		url := fmt.Sprintf("%s/api/pods/%s/logs", cfg.Server, m.podID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return logsUpdated{logs: []string{"Error: " + err.Error()}, status: "error"}
		}
		req.Header.Set("Authorization", "Bearer "+cfg.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return logsUpdated{logs: []string{"Error: " + err.Error()}, status: "error"}
		}
		defer resp.Body.Close()

		var result logsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return logsUpdated{logs: []string{"Error: " + err.Error()}, status: "error"}
		}

		return logsUpdated{logs: result.Logs, status: result.Status}
	}
}

func (m LogViewPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case pollLogsMsg:
		// Keep polling while building
		if m.status == "building" {
			return m, tea.Batch(m.fetchLogs(), m.schedulePoll())
		}
		return m, nil

	case logsUpdated:
		m.logs = tmsg.logs
		m.status = tmsg.status
		m.updateViewport()
		return m, nil

	case tea.KeyPressMsg:
		if key.Matches(tmsg, m.keyBack) {
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			}
		}
		if key.Matches(tmsg, m.keyDeploy) {
			// Redeploy - trigger deploy and restart polling
			m.status = "building"
			m.logs = []string{"Starting new deployment..."}
			m.updateViewport()
			return m, tea.Batch(
				m.triggerDeploy(),
				m.schedulePoll(),
			)
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(tmsg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		m.viewport.SetWidth(m.width)
		m.viewport.SetHeight(m.height - 4)
		m.updateViewport()
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(tmsg)
	return m, cmd
}

func (m LogViewPage) triggerDeploy() tea.Cmd {
	return func() tea.Msg {
		cfg, _ := config.Load()
		url := fmt.Sprintf("%s/api/pods/%s/deploy", cfg.Server, m.podID)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
		http.DefaultClient.Do(req)
		return pollLogsMsg{}
	}
}

func (m *LogViewPage) updateViewport() {
	content := strings.Join(m.logs, "\n")
	m.viewport.SetContent(content)
	if m.status == "building" {
		m.viewport.GotoBottom()
	}
}

func (m LogViewPage) View() tea.View {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	header := titleStyle.Render(fmt.Sprintf("Build Logs: %s", m.podTitle))

	var statusText string
	switch m.status {
	case "building":
		statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("● building...")
	case "running":
		statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("● running")
	case "failed":
		statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("● failed")
	default:
		statusText = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("● " + m.status)
	}

	headerLine := fmt.Sprintf("%s  %s", header, statusText)
	help := styles.MutedStyle().Render("esc: back  D: redeploy  ↑↓: scroll")

	content := lipgloss.JoinVertical(lipgloss.Left,
		headerLine,
		"",
		m.viewport.View(),
		"",
		help,
	)

	return tea.NewView(content)
}

func (m LogViewPage) Breadcrumbs() []string {
	return []string{"Build Logs", m.podTitle}
}
