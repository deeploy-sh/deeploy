package components

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/config"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/theme"
)

// ThemeSwitcherCloseMsg is sent when the theme switcher should close
type ThemeSwitcherCloseMsg struct {
	Selected bool // true if a theme was selected, false if cancelled
}

// OpenThemeSwitcherMsg is sent to open the theme switcher from palette
type OpenThemeSwitcherMsg struct{}

// themeItem represents a theme in the list
type themeItem struct {
	name     string
	isActive bool
}

func (i themeItem) Title() string       { return i.name }
func (i themeItem) FilterValue() string { return i.name }
func (i themeItem) Prefix() string {
	if i.isActive {
		return "●"
	}
	return " "
}

// ThemeSwitcher is an overlay component for selecting themes with live preview
type ThemeSwitcher struct {
	list          ScrollList
	originalTheme string
}

// NewThemeSwitcher creates a new theme switcher overlay
func NewThemeSwitcher() ThemeSwitcher {
	themes := theme.ThemeNames()
	currentTheme := theme.Current.Name()

	// Create list items
	items := make([]ScrollItem, len(themes))
	for i, t := range themes {
		items[i] = themeItem{name: t, isActive: t == currentTheme}
	}

	card := CardProps{Width: 50, Padding: []int{1, 1}}
	l := NewScrollList(items, card.InnerWidth(), 15)

	// Find and select current theme
	for i, t := range themes {
		if t == currentTheme {
			l.Select(i)
			break
		}
	}

	return ThemeSwitcher{
		list:          l,
		originalTheme: currentTheme,
	}
}

func (m ThemeSwitcher) Init() tea.Cmd {
	return nil
}

// OriginalTheme returns the theme that was active when the switcher was opened
func (m ThemeSwitcher) OriginalTheme() string {
	return m.originalTheme
}

func (m ThemeSwitcher) Update(msg tea.Msg) (ThemeSwitcher, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		// Select theme and save
		case tea.KeyEnter:
			// Get selected theme
			if item := m.list.SelectedItem(); item != nil {
				if ti, ok := item.(themeItem); ok {
					// Save to config
					cfg, err := config.Load()
					if err != nil {
						cfg = &config.Config{}
					}
					cfg.Theme = ti.name
					_ = config.Save(cfg)

					return m, func() tea.Msg {
						return ThemeSwitcherCloseMsg{Selected: true}
					}
				}
			}

		// Cancel and revert
		case tea.KeyEscape:
			// Revert to original theme
			theme.SetTheme(m.originalTheme)
			return m, func() tea.Msg {
				return ThemeSwitcherCloseMsg{Selected: false}
			}

		case tea.KeyUp:
			prevIndex := m.list.Index()
			m.list.CursorUp()
			if m.list.Index() != prevIndex {
				if item := m.list.SelectedItem(); item != nil {
					if ti, ok := item.(themeItem); ok {
						theme.SetTheme(ti.name)
					}
				}
			}
			return m, nil

		case tea.KeyDown:
			prevIndex := m.list.Index()
			m.list.CursorDown()
			if m.list.Index() != prevIndex {
				if item := m.list.SelectedItem(); item != nil {
					if ti, ok := item.(themeItem); ok {
						theme.SetTheme(ti.name)
					}
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func (m ThemeSwitcher) View() string {
	card := CardProps{Width: 50, Padding: []int{1, 1}}
	w := card.InnerWidth()

	// Custom title
	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render("Select Theme")

	// List with background
	list := lipgloss.NewStyle().
		Width(w).
		Height(m.list.Height()).
		Background(styles.ColorBackgroundPanel()).
		Render(m.list.View())

	content := lipgloss.JoinVertical(lipgloss.Left, title, list)

	return Card(card).Render(content)
}
