package pages

import (
	"encoding/json"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectListPage struct {
	stack         []tea.Model
	width         int
	height        int
	projects      []repo.ProjectDTO
	selectedIndex int
	err           error
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectListPage() ProjectListPage {
	return ProjectListPage{}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectListPage) Init() tea.Cmd {
	return getProjects
}

func (p ProjectListPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("FROM PROJECTS LIST: ", msg)
		if msg.Type == tea.KeyEsc {
			if len(p.stack) == 0 {
				return p, func() tea.Msg {
					return messages.ChangePageMsg{Page: NewDashboard()}
				}
			}
			return p, func() tea.Msg {
				return messages.ProjectPopPageMsg{}
			}
		}
		switch msg.String() {
		case "n":
			return p, func() tea.Msg {
				return messages.ProjectPushPageMsg{Page: NewProjectFormPage(nil)}
			}
		case "e":
			return p, func() tea.Msg {
				return messages.ProjectPushPageMsg{Page: NewProjectFormPage(&p.projects[p.selectedIndex])}
			}
		case "d":
			return p, func() tea.Msg {
				return messages.ProjectPushPageMsg{Page: NewProjectDeletePage(&p.projects[p.selectedIndex])}
			}
		case "down", "j":
			if p.selectedIndex == len(p.projects)-1 {
				p.selectedIndex = 0
			} else {
				p.selectedIndex++
			}
		case "up", "k":
			if p.selectedIndex == 0 {
				p.selectedIndex = len(p.projects) - 1
			} else {
				p.selectedIndex--
			}
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	case messages.ProjectErrMsg:
		p.err = msg
	case messages.ProjectsInitDataMsg:
		p.projects = msg
		return p, nil
	case messages.ProjectCreatedMsg:
		p.projects = append(p.projects, repo.ProjectDTO(msg))
		return p, nil
	case messages.ProjectUpdatedMsg:
		project := msg
		for i, v := range p.projects {
			if v.ID == project.ID {
				p.projects[i] = repo.ProjectDTO(project)
				break
			}
		}

	case messages.ProjectDeleteMsg:
		project := msg
		var index int
		for i, v := range p.projects {
			if v.ID == project.ID {
				index = i
				break
			}
		}
		p.projects = append(p.projects[:index], p.projects[index+1:]...)
		p.selectedIndex--
		return p, nil
	}
	return p, nil
}

func (p ProjectListPage) View() string {
	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	log.Println(p.projects)
	var cards []string
	if p.err != nil {
		cards = append(cards, components.ErrorCard(30).Render(p.err.Error()))
	} else {
		for i, project := range p.projects {
			props := components.CardProps{
				Width:   30,
				Padding: []int{0, 1},
			}
			if p.selectedIndex == i {
				props.BorderForeground = styles.ColorPrimary
			}
			cards = append(cards, components.Card(props).Render(project.Title))
		}
	}

	projectsView := lipgloss.JoinVertical(0.5, cards...)

	if len(cards) == 0 {
		projectsView = components.Card(components.CardProps{Width: 30}).Align(lipgloss.Position(0.5)).Render(styles.FocusedStyle.Render("No projects yet"))
	}

	view := lipgloss.JoinVertical(0.5, logo, projectsView)

	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)

	return layout
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func getProjects() tea.Msg {
	config, err := config.LoadConfig()
	if err != nil {
		return messages.ChangePageMsg{Page: NewConnectPage()}
	}

	r, err := http.NewRequest("GET", config.Server+"/api/projects", nil)
	if err != nil {
		return messages.ProjectErrMsg(err)
	}
	r.Header.Set("Authorization", "Bearer "+config.Token)

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return messages.ProjectErrMsg(err)
	}
	if res.StatusCode == http.StatusUnauthorized {
		return messages.ChangePageMsg{Page: NewConnectPage()}
	}
	defer res.Body.Close()

	var projects []repo.ProjectDTO
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		return messages.ProjectErrMsg(err)
	}

	log.Println("Projects loaded:", len(projects)) // Debugging output, replace with proper logging if needed
	return messages.ProjectsInitDataMsg(projects)
}
