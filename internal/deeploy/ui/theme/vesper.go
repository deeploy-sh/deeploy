package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// vesper implements the Vesper theme
// Minimal dark theme with warm orange accents
type vesper struct{}

// Vesper returns the Vesper theme
func Vesper() Theme { return vesper{} }

func (t vesper) Name() string { return "vesper" }

// Backgrounds - deep black
func (t vesper) Background() color.Color        { return lipgloss.Color("#101010") }
func (t vesper) BackgroundPanel() color.Color   { return lipgloss.Color("#101010") }
func (t vesper) BackgroundElement() color.Color { return lipgloss.Color("#1C1C1C") }

// Text
func (t vesper) Foreground() color.Color      { return lipgloss.Color("#FFF") }
func (t vesper) ForegroundMuted() color.Color { return lipgloss.Color("#A0A0A0") }
func (t vesper) ForegroundDim() color.Color   { return lipgloss.Color("#282828") }

// Semantic colors - warm
func (t vesper) Primary() color.Color { return lipgloss.Color("#FFC799") } // Orange/peach
func (t vesper) Success() color.Color { return lipgloss.Color("#99FFE4") } // Mint
func (t vesper) Warning() color.Color { return lipgloss.Color("#FFC799") } // Orange
func (t vesper) Error() color.Color   { return lipgloss.Color("#FF8080") } // Soft red

// AccentBorder - peach
func (t vesper) AccentBorder() color.Color { return lipgloss.Color("#FFC799") }
