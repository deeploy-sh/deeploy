package pages

import (
	"encoding/json"
	"log"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

// /////////////////////////////////////////////////////////////////////////////
// KeyMap
// /////////////////////////////////////////////////////////////////////////////

type formKeyMap struct {
	Save   key.Binding
	Cancel key.Binding
}

func (k formKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Cancel}
}

func (k formKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Save, k.Cancel}}
}

func newFormKeyMap() formKeyMap {
	return formKeyMap{
		Save:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "save")),
		Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectFormPage struct {
	titleInput textinput.Model
	keys       formKeyMap
	help       help.Model
	project    *repo.Project
	width      int
	height     int
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectFormPage(project *repo.Project) ProjectFormPage {
	titleInput := textinput.New()
	titleInput.Focus()
	titleInput.Placeholder = "Title"
	if project != nil {
		titleInput.SetValue(project.Title)
	}

	projectFormPage := ProjectFormPage{
		titleInput: titleInput,
		keys:       newFormKeyMap(),
		help:       styles.NewHelpModel(),
	}
	if project != nil {
		projectFormPage.project = project
	}
	return projectFormPage
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectFormPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p ProjectFormPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEscape:
			return p, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewDashboard()}
			}
		}
		switch {
		case key.Matches(msg, p.keys.Save):
			if len(p.titleInput.Value()) > 0 {
				return p, tea.Batch(
					p.Submit,
					func() tea.Msg {
						return messages.ChangePageMsg{Page: NewDashboard()}
					},
				)
			}
		}

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	}

	p.titleInput, cmd = p.titleInput.Update(msg)
	return p, cmd
}

func (p ProjectFormPage) View() tea.View {
	card := components.Card(components.CardProps{
		Width:   40,
		Padding: []int{2, 3},
	}).Render(p.titleInput.View())

	helpView := p.help.View(p.keys)
	contentHeight := p.height - 1 // 1 für help

	// Card vertikal zentrieren
	centeredCard := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, centeredCard, helpView))
}

func (p ProjectFormPage) Breadcrumbs() []string {
	if p.project != nil {
		return []string{"Projects", "Edit"}
	}
	return []string{"Projects", "New"}
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectFormPage) HasFocusedInput() bool {
	return p.titleInput.Focused()
}

func (p ProjectFormPage) Submit() tea.Msg {
	if p.project != nil {
		return p.UpdateProject()
	}
	return p.CreateProject()
}

func (p ProjectFormPage) CreateProject() tea.Msg {
	postData := struct {
		Title string
	}{
		Title: p.titleInput.Value(),
	}

	res, err := utils.Request("POST", "/projects", postData)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var project repo.Project
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return nil
	}

	return messages.ProjectCreatedMsg(project)
}

func (p ProjectFormPage) UpdateProject() tea.Msg {
	postData := p.project
	postData.Title = p.titleInput.Value()

	res, err := utils.Request("PUT", "/projects", postData)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var project repo.Project
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return nil
	}

	return messages.ProjectUpdatedMsg(project)
}
