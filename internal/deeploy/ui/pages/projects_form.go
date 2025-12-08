package pages

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/api"
	"github.com/deeploy-sh/deeploy/internal/deeploy/msg"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

type formKeyMap struct {
	Save   key.Binding
	Cancel key.Binding
}

func (k formKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Cancel}
}

func (k formKeyMap) FullHelp() [][]key.Binding {
	return nil
}

func (p ProjectFormPage) HelpKeys() help.KeyMap {
	return p.keys
}

func newFormKeyMap() formKeyMap {
	return formKeyMap{
		Save:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
		Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

type ProjectFormPage struct {
	titleInput textinput.Model
	keys       formKeyMap
	project    *repo.Project
	width      int
	height     int
}

func NewProjectFormPage(project *repo.Project) ProjectFormPage {
	card := components.CardProps{Width: 40, Padding: []int{1, 2}, Accent: true}
	titleInput := components.NewTextInput(card.InnerWidth())
	titleInput.Focus()
	titleInput.Placeholder = "Title"
	if project != nil {
		titleInput.SetValue(project.Title)
	}

	projectFormPage := ProjectFormPage{
		titleInput: titleInput,
		keys:       newFormKeyMap(),
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
		case key.Matches(tmsg, p.keys.Save):
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

	card := components.Card(components.CardProps{
		Width:   40,
		Padding: []int{1, 2},
		Accent:  true,
	}).Render(content)

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
