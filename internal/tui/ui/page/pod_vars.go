package page

import (
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type podVars struct {
	pod      *model.Pod
	project  *model.Project
	textarea textarea.Model
	envVars  []model.PodEnvVar
	keySave  key.Binding
	keyBack  key.Binding
	width    int
	height   int
}

func (m podVars) HelpKeys() []key.Binding {
	return []key.Binding{m.keySave, m.keyBack}
}

func NewPodVars(pod *model.Pod, project *model.Project) podVars {
	ta := textarea.New()
	ta.Placeholder = "DATABASE_URL=postgres://..."
	ta.Prompt = ""
	ta.SetWidth(60)
	ta.SetHeight(10)
	ta.Focus()

	return podVars{
		pod:      pod,
		project:  project,
		textarea: ta,
		keySave:  key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyBack:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (m podVars) Init() tea.Cmd {
	return tea.Batch(api.FetchPodEnvVars(m.pod.ID), textarea.Blink)
}

func (m podVars) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.PodEnvVarsLoaded:
		m.envVars = tmsg.EnvVars
		m.textarea.SetValue(m.envVarsToText())
		return m, nil

	case msg.PodEnvVarsUpdated:
		podID := m.pod.ID
		return m, tea.Batch(
			func() tea.Msg { return msg.ShowStatus{Text: "Saved. Restart or deploy to apply.", Type: msg.StatusSuccess} },
			func() tea.Msg { return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewPodDetail(s, podID) }} },
		)

	case tea.KeyPressMsg:
		if key.Matches(tmsg, m.keyBack) {
			podID := m.pod.ID
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewPodDetail(s, podID)
					},
				}
			}
		}

		if key.Matches(tmsg, m.keySave) {
			return m.save()
		}

		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(tmsg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(tmsg)
	return m, cmd
}

func (m podVars) envVarsToText() string {
	var lines []string
	for _, v := range m.envVars {
		lines = append(lines, v.Key+"="+v.Value)
	}
	return strings.Join(lines, "\n")
}

func (m podVars) textToEnvVars() []model.PodEnvVar {
	var vars []model.PodEnvVar
	lines := strings.Split(m.textarea.Value(), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			continue
		}

		vars = append(vars, model.PodEnvVar{
			Key:   key,
			Value: value,
		})
	}

	return vars
}

func (m *podVars) save() (tea.Model, tea.Cmd) {
	vars := m.textToEnvVars()
	return m, tea.Batch(
		func() tea.Msg { return msg.StartLoading{Text: "Saving"} },
		api.UpdatePodEnvVars(m.pod.ID, vars),
	)
}

func (m podVars) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Environment Variables"))
	b.WriteString("\n")
	b.WriteString(styles.MutedStyle().Render("One KEY=value per line."))
	b.WriteString("\n\n")

	b.WriteString(m.textarea.View())

	b.WriteString("\n\n")
	b.WriteString(styles.MutedStyle().Render("Values are encrypted at rest."))

	card := styles.Card(styles.CardProps{
		Width:   styles.CardWidthLG,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(b.String())

	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m podVars) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title, "Env Vars"}
}
