package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// ayu implements the Ayu Dark theme
// A simple, bright and elegant theme
type ayu struct{}

// Ayu returns the Ayu Dark theme
func Ayu() Theme { return ayu{} }

func (t ayu) Name() string { return "ayu" }

// Backgrounds - deep dark blues
func (t ayu) Background() color.Color        { return lipgloss.Color("#0B0E14") }
func (t ayu) BackgroundPanel() color.Color   { return lipgloss.Color("#0D1017") }
func (t ayu) BackgroundElement() color.Color { return lipgloss.Color("#0F131A") }

// Text
func (t ayu) Foreground() color.Color      { return lipgloss.Color("#BFBDB6") }
func (t ayu) ForegroundMuted() color.Color { return lipgloss.Color("#565B66") }
func (t ayu) ForegroundDim() color.Color   { return lipgloss.Color("#11151C") }

// Semantic colors
func (t ayu) Primary() color.Color { return lipgloss.Color("#59C2FF") } // Entity blue
func (t ayu) Success() color.Color { return lipgloss.Color("#7FD962") } // Green
func (t ayu) Warning() color.Color { return lipgloss.Color("#E6B450") } // Accent yellow
func (t ayu) Error() color.Color   { return lipgloss.Color("#D95757") } // Error red

// AccentBorder - orange func
func (t ayu) AccentBorder() color.Color { return lipgloss.Color("#FFB454") }
