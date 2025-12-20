package page

import (
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

type projectForm struct {
	titleInput textinput.Model
	keySave    key.Binding
	keyCancel  key.Binding
	project    *model.Project
	width      int
	height     int
}

func (p projectForm) HelpKeys() []key.Binding {
	return []key.Binding{p.keySave, p.keyCancel}
}

func NewProjectForm(project *model.Project) projectForm {
	card := styles.CardProps{Width: 40, Padding: []int{1, 2}, Accent: true}
	titleInput := components.NewTextInput(card.InnerWidth())
	titleInput.Focus()
	titleInput.Placeholder = "Title"
	if project != nil {
		titleInput.SetValue(project.Title)
	}

	projectForm := projectForm{
		titleInput: titleInput,
		keySave:    key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
	if project != nil {
		projectForm.project = project
	}
	return projectForm
}

func (p projectForm) Init() tea.Cmd {
	return textinput.Blink
}

func (p projectForm) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case msg.ProjectCreated:
		return p, tea.Batch(
			api.LoadData(),
			func() tea.Msg { return msg.ShowStatus{Text: "Project created", Type: msg.StatusSuccess} },
			func() tea.Msg { return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }} },
		)

	case msg.ProjectUpdated:
		projectID := p.project.ID
		return p, tea.Batch(
			api.LoadData(),
			func() tea.Msg { return msg.ShowStatus{Text: "Project saved", Type: msg.StatusSuccess} },
			func() tea.Msg { return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetail(s, projectID) }} },
		)

	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			if p.project != nil {
				projectID := p.project.ID
				return p, func() tea.Msg {
					return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewProjectDetail(s, projectID) }}
				}
			}
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
			}
		}
		switch {
		case key.Matches(tmsg, p.keySave):
			if len(p.titleInput.Value()) > 0 {
				if p.project != nil {
					p.project.Title = p.titleInput.Value()
					return p, tea.Batch(
						func() tea.Msg { return msg.StartLoading{Text: "Updating project"} },
						api.UpdateProject(p.project),
					)
				}
				return p, tea.Batch(
					func() tea.Msg { return msg.StartLoading{Text: "Creating project"} },
					api.CreateProject(p.titleInput.Value()),
				)
			}
		}

	case tea.WindowSizeMsg:
		p.width = tmsg.Width
		p.height = tmsg.Height
		return p, nil
	}

	p.titleInput, cmd = p.titleInput.Update(tmsg)
	return p, cmd
}

func (p projectForm) View() tea.View {
	var titleText string
	if p.project != nil {
		titleText = "Edit Project"
	} else {
		titleText = "New Project"
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Render(titleText)

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		p.titleInput.View(),
	)

	card := styles.Card(styles.CardProps{
		Width:   40,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p projectForm) Breadcrumbs() []string {
	if p.project != nil {
		return []string{"Projects", "Edit"}
	}
	return []string{"Projects", "New"}
}

func (p projectForm) HasFocusedInput() bool {
	return p.titleInput.Focused()
}
