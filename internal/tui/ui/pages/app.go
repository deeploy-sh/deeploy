package pages

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/tui/api"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/components"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/theme"
)

const headerHeight = 1
const footerHeight = 1

type HelpProvider interface {
	HelpKeys() []key.Binding
}

type app struct {
	currentPage      tea.Model
	palette          *components.Palette
	themeSwitcher    *components.ThemeSwitcher
	projects         []model.Project
	pods             []model.Pod
	width            int
	height           int
	heartbeatStarted bool
	offline          bool
	bootstrapped     bool
	statusText       string
	statusType       msg.StatusType
}

func (m *app) Projects() []model.Project {
	return m.projects
}

func (m *app) Pods() []model.Pod {
	return m.pods
}

func (m *app) clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return msg.ClearStatus{}
	})
}

func NewApp() tea.Model {
	cfg, err := config.Load()
	if err == nil && cfg.Theme != "" {
		theme.SetTheme(cfg.Theme)
	}

	return &app{
		currentPage: NewBootstrap(),
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(
		m.currentPage.Init(),
		tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
			return api.CheckConnection()()
		}),
	)
}

func (m app) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.ConnectionResult:
		switch {
		case tmsg.NeedsSetup:
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewConnectPage(nil) },
				}
			}
		case tmsg.NeedsAuth:
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewAuthPage("") },
				}
			}
		case tmsg.Offline:
			m.offline = true
			if !m.bootstrapped {
				bp, ok := m.currentPage.(*bootstrap)
				if ok {
					bp.offline = true
				}
				return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
					return api.CheckConnection()()
				})
			}
		case tmsg.Online:
			m.offline = false
		}

		if m.heartbeatStarted {
			return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return api.CheckConnection()()
			})
		}

		m.heartbeatStarted = true
		m.bootstrapped = true

		// Load data - Dashboard will be created in DataLoaded handler
		return m, tea.Batch(
			api.LoadData(),
			tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
				return api.CheckConnection()()
			}),
		)

	case tea.KeyPressMsg:
		if m.offline && tmsg.String() != "ctrl+c" {
			return m, nil
		}
		if tmsg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if tmsg.String() == "alt+p" && m.bootstrapped && m.palette == nil {
			if m.themeSwitcher != nil {
				theme.SetTheme(m.themeSwitcher.OriginalTheme())
				m.themeSwitcher = nil
			}
			palette := components.NewPalette(m.getPaletteItems())
			palette.SetSize(50, 20)
			m.palette = &palette
			return m, palette.Init()
		}

		if tmsg.Code == tea.KeyEscape {
			if m.themeSwitcher != nil {
				theme.SetTheme(m.themeSwitcher.OriginalTheme())
				m.themeSwitcher = nil
				return m, nil
			}
			if m.palette != nil {
				m.palette = nil
				return m, nil
			}
		}

		if m.themeSwitcher != nil {
			var cmd tea.Cmd
			*m.themeSwitcher, cmd = m.themeSwitcher.Update(tmsg)
			return m, cmd
		}

		if m.palette != nil {
			var cmd tea.Cmd
			*m.palette, cmd = m.palette.Update(tmsg)
			return m, cmd
		}

		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = tmsg.Width
		m.height = tmsg.Height

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

	case msg.ChangePage:
		m.currentPage = tmsg.PageFactory(&m)
		m.palette = nil

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

	case msg.DataLoaded:
		m.projects = tmsg.Projects
		m.pods = tmsg.Pods

		// Forward to current page so it can update its list
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)

		// If still on bootstrap, switch to dashboard now
		_, onBootstrap := m.currentPage.(*bootstrap)
		if onBootstrap {
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			}
		}
		return m, cmd

	case msg.Error:
		// If unauthorized, redirect to auth page
		if tmsg.Err == errs.ErrUnauthorized {
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewAuthPage("") },
				}
			}
		}
		// Show error in status line
		m.statusText = tmsg.Err.Error()
		m.statusType = msg.StatusError
		// Forward to current page
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)
		return m, tea.Batch(cmd, m.clearStatusAfter(5*time.Second))

	case msg.AuthSuccess:
		return m, tea.Batch(
			api.LoadData(),
			func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			},
		)

	case msg.ShowStatus:
		m.statusText = tmsg.Text
		m.statusType = tmsg.Type
		return m, m.clearStatusAfter(3 * time.Second)

	case msg.ClearStatus:
		m.statusText = ""
		return m, nil

	case msg.ThemeSwitcherClose:
		m.themeSwitcher = nil
		return m, nil

	case msg.OpenThemeSwitcher:
		m.palette = nil
		switcher := components.NewThemeSwitcher()
		m.themeSwitcher = &switcher
		return m, switcher.Init()

	default:
		if m.palette != nil {
			var cmd tea.Cmd
			*m.palette, cmd = m.palette.Update(tmsg)
			return m, cmd
		}

		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)
		return m, cmd
	}
}

type PageInfo interface {
	Breadcrumbs() []string
}

func (m app) getPaletteItems() []components.PaletteItem {
	items := []components.PaletteItem{
		{
			ItemTitle:   "Dashboard",
			Description: "Go to dashboard",
			Category:    "action",
			Action: func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewDashboard(s) },
				}
			},
		},
		{
			ItemTitle:   "New Project",
			Description: "Create a new project",
			Category:    "action",
			Action: func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(m msg.Store) tea.Model { return NewProjectFormPage(nil) },
				}
			},
		},
		{
			ItemTitle:   "Change Theme",
			Description: "Switch color theme",
			Category:    "settings",
			Action: func() tea.Msg {
				return msg.OpenThemeSwitcher{}
			},
		},
		{
			ItemTitle:   "Git Tokens",
			Description: "Manage Git tokens for private repos",
			Category:    "settings",
			Action: func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewGitTokensPage() },
				}
			},
		},
	}

	for _, p := range m.projects {
		project := p
		items = append(items, components.PaletteItem{
			ItemTitle:   project.Title,
			Description: "",
			Category:    "project",
			Action: func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewProjectDetailPage(s, project.ID) },
				}
			},
		})
	}

	for _, p := range m.pods {
		pod := p
		items = append(items, components.PaletteItem{
			ItemTitle:   pod.Title,
			Description: "",
			Category:    "pod",
			Action: func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						// Find project for this pod
						var project *model.Project
						for _, pr := range s.Projects() {
							if pr.ID == pod.ProjectID {
								proj := pr
								project = &proj
								break
							}
						}
						return NewPodDetailPage(&pod, project)
					},
				}
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

	var status string
	var statusStyle lipgloss.Style
	if m.offline {
		status = "● reconnecting"
		statusStyle = styles.OfflineStyle()
	} else {
		status = "● online"
		statusStyle = styles.OnlineStyle()
	}

	logo := "⚡ deeploy.sh"
	breadcrumbParts := []string{logo}
	p, ok := m.currentPage.(PageInfo)
	if ok {
		breadcrumbParts = append(breadcrumbParts, p.Breadcrumbs()...)
	}
	breadcrumbs := strings.Join(breadcrumbParts, styles.MutedStyle().Render("  >  "))

	gap := max(m.width-lipgloss.Width(breadcrumbs)-lipgloss.Width(status)-2, 1)
	headerContent := breadcrumbs + strings.Repeat(" ", gap) + statusStyle.Render(status)
	header := lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1).
		Render(headerContent)

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
		helpText := components.RenderHelpFooter(hp.HelpKeys())

		// Add status message if present
		var statusMsg string
		if m.statusText != "" {
			var statusStyle lipgloss.Style
			var icon string
			switch m.statusType {
			case msg.StatusSuccess:
				statusStyle = styles.SuccessStyle()
				icon = "✓"
			case msg.StatusError:
				statusStyle = styles.ErrorStyle()
				icon = "✗"
			default:
				statusStyle = styles.MutedStyle()
				icon = "●"
			}
			statusMsg = statusStyle.Render(icon + " " + m.statusText)
		}

		hs := lipgloss.NewStyle().Padding(0, 1)
		if statusMsg != "" {
			footerGap := max(m.width-lipgloss.Width(helpText)-lipgloss.Width(statusMsg)-4, 1)
			helpView = hs.Render(helpText + strings.Repeat(" ", footerGap) + statusMsg)
		} else {
			helpView = hs.Render(helpText)
		}
	}

	base := lipgloss.JoinVertical(lipgloss.Left, header, contentArea, helpView)

	if m.themeSwitcher != nil {
		switcherCard := m.themeSwitcher.View()

		switcherWidth := lipgloss.Width(switcherCard)
		switcherHeight := lipgloss.Height(switcherCard)
		switcherX := (m.width - switcherWidth) / 2
		switcherY := (m.height - switcherHeight) / 2

		baseLayer := lipgloss.NewLayer(base)
		switcherLayer := lipgloss.NewLayer(switcherCard).
			X(switcherX).
			Y(switcherY).
			Z(1)

		canvas := lipgloss.NewCanvas(baseLayer, switcherLayer)
		return tea.NewView(canvas.Render())
	}

	if m.palette != nil {
		paletteCard := m.palette.View()

		paletteWidth := lipgloss.Width(paletteCard)
		paletteX := (m.width - paletteWidth) / 2
		paletteY := m.height * 3 / 10

		baseLayer := lipgloss.NewLayer(base)
		paletteLayer := lipgloss.NewLayer(paletteCard).
			X(paletteX).
			Y(paletteY).
			Z(1)

		canvas := lipgloss.NewCanvas(baseLayer, paletteLayer)
		return tea.NewView(canvas.Render())
	}

	x := tea.NewView(base)
	x.BackgroundColor = styles.ColorBackground()
	return x
}
