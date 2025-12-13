package pages

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/server/repo"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type PodEditPage struct {
	pod             *repo.Pod
	project         *repo.Project
	gitTokens       []api.GitToken
	titleInput      textinput.Model
	descInput       textinput.Model
	repoURLInput    textinput.Model
	branchInput     textinput.Model
	dockerfileInput textinput.Model
	selectedToken   int
	focusedField    int
	loading         bool
	keySave         key.Binding
	keyBack         key.Binding
	keyTab          key.Binding
	width           int
	height          int
}

const (
	fieldTitle = iota
	fieldDesc
	fieldRepoURL
	fieldBranch
	fieldDockerfile
	fieldGitToken
)

func (m PodEditPage) HelpKeys() []key.Binding {
	return []key.Binding{m.keySave, m.keyTab, m.keyBack}
}

func NewPodEditPage(pod *repo.Pod, project *repo.Project) PodEditPage {
	card := styles.CardProps{Width: 70, Padding: []int{1, 2}, Accent: true}
	inputWidth := card.InnerWidth()

	titleInput := components.NewTextInput(inputWidth)
	titleInput.Placeholder = "My Pod"
	titleInput.SetValue(pod.Title)
	titleInput.Focus()

	descInput := components.NewTextInput(inputWidth)
	descInput.Placeholder = "Description (optional)"
	descInput.SetValue(pod.Description)

	repoInput := components.NewTextInput(inputWidth)
	repoInput.Placeholder = "https://github.com/user/repo"
	if pod.RepoURL != nil {
		repoInput.SetValue(*pod.RepoURL)
	}

	branchInput := components.NewTextInput(inputWidth)
	branchInput.Placeholder = "main"
	if pod.Branch != "" {
		branchInput.SetValue(pod.Branch)
	} else {
		branchInput.SetValue("main")
	}

	dockerfileInput := components.NewTextInput(inputWidth)
	dockerfileInput.Placeholder = "Dockerfile"
	if pod.DockerfilePath != "" {
		dockerfileInput.SetValue(pod.DockerfilePath)
	} else {
		dockerfileInput.SetValue("Dockerfile")
	}

	return PodEditPage{
		pod:             pod,
		project:         project,
		titleInput:      titleInput,
		descInput:       descInput,
		repoURLInput:    repoInput,
		branchInput:     branchInput,
		dockerfileInput: dockerfileInput,
		selectedToken:   0,
		focusedField:    fieldTitle,
		keySave:         key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyBack:         key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		keyTab:          key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
	}
}

func (m PodEditPage) Init() tea.Cmd {
	return tea.Batch(api.FetchGitTokens(), textinput.Blink)
}

func (m PodEditPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.GitTokensLoaded:
		if tokens, ok := tmsg.Tokens.([]api.GitToken); ok {
			m.gitTokens = tokens
			// Find current token selection
			if m.pod.GitTokenID != nil {
				for i, t := range m.gitTokens {
					if t.ID == *m.pod.GitTokenID {
						m.selectedToken = i + 1
						break
					}
				}
			}
		}
		return m, nil

	case msg.PodUpdated:
		m.loading = false
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDetailPage(pod, project)
				},
			}
		}

	case tea.KeyPressMsg:
		return m.handleKeyPress(tmsg)

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	// Update focused input for blink messages
	var cmd tea.Cmd
	switch m.focusedField {
	case fieldTitle:
		m.titleInput, cmd = m.titleInput.Update(tmsg)
	case fieldDesc:
		m.descInput, cmd = m.descInput.Update(tmsg)
	case fieldRepoURL:
		m.repoURLInput, cmd = m.repoURLInput.Update(tmsg)
	case fieldBranch:
		m.branchInput, cmd = m.branchInput.Update(tmsg)
	case fieldDockerfile:
		m.dockerfileInput, cmd = m.dockerfileInput.Update(tmsg)
	}
	return m, cmd
}

func (m *PodEditPage) handleKeyPress(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDetailPage(pod, project)
				},
			}
		}

	case key.Matches(tmsg, m.keySave):
		return m.save()

	case key.Matches(tmsg, m.keyTab):
		m.focusedField = (m.focusedField + 1) % 6
		m.blurAll()
		var cmd tea.Cmd
		switch m.focusedField {
		case fieldTitle:
			cmd = m.titleInput.Focus()
		case fieldDesc:
			cmd = m.descInput.Focus()
		case fieldRepoURL:
			cmd = m.repoURLInput.Focus()
		case fieldBranch:
			cmd = m.branchInput.Focus()
		case fieldDockerfile:
			cmd = m.dockerfileInput.Focus()
		}
		return m, cmd

	case tmsg.Code == tea.KeyUp:
		if m.focusedField == fieldGitToken && m.selectedToken > 0 {
			m.selectedToken--
		}
		return m, nil

	case tmsg.Code == tea.KeyDown:
		if m.focusedField == fieldGitToken && m.selectedToken < len(m.gitTokens) {
			m.selectedToken++
		}
		return m, nil
	}

	// Update focused input
	var cmd tea.Cmd
	switch m.focusedField {
	case fieldTitle:
		m.titleInput, cmd = m.titleInput.Update(tmsg)
	case fieldDesc:
		m.descInput, cmd = m.descInput.Update(tmsg)
	case fieldRepoURL:
		m.repoURLInput, cmd = m.repoURLInput.Update(tmsg)
	case fieldBranch:
		m.branchInput, cmd = m.branchInput.Update(tmsg)
	case fieldDockerfile:
		m.dockerfileInput, cmd = m.dockerfileInput.Update(tmsg)
	}
	return m, cmd
}

func (m *PodEditPage) blurAll() {
	m.titleInput.Blur()
	m.descInput.Blur()
	m.repoURLInput.Blur()
	m.branchInput.Blur()
	m.dockerfileInput.Blur()
}

func (m *PodEditPage) save() (tea.Model, tea.Cmd) {
	title := strings.TrimSpace(m.titleInput.Value())
	if title == "" {
		return m, nil
	}

	m.pod.Title = title

	m.pod.Description = strings.TrimSpace(m.descInput.Value())

	repoURL := strings.TrimSpace(m.repoURLInput.Value())
	if repoURL != "" {
		m.pod.RepoURL = &repoURL
	} else {
		m.pod.RepoURL = nil
	}

	m.pod.Branch = m.branchInput.Value()
	if m.pod.Branch == "" {
		m.pod.Branch = "main"
	}

	m.pod.DockerfilePath = m.dockerfileInput.Value()
	if m.pod.DockerfilePath == "" {
		m.pod.DockerfilePath = "Dockerfile"
	}

	if m.selectedToken > 0 && m.selectedToken <= len(m.gitTokens) {
		tokenID := m.gitTokens[m.selectedToken-1].ID
		m.pod.GitTokenID = &tokenID
	} else {
		m.pod.GitTokenID = nil
	}

	m.loading = true
	return m, api.UpdatePod(m.pod)
}

func (m PodEditPage) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Edit Pod"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorMuted())
	activeLabel := lipgloss.NewStyle().Foreground(styles.ColorPrimary())

	// Title
	if m.focusedField == fieldTitle {
		b.WriteString(activeLabel.Render("Title"))
	} else {
		b.WriteString(labelStyle.Render("Title"))
	}
	b.WriteString("\n")
	b.WriteString(m.titleInput.View())
	b.WriteString("\n\n")

	// Description
	if m.focusedField == fieldDesc {
		b.WriteString(activeLabel.Render("Description"))
	} else {
		b.WriteString(labelStyle.Render("Description"))
	}
	b.WriteString("\n")
	b.WriteString(m.descInput.View())
	b.WriteString("\n\n")

	// Repo URL
	if m.focusedField == fieldRepoURL {
		b.WriteString(activeLabel.Render("Repo URL"))
	} else {
		b.WriteString(labelStyle.Render("Repo URL"))
	}
	b.WriteString("\n")
	b.WriteString(m.repoURLInput.View())
	b.WriteString("\n\n")

	// Branch
	if m.focusedField == fieldBranch {
		b.WriteString(activeLabel.Render("Branch"))
	} else {
		b.WriteString(labelStyle.Render("Branch"))
	}
	b.WriteString("\n")
	b.WriteString(m.branchInput.View())
	b.WriteString("\n\n")

	// Dockerfile
	if m.focusedField == fieldDockerfile {
		b.WriteString(activeLabel.Render("Dockerfile"))
	} else {
		b.WriteString(labelStyle.Render("Dockerfile"))
	}
	b.WriteString("\n")
	b.WriteString(m.dockerfileInput.View())
	b.WriteString("\n\n")

	// Git Token
	if m.focusedField == fieldGitToken {
		b.WriteString(activeLabel.Render("Git Token"))
	} else {
		b.WriteString(labelStyle.Render("Git Token"))
	}
	b.WriteString("\n")

	cursor := "  "
	if m.focusedField == fieldGitToken && m.selectedToken == 0 {
		cursor = "> "
	}
	b.WriteString(fmt.Sprintf("%s(none - public repo)\n", cursor))

	for i, t := range m.gitTokens {
		cursor = "  "
		if m.focusedField == fieldGitToken && m.selectedToken == i+1 {
			cursor = "> "
		}
		b.WriteString(fmt.Sprintf("%s%s [%s]\n", cursor, t.Name, t.Provider))
	}

	if m.loading {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render("Saving..."))
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

func (m PodEditPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title, "Edit"}
}
