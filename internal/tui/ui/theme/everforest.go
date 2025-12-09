package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// everforest implements the Everforest Dark theme
// A comfortable, green-based color scheme
type everforest struct{}

// Everforest returns the Everforest Dark theme
func Everforest() Theme { return everforest{} }

func (t everforest) Name() string { return "everforest" }

// Backgrounds - forest greens
func (t everforest) Background() color.Color        { return lipgloss.Color("#2d353b") }
func (t everforest) BackgroundPanel() color.Color   { return lipgloss.Color("#333c43") }
func (t everforest) BackgroundElement() color.Color { return lipgloss.Color("#3d484d") }

// Text
func (t everforest) Foreground() color.Color      { return lipgloss.Color("#d3c6aa") }
func (t everforest) ForegroundMuted() color.Color { return lipgloss.Color("#7a8478") }
func (t everforest) ForegroundDim() color.Color   { return lipgloss.Color("#475258") }

// Semantic colors
func (t everforest) Primary() color.Color { return lipgloss.Color("#a7c080") } // Green
func (t everforest) Success() color.Color { return lipgloss.Color("#83c092") } // Aqua
func (t everforest) Warning() color.Color { return lipgloss.Color("#dbbc7f") } // Yellow
func (t everforest) Error() color.Color   { return lipgloss.Color("#e67e80") } // Red

// AccentBorder - soft green
func (t everforest) AccentBorder() color.Color { return lipgloss.Color("#a7c080") }
