package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// solarized implements the Solarized Dark theme
// Precision colors for machines and people
type solarized struct{}

// Solarized returns the Solarized Dark theme
func Solarized() Theme { return solarized{} }

func (t solarized) Name() string { return "solarized" }

// Backgrounds - base03/base02
func (t solarized) Background() color.Color        { return lipgloss.Color("#002b36") }
func (t solarized) BackgroundPanel() color.Color   { return lipgloss.Color("#073642") }
func (t solarized) BackgroundElement() color.Color { return lipgloss.Color("#073642") }

// Text - base0/base1
func (t solarized) Foreground() color.Color      { return lipgloss.Color("#839496") }
func (t solarized) ForegroundMuted() color.Color { return lipgloss.Color("#586e75") }
func (t solarized) ForegroundDim() color.Color   { return lipgloss.Color("#073642") }

// Semantic colors
func (t solarized) Primary() color.Color { return lipgloss.Color("#268bd2") } // Blue
func (t solarized) Success() color.Color { return lipgloss.Color("#859900") } // Green
func (t solarized) Warning() color.Color { return lipgloss.Color("#b58900") } // Yellow
func (t solarized) Error() color.Color   { return lipgloss.Color("#dc322f") } // Red

// AccentBorder - cyan
func (t solarized) AccentBorder() color.Color { return lipgloss.Color("#2aa198") }
