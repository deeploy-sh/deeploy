package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// matrix implements the Matrix hacker theme
// Green on black, classic terminal aesthetic
type matrix struct{}

// Matrix returns the Matrix theme
func Matrix() Theme { return matrix{} }

func (t matrix) Name() string { return "matrix" }

// Backgrounds - deep blacks
func (t matrix) Background() color.Color        { return lipgloss.Color("#0a0e0a") }
func (t matrix) BackgroundPanel() color.Color   { return lipgloss.Color("#0e130d") }
func (t matrix) BackgroundElement() color.Color { return lipgloss.Color("#141c12") }

// Text - matrix green
func (t matrix) Foreground() color.Color      { return lipgloss.Color("#62ff94") }
func (t matrix) ForegroundMuted() color.Color { return lipgloss.Color("#8ca391") }
func (t matrix) ForegroundDim() color.Color   { return lipgloss.Color("#1e2a1b") }

// Semantic colors
func (t matrix) Primary() color.Color { return lipgloss.Color("#2eff6a") } // Rain green
func (t matrix) Success() color.Color { return lipgloss.Color("#62ff94") } // Bright green
func (t matrix) Warning() color.Color { return lipgloss.Color("#e6ff57") } // Alert yellow
func (t matrix) Error() color.Color   { return lipgloss.Color("#ff4b4b") } // Alert red

// AccentBorder - neon green
func (t matrix) AccentBorder() color.Color { return lipgloss.Color("#2eff6a") }
