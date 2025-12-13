package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// aura implements the Aura theme
// A beautiful dark theme with purple neon accents
type aura struct{}

// Aura returns the Aura theme
func Aura() Theme { return aura{} }

func (t aura) Name() string { return "aura" }

// Backgrounds - deep dark
func (t aura) Background() color.Color        { return lipgloss.Color("#0f0f0f") }
func (t aura) BackgroundPanel() color.Color   { return lipgloss.Color("#15141b") }
func (t aura) BackgroundElement() color.Color { return lipgloss.Color("#15141b") }

// Text
func (t aura) Foreground() color.Color      { return lipgloss.Color("#edecee") }
func (t aura) ForegroundMuted() color.Color { return lipgloss.Color("#6d6d6d") }
func (t aura) ForegroundDim() color.Color   { return lipgloss.Color("#2d2d2d") }

// Semantic colors - neon
func (t aura) Primary() color.Color { return lipgloss.Color("#a277ff") } // Purple
func (t aura) Success() color.Color { return lipgloss.Color("#61ffca") } // Cyan/green
func (t aura) Warning() color.Color { return lipgloss.Color("#ffca85") } // Orange
func (t aura) Error() color.Color   { return lipgloss.Color("#ff6767") } // Red

// AccentBorder - pink
func (t aura) AccentBorder() color.Color { return lipgloss.Color("#f694ff") }
