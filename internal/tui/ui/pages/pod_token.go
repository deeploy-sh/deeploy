package pages

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

type PodTokenPage struct {
	pod       *model.Pod
	project   *model.Project
	gitTokens []model.GitToken
	selected  int // 0 = none, 1+ = token index
	keySelect key.Binding
	keyBack   key.Binding
	width     int
	height    int
}

func (m PodTokenPage) HelpKeys() []key.Binding {
	return []key.Binding{m.keySelect, m.keyBack}
}

func NewPodTokenPage(pod *model.Pod, project *model.Project, gitTokens []model.GitToken) PodTokenPage {
	// Find current selection
	selected := 0
	if pod.GitTokenID != nil {
		for i, t := range gitTokens {
			if t.ID == *pod.GitTokenID {
				selected = i + 1
				break
			}
		}
	}

	return PodTokenPage{
		pod:       pod,
		project:   project,
		gitTokens: gitTokens,
		selected:  selected,
		keySelect: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		keyBack:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}
}

func (m PodTokenPage) Init() tea.Cmd {
	return nil
}

func (m PodTokenPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.PodUpdated:
		pod := m.pod
		project := m.project
		return m, tea.Batch(
			api.LoadData(),
			func() tea.Msg { return msg.ShowStatus{Text: "Token updated", Type: msg.StatusSuccess} },
			func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model {
					return NewPodDetailPage(pod, project, s.GitTokens())
				}}
			},
		)

	case tea.KeyPressMsg:
		return m.handleKeyPress(tmsg)

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height
		return m, nil
	}

	return m, nil
}

func (m *PodTokenPage) handleKeyPress(tmsg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(tmsg, m.keyBack):
		pod := m.pod
		project := m.project
		return m, func() tea.Msg {
			return msg.ChangePage{
				PageFactory: func(s msg.Store) tea.Model {
					return NewPodDetailPage(pod, project, s.GitTokens())
				},
			}
		}

	case key.Matches(tmsg, m.keySelect):
		return m.selectToken()

	case tmsg.Code == tea.KeyUp:
		if m.selected > 0 {
			m.selected--
		}
		return m, nil

	case tmsg.Code == tea.KeyDown:
		if m.selected < len(m.gitTokens) {
			m.selected++
		}
		return m, nil
	}

	return m, nil
}

func (m *PodTokenPage) selectToken() (tea.Model, tea.Cmd) {
	if m.selected == 0 {
		m.pod.GitTokenID = nil
	} else if m.selected <= len(m.gitTokens) {
		tokenID := m.gitTokens[m.selected-1].ID
		m.pod.GitTokenID = &tokenID
	}

	return m, api.UpdatePod(m.pod)
}

func (m PodTokenPage) View() tea.View {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary())
	b.WriteString(titleStyle.Render("Select Git Token"))
	b.WriteString("\n\n")

	cursorStyle := lipgloss.NewStyle().Foreground(styles.ColorPrimary())
	selectedStyle := lipgloss.NewStyle().Foreground(styles.ColorPrimary())
	normalStyle := lipgloss.NewStyle()

	// None option
	cursor := "  "
	style := normalStyle
	if m.selected == 0 {
		cursor = "> "
		style = selectedStyle
	}
	b.WriteString(cursorStyle.Render(cursor))
	b.WriteString(style.Render("(none - public repo)"))
	b.WriteString("\n")

	// Token options
	for i, t := range m.gitTokens {
		cursor = "  "
		style = normalStyle
		if m.selected == i+1 {
			cursor = "> "
			style = selectedStyle
		}
		b.WriteString(cursorStyle.Render(cursor))
		b.WriteString(style.Render(fmt.Sprintf("%s [%s]", t.Name, t.Provider)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if len(m.gitTokens) == 0 {
		b.WriteString(styles.MutedStyle().Render("No git tokens configured."))
		b.WriteString("\n")
	}
	b.WriteString(styles.MutedStyle().Render("Manage tokens: Alt+P > Git Tokens"))

	card := styles.Card(styles.CardProps{
		Width:   50,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(b.String())

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (m PodTokenPage) Breadcrumbs() []string {
	return []string{"Projects", m.project.Title, "Pods", m.pod.Title, "Token"}
}
