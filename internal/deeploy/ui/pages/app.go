package pages

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
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
const footerHeight = 1

type Store interface {
	Projects() []repo.Project
	Pods() []repo.Pod
}

type HelpProvider interface {
	HelpKeys() help.KeyMap
}

type app struct {
	currentPage      tea.Model
	palette          *components.Palette
	projects         []repo.Project
	pods             []repo.Pod
	width            int
	height           int
	heartbeatStarted bool
	offline          bool
	bootstrapped     bool
	help             help.Model
}

func (m *app) Projects() []repo.Project {
	return m.projects
}

func (m *app) Pods() []repo.Pod {
	return m.pods
}

func NewApp() tea.Model {
	return &app{
		currentPage: NewBootstrap(),
		help:        styles.NewHelpModel(),
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
			initData,
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
			Height: m.height - headerHeight - footerHeight,
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
			Height: m.height - headerHeight - footerHeight,
		}
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(pageMsg)

		return m, tea.Batch(
			m.currentPage.Init(),
			cmd,
		)

	case InitDataMsg:
		m.projects = msg.projects
		m.pods = msg.pods
		return m, nil

	default:
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(msg)
		return m, cmd
	}
}

// InitDataMsg holds loaded projects for the palette
type InitDataMsg struct {
	projects []repo.Project
	pods     []repo.Pod
	// later we add settings, domains, etc.
}

func initData() tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return nil
	}

	projects, err := initProjects(*cfg)
	if err != nil {
		log.Fatal("something went wrong: ", err)
	}

	pods, err := initPods(*cfg)
	if err != nil {
		log.Fatal("something went wrong: ", err)
	}

	return InitDataMsg{
		projects: projects,
		pods:     pods,
	}
}

func initProjects(cfg config.Config) ([]repo.Project, error) {
	req, err := http.NewRequest("GET", cfg.Server+"/api/projects", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var projects []repo.Project
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func initPods(cfg config.Config) ([]repo.Pod, error) {
	req, err := http.NewRequest("GET", cfg.Server+"/api/pods", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var pods []repo.Pod
	err = json.NewDecoder(res.Body).Decode(&pods)
	if err != nil {
		return nil, err
	}

	return pods, nil
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

	for _, p := range m.pods {
		pod := p // Capture for closure
		items = append(items, components.PaletteItem{
			Title:       pod.Title,
			Description: pod.Description,
			Category:    "pod",
			Action: func() tea.Msg {
				return nil
				// return messages.ChangePageMsg{Page: NewProjectDetailPage(project.ID)}
			},
		})
	}

	log.Println(m.pods)
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
	contentHeight := m.height - headerHeight - footerHeight

	contentArea := lipgloss.Place(
		m.width,
		contentHeight,
		lipgloss.Left,
		lipgloss.Top,
		content,
	)

	var helpView string
	hp, ok := m.currentPage.(HelpProvider)
	if ok {
		hs := lipgloss.NewStyle().Padding(0, 1)
		helpView = hs.Render(m.help.View(hp.HelpKeys()))
	}

	base := lipgloss.JoinVertical(lipgloss.Left, header, contentArea, helpView)

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
