package components

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
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
	name string
}

func (i themeItem) FilterValue() string { return i.name }

// themeDelegate handles rendering of theme items
type themeDelegate struct {
	activeTheme string
}

func (d themeDelegate) Height() int                             { return 1 }
func (d themeDelegate) Spacing() int                            { return 0 }
func (d themeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d themeDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(themeItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	isActive := item.name == d.activeTheme

	// Dot for active theme
	dot := " "
	if isActive {
		dot = "●"
	}
	content := fmt.Sprintf(" %s %s", dot, item.name)

	lineStyle := lipgloss.NewStyle().Width(m.Width())

	var line string
	if isSelected {
		line = lineStyle.
			Background(styles.ColorPrimary()).
			Foreground(styles.ColorBackground()).
			Bold(true).
			Render(content)
	} else {
		line = lineStyle.
			Background(styles.ColorBackgroundPanel()).
			Foreground(styles.ColorForeground()).
			Render(content)
	}

	fmt.Fprint(w, line)
}

// ThemeSwitcher is an overlay component for selecting themes with live preview
type ThemeSwitcher struct {
	list          list.Model
	originalTheme string
}

// NewThemeSwitcher creates a new theme switcher overlay
func NewThemeSwitcher() ThemeSwitcher {
	themes := theme.ThemeNames()
	currentTheme := theme.Current.Name()

	// Create list items
	items := make([]list.Item, len(themes))
	for i, t := range themes {
		items[i] = themeItem{name: t}
	}

	delegate := themeDelegate{activeTheme: currentTheme}

	card := CardProps{Width: 50, Padding: []int{1, 1}}
	l := list.New(items, delegate, card.InnerWidth(), 15)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.InfiniteScrolling = true

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
			if item, ok := m.list.SelectedItem().(themeItem); ok {
				// Save to config
				cfg, err := config.Load()
				if err != nil {
					cfg = &config.Config{}
				}
				cfg.Theme = item.name
				_ = config.Save(cfg)

				return m, func() tea.Msg {
					return ThemeSwitcherCloseMsg{Selected: true}
				}
			}

		// Cancel and revert
		case tea.KeyEscape:
			// Revert to original theme
			theme.SetTheme(m.originalTheme)
			return m, func() tea.Msg {
				return ThemeSwitcherCloseMsg{Selected: false}
			}
		}
	}

	// Update list and apply live preview
	var cmd tea.Cmd
	prevIndex := m.list.Index()
	m.list, cmd = m.list.Update(msg)

	// If selection changed, apply live preview
	if m.list.Index() != prevIndex {
		if item, ok := m.list.SelectedItem().(themeItem); ok {
			theme.SetTheme(item.name)
		}
	}

	return m, cmd
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
