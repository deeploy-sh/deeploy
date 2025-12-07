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

	// Dot column - fixed width, always takes same space
	var dot string
	if isActive {
		dot = lipgloss.NewStyle().
			Foreground(t.Primary()).
			Render("●")
	} else {
		dot = " " // Space placeholder to keep alignment
	}

	// Format: "●  themename" - dot is in its own column, name always starts at same position
	content := fmt.Sprintf("%s  %s", dot, themeName)

	// Width: Card(52) - Padding(4) - Border(1) = 47
	const listWidth = 47

	var line string
	if isFocused {
		// Focused: background highlight with theme colors
		line = lipgloss.NewStyle().
			Background(t.BackgroundElement()).
			Foreground(t.Primary()).
			Bold(true).
			Width(listWidth).
			Render(content)
	} else {
		// Normal item
		line = lipgloss.NewStyle().
			Width(listWidth).
			Render(content)
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

	// Card dimensions: Width=52, Padding=2 left/right, Border=1
	// Inner width: 52 - 4 - 1 = 47
	const listWidth = 47
	const listHeight = 18 // More items visible

	// Create list
	l := list.New(items, delegate, listWidth, listHeight)
	l.Title = "Select Theme"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowTitle(true)

	// Style the list - align title with the theme names (after dot column)
	l.Styles.TitleBar = lipgloss.NewStyle().
		Padding(0, 0, 1, 0)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(styles.ColorPrimary()).
		Bold(true).
		PaddingLeft(3) // Align with theme names (after "●  ")

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
		width:         listWidth,
		height:        listHeight,
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

	// List (with native title)
	b.WriteString(m.list.View())

	// Help text - aligned with content
	b.WriteString("\n")
	helpStyle := styles.MutedStyle().
		Italic(true).
		PaddingLeft(3) // Align with theme names
	b.WriteString(helpStyle.Render("↑/↓ navigate • enter • esc"))

	return b.String()
}
