package components

import (
	"fmt"
	"io"
	"strings"

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
	activeTheme string // The theme currently saved in config
}

func (d themeDelegate) Height() int                             { return 1 }
func (d themeDelegate) Spacing() int                            { return 0 }
func (d themeDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d themeDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(themeItem)
	if !ok {
		return
	}

	themeName := item.name
	t, exists := theme.Available[themeName]
	if !exists {
		return
	}

	isFocused := index == m.Index()
	isActive := themeName == d.activeTheme

	// Build the line
	var dot string
	if isActive {
		// Active theme gets a colored dot with the theme's primary color
		dot = lipgloss.NewStyle().
			Foreground(t.Primary()).
			Render("●")
	} else {
		dot = ""
	}

	// Theme name with padding
	name := fmt.Sprintf("%s %s", dot, themeName)

	var line string
	if isFocused {
		// Focused: use theme's BackgroundElement as subtle highlight
		// with theme's Primary as text color
		line = lipgloss.NewStyle().
			Background(t.BackgroundElement()).
			Foreground(t.Primary()).
			Bold(true).
			Width(36).
			// Padding(0, 1).
			Render(name)
	} else {
		// Normal item
		line = lipgloss.NewStyle().
			Width(36).
			// Padding(0, 1).
			Render(name)
	}

	fmt.Fprint(w, line)
}

// ThemeSwitcher is an overlay component for selecting themes with live preview
type ThemeSwitcher struct {
	list          list.Model
	originalTheme string
	width         int
	height        int
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

	// Create delegate
	delegate := themeDelegate{activeTheme: currentTheme}

	// Create list
	l := list.New(items, delegate, 40, 14)
	l.Title = "Select Theme"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(true) // We'll render our own title

	// Style the list
	l.Styles.TitleBar = lipgloss.NewStyle()
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(styles.ColorPrimary()).
		Bold(true)

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
		width:         40,
		height:        20,
	}
}

func (m ThemeSwitcher) Init() tea.Cmd {
	return nil
}

func (m ThemeSwitcher) Update(msg tea.Msg) (ThemeSwitcher, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		// Select theme and save
		case msg.Code == tea.KeyEnter:
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
		case msg.Code == tea.KeyEscape:
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

func (m *ThemeSwitcher) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height-4) // Account for title and help
}

func (m ThemeSwitcher) View() string {
	var b strings.Builder

	// Title - centered
	// title := "Select Theme"
	// titleStyle := lipgloss.NewStyle().
	// 	Foreground(styles.ColorPrimary()).
	// 	Bold(true).
	// 	Width(38).
	// 	Align(lipgloss.Center)
	//
	// b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// List
	b.WriteString(m.list.View())

	// Help text
	b.WriteString("\n")
	helpStyle := styles.MutedStyle().
		Italic(true).
		Width(38).
		Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("↑/↓ navigate • enter select • esc cancel"))

	return b.String()
}
