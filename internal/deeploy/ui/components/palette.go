package components

import (
	"log"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
)

// PaletteItem represents a single item in the command palette
type PaletteItem struct {
	Title       string
	Description string
	Category    string         // "project", "pod", "action"
	Action      func() tea.Msg // Action to execute on selection
}

// FilterValue returns the filterable content for this item
func (i PaletteItem) FilterValue() string {
	return i.Title + " " + i.Description
}

// Palette is a command palette component
type Palette struct {
	textinput textinput.Model
	items     []PaletteItem
	filtered  []PaletteItem
	cursor    int
	width     int
	height    int
}

func NewPalette(items []PaletteItem) Palette {
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 100

	p := Palette{
		textinput: ti,
		items:     items,
		filtered:  items,
		cursor:    0,
	}

	return p
}

func (m Palette) Init() tea.Cmd {
	return textinput.Blink
}

func (m Palette) Update(msg tea.Msg) (Palette, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case msg.Code == tea.KeyUp || msg.String() == "ctrl+p":
			log.Println(len(m.filtered))
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.filtered) - 1
			}
			return m, nil
		case msg.Code == tea.KeyDown || msg.String() == "ctrl+n":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			return m, nil
		case msg.Code == tea.KeyEnter:
			if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
				item := m.filtered[m.cursor]
				if item.Action != nil {
					return m, item.Action
				}
			}
			return m, nil
		}
	}

	m.textinput, cmd = m.textinput.Update(msg)

	// Filter items based on input
	m.filterItems()

	return m, cmd
}

// filterItems filters the items based on the current input
func (m *Palette) filterItems() {
	query := strings.ToLower(m.textinput.Value())
	if query == "" {
		m.filtered = m.items
		return
	}

	filtered := make([]PaletteItem, 0)
	for _, item := range m.items {
		if strings.Contains(strings.ToLower(item.FilterValue()), query) {
			filtered = append(filtered, item)
		}
	}
	m.filtered = filtered

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
}

func (m *Palette) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m Palette) View() string {
	var b strings.Builder

	inputStyle := lipgloss.NewStyle().
		Width(m.width-4).
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary)

	b.WriteString(inputStyle.Render(m.textinput.View()))
	b.WriteString("\n\n")

	maxVisible := min(8, len(m.filtered))
	if maxVisible == 0 {
		b.WriteString(styles.MutedStyle.Render("  No results"))
	} else {
		for i := 0; i < maxVisible; i++ {
			item := m.filtered[i]

			var categoryBadge string
			switch item.Category {
			case "project":
				categoryBadge = lipgloss.NewStyle().
					Foreground(lipgloss.Color("33")).
					Render("[P]")
			case "pod":
				categoryBadge = lipgloss.NewStyle().
					Foreground(lipgloss.Color("35")).
					Render("[D]")
			case "action":
				categoryBadge = lipgloss.NewStyle().
					Foreground(lipgloss.Color("208")).
					Render("[A]")
			default:
				categoryBadge = "   "
			}

			var line string
			if i == m.cursor {
				line = lipgloss.NewStyle().
					Foreground(styles.ColorPrimary).
					Bold(true).
					Render("> " + categoryBadge + " " + item.Title)
			} else {
				line = "  " + categoryBadge + " " + item.Title
			}

			if item.Description != "" && i == m.cursor {
				line += styles.MutedStyle.Render(" - " + item.Description)
			}

			b.WriteString(line + "\n")
		}

		if len(m.filtered) > maxVisible {
			b.WriteString(styles.MutedStyle.Render(
				"\n  ..." + string(rune(len(m.filtered)-maxVisible)) + " more"))
		}
	}

	return b.String()
}

// // SelectedItem returns the currently selected item, or nil if none
// func (m Palette) SelectedItem() *PaletteItem {
// 	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
// 		return nil
// 	}
// 	return &m.filtered[m.cursor]
// }
