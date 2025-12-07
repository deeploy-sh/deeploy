package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// catppuccin implements the Catppuccin Mocha theme
// A soothing pastel theme with warm colors
type catppuccin struct{}

// Catppuccin returns the Catppuccin Mocha theme
func Catppuccin() Theme { return catppuccin{} }

func (t catppuccin) Name() string { return "catppuccin" }

// Backgrounds - deep purple-blue darks
func (t catppuccin) Background() color.Color        { return lipgloss.Color("#1e1e2e") }
func (t catppuccin) BackgroundPanel() color.Color   { return lipgloss.Color("#181825") }
func (t catppuccin) BackgroundElement() color.Color { return lipgloss.Color("#11111b") }

// Text
func (t catppuccin) Foreground() color.Color      { return lipgloss.Color("#cdd6f4") }
func (t catppuccin) ForegroundMuted() color.Color { return lipgloss.Color("#bac2de") }
func (t catppuccin) ForegroundDim() color.Color   { return lipgloss.Color("#6c7086") }

// Semantic colors
func (t catppuccin) Primary() color.Color { return lipgloss.Color("#89b4fa") } // Blue
func (t catppuccin) Success() color.Color { return lipgloss.Color("#a6e3a1") } // Green
func (t catppuccin) Warning() color.Color { return lipgloss.Color("#f9e2af") } // Yellow
func (t catppuccin) Error() color.Color   { return lipgloss.Color("#f38ba8") } // Red

// AccentBorder - blue accent
func (t catppuccin) AccentBorder() color.Color { return lipgloss.Color("#89b4fa") }
