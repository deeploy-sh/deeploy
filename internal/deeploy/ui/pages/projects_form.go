package pages

import (
	"encoding/json"
	"log"

	"github.com/deeploy-sh/deeploy/internal/shared/repo"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectFormPage struct {
	titleInput textinput.Model
	project    *repo.ProjectDTO
	width      int
	height     int
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectFormPage(project *repo.ProjectDTO) ProjectFormPage {
	titleInput := textinput.New()
	titleInput.Focus()
	titleInput.Placeholder = "Title"
	if project != nil {
		titleInput.SetValue(project.Title)
	}

	projectFormPage := ProjectFormPage{
		titleInput: titleInput,
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
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(p.titleInput.Value()) > 0 {
				return p, tea.Batch(
					p.Submit,
					func() tea.Msg { return messages.ProjectPopPageMsg{} },
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

func (p ProjectFormPage) View() string {
	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")

	card := components.Card(components.CardProps{
		Width:   p.width / 2,
		Padding: []int{2, 3},
	}).Render(p.titleInput.View())
	view := lipgloss.JoinVertical(0.5, logo, card)

	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)

	return layout
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

	res, err := utils.Request(utils.RequestProps{
		Method: "POST",
		URL:    "/projects",
		Data:   postData,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var project repo.ProjectDTO
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return nil
	}

	return messages.ProjectCreatedMsg(project)
}

func (p ProjectFormPage) UpdateProject() tea.Msg {
	postData := p.project
	postData.Title = p.titleInput.Value()

	res, err := utils.Request(utils.RequestProps{
		Method: "PUT",
		URL:    "/projects",
		Data:   postData,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var project repo.ProjectDTO
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return nil
	}

	return messages.ProjectUpdatedMsg(project)
}
