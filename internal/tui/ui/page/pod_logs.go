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

// podLogs displays streaming build/container logs with auto-scroll.
// Uses bubbles/viewport for proper scrolling and dimension constraints.
type podLogs struct {
	store     msg.Store
	pod       *model.Pod
	project   *model.Project
	viewport  viewport.Model // handles scrolling, truncation, rendering
	logs      []string       // raw log lines from API
	status    string         // building, running, failed
	keyBack   key.Binding
	keyDeploy key.Binding
	width     int
	height    int
	cardProps styles.CardProps
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

	return podLogs{
		store:     s,
		pod:       &pod,
		project:   &project,
		viewport:  viewport.New(),
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
	var cmd tea.Cmd

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

		// viewport handles up/down/pgup/pgdown/home/end natively
		m.viewport, cmd = m.viewport.Update(tmsg)
		return m, cmd

	case tea.MouseWheelMsg:
		// viewport handles mouse scroll natively
		m.viewport, cmd = m.viewport.Update(tmsg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		m.cardProps = styles.CardProps{Width: m.width, Padding: []int{1, 1}}

		// viewport height = available - card padding (2) - header area (2)
		viewportHeight := m.height - 4
		m.viewport.SetWidth(m.cardProps.InnerWidth())
		m.viewport.SetHeight(viewportHeight)
		m.updateViewport()
		return m, nil
	}

	return m, nil
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

// updateViewport syncs logs to viewport with "follow mode":
// - If user was at bottom, stay at bottom (follow new logs)
// - If user scrolled up, stay there (let them read)
func (m *podLogs) updateViewport() {
	wasAtBottom := m.viewport.AtBottom()
	m.viewport.SetContent(strings.Join(m.logs, "\n"))
	if wasAtBottom {
		m.viewport.GotoBottom()
	}
}

func (m podLogs) View() tea.View {
	if m.height == 0 {
		return tea.NewView("Loading...")
	}

	bg := styles.ColorBackgroundPanel()
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary()).Background(bg)
	header := titleStyle.Render(fmt.Sprintf("Build Logs: %s", m.pod.Title))

	var statusText string
	switch m.status {
	case "building":
		statusText = styles.WarningStyle().Background(bg).Render("● building...")
	case "running":
		statusText = styles.SuccessStyle().Background(bg).Render("● running")
	case "failed":
		statusText = styles.ErrorStyle().Background(bg).Render("● failed")
	default:
		statusText = styles.MutedStyle().Background(bg).Render("● " + m.status)
	}

	spacer := lipgloss.NewStyle().Background(bg).Render("  ")
	headerText := header + spacer + statusText

	// Extend header to full width with background
	headerLine := lipgloss.NewStyle().
		Width(m.cardProps.InnerWidth()).
		Background(bg).
		Render(headerText)

	// Viewport content
	logsContent := m.viewport.View()

	content := lipgloss.JoinVertical(lipgloss.Left,
		headerLine,
		"",
		logsContent,
	)

	card := styles.Card(m.cardProps).Render(content)

	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m podLogs) Breadcrumbs() []string {
	return []string{"Build Logs", m.pod.Title}
}

func (m podLogs) HelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "redeploy")),
		key.NewBinding(key.WithKeys("up", "down"), key.WithHelp("↑↓", "scroll")),
	}
}
