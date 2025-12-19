package components

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

// ScrollItem interface für Items in der Liste
type ScrollItem interface {
	Title() string
	FilterValue() string
}

// PrefixedItem optionales Interface für Items mit Prefix (●, [P], etc.)
type PrefixedItem interface {
	Prefix() string
}

// SuffixedItem optionales Interface für Items mit Suffix (count, status, etc.)
type SuffixedItem interface {
	Suffix() string
}

// ScrollListConfig für NewScrollList
type ScrollListConfig struct {
	Width       int
	Height      int
	WithInput   bool
	Placeholder string
}

// ScrollList ist eine einfache scrollbare Liste mit echtem Zeile-für-Zeile Scrolling
type ScrollList struct {
	allItems  []ScrollItem // original items (für filter)
	items     []ScrollItem // filtered items
	cursor    int          // selected item
	viewStart int          // erstes sichtbares item
	width     int
	height    int // anzahl sichtbarer items
	input     *textinput.Model
}

func NewScrollList(items []ScrollItem, cfg ScrollListConfig) ScrollList {
	l := ScrollList{
		allItems: items,
		items:    items,
		width:    cfg.Width,
		height:   cfg.Height,
	}

	if cfg.WithInput {
		ti := NewTextInput(20)
		ti.Placeholder = cfg.Placeholder
		if ti.Placeholder == "" {
			ti.Placeholder = "Type to search..."
		}
		ti.Focus()
		ti.CharLimit = 100
		l.input = &ti
	}

	return l
}

func (m *ScrollList) CursorUp() {
	m.cursor--
	if m.cursor < 0 {
		// Wrap to end
		m.cursor = len(m.items) - 1
		m.viewStart = max(0, len(m.items)-m.height)
	} else if m.cursor < m.viewStart {
		// Scroll up
		m.viewStart = m.cursor
	}
}

func (m *ScrollList) CursorDown() {
	m.cursor++
	if m.cursor >= len(m.items) {
		// Wrap to start
		m.cursor = 0
		m.viewStart = 0
	} else if m.cursor >= m.viewStart+m.height {
		// Scroll down
		m.viewStart = m.cursor - m.height + 1
	}
}

func (m ScrollList) SelectedItem() ScrollItem {
	if m.cursor >= 0 && m.cursor < len(m.items) {
		return m.items[m.cursor]
	}
	return nil
}

func (m ScrollList) Index() int          { return m.cursor }
func (m ScrollList) Width() int          { return m.width }
func (m ScrollList) Height() int         { return m.height }
func (m ScrollList) Items() []ScrollItem { return m.items }

func (m *ScrollList) SetWidth(w int)  { m.width = w }
func (m *ScrollList) SetHeight(h int) { m.height = h }

func (m *ScrollList) Select(index int) {
	if index >= 0 && index < len(m.items) {
		m.cursor = index
		// Adjust viewStart to make cursor visible
		if m.cursor < m.viewStart {
			m.viewStart = m.cursor
		} else if m.cursor >= m.viewStart+m.height {
			m.viewStart = m.cursor - m.height + 1
		}
	}
}

func (m *ScrollList) SetItems(items []ScrollItem) {
	m.allItems = items
	m.items = items
	if m.cursor >= len(items) {
		m.cursor = max(0, len(items)-1)
	}
	if m.viewStart > m.cursor {
		m.viewStart = m.cursor
	}
}

func (m ScrollList) Init() tea.Cmd {
	if m.input != nil {
		return textinput.Blink
	}
	return nil
}

func (m *ScrollList) filter() {
	if m.input == nil {
		return
	}

	query := strings.ToLower(m.input.Value())
	if query == "" {
		m.items = m.allItems
	} else {
		var filtered []ScrollItem
		for _, item := range m.allItems {
			if strings.Contains(strings.ToLower(item.FilterValue()), query) {
				filtered = append(filtered, item)
			}
		}
		m.items = filtered
	}

	// Reset cursor wenn nötig
	if m.cursor >= len(m.items) {
		m.cursor = max(0, len(m.items)-1)
	}
	m.viewStart = 0
}

func (m ScrollList) Update(msg tea.Msg) (ScrollList, tea.Cmd) {
	var cmd tea.Cmd

	// Input updaten wenn vorhanden (für Blink und Typing)
	if m.input != nil {
		*m.input, cmd = m.input.Update(msg)
		m.filter()
	}

	// Navigation (vim-style + mouse)
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()

		// Mit Input: Ctrl+P/N, Tab/Shift+Tab, Pfeiltasten
		// Ohne Input: zusätzlich j/k
		isUp := msg.Code == tea.KeyUp || key == "ctrl+p" || key == "shift+tab" || (m.input == nil && key == "k")
		isDown := msg.Code == tea.KeyDown || key == "ctrl+n" || key == "tab" || (m.input == nil && key == "j")

		switch {
		case isUp:
			m.CursorUp()
		case isDown:
			m.CursorDown()
		}

	case tea.MouseWheelMsg:
		if msg.Button == tea.MouseWheelUp {
			m.CursorUp()
		} else if msg.Button == tea.MouseWheelDown {
			m.CursorDown()
		}
	}

	return m, cmd
}

func (m ScrollList) HasInput() bool {
	return m.input != nil
}

func (m ScrollList) InputView() string {
	if m.input == nil {
		return ""
	}
	return lipgloss.NewStyle().
		Width(m.width).
		Background(styles.ColorBackgroundPanel()).
		PaddingLeft(1).
		PaddingBottom(1).
		Render(m.input.View())
}

func (m ScrollList) View() string {
	var lines []string

	// Items rendern
	if len(m.items) == 0 {
		empty := lipgloss.NewStyle().
			Width(m.width).
			Background(styles.ColorBackgroundPanel()).
			Render(styles.MutedStyle().Render(" No results"))
		lines = append(lines, empty)
	} else {
		end := min(m.viewStart+m.height, len(m.items))
		for i := m.viewStart; i < end; i++ {
			item := m.items[i]
			selected := i == m.cursor

			// Get prefix if item implements PrefixedItem
			prefix := ""
			pi, ok := item.(PrefixedItem)
			if ok {
				prefix = pi.Prefix() + " "
			}

			// Get suffix if item implements SuffixedItem
			suffix := ""
			si, ok := item.(SuffixedItem)
			if ok {
				suffix = si.Suffix()
			}

			title := item.Title()
			// Calculate space-between padding (1 leading + 1 trailing space)
			usedWidth := 2 + len(prefix) + len(title) + len(suffix)
			padding := max(1, m.width-usedWidth)

			content := " " + prefix + title + strings.Repeat(" ", padding) + suffix + " "
			lineStyle := lipgloss.NewStyle().Width(m.width)

			var line string
			if selected {
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

			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n")
}
