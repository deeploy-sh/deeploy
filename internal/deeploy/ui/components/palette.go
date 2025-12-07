package components

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
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
	textinput textinput.Model
	items     []PaletteItem
	list      ScrollList
	width     int
	height    int
}

func NewPalette(items []PaletteItem) Palette {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Prompt = ""
	ti.SetWidth(20)

	// Style with panel background
	bgStyle := lipgloss.NewStyle().Background(styles.ColorBackgroundPanel())
	inputStyles := textinput.Styles{
		Focused: textinput.StyleState{
			Text:        bgStyle.Foreground(styles.ColorForeground()),
			Placeholder: bgStyle.Foreground(styles.ColorMuted()),
		},
		Blurred: textinput.StyleState{
			Text:        bgStyle.Foreground(styles.ColorForeground()),
			Placeholder: bgStyle.Foreground(styles.ColorMuted()),
		},
		Cursor: textinput.CursorStyle{
			Blink: true,
		},
	}
	ti.SetStyles(inputStyles)

	// Convert to ScrollItems
	scrollItems := make([]ScrollItem, len(items))
	for i, item := range items {
		scrollItems[i] = item
	}

	card := CardProps{Width: 60, Padding: []int{1, 1}, Accent: true}
	list := NewScrollList(scrollItems, card.InnerWidth(), 10)

	p := Palette{
		textinput: ti,
		items:     items,
		list:      list,
	}

	return p
}

func (m Palette) Init() tea.Cmd {
	return textinput.Blink
}

func (m Palette) Update(msg tea.Msg) (Palette, tea.Cmd) {
	var cmd tea.Cmd

	// Always update textinput first (for blink, typing, etc.)
	m.textinput, cmd = m.textinput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case msg.Code == tea.KeyUp || msg.String() == "ctrl+p":
			m.list.CursorUp()
			return m, cmd
		case msg.Code == tea.KeyDown || msg.String() == "ctrl+n":
			m.list.CursorDown()
			return m, cmd
		case msg.Code == tea.KeyEnter:
			if item := m.list.SelectedItem(); item != nil {
				if pi, ok := item.(PaletteItem); ok && pi.Action != nil {
					return m, pi.Action
				}
			}
			return m, cmd
		}
	}

	// Filter items based on input
	m.filterItems()

	return m, cmd
}

// filterItems filters the items based on the current input
func (m *Palette) filterItems() {
	query := strings.ToLower(m.textinput.Value())

	var filtered []ScrollItem
	if query == "" {
		filtered = make([]ScrollItem, len(m.items))
		for i, item := range m.items {
			filtered[i] = item
		}
	} else {
		for _, item := range m.items {
			if strings.Contains(strings.ToLower(item.FilterValue()), query) {
				filtered = append(filtered, item)
			}
		}
	}

	m.list.SetItems(filtered)
}

func (m *Palette) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Palette) View() string {
	card := CardProps{Width: 60, Padding: []int{1, 1}, Accent: true}
	w := card.InnerWidth()

	// Title (like other lists)
	title := lipgloss.NewStyle().
		Bold(true).
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		Foreground(styles.ColorPrimary()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render("Command Palette")

	// Search input
	input := lipgloss.NewStyle().
		Width(w).
		Background(styles.ColorBackgroundPanel()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render(m.textinput.View())

	// List with background
	var listContent string
	if len(m.list.Items()) == 0 {
		listContent = styles.MutedStyle().Render(" No results")
	} else {
		listContent = m.list.View()
	}

	list := lipgloss.NewStyle().
		Width(w).
		Height(m.list.Height()).
		Background(styles.ColorBackgroundPanel()).
		Render(listContent)

	content := lipgloss.JoinVertical(lipgloss.Left, title, input, list)

	return Card(card).Render(content)
}
