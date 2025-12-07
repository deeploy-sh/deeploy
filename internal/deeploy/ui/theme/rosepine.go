package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// rosepine implements the Rose Pine theme
// All natural pine, faux fur and a bit of soho vibes
type rosepine struct{}

// RosePine returns the Rose Pine theme
func RosePine() Theme { return rosepine{} }

func (t rosepine) Name() string { return "rose-pine" }

// Backgrounds - deep purples
func (t rosepine) Background() color.Color        { return lipgloss.Color("#191724") }
func (t rosepine) BackgroundPanel() color.Color   { return lipgloss.Color("#1f1d2e") }
func (t rosepine) BackgroundElement() color.Color { return lipgloss.Color("#26233a") }

// Text
func (t rosepine) Foreground() color.Color      { return lipgloss.Color("#e0def4") }
func (t rosepine) ForegroundMuted() color.Color { return lipgloss.Color("#6e6a86") }
func (t rosepine) ForegroundDim() color.Color   { return lipgloss.Color("#403d52") }

// Semantic colors
func (t rosepine) Primary() color.Color { return lipgloss.Color("#9ccfd8") } // Foam
func (t rosepine) Success() color.Color { return lipgloss.Color("#31748f") } // Pine
func (t rosepine) Warning() color.Color { return lipgloss.Color("#f6c177") } // Gold
func (t rosepine) Error() color.Color   { return lipgloss.Color("#eb6f92") } // Love

// AccentBorder - rose
func (t rosepine) AccentBorder() color.Color { return lipgloss.Color("#ebbcba") }
