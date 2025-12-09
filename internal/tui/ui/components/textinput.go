package components

import (
	"charm.land/bubbles/v2/textinput"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

// NewTextInput creates a styled textinput with panel background.
// Optional width parameter sets the input width.
func NewTextInput(width ...int) textinput.Model {
	ti := textinput.New()
	ti.Prompt = ""

	bg := lipgloss.NewStyle().Background(styles.ColorBackgroundPanel())
	ti.SetStyles(textinput.Styles{
		Focused: textinput.StyleState{
			Text:        bg.Foreground(styles.ColorForeground()),
			Placeholder: bg.Foreground(styles.ColorMuted()),
		},
		Blurred: textinput.StyleState{
			Text:        bg.Foreground(styles.ColorForeground()),
			Placeholder: bg.Foreground(styles.ColorMuted()),
		},
		Cursor: textinput.CursorStyle{
			Blink: true,
		},
	})

	if len(width) > 0 {
		ti.SetWidth(width[0])
	}
	return ti
}
