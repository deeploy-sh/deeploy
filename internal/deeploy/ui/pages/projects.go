package pages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectPage struct {
	stack  []tea.Model
	width  int
	height int
}

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectPage() ProjectPage {
	return ProjectPage{
		stack: make([]tea.Model, 0),
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectPage) Init() tea.Cmd {
	return nil
}

func (p ProjectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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

		// Pass current page's KeyMsg
		if len(p.stack) == 0 {
			return p, nil
		}
		currentPage := p.stack[len(p.stack)-1]
		updatedPage, cmd := currentPage.Update(msg)
		p.stack[len(p.stack)-1] = updatedPage
		return p, cmd

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height

		// If no pages yet, create first one
		if len(p.stack) == 0 {

			page := NewProjectListPage()

			// Add first page to stack
			p.stack = append(p.stack, page)

			// Update page with window size and initialize it
			updatedPage, cmd := page.Update(msg)
			p.stack[len(p.stack)-1] = updatedPage
			return p, tea.Batch(cmd, updatedPage.Init())
		}

		// Update current page's window size
		currentPage := p.stack[len(p.stack)-1]
		updatedPage, cmd := currentPage.Update(msg)
		p.stack[len(p.stack)-1] = updatedPage
		return p, cmd

	case messages.ProjectPushPageMsg:
		newPage := msg.Page

		p.stack = append(p.stack, newPage)

		// Batch window size and init commands together
		// This prevents double rendering by ensuring both happen in sequence
		return p, tea.Batch(
			func() tea.Msg {
				return tea.WindowSizeMsg{
					Width:  p.width,
					Height: p.height,
				}
			},
			newPage.Init(),
		)

	case messages.ProjectPopPageMsg:
		if len(p.stack) > 1 {
			p.stack = p.stack[:len(p.stack)-1]
			return p, nil
		}

	default:
		// Forward all other messages to the current page
		if len(p.stack) > 0 {
			currentPage := p.stack[len(p.stack)-1]
			updatedPage, cmd := currentPage.Update(msg)
			p.stack[len(p.stack)-1] = updatedPage
			return p, cmd
		}
	}
	return p, nil
}

func (p ProjectPage) View() string {
	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")

	if len(p.stack) == 0 {
		return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(0.5, logo, "Loading..."))
	}

	main := p.stack[len(p.stack)-1].View()
	view := lipgloss.JoinVertical(0.5, logo, main)
	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}
