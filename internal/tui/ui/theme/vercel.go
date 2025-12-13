package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// vercel implements the Vercel theme
// Minimal black and white with blue accents
type vercel struct{}

// Vercel returns the Vercel theme
func Vercel() Theme { return vercel{} }

func (t vercel) Name() string { return "vercel" }

// Backgrounds - pure blacks
func (t vercel) Background() color.Color        { return lipgloss.Color("#000000") }
func (t vercel) BackgroundPanel() color.Color   { return lipgloss.Color("#1A1A1A") }
func (t vercel) BackgroundElement() color.Color { return lipgloss.Color("#292929") }

// Text
func (t vercel) Foreground() color.Color      { return lipgloss.Color("#EDEDED") }
func (t vercel) ForegroundMuted() color.Color { return lipgloss.Color("#878787") }
func (t vercel) ForegroundDim() color.Color   { return lipgloss.Color("#1F1F1F") }

// Semantic colors
func (t vercel) Primary() color.Color { return lipgloss.Color("#0070F3") } // Vercel blue
func (t vercel) Success() color.Color { return lipgloss.Color("#46A758") } // Green
func (t vercel) Warning() color.Color { return lipgloss.Color("#FFB224") } // Amber
func (t vercel) Error() color.Color   { return lipgloss.Color("#E5484D") } // Red

// AccentBorder - blue
func (t vercel) AccentBorder() color.Color { return lipgloss.Color("#0070F3") }
