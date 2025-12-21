package page

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type podDetail struct {
	store       msg.Store
	pod         *model.Pod
	project     *model.Project
	domains     []model.PodDomain
	envVarCount int
	keyDeploy   key.Binding
	keyStop     key.Binding
	keyRestart  key.Binding
	keyLogs     key.Binding
	keyEdit     key.Binding
	keyDomains  key.Binding
	keyVars     key.Binding
	keyToken    key.Binding
	keyBack     key.Binding
	width       int
	height      int
}

func (m podDetail) HelpKeys() []key.Binding {
	return []key.Binding{m.keyDeploy, m.keyStop, m.keyRestart, m.keyLogs, m.keyEdit, m.keyDomains, m.keyVars, m.keyToken, m.keyBack}
}

func NewPodDetail(s msg.Store, podID string) podDetail {
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

	return podDetail{
		store:      s,
		pod:        &pod,
		project:    &project,
		keyDeploy:  key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "deploy")),
		keyStop:    key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "stop")),
		keyRestart: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "restart")),
		keyLogs:    key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logs")),
		keyEdit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
		keyDomains: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "domains")),
		keyVars:    key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "env vars")),
		keyToken:   key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "token")),
		keyBack:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

func (m podDetail) Init() tea.Cmd {
	return tea.Batch(
		api.FetchPodDomains(m.pod.ID),
		api.FetchPodEnvVars(m.pod.ID),
	)
}

func (m podDetail) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.PodDomainsLoaded:
		m.domains = tmsg.Domains
		return m, nil

	case msg.PodEnvVarsLoaded:
		m.envVarCount = len(tmsg.EnvVars)
		return m, nil

	case msg.PodDeployed:
		return m, api.LoadData()

	case msg.PodStopped:
		return m, tea.Batch(
			api.LoadData(),
			func() tea.Msg { return msg.ShowStatus{Text: "Pod stopped", Type: msg.StatusSuccess} },
		)

	case msg.PodRestarted:
		return m, tea.Batch(
			api.LoadData(),
			func() tea.Msg { return msg.ShowStatus{Text: "Pod restarted", Type: msg.StatusSuccess} },
		)

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

func (m podDetail) handleKeyPress(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		projectID := m.project.ID
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewProjectDetail(s, projectID)
				},
			}
		}

	case key.Matches(tmsg, m.keyDeploy):
		podID := m.pod.ID
		return m, tea.Batch(
			func() tea.Msg { return msg.StartLoading{Text: "Deploying"} },
			api.DeployPod(podID),
			func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewPodLogs(s, podID)
					},
				}
			},
		)

	case key.Matches(tmsg, m.keyStop):
		return m, tea.Batch(
			func() tea.Msg { return msg.StartLoading{Text: "Stopping"} },
			api.StopPod(m.pod.ID),
		)

	case key.Matches(tmsg, m.keyRestart):
		return m, tea.Batch(
			func() tea.Msg { return msg.StartLoading{Text: "Restarting"} },
			api.RestartPod(m.pod.ID),
		)

	case key.Matches(tmsg, m.keyLogs):
		podID := m.pod.ID
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodLogs(s, podID)
				},
			}
		}

	case key.Matches(tmsg, m.keyEdit):
		pod := m.pod
		projectID := m.project.ID
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodForm(projectID, pod)
				},
			}
		}

	case key.Matches(tmsg, m.keyDomains):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDomains(pod, project)
				},
			}
		}

	case key.Matches(tmsg, m.keyVars):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodVars(pod, project)
				},
			}
		}

	case key.Matches(tmsg, m.keyToken):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodToken(pod, project, s.GitTokens())
				},
			}
		}
	}

	return m, nil
}

func (m podDetail) View() tea.View {
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

	// Git Token
	b.WriteString(labelStyle.Render("Git Token"))
	b.WriteString("\n")
	if m.pod.GitTokenID != nil {
		tokenFound := false
		for _, t := range m.store.GitTokens() {
			if t.ID == *m.pod.GitTokenID {
				b.WriteString(fmt.Sprintf("%s [%s]", t.Name, t.Provider))
				tokenFound = true
				break
			}
		}
		if !tokenFound {
			b.WriteString(styles.MutedStyle().Render("(unknown)"))
		}
	} else {
		b.WriteString(styles.MutedStyle().Render("(none)"))
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

	card := styles.Card(styles.CardProps{
		Width:   styles.CardWidthLG,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(b.String())

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m podDetail) renderStatus() string {
	// Prefer live container state over DB status
	status := m.pod.ContainerState
	if status == "" {
		status = m.pod.Status
	}
	if status == "" {
		status = "idle"
	}
	style := lipgloss.NewStyle()
	switch status {
	case "running":
		style = style.Foreground(lipgloss.Color("10"))
	case "failed", "exited", "dead":
		style = style.Foreground(lipgloss.Color("9"))
	case "building", "restarting":
		style = style.Foreground(lipgloss.Color("11"))
	default:
		style = style.Foreground(lipgloss.Color("8"))
	}
	return style.Render("[" + status + "]")
}

func (m podDetail) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title}
}
