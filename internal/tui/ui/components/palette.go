package components

import (
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

// PaletteItem represents a single item in the command palette
type PaletteItem struct {
	ItemTitle   string
	Description string
	Category    string         // "project", "pod", "action", "settings"
	Action      func() tea.Msg // Action to execute on selection
}

func (i PaletteItem) Title() string       { return i.ItemTitle }
func (i PaletteItem) FilterValue() string { return i.ItemTitle + " " + i.Description }

func (i PaletteItem) Prefix() string {
	switch i.Category {
	case "project":
		return "[P]"
	case "pod":
		return "[D]"
	case "action":
		return "[A]"
	case "settings":
		return "[S]"
	default:
		return "   "
	}
}

// Palette is a command palette component
type Palette struct {
	list  ScrollList
	width int
}

func NewPalette(items []PaletteItem) Palette {
	scrollItems := make([]ScrollItem, len(items))
	for i, item := range items {
		scrollItems[i] = item
	}

	card := styles.CardProps{Width: styles.CardWidthMD, Padding: []int{1, 1}, Accent: true}
	list := NewScrollList(scrollItems, ScrollListConfig{
		Width:       card.InnerWidth(),
		Height:      10,
		WithInput:   true,
		Placeholder: "Type to search...",
	})

	return Palette{list: list, width: styles.CardWidthMD}
}

func (m Palette) Init() tea.Cmd {
	return m.list.Init()
}

func (m Palette) Update(msg tea.Msg) (Palette, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Enter executes action
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if ok && keyMsg.Code == tea.KeyEnter {
		item := m.list.SelectedItem()
		if item != nil {
			pi, ok := item.(PaletteItem)
			if ok && pi.Action != nil {
				return m, pi.Action
			}
		}
	}

	return m, cmd
}

func (m *Palette) SetSize(width, height int) {
	m.width = width
}

func (m Palette) View() string {
	card := styles.CardProps{Width: styles.CardWidthMD, Padding: []int{1, 1}, Accent: true}
	w := card.InnerWidth()

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render("Command Palette")

	// Input
	input := m.list.InputView()

	// List
	list := lipgloss.NewStyle().
		Width(w).
		Height(m.list.Height()).
		Background(styles.ColorBackgroundPanel()).
		Render(m.list.View())

	content := lipgloss.JoinVertical(lipgloss.Left, title, input, list)

	return styles.Card(card).Render(content)
}
