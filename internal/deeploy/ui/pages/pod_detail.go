package pages

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/api"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type podDetailMode int

const (
	modeView podDetailMode = iota
	modeEditRepo
)

type PodDetailPage struct {
	pod            *repo.Pod
	project        *repo.Project
	logs           []string
	domains        []api.PodDomain
	gitTokens      []api.GitToken
	loading        bool
	mode           podDetailMode
	repoURLInput   textinput.Model
	branchInput    textinput.Model
	dockerfileInput textinput.Model
	selectedToken  int
	focusedField   int
	keyDeploy      key.Binding
	keyStop        key.Binding
	keyRestart     key.Binding
	keyLogs        key.Binding
	keyEdit        key.Binding
	keyDomains     key.Binding
	keySave        key.Binding
	keyBack        key.Binding
	keyTab         key.Binding
	width          int
	height         int
}

func (m PodDetailPage) HelpKeys() []key.Binding {
	if m.mode == modeEditRepo {
		return []key.Binding{m.keySave, m.keyTab, m.keyBack}
	}
	return []key.Binding{m.keyDeploy, m.keyStop, m.keyRestart, m.keyLogs, m.keyEdit, m.keyDomains, m.keyBack}
}

func NewPodDetailPage(pod *repo.Pod, project *repo.Project) PodDetailPage {
	repoInput := components.NewTextInput(50)
	repoInput.Placeholder = "https://github.com/user/repo"
	if pod.RepoURL != nil {
		repoInput.SetValue(*pod.RepoURL)
	}

	branchInput := components.NewTextInput(30)
	branchInput.Placeholder = "main"
	if pod.Branch != "" {
		branchInput.SetValue(pod.Branch)
	} else {
		branchInput.SetValue("main")
	}

	dockerfileInput := components.NewTextInput(30)
	dockerfileInput.Placeholder = "Dockerfile"
	if pod.DockerfilePath != "" {
		dockerfileInput.SetValue(pod.DockerfilePath)
	} else {
		dockerfileInput.SetValue("Dockerfile")
	}

	return PodDetailPage{
		pod:             pod,
		project:         project,
		repoURLInput:    repoInput,
		branchInput:     branchInput,
		dockerfileInput: dockerfileInput,
		selectedToken:   0,
		keyDeploy:       key.NewBinding(key.WithKeys("D"), key.WithHelp("D", "deploy")),
		keyStop:         key.NewBinding(key.WithKeys("S"), key.WithHelp("S", "stop")),
		keyRestart:      key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "restart")),
		keyLogs:         key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "refresh logs")),
		keyEdit:         key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit config")),
		keyDomains:      key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "domains")),
		keySave:         key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyBack:         key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		keyTab:          key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	}
}

func (m PodDetailPage) Init() tea.Cmd {
	return tea.Batch(
		api.FetchPodLogs(m.pod.ID),
		api.FetchPodDomains(m.pod.ID),
		api.FetchGitTokens(),
		textinput.Blink,
	)
}

func (m PodDetailPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.GitTokensLoaded:
		if tokens, ok := tmsg.Tokens.([]api.GitToken); ok {
			m.gitTokens = tokens
		}
		return m, nil

	case msg.PodLogsLoaded:
		m.logs = tmsg.Logs
		m.loading = false
		return m, nil

	case msg.PodDomainsLoaded:
		if domains, ok := tmsg.Domains.([]api.PodDomain); ok {
			m.domains = domains
		}
		return m, nil

	case msg.PodDeployed, msg.PodStopped, msg.PodRestarted:
		m.loading = false
		return m, tea.Batch(api.LoadData(), api.FetchPodLogs(m.pod.ID))

	case msg.PodUpdated:
		m.mode = modeView
		return m, api.LoadData()

	case msg.DataLoaded:
		// Find updated pod
		for _, p := range tmsg.Pods {
			if p.ID == m.pod.ID {
				pod := p
				m.pod = &pod
				break
			}
		}
		return m, nil

	case tea.KeyPressMsg:
		if m.mode == modeEditRepo {
			return m.handleEditMode(tmsg)
		}
		return m.handleViewMode(tmsg)

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	return m, nil
}

func (m PodDetailPage) handleViewMode(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
						return NewLogViewPage(podID, podTitle)
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
		m.loading = true
		return m, api.FetchPodLogs(m.pod.ID)
	case key.Matches(tmsg, m.keyEdit):
		m.mode = modeEditRepo
		m.focusedField = 0
		m.repoURLInput.Focus()
		return m, textinput.Blink
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
	}
	return m, nil
}

func (m PodDetailPage) handleEditMode(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		m.mode = modeView
		m.repoURLInput.Blur()
		m.branchInput.Blur()
		m.dockerfileInput.Blur()
		return m, nil

	case key.Matches(tmsg, m.keySave):
		repoURL := m.repoURLInput.Value()
		m.pod.RepoURL = &repoURL
		m.pod.Branch = m.branchInput.Value()
		if m.pod.Branch == "" {
			m.pod.Branch = "main"
		}
		m.pod.DockerfilePath = m.dockerfileInput.Value()
		if m.pod.DockerfilePath == "" {
			m.pod.DockerfilePath = "Dockerfile"
		}
		// Set git token if selected
		if m.selectedToken > 0 && m.selectedToken <= len(m.gitTokens) {
			tokenID := m.gitTokens[m.selectedToken-1].ID
			m.pod.GitTokenID = &tokenID
		} else {
			m.pod.GitTokenID = nil
		}
		return m, api.UpdatePod(m.pod)

	case key.Matches(tmsg, m.keyTab):
		m.focusedField = (m.focusedField + 1) % 4
		m.repoURLInput.Blur()
		m.branchInput.Blur()
		m.dockerfileInput.Blur()
		switch m.focusedField {
		case 0:
			m.repoURLInput.Focus()
		case 1:
			m.branchInput.Focus()
		case 2:
			m.dockerfileInput.Focus()
		}
		return m, nil

	case tmsg.Code == tea.KeyUp:
		if m.focusedField == 3 && m.selectedToken > 0 {
			m.selectedToken--
		}
		return m, nil

	case tmsg.Code == tea.KeyDown:
		if m.focusedField == 3 && m.selectedToken < len(m.gitTokens) {
			m.selectedToken++
		}
		return m, nil
	}

	var cmd tea.Cmd
	switch m.focusedField {
	case 0:
		m.repoURLInput, cmd = m.repoURLInput.Update(tmsg)
	case 1:
		m.branchInput, cmd = m.branchInput.Update(tmsg)
	case 2:
		m.dockerfileInput, cmd = m.dockerfileInput.Update(tmsg)
	}
	return m, cmd
}

func (m PodDetailPage) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render(m.pod.Title))
	b.WriteString("  ")
	b.WriteString(m.renderStatus())
	b.WriteString("\n\n")

	if m.mode == modeEditRepo {
		b.WriteString(m.renderEditMode())
	} else {
		b.WriteString(m.renderViewMode())
	}

	card := components.Card(components.CardProps{
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

func (m PodDetailPage) renderEditMode() string {
	var b strings.Builder
	labelStyle := lipgloss.NewStyle().Width(12).Foreground(styles.ColorMuted())
	activeLabel := lipgloss.NewStyle().Width(12).Foreground(styles.ColorPrimary())

	b.WriteString("Configure Repository\n\n")

	// Repo URL
	if m.focusedField == 0 {
		b.WriteString(activeLabel.Render("Repo URL:"))
	} else {
		b.WriteString(labelStyle.Render("Repo URL:"))
	}
	b.WriteString(m.repoURLInput.View())
	b.WriteString("\n\n")

	// Branch
	if m.focusedField == 1 {
		b.WriteString(activeLabel.Render("Branch:"))
	} else {
		b.WriteString(labelStyle.Render("Branch:"))
	}
	b.WriteString(m.branchInput.View())
	b.WriteString("\n\n")

	// Dockerfile
	if m.focusedField == 2 {
		b.WriteString(activeLabel.Render("Dockerfile:"))
	} else {
		b.WriteString(labelStyle.Render("Dockerfile:"))
	}
	b.WriteString(m.dockerfileInput.View())
	b.WriteString("\n\n")

	// Git Token
	if m.focusedField == 3 {
		b.WriteString(activeLabel.Render("Git Token:"))
	} else {
		b.WriteString(labelStyle.Render("Git Token:"))
	}
	b.WriteString("\n")

	// Token selection
	cursor := "  "
	if m.focusedField == 3 && m.selectedToken == 0 {
		cursor = "> "
	}
	b.WriteString(fmt.Sprintf("%s(none - public repo)\n", cursor))

	for i, t := range m.gitTokens {
		cursor = "  "
		if m.focusedField == 3 && m.selectedToken == i+1 {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s [%s]\n", cursor, t.Name, t.Provider))
	}

	return b.String()
}

func (m PodDetailPage) renderViewMode() string {
	var b strings.Builder
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorMuted())

	// Repo Info
	b.WriteString(labelStyle.Render("Repo: "))
	if m.pod.RepoURL != nil && *m.pod.RepoURL != "" {
		b.WriteString(*m.pod.RepoURL)
		b.WriteString(" @ ")
		b.WriteString(m.pod.Branch)
	} else {
		b.WriteString(styles.MutedStyle().Render("(not configured - press 'e' to edit)"))
	}
	b.WriteString("\n")

	// Dockerfile
	b.WriteString(labelStyle.Render("Dockerfile: "))
	if m.pod.DockerfilePath != "" {
		b.WriteString(m.pod.DockerfilePath)
	} else {
		b.WriteString("Dockerfile")
	}
	b.WriteString("\n")

	// Domains
	b.WriteString(labelStyle.Render("Domains: "))
	if len(m.domains) > 0 {
		for i, d := range m.domains {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(d.Domain)
		}
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("(none - press 'o' to add)"))
	}
	b.WriteString("\n\n")

	// Logs
	b.WriteString(labelStyle.Render("Recent Logs:"))
	b.WriteString("\n")
	if len(m.logs) > 0 {
		maxLogs := 8
		start := 0
		if len(m.logs) > maxLogs {
			start = len(m.logs) - maxLogs
		}
		for _, line := range m.logs[start:] {
			b.WriteString(styles.MutedStyle().Render("  " + line))
			b.WriteString("\n")
		}
	} else {
		b.WriteString(styles.MutedStyle().Render("  (no logs yet)"))
		b.WriteString("\n")
	}

	if m.loading {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Loading..."))
	}

	return b.String()
}

func (m PodDetailPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title}
}
