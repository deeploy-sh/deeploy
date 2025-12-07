package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// onedark implements the One Dark theme
// Inspired by Atom's iconic One Dark theme
type onedark struct{}

// OneDark returns the One Dark theme
func OneDark() Theme { return onedark{} }

func (t onedark) Name() string { return "one-dark" }

// Backgrounds - cool grays
func (t onedark) Background() color.Color        { return lipgloss.Color("#282c34") }
func (t onedark) BackgroundPanel() color.Color   { return lipgloss.Color("#21252b") }
func (t onedark) BackgroundElement() color.Color { return lipgloss.Color("#353b45") }

// Text
func (t onedark) Foreground() color.Color      { return lipgloss.Color("#abb2bf") }
func (t onedark) ForegroundMuted() color.Color { return lipgloss.Color("#5c6370") }
func (t onedark) ForegroundDim() color.Color   { return lipgloss.Color("#393f4a") }

// Semantic colors
func (t onedark) Primary() color.Color { return lipgloss.Color("#61afef") } // Blue
func (t onedark) Success() color.Color { return lipgloss.Color("#98c379") } // Green
func (t onedark) Warning() color.Color { return lipgloss.Color("#e5c07b") } // Yellow
func (t onedark) Error() color.Color   { return lipgloss.Color("#e06c75") } // Red

// AccentBorder - purple
func (t onedark) AccentBorder() color.Color { return lipgloss.Color("#c678dd") }
