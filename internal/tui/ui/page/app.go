package page

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/shared/model"
	"github.com/deeploy-sh/deeploy/internal/shared/version"
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
	gitTokens        []model.GitToken
	width            int
	height           int
	heartbeatStarted bool
	offline          bool
	bootstrapped     bool
	statusText       string
	statusType       msg.StatusType
	serverVersion    string // From /health endpoint
	latestVersion    string // From GitHub API
	// Security: true if using HTTPS, false if using plain HTTP
	secureConnection bool
	// Loading state
	isLoading   bool
	loadingText string
	spinner     spinner.Model
}

func (m *app) Projects() []model.Project {
	return m.projects
}

func (m *app) Pods() []model.Pod {
	return m.pods
}

func (m *app) GitTokens() []model.GitToken {
	return m.gitTokens
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

	// Check if connection is secure (HTTPS)
	secureConnection := false
	if err == nil && cfg.Server != "" {
		secureConnection = strings.HasPrefix(cfg.Server, "https://")
	}

	// Initialize spinner for loading state
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.ColorPrimary())

	return &app{
		currentPage:      NewBootstrap(),
		secureConnection: secureConnection,
		spinner:          s,
	}
}

func (m app) Init() tea.Cmd {
	return tea.Batch(
		m.currentPage.Init(),
		tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
			return api.CheckConnection()()
		}),
		api.CheckLatestVersion(), // Check for updates on startup
	)
}

func (m app) Update(tmsg tea.Msg) (tea.Model, tea.Cmd) {
	switch tmsg := tmsg.(type) {
	case msg.ConnectionResult:
		switch {
		case tmsg.NeedsSetup:
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewConnect(nil) },
				}
			}
		case tmsg.NeedsAuth:
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewAuth("") },
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
			m.serverVersion = tmsg.ServerVersion
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
		// Allow quit even during loading
		if tmsg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Block all input during loading
		if m.isLoading {
			return m, nil
		}

		// Allow palette (alt+p) even when offline so user can change server
		if m.offline && tmsg.String() != "alt+p" {
			return m, nil
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
		m.gitTokens = tmsg.GitTokens

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
		m.isLoading = false
		// If unauthorized, redirect to auth page
		if tmsg.Err == errs.ErrUnauthorized {
			return m, func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewAuth("") },
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

	case msg.StartLoading:
		m.isLoading = true
		m.loadingText = tmsg.Text
		m.statusText = ""
		return m, m.spinner.Tick

	case msg.ProjectCreated, msg.ProjectUpdated, msg.ProjectDeleted,
		msg.PodCreated, msg.PodUpdated, msg.PodDeleted,
		msg.PodDeployed, msg.PodStopped, msg.PodRestarted,
		msg.GitTokenCreated, msg.GitTokenDeleted,
		msg.PodDomainCreated, msg.PodDomainUpdated, msg.PodDomainDeleted,
		msg.PodEnvVarsUpdated:
		m.isLoading = false
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)
		return m, cmd

	case msg.ThemeSwitcherClose:
		m.themeSwitcher = nil
		return m, nil

	case msg.OpenThemeSwitcher:
		m.palette = nil
		switcher := components.NewThemeSwitcher()
		m.themeSwitcher = &switcher
		return m, switcher.Init()

	case msg.LatestVersionResult:
		if tmsg.Error == nil && tmsg.Version != "" {
			m.latestVersion = tmsg.Version
		}
		// Check again in 1 hour
		return m, tea.Tick(1*time.Hour, func(t time.Time) tea.Msg {
			return api.CheckLatestVersion()()
		})

	case msg.ServerDomainSet:
		m.isLoading = false
		m.secureConnection = true
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)
		return m, cmd

	case msg.ServerDomainDeleted:
		m.isLoading = false
		m.secureConnection = false
		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)
		return m, cmd

	default:
		var cmds []tea.Cmd

		// Update spinner when loading
		if m.isLoading {
			var spinnerCmd tea.Cmd
			m.spinner, spinnerCmd = m.spinner.Update(tmsg)
			cmds = append(cmds, spinnerCmd)
		}

		if m.palette != nil {
			var cmd tea.Cmd
			*m.palette, cmd = m.palette.Update(tmsg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		var cmd tea.Cmd
		m.currentPage, cmd = m.currentPage.Update(tmsg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
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
					PageFactory: func(m msg.Store) tea.Model { return NewProjectForm(nil) },
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
					PageFactory: func(s msg.Store) tea.Model { return NewGitTokens(s.GitTokens()) },
				}
			},
		},
		{
			ItemTitle:   "About / Updates",
			Description: "Version info and updates",
			Category:    "settings",
			Action: func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model {
						return NewInfo(version.Version, m.serverVersion, m.latestVersion)
					},
				}
			},
		},
		{
			ItemTitle:   "Domain",
			Description: "Setup HTTPS with custom domain",
			Category:    "settings",
			Action: func() tea.Msg {
				return msg.ChangePage{
					PageFactory: func(s msg.Store) tea.Model { return NewServerDomain() },
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
					PageFactory: func(s msg.Store) tea.Model { return NewProjectDetail(s, project.ID) },
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
						return NewPodDetail(s, pod.ID)
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
		v := m.currentPage.View()
		v.MouseMode = tea.MouseModeCellMotion
		return v
	}

	var status string
	var statusStyle lipgloss.Style
	if m.offline {
		status = "● reconnecting"
		statusStyle = styles.OfflineStyle()
	} else if !m.secureConnection {
		// Show warning for insecure HTTP connection with hint
		status = "⚠ insecure (alt+p > Domain)"
		statusStyle = styles.WarningStyle()
	} else {
		status = "● online"
		statusStyle = styles.OnlineStyle()
	}

	// Build version info
	tuiVersion := version.Version
	serverVersion := m.serverVersion
	if serverVersion == "" {
		serverVersion = "..."
	}

	var versionInfo string
	if tuiVersion == serverVersion {
		versionInfo = fmt.Sprintf(" %s", tuiVersion)
	} else {
		versionInfo = fmt.Sprintf(" TUI %s · Server %s", tuiVersion, serverVersion)
	}

	// Check if update available
	updateAvailable := false
	if m.latestVersion != "" {
		if tuiVersion != "dev" && tuiVersion != m.latestVersion {
			updateAvailable = true
		}
		if serverVersion != "dev" && serverVersion != "..." && serverVersion != m.latestVersion {
			updateAvailable = true
		}
	}
	if updateAvailable {
		versionInfo += " ⬆"
	}

	status = statusStyle.Render(status) + styles.MutedStyle().Render(versionInfo)

	logo := lipgloss.NewStyle().Bold(true).Render("deeploy")
	breadcrumbParts := []string{logo}
	p, ok := m.currentPage.(PageInfo)
	if ok {
		breadcrumbParts = append(breadcrumbParts, p.Breadcrumbs()...)
	}
	breadcrumbs := strings.Join(breadcrumbParts, styles.MutedStyle().Render(" > "))

	gap := max(m.width-lipgloss.Width(breadcrumbs)-lipgloss.Width(status)-2, 1)
	headerContent := breadcrumbs + strings.Repeat(" ", gap) + status
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
		keys := hp.HelpKeys()
		helpText := components.RenderHelpFooter(keys)

		// Add loading spinner or status message
		var statusMsg string
		if m.isLoading {
			statusMsg = styles.PrimaryStyle().Render(strings.TrimSpace(m.spinner.View()) + " " + m.loadingText)
		} else if m.statusText != "" {
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
			footerGap := max(m.width-lipgloss.Width(helpText)-lipgloss.Width(statusMsg)-2, 1)
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
		v := tea.NewView(canvas.Render())
		v.MouseMode = tea.MouseModeCellMotion
		return v
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
		v := tea.NewView(canvas.Render())
		v.MouseMode = tea.MouseModeCellMotion
		return v
	}

	v := tea.NewView(base)
	v.BackgroundColor = styles.ColorBackground()
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
