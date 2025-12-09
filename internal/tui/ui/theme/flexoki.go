package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// flexoki implements the Flexoki theme
// An inky color scheme inspired by traditional printing
type flexoki struct{}

// Flexoki returns the Flexoki theme
func Flexoki() Theme { return flexoki{} }

func (t flexoki) Name() string { return "flexoki" }

// Backgrounds - warm ink blacks
func (t flexoki) Background() color.Color        { return lipgloss.Color("#100F0F") }
func (t flexoki) BackgroundPanel() color.Color   { return lipgloss.Color("#1C1B1A") }
func (t flexoki) BackgroundElement() color.Color { return lipgloss.Color("#282726") }

// Text - paper tones
func (t flexoki) Foreground() color.Color      { return lipgloss.Color("#CECDC3") }
func (t flexoki) ForegroundMuted() color.Color { return lipgloss.Color("#6F6E69") }
func (t flexoki) ForegroundDim() color.Color   { return lipgloss.Color("#403E3C") }

// Semantic colors
func (t flexoki) Primary() color.Color { return lipgloss.Color("#DA702C") } // Orange
func (t flexoki) Success() color.Color { return lipgloss.Color("#879A39") } // Green
func (t flexoki) Warning() color.Color { return lipgloss.Color("#D0A215") } // Yellow
func (t flexoki) Error() color.Color   { return lipgloss.Color("#D14D41") } // Red

// AccentBorder - cyan
func (t flexoki) AccentBorder() color.Color { return lipgloss.Color("#3AA99F") }
