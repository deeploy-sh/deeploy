package page

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type projectToDelete struct {
	ID    string
	Title string
}

type projectDelete struct {
	project    projectToDelete
	podCount   int
	input      textinput.Model
	keyConfirm key.Binding
	keyCancel  key.Binding
	width      int
	height     int
}

func (p projectDelete) HelpKeys() []key.Binding {
	return []key.Binding{p.keyConfirm, p.keyCancel}
}

func NewProjectDelete(s msg.Store, project *model.Project) projectDelete {
	podCount := 0
	for _, p := range s.Pods() {
		if p.ProjectID == project.ID {
			podCount++
		}
	}

	card := styles.CardProps{Width: 40, Padding: []int{1, 2}, Accent: true}
	ti := components.NewTextInput(card.InnerWidth())
	ti.Placeholder = project.Title
	ti.Focus()
	ti.CharLimit = 100

	return projectDelete{
		project:    projectToDelete{ID: project.ID, Title: project.Title},
		podCount:   podCount,
		input:      ti,
		keyConfirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

func (p projectDelete) Init() tea.Cmd {
	return textinput.Blink
}

func (p projectDelete) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case msg.ProjectDeleted:
		return p, tea.Batch(
			api.LoadData(),
			func() tea.Msg { return msg.ShowStatus{Text: "Project deleted", Type: msg.StatusSuccess} },
			func() tea.Msg { return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }} },
		)

	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
			}
		case tea.KeyEnter:
			if p.podCount > 0 {
				return p, nil
			}
			// Only delete if input matches project title exactly
			if p.input.Value() != p.project.Title {
				return p, nil
			}
			return p, api.DeleteProject(p.project.ID)
		}
	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}

	p.input, cmd = p.input.Update(tmsg)
	return p, cmd
}

func (p projectDelete) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Render(fmt.Sprintf("Delete Project (%v)", p.project.Title))

	var hint string
	if p.podCount > 0 {
		hint = styles.MutedStyle().
			PaddingTop(1).
			PaddingBottom(1).
			Render(fmt.Sprintf("Delete all %d pods first", p.podCount))
	} else {
		hint = styles.MutedStyle().
			PaddingTop(1).
			PaddingBottom(1).
			Render("Type '" + p.project.Title + "' to confirm")
	}

	content := lipgloss.JoinVertical(lipgloss.Left, title, hint, p.input.View())

	card := styles.Card(styles.CardProps{
		Width:   40,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p projectDelete) Breadcrumbs() []string {
	return []string{"Projects", "Delete"}
}
