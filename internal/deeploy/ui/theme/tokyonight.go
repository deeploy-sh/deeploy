package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// tokyonight implements the Tokyo Night theme
// A clean, dark theme with blue accents
type tokyonight struct{}

// TokyoNight returns the Tokyo Night theme
func TokyoNight() Theme { return tokyonight{} }

func (t tokyonight) Name() string { return "tokyonight" }

// Backgrounds - deep blue-tinted darks
func (t tokyonight) Background() color.Color        { return lipgloss.Color("#1a1b26") }
func (t tokyonight) BackgroundPanel() color.Color   { return lipgloss.Color("#1e2030") }
func (t tokyonight) BackgroundElement() color.Color { return lipgloss.Color("#222436") }

// Text
func (t tokyonight) Foreground() color.Color      { return lipgloss.Color("#c8d3f5") }
func (t tokyonight) ForegroundMuted() color.Color { return lipgloss.Color("#828bb8") }
func (t tokyonight) ForegroundDim() color.Color   { return lipgloss.Color("#545c7e") }

// Semantic colors
func (t tokyonight) Primary() color.Color { return lipgloss.Color("#82aaff") } // Blue
func (t tokyonight) Success() color.Color { return lipgloss.Color("#c3e88d") } // Green
func (t tokyonight) Warning() color.Color { return lipgloss.Color("#ffc777") } // Yellow
func (t tokyonight) Error() color.Color   { return lipgloss.Color("#ff757f") } // Red

// AccentBorder - blue accent
func (t tokyonight) AccentBorder() color.Color { return lipgloss.Color("#82aaff") }
