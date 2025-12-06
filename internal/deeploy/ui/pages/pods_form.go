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
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type PodFormPage struct {
	titleInput textinput.Model
	keys       formKeyMap
	projectID  string
	pod        *repo.Pod
	width      int
	height     int
}

func (p PodFormPage) HelpKeys() help.KeyMap {
	return p.keys
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewPodFormPage(projectID string, pod *repo.Pod) PodFormPage {
	titleInput := textinput.New()
	titleInput.Focus()
	titleInput.Placeholder = "Title"
	if pod != nil {
		titleInput.SetValue(pod.Title)
	}

	podFormPage := PodFormPage{
		titleInput: titleInput,
		keys:       newFormKeyMap(),
		projectID:  projectID,
	}
	if pod != nil {
		podFormPage.pod = pod
	}
	return podFormPage
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p PodFormPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p PodFormPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEscape:
			projectID := p.projectID
			return p, func() tea.Msg {
				return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
			}
		}
		switch {
		case key.Matches(msg, p.keys.Save):
			if len(p.titleInput.Value()) > 0 {
				projectID := p.projectID
				return p, tea.Batch(
					p.Submit,
					func() tea.Msg {
						return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewProjectDetailPage(s, projectID) }}
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

func (p PodFormPage) View() tea.View {
	card := components.Card(components.CardProps{
		Width:   40,
		Padding: []int{2, 3},
	}).Render(p.titleInput.View())

	contentHeight := p.height

	// Card vertikal zentrieren
	centeredCard := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centeredCard)
}

func (p PodFormPage) Breadcrumbs() []string {
	if p.pod != nil {
		return []string{"Projects", "Pods", "Edit"}
	}
	return []string{"Projects", "Pods", "New"}
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func (p PodFormPage) HasFocusedInput() bool {
	return p.titleInput.Focused()
}

func (p PodFormPage) Submit() tea.Msg {
	if p.pod != nil {
		return p.UpdatePod()
	}
	return p.CreatePod()
}

func (p PodFormPage) CreatePod() tea.Msg {
	postData := struct {
		Title     string `json:"title"`
		ProjectID string `json:"project_id"`
	}{
		Title:     p.titleInput.Value(),
		ProjectID: p.projectID,
	}

	res, err := utils.Request("POST", "/pods", postData)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var pod repo.Pod
	err = json.NewDecoder(res.Body).Decode(&pod)
	if err != nil {
		return nil
	}

	return messages.PodCreatedMsg(pod)
}

func (p PodFormPage) UpdatePod() tea.Msg {
	postData := p.pod
	postData.Title = p.titleInput.Value()

	res, err := utils.Request("PUT", "/pods", postData)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var pod repo.Pod
	err = json.NewDecoder(res.Body).Decode(&pod)
	if err != nil {
		return nil
	}

	return messages.PodUpdatedMsg(pod)
}
