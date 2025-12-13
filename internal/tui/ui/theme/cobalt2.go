package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// cobalt2 implements the Cobalt2 theme
// A bold blue theme with vibrant accent colors
type cobalt2 struct{}

// Cobalt2 returns the Cobalt2 theme
func Cobalt2() Theme { return cobalt2{} }

func (t cobalt2) Name() string { return "cobalt2" }

// Backgrounds - deep blue
func (t cobalt2) Background() color.Color        { return lipgloss.Color("#193549") }
func (t cobalt2) BackgroundPanel() color.Color   { return lipgloss.Color("#122738") }
func (t cobalt2) BackgroundElement() color.Color { return lipgloss.Color("#1f4662") }

// Text
func (t cobalt2) Foreground() color.Color      { return lipgloss.Color("#ffffff") }
func (t cobalt2) ForegroundMuted() color.Color { return lipgloss.Color("#adb7c9") }
func (t cobalt2) ForegroundDim() color.Color   { return lipgloss.Color("#1f4662") }

// Semantic colors
func (t cobalt2) Primary() color.Color { return lipgloss.Color("#0088ff") } // Blue
func (t cobalt2) Success() color.Color { return lipgloss.Color("#9eff80") } // Green
func (t cobalt2) Warning() color.Color { return lipgloss.Color("#ffc600") } // Yellow
func (t cobalt2) Error() color.Color   { return lipgloss.Color("#ff0088") } // Red/Pink

// AccentBorder - mint cyan
func (t cobalt2) AccentBorder() color.Color { return lipgloss.Color("#2affdf") }
