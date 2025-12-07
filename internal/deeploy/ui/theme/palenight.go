package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// palenight implements the Palenight theme
// Material Design inspired with soft purple tones
type palenight struct{}

// Palenight returns the Palenight theme
func Palenight() Theme { return palenight{} }

func (t palenight) Name() string { return "palenight" }

// Backgrounds - soft purple-grey
func (t palenight) Background() color.Color        { return lipgloss.Color("#292d3e") }
func (t palenight) BackgroundPanel() color.Color   { return lipgloss.Color("#1e2132") }
func (t palenight) BackgroundElement() color.Color { return lipgloss.Color("#32364a") }

// Text
func (t palenight) Foreground() color.Color      { return lipgloss.Color("#a6accd") }
func (t palenight) ForegroundMuted() color.Color { return lipgloss.Color("#676e95") }
func (t palenight) ForegroundDim() color.Color   { return lipgloss.Color("#32364a") }

// Semantic colors
func (t palenight) Primary() color.Color { return lipgloss.Color("#82aaff") } // Blue
func (t palenight) Success() color.Color { return lipgloss.Color("#c3e88d") } // Green
func (t palenight) Warning() color.Color { return lipgloss.Color("#ffcb6b") } // Yellow
func (t palenight) Error() color.Color   { return lipgloss.Color("#f07178") } // Red

// AccentBorder - purple
func (t palenight) AccentBorder() color.Color { return lipgloss.Color("#c792ea") }
