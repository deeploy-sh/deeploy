package components

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/config"
	"github.com/deeploy-sh/deeploy/internal/tui/msg"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/theme"
)

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

type ThemeSwitcher struct {
	list          ScrollList
	originalTheme string
}

func NewThemeSwitcher() ThemeSwitcher {
	themes := theme.ThemeNames()
	currentTheme := theme.Current.Name()

	items := make([]ScrollItem, len(themes))
	for i, t := range themes {
		items[i] = themeItem{name: t, isActive: t == currentTheme}
	}

	card := styles.CardProps{Width: 50, Padding: []int{1, 1}}
	l := NewScrollList(items, ScrollListConfig{
		Width:  card.InnerWidth(),
		Height: 15,
	})

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

func (m ThemeSwitcher) OriginalTheme() string {
	return m.originalTheme
}

func (m ThemeSwitcher) Update(tmsg tea.Msg) (ThemeSwitcher, tea.Cmd) {
	prevIndex := m.list.Index()

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(tmsg)

	if m.list.Index() != prevIndex {
		if item := m.list.SelectedItem(); item != nil {
			if ti, ok := item.(themeItem); ok {
				theme.SetTheme(ti.name)
			}
		}
	}

	switch tmsg := tmsg.(type) {
	case tea.KeyPressMsg:
		switch tmsg.Code {
		case tea.KeyEnter:
			if item := m.list.SelectedItem(); item != nil {
				if ti, ok := item.(themeItem); ok {
					cfg, err := config.Load()
					if err != nil {
						cfg = &config.Config{}
					}
					cfg.Theme = ti.name
					_ = config.Save(cfg)
					return m, func() tea.Msg {
						return msg.ThemeSwitcherClose{Theme: ti.name}
					}
				}
			}
		case tea.KeyEscape:
			theme.SetTheme(m.originalTheme)
			return m, func() tea.Msg {
				return msg.ThemeSwitcherClose{Theme: m.originalTheme}
			}
		}
	}

	return m, cmd
}

func (m ThemeSwitcher) View() string {
	card := styles.CardProps{Width: 50, Padding: []int{1, 1}}
	w := card.InnerWidth()

	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render("Select Theme")

	list := lipgloss.NewStyle().
		Width(w).
		Height(m.list.Height()).
		Background(styles.ColorBackgroundPanel()).
		Render(m.list.View())

	content := lipgloss.JoinVertical(lipgloss.Left, title, list)

	return styles.Card(card).Render(content)
}
