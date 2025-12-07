package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// monokai implements the classic Monokai theme
// The iconic theme from Sublime Text
type monokai struct{}

// Monokai returns the Monokai theme
func Monokai() Theme { return monokai{} }

func (t monokai) Name() string { return "monokai" }

// Backgrounds - warm dark
func (t monokai) Background() color.Color        { return lipgloss.Color("#272822") }
func (t monokai) BackgroundPanel() color.Color   { return lipgloss.Color("#1e1f1c") }
func (t monokai) BackgroundElement() color.Color { return lipgloss.Color("#3e3d32") }

// Text
func (t monokai) Foreground() color.Color      { return lipgloss.Color("#f8f8f2") }
func (t monokai) ForegroundMuted() color.Color { return lipgloss.Color("#75715e") }
func (t monokai) ForegroundDim() color.Color   { return lipgloss.Color("#3e3d32") }

// Semantic colors
func (t monokai) Primary() color.Color { return lipgloss.Color("#66d9ef") } // Cyan
func (t monokai) Success() color.Color { return lipgloss.Color("#a6e22e") } // Green
func (t monokai) Warning() color.Color { return lipgloss.Color("#e6db74") } // Yellow
func (t monokai) Error() color.Color   { return lipgloss.Color("#f92672") } // Pink/Red

// AccentBorder - purple
func (t monokai) AccentBorder() color.Color { return lipgloss.Color("#ae81ff") }
