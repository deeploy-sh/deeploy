package pages

import (
	"log"

	"github.com/deeploy-sh/deeploy/internal/shared/repo"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	confirmNo = iota
	confirmYes
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectDeletePage struct {
	project  *repo.ProjectDTO
	decision int
	width    int
	height   int
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectDeletePage(project *repo.ProjectDTO) ProjectDeletePage {
	return ProjectDeletePage{
		project: project,
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectDeletePage) Init() tea.Cmd {
	return textinput.Blink
}

func (p ProjectDeletePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			p.decision = confirmNo
			return p, nil
		case "right", "l":
			p.decision = confirmYes
			return p, nil
		case "tab":
			if p.decision == confirmNo {
				p.decision = confirmYes
			} else {
				p.decision = confirmNo
			}
		case "enter":
			if p.decision == confirmNo {
				return p, func() tea.Msg {
					return messages.ProjectPopPageMsg{}
				}
			}
			return p, tea.Batch(
				p.DeleteProject,
				func() tea.Msg {
					return messages.ProjectPopPageMsg{}
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

func (p ProjectDeletePage) View() string {
	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")

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

	// Help text
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("← →/h l to move • enter to confirm")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, noButton, yesButton)
	content := lipgloss.JoinVertical(0.5,
		title,
		buttons,
		help,
	)

	card := components.Card(components.CardProps{
		// Width:   p.width / 2,
		Padding: []int{2, 1},
	}).Render(content)

	view := lipgloss.JoinVertical(0.5, logo, card)
	return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectDeletePage) DeleteProject() tea.Msg {
	res, err := utils.Request(utils.RequestProps{
		Method: "DELETE",
		URL:    "/projects/" + p.project.ID,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	log.Println(res)
	return messages.ProjectDeleteMsg(p.project)
}
