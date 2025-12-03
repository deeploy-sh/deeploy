package pages

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/messages"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/components"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/utils"
	"github.com/deeploy-sh/deeploy/internal/deeployd/repo"
)

const headerHeight = 1

type app struct {
	currentPage      tea.Model
	palette          *components.Palette
	projects         []repo.Project
	width            int
	height           int
	heartbeatStarted bool
	offline          bool
	bootstrapped     bool
}

func NewApp() tea.Model {
	return &app{
		currentPage: NewBootstrap(),
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(
		m.currentPage.Init(),
		// INFO: use tick here to show bootstrap(logo) min. 1 second
		tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
			return utils.CheckConnection()
		}),
	)
}

func (m app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case utils.ConnectionResultMsg:
		switch {
		case msg.NeedsSetup:
			return m, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewConnectPage(nil)}
			}
		case msg.NeedsAuth:
			return m, func() tea.Msg {
				return messages.ChangePageMsg{Page: NewAuthPage("")}
			}
		case msg.Offline:
			m.offline = true
			if !m.bootstrapped {
				bp, ok := m.currentPage.(*bootstrap)
				if ok {
					bp.offline = true
				}
				return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
					return utils.CheckConnection()
				})
			}
		case msg.Online:
			m.offline = false
		}

		if m.heartbeatStarted {
			return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return utils.CheckConnection()
			})
		}

		m.heartbeatStarted = true
		m.bootstrapped = true

		return m, tea.Batch(
			func() tea.Msg {
				return messages.ChangePageMsg{Page: NewDashboard()}
			},
			tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return utils.CheckConnection()
			}),
			loadAppProjects,
		)

	case tea.KeyPressMsg:
		if m.offline && msg.String() != "ctrl+c" {
			return m, nil // block app
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Ctrl+K toggles palette
		if msg.String() == "ctrl+k" && m.bootstrapped {
			if m.palette != nil {
				m.palette = nil
			} else {
				palette := components.NewPalette(m.getPaletteItems())
				palette.SetSize(50, 20)
				m.palette = &palette
				return m, palette.Init()
			}
			return m, nil
		}

		// Esc closes palette
		if msg.Code == tea.KeyEscape && m.palette != nil {
			m.palette = nil
			return m, nil
		}

		// Forward to palette if open
		if m.palette != nil {
			var cmd tea.Cmd
			*m.palette, cmd = m.palette.Update(msg)
			return m, cmd
		}

		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		pageMsg := tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		}
		var cmd tea.Cmd
		if m.currentPage == nil {
			return m, nil
		}
		m.currentPage, cmd = m.currentPage.Update(pageMsg)
		return m, cmd

	case messages.ChangePageMsg:
		m.currentPage = msg.Page
		m.palette = nil // Close palette on page change

		pageMsg := tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height - headerHeight,
		}
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(pageMsg)

		return m, tea.Batch(
			m.currentPage.Init(),
			cmd,
		)

	case appProjectsMsg:
		m.projects = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(msg)
		return m, cmd
	}
}

// appProjectsMsg holds loaded projects for the palette
type appProjectsMsg []repo.Project

// loadAppProjects loads projects for the command palette
func loadAppProjects() tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return nil
	}

	r, err := http.NewRequest("GET", cfg.Server+"/api/projects", nil)
	if err != nil {
		return nil
	}
	r.Header.Set("Authorization", "Bearer "+cfg.Token)

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var projects []repo.Project
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		return nil
	}

	return appProjectsMsg(projects)
}

// PageInfo interface for pages to provide breadcrumbs
type PageInfo interface {
	Breadcrumbs() []string
}

// getPaletteItems returns the items for the command palette
func (m app) getPaletteItems() []components.PaletteItem {
	items := []components.PaletteItem{
		// Actions
		{
			Title:       "Dashboard",
			Description: "Go to dashboard",
			Category:    "action",
			Action: func() tea.Msg {
				return messages.ChangePageMsg{Page: NewDashboard()}
			},
		},
		{
			Title:       "New Project",
			Description: "Create a new project",
			Category:    "action",
			Action: func() tea.Msg {
				return messages.ChangePageMsg{Page: NewProjectFormPage(nil)}
			},
		},
	}

	// Add projects dynamically
	for _, p := range m.projects {
		project := p // Capture for closure
		items = append(items, components.PaletteItem{
			Title:       project.Title,
			Description: project.Description,
			Category:    "project",
			Action: func() tea.Msg {
				return messages.ChangePageMsg{Page: NewProjectDetailPage(project.ID)}
			},
		})
	}

	return items
}

func (m app) View() tea.View {
	_, ok := m.currentPage.(*bootstrap)
	if ok {
		return m.currentPage.View()
	}

	// Status
	var status string
	var statusStyle lipgloss.Style
	if m.offline {
		status = "● reconnecting"
		statusStyle = styles.OfflineStyle
	} else {
		status = "● online"
		statusStyle = styles.OnlineStyle
	}

	// Breadcrumbs
	logo := "⚡ deeploy.sh"
	breadcrumbParts := []string{logo}
	if p, ok := m.currentPage.(PageInfo); ok {
		breadcrumbParts = append(breadcrumbParts, p.Breadcrumbs()...)
	}
	breadcrumbs := strings.Join(breadcrumbParts, styles.MutedStyle.Render("  >  "))

	// Header - minimal, ohne Border
	gap := max(m.width-lipgloss.Width(breadcrumbs)-lipgloss.Width(status)-2, 1)
	headerContent := breadcrumbs + strings.Repeat(" ", gap) + statusStyle.Render(status)
	header := lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1).
		Render(headerContent)

	// Content - Pages render their own help footer
	content := m.currentPage.View().Content
	contentHeight := m.height - headerHeight

	contentArea := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	)

	base := lipgloss.JoinVertical(lipgloss.Left, header, contentArea)

	// Render palette overlay if open
	if m.palette != nil {
		paletteContent := m.palette.View()
		paletteCard := components.Card(components.CardProps{
			Width:   54,
			Padding: []int{1, 2},
		}).Render(paletteContent)

		// Calculate palette position (centered horizontally, 30% from top)
		paletteWidth := lipgloss.Width(paletteCard)
		paletteX := (m.width - paletteWidth) / 2
		paletteY := m.height * 3 / 10

		// Create layers with proper z-ordering
		baseLayer := lipgloss.NewLayer(base)
		paletteLayer := lipgloss.NewLayer(paletteCard).
			X(paletteX).
			Y(paletteY).
			Z(1)

		// Compose with Canvas
		canvas := lipgloss.NewCanvas(baseLayer, paletteLayer)
		return tea.NewView(canvas.Render())
	}

	return tea.NewView(base)
}
