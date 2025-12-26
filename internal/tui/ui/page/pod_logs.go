package page

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
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

// logsViewport is a simple scrollable viewport for log lines
type logsViewport struct {
	lines  []string
	width  int
	height int
	offset int
}

func (v *logsViewport) setSize(w, h int) {
	v.width = w
	v.height = h
}

func (v *logsViewport) setLines(lines []string) {
	v.lines = lines
}

func (v *logsViewport) scrollUp(n int) {
	v.offset -= n
	if v.offset < 0 {
		v.offset = 0
	}
}

func (v *logsViewport) scrollDown(n int) {
	v.offset += n
	max := len(v.lines) - v.height
	if max < 0 {
		max = 0
	}
	if v.offset > max {
		v.offset = max
	}
}

func (v *logsViewport) gotoBottom() {
	v.offset = len(v.lines) - v.height
	if v.offset < 0 {
		v.offset = 0
	}
}

func (v *logsViewport) isAtBottom() bool {
	return v.offset >= len(v.lines)-v.height
}

func (v *logsViewport) view() string {
	if v.height <= 0 {
		return ""
	}

	lines := make([]string, v.height)
	for i := 0; i < v.height; i++ {
		idx := v.offset + i
		if idx < len(v.lines) {
			lines[i] = v.lines[idx]
		}
	}
	return strings.Join(lines, "\n")
}

type podLogs struct {
	store     msg.Store
	pod       *model.Pod
	project   *model.Project
	viewport  logsViewport
	logs      []string
	status    string
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
		viewport:  logsViewport{},
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

		// Keyboard scroll
		switch tmsg.String() {
		case "up", "k":
			m.viewport.scrollUp(1)
		case "down", "j":
			m.viewport.scrollDown(1)
		}
		return m, nil

	case tea.MouseWheelMsg:
		// Mouse scroll (3 lines per event)
		if tmsg.Button == tea.MouseWheelUp {
			m.viewport.scrollUp(3)
		} else if tmsg.Button == tea.MouseWheelDown {
			m.viewport.scrollDown(3)
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height

		m.cardProps = styles.CardProps{Width: m.width, Height: m.height, Padding: []int{1, 1}}

		const logsHeaderHeight = 2 // header + empty line
		m.viewport.setSize(m.cardProps.InnerWidth(), m.cardProps.InnerHeight()-logsHeaderHeight)
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

func (m *podLogs) updateViewport() {
	wasAtBottom := m.viewport.isAtBottom()
	m.viewport.setLines(m.logs)
	if wasAtBottom {
		m.viewport.gotoBottom()
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

	// Viewport content - just the lines, lipgloss handles overflow
	logsContent := m.viewport.view()

	content := lipgloss.JoinVertical(lipgloss.Left,
		headerLine,
		"",
		logsContent,
	)

	card := styles.Card(m.cardProps).Render(content)

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

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
