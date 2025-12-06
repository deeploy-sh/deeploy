package pages

import (
	"log"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
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

type deleteKeyMap struct {
	Select  key.Binding
	Confirm key.Binding
	Cancel  key.Binding
}

func (k deleteKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Confirm, k.Cancel}
}

func (k deleteKeyMap) FullHelp() [][]key.Binding {
	return nil
}

func (p ProjectDeletePage) HelpKeys() help.KeyMap {
	return p.keys
}

func newDeleteKeyMap() deleteKeyMap {
	return deleteKeyMap{
		Select:  key.NewBinding(key.WithKeys("left", "right", "h", "l", "tab"), key.WithHelp("←/→", "select")),
		Confirm: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
		Cancel:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	}
}

const (
	confirmNo = iota
	confirmYes
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectDeletePage struct {
	project  *repo.Project
	decision int
	keys     deleteKeyMap
	width    int
	height   int
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectDeletePage(project *repo.Project) ProjectDeletePage {
	return ProjectDeletePage{
		project: project,
		keys:    newDeleteKeyMap(),
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectDeletePage) Init() tea.Cmd {
	return nil
}

func (p ProjectDeletePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyLeft, 'h':
			p.decision = confirmNo
			return p, nil
		case tea.KeyRight, 'l':
			p.decision = confirmYes
			return p, nil
		case tea.KeyTab:
			if p.decision == confirmNo {
				p.decision = confirmYes
			} else {
				p.decision = confirmNo
			}
		case tea.KeyEscape:
			return p, func() tea.Msg {
				return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewDashboard(s) }}
			}
		case tea.KeyEnter:
			if p.decision == confirmNo {
				return p, func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewDashboard(s) }}
				}
			}
			return p, tea.Batch(
				p.DeleteProject,
				func() tea.Msg {
					return ChangePageMsg{PageFactory: func(s Store) tea.Model { return NewDashboard(s) }}
				},
			)
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	}
	return p, nil
}

func (p ProjectDeletePage) View() tea.View {
	title := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 0, 1, 0).
		Render("Delete " + p.project.Title + "?")

	// Button Styles
	baseButton := lipgloss.NewStyle().
		Padding(0, 3).
		Width(1).
		MarginRight(1)

	activeButton := baseButton.
		Background(styles.ColorPrimary).
		Foreground(lipgloss.Color("0"))

	inactiveButton := baseButton.
		Background(lipgloss.Color("237"))

	// Render Buttons based on decision
	var noButton, yesButton string
	if p.decision == confirmNo {
		noButton = activeButton.Render("NO")
		yesButton = inactiveButton.Render("YES")
	} else {
		noButton = inactiveButton.Render("NO")
		yesButton = activeButton.Render("YES")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, noButton, yesButton)
	content := lipgloss.JoinVertical(0.5, title, buttons)

	card := components.Card(components.CardProps{
		Padding: []int{2, 1},
	}).Render(content)

	contentHeight := p.height

	// Card vertikal zentrieren
	centeredCard := lipgloss.Place(p.width, contentHeight,
		lipgloss.Center, lipgloss.Center, card)

	return tea.NewView(centeredCard)
}

func (p ProjectDeletePage) Breadcrumbs() []string {
	return []string{"Projects", "Delete"}
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectDeletePage) DeleteProject() tea.Msg {
	res, err := utils.Request("DELETE", "/projects/"+p.project.ID, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	return messages.ProjectDeleteMsg(p.project)
}
