package page

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
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
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

type podLogs struct {
	store     msg.Store
	pod       *model.Pod
	project   *model.Project
	viewport  viewport.Model
	logs      []string
	status    string
	keyBack   key.Binding
	keyDeploy key.Binding
	width     int
	height    int
}

func NewPodLogs(s msg.Store, podID string) podLogs {
	var pod model.Pod
	for _, p := range s.Pods() {
		if p.ID == podID {
			pod = p
			break
		}
	}

	var project model.Project
	for _, pr := range s.Projects() {
		if pr.ID == pod.ProjectID {
			project = pr
			break
		}
	}

	vp := viewport.New()
	return podLogs{
		store:     s,
		pod:       &pod,
		project:   &project,
		viewport:  vp,
		status:    "building",
		keyBack:   key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc/q", "back")),
		keyDeploy: key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "redeploy")),
	}
}

func (m podLogs) Init() tea.Cmd {
	return tea.Batch(
		m.fetchLogs(),
		m.schedulePoll(),
	)
}

func (m podLogs) schedulePoll() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return pollLogsMsg{}
	})
}

func (m podLogs) fetchLogs() tea.Cmd {
	return func() tea.Msg {
		cfg, err := config.Load()
		if err != nil {
			return logsUpdated{logs: []string{"Error: " + err.Error()}, status: "error"}
		}

		url := fmt.Sprintf("%s/api/pods/%s/logs", cfg.Server, m.pod.ID)
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
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return logsUpdated{logs: []string{"Error: " + err.Error()}, status: "error"}
		}

		return logsUpdated{logs: result.Logs, status: result.Status}
	}
}

func (m podLogs) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
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
			podID := m.pod.ID
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewPodDetail(s, podID) },
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
		// Card width: responsive, max 120
		cardWidth := m.width - 8
		if cardWidth > 120 {
			cardWidth = 120
		}
		// Inner width: card width minus padding (2 on each side) and accent border (1)
		innerWidth := cardWidth - 5
		// Height: total minus card padding (1 top, 1 bottom), header, help, spacing
		innerHeight := m.height - 10
		m.viewport.SetWidth(innerWidth)
		m.viewport.SetHeight(innerHeight)
		m.updateViewport()
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(tmsg)
	return m, cmd
}

func (m podLogs) triggerDeploy() tea.Cmd {
	return func() tea.Msg {
		cfg, _ := config.Load()
		url := fmt.Sprintf("%s/api/pods/%s/deploy", cfg.Server, m.pod.ID)
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
		http.DefaultClient.Do(req)
		return pollLogsMsg{}
	}
}

func (m *podLogs) updateViewport() {
	content := strings.Join(m.logs, "\n")
	m.viewport.SetContent(content)
	if m.status == "building" {
		m.viewport.GotoBottom()
	}
}

func (m podLogs) View() tea.View {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	header := titleStyle.Render(fmt.Sprintf("Build Logs: %s", m.pod.Title))

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

	// Card width: responsive, max 120
	cardWidth := m.width - 8
	if cardWidth > 120 {
		cardWidth = 120
	}

	card := styles.Card(styles.CardProps{
		Width:   cardWidth,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m podLogs) Breadcrumbs() []string {
	return []string{"Build Logs", m.pod.Title}
}
