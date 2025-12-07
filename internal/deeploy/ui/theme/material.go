package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// material implements the Material Dark theme
// Based on Google's Material Design color palette
type material struct{}

// Material returns the Material Dark theme
func Material() Theme { return material{} }

func (t material) Name() string { return "material" }

// Backgrounds - blue-grey
func (t material) Background() color.Color        { return lipgloss.Color("#263238") }
func (t material) BackgroundPanel() color.Color   { return lipgloss.Color("#1e272c") }
func (t material) BackgroundElement() color.Color { return lipgloss.Color("#37474f") }

// Text
func (t material) Foreground() color.Color      { return lipgloss.Color("#eeffff") }
func (t material) ForegroundMuted() color.Color { return lipgloss.Color("#546e7a") }
func (t material) ForegroundDim() color.Color   { return lipgloss.Color("#37474f") }

// Semantic colors
func (t material) Primary() color.Color { return lipgloss.Color("#82aaff") } // Blue
func (t material) Success() color.Color { return lipgloss.Color("#c3e88d") } // Green
func (t material) Warning() color.Color { return lipgloss.Color("#ffcb6b") } // Yellow
func (t material) Error() color.Color   { return lipgloss.Color("#f07178") } // Red

// AccentBorder - cyan
func (t material) AccentBorder() color.Color { return lipgloss.Color("#89ddff") }
