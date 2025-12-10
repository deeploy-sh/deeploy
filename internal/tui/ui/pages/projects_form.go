package pages

import (
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

type ProjectFormPage struct {
	titleInput textinput.Model
	keySave    key.Binding
	keyCancel  key.Binding
	project    *repo.Project
	width      int
	height     int
}

func (p ProjectFormPage) HelpKeys() []key.Binding {
	return []key.Binding{p.keySave, p.keyCancel}
}

func NewProjectFormPage(project *repo.Project) ProjectFormPage {
	titleInput := components.NewTextInput(styles.CardInner(styles.CardSmall))
	titleInput.Focus()
	titleInput.Placeholder = "Title"
	if project != nil {
		titleInput.SetValue(project.Title)
	}

	projectFormPage := ProjectFormPage{
		titleInput: titleInput,
		keySave:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
		keyCancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
	if project != nil {
		projectFormPage.project = project
	}
	return projectFormPage
}

func (p ProjectFormPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p ProjectFormPage) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEscape:
			return p, func() tea.Msg {
				return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
			}
		}
		switch {
		case key.Matches(tmsg, p.keySave):
			if len(p.titleInput.Value()) > 0 {
				var apiCmd tea.Cmd
				if p.project != nil {
					p.project.Title = p.titleInput.Value()
					apiCmd = api.UpdateProject(p.project)
				} else {
					apiCmd = api.CreateProject(p.titleInput.Value())
				}
				return p, tea.Batch(
					apiCmd,
					func() tea.Msg {
						return msg.ChangePage{PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) }}
					},
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

func (p ProjectFormPage) View() tea.View {
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

	card := styles.Card(styles.CardSmall, true).Render(content)

	centered := lipgloss.Place(p.width, p.height,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centered)
}

func (p ProjectFormPage) Breadcrumbs() []string {
	if p.project != nil {
		return []string{"Projects", "Edit"}
	}
	return []string{"Projects", "New"}
}

func (p ProjectFormPage) HasFocusedInput() bool {
	return p.titleInput.Focused()
}
