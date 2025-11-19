package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type HasInputView interface {
	HasFocusedInput() bool
}

type app struct {
	currentPage tea.Model
	width       int
	height      int
}

// /////////////////////////////////////////////////////////////////////////////
// Constructors
// /////////////////////////////////////////////////////////////////////////////

func NewApp() app {
	return app{}
}

// We wait for window size before creating pages
func (a app) Init() tea.Cmd {
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}

		currentPage := a.currentPage
		if page, ok := currentPage.(HasInputView); ok && page.HasFocusedInput() {
			// this disable "q"
		} else if msg.String() == "q" {
			return a, tea.Quit
		}

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		// If no pages yet, create first one
		if a.currentPage == nil {
			// config, err := config.LoadConfig()
			// config, _ := config.LoadConfig()
			// log.Println(config)
			var page tea.Model

			page = NewBootstrap()
			// No config = show login, has config = show dashboard
			// if err != nil || config.Server == "" || config.Token == "" {
			// 	page = NewConnectPage()
			// } else {
			// 	page = NewDashboard()
			// }

			// Add first page to stack
			a.currentPage = page

			// Update page with window size and initialize it
			updatedPage, cmd := page.Update(msg)
			a.currentPage = updatedPage
			return a, tea.Batch(cmd, updatedPage.Init())
		}

		// Update current page's window size
		currentPage := a.currentPage
		updatedPage, cmd := currentPage.Update(msg)
		a.currentPage = updatedPage
		return a, cmd

	case messages.ChangePageMsg:
		newPage := msg.Page

		a.currentPage = newPage

		// Batch window size and init commands together
		// This prevents double rendering by ensuring both happen in sequence
		return a, tea.Batch(
			func() tea.Msg {
				return tea.WindowSizeMsg{
					Width:  a.width,
					Height: a.height,
				}
			},
			newPage.Init(),
		)

	// All other messages go to current page
	default:
		if a.currentPage == nil {
			return a, nil
		}
		currentPage := a.currentPage
		updatedPage, cmd := currentPage.Update(msg)
		a.currentPage = updatedPage
		return a, cmd
	}
}

type FooterMenuItem struct {
	Key  string
	Desc string
}

func (a app) View() string {
	if a.currentPage == nil {
		return "Loading..."
	}

	main := a.currentPage.View()

	footerMenuItems := []FooterMenuItem{
		{Key: ":", Desc: "menu"},
		{Key: "esc", Desc: "back"},
		{Key: "q", Desc: "quit"},
	}

	var footer strings.Builder

	for i, v := range footerMenuItems {
		footer.WriteString(styles.FocusedStyle.Render(v.Key))
		footer.WriteString(" ")
		footer.WriteString(v.Desc)
		if len(footerMenuItems)-1 != i {
			footer.WriteString(" • ")
		}
	}

	footerCard := components.Card(components.CardProps{
		Width:   a.width,
		Padding: []int{0, 1},
	}).Render(footer.String())

	horizontal := lipgloss.JoinHorizontal(0.5, main)
	view := lipgloss.JoinVertical(lipgloss.Bottom, horizontal, footerCard)

	return view
}
