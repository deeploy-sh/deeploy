package components

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
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

// ScrollList ist eine einfache scrollbare Liste mit echtem Zeile-für-Zeile Scrolling
type ScrollList struct {
	items     []ScrollItem
	cursor    int // selected item
	viewStart int // erstes sichtbares item
	width     int
	height    int // anzahl sichtbarer items
}

func NewScrollList(items []ScrollItem, width, height int) ScrollList {
	return ScrollList{
		items:  items,
		width:  width,
		height: height,
	}
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
	m.items = items
	if m.cursor >= len(items) {
		m.cursor = max(0, len(items)-1)
	}
	if m.viewStart > m.cursor {
		m.viewStart = m.cursor
	}
}

func (m ScrollList) Update(msg tea.Msg) (ScrollList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyUp:
			m.CursorUp()
		case tea.KeyDown:
			m.CursorDown()
		}
	}
	return m, nil
}

func (m ScrollList) View() string {
	if len(m.items) == 0 {
		return ""
	}

	var lines []string
	end := min(m.viewStart+m.height, len(m.items))

	for i := m.viewStart; i < end; i++ {
		item := m.items[i]
		selected := i == m.cursor

		// Get prefix if item implements PrefixedItem
		prefix := ""
		if pi, ok := item.(PrefixedItem); ok {
			prefix = pi.Prefix() + " "
		}

		content := " " + prefix + item.Title()
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

	return strings.Join(lines, "\n")
}
