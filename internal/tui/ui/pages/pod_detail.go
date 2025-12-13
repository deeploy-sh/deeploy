package pages

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type PodDetailPage struct {
	pod         *repo.Pod
	project     *repo.Project
	domains     []api.PodDomain
	envVarCount int
	loading     bool
	keyDeploy   key.Binding
	keyStop     key.Binding
	keyRestart  key.Binding
	keyLogs     key.Binding
	keyEdit     key.Binding
	keyDomains  key.Binding
	keyVars     key.Binding
	keyBack     key.Binding
	width       int
	height      int
}

func (m PodDetailPage) HelpKeys() []key.Binding {
	return []key.Binding{m.keyDeploy, m.keyStop, m.keyRestart, m.keyLogs, m.keyEdit, m.keyDomains, m.keyVars, m.keyBack}
}

func NewPodDetailPage(pod *repo.Pod, project *repo.Project) PodDetailPage {
	return PodDetailPage{
		pod:        pod,
		project:    project,
		keyDeploy:  key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "deploy")),
		keyStop:    key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "stop")),
		keyRestart: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "restart")),
		keyLogs:    key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logs")),
		keyEdit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		keyDomains: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "domains")),
		keyVars:    key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "env vars")),
		keyBack:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

func (m PodDetailPage) Init() tea.Cmd {
	return tea.Batch(
		api.FetchPodDomains(m.pod.ID),
		api.FetchPodEnvVars(m.pod.ID),
	)
}

func (m PodDetailPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.PodDomainsLoaded:
		if domains, ok := tmsg.Domains.([]api.PodDomain); ok {
			m.domains = domains
		}
		return m, nil

	case msg.PodEnvVarsLoaded:
		if envVars, ok := tmsg.EnvVars.([]model.PodEnvVar); ok {
			m.envVarCount = len(envVars)
		}
		return m, nil

	case msg.PodDeployed, msg.PodStopped, msg.PodRestarted:
		m.loading = false
		return m, api.LoadData()

	case msg.DataLoaded:
		for _, p := range tmsg.Pods {
			if p.ID == m.pod.ID {
				pod := p
				m.pod = &pod
				break
			}
		}
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKeyPress(tmsg)

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	return m, nil
}

func (m PodDetailPage) handleKeyPress(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		projectID := m.project.ID
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewProjectDetailPage(s, projectID)
				},
			}
		}

	case key.Matches(tmsg, m.keyDeploy):
		podID := m.pod.ID
		podTitle := m.pod.Title
		return m, tea.Batch(
			api.DeployPod(podID),
			func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewPodLogsPage(podID, podTitle)
					},
				}
			},
		)

	case key.Matches(tmsg, m.keyStop):
		m.loading = true
		return m, api.StopPod(m.pod.ID)

	case key.Matches(tmsg, m.keyRestart):
		m.loading = true
		return m, api.RestartPod(m.pod.ID)

	case key.Matches(tmsg, m.keyLogs):
		podID := m.pod.ID
		podTitle := m.pod.Title
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodLogsPage(podID, podTitle)
				},
			}
		}

	case key.Matches(tmsg, m.keyEdit):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodEditPage(pod, project)
				},
			}
		}

	case key.Matches(tmsg, m.keyDomains):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDomainsPage(pod, project)
				},
			}
		}

	case key.Matches(tmsg, m.keyVars):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodVarsPage(pod, project)
				},
			}
		}
	}

	return m, nil
}

func (m PodDetailPage) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	statusStyle := lipgloss.NewStyle().Background(styles.ColorBackgroundPanel())

	b.WriteString(titleStyle.Render(m.pod.Title))
	b.WriteString(statusStyle.Render(" " + m.renderStatus()))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorMuted())

	// Repo Info
	b.WriteString(labelStyle.Render("Repo"))
	b.WriteString("\n")
	if m.pod.RepoURL != nil && *m.pod.RepoURL != "" {
		b.WriteString(*m.pod.RepoURL)
		b.WriteString(" @ ")
		b.WriteString(m.pod.Branch)
	} else {
		b.WriteString(styles.MutedStyle().Render("(not configured)"))
	}
	b.WriteString("\n\n")

	// Dockerfile
	b.WriteString(labelStyle.Render("Dockerfile"))
	b.WriteString("\n")
	if m.pod.DockerfilePath != "" {
		b.WriteString(m.pod.DockerfilePath)
	} else {
		b.WriteString("Dockerfile")
	}
	b.WriteString("\n\n")

	// Domains
	b.WriteString(labelStyle.Render("Domains"))
	b.WriteString("\n")
	if len(m.domains) > 0 {
		b.WriteString(fmt.Sprintf("%d configured", len(m.domains)))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("(none)"))
	}
	b.WriteString("\n\n")

	// Env Vars
	b.WriteString(labelStyle.Render("Env Vars"))
	b.WriteString("\n")
	if m.envVarCount > 0 {
		b.WriteString(fmt.Sprintf("%d configured", m.envVarCount))
	} else {
		b.WriteString(styles.MutedStyle().Render("(none)"))
	}

	if m.loading {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Loading..."))
	}

	card := styles.Card(styles.CardProps{
		Width:   70,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(b.String())

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m PodDetailPage) renderStatus() string {
	status := m.pod.Status
	if status == "" {
		status = "idle"
	}
	style := lipgloss.NewStyle()
	switch status {
	case "running":
		style = style.Foreground(lipgloss.Color("10"))
	case "failed":
		style = style.Foreground(lipgloss.Color("9"))
	case "building":
		style = style.Foreground(lipgloss.Color("11"))
	default:
		style = style.Foreground(lipgloss.Color("8"))
	}
	return style.Render("[" + status + "]")
}

func (m PodDetailPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title}
}
