package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// zenburn implements the Zenburn theme
// Low-contrast theme designed for long coding sessions
type zenburn struct{}

// Zenburn returns the Zenburn theme
func Zenburn() Theme { return zenburn{} }

func (t zenburn) Name() string { return "zenburn" }

// Backgrounds - warm grey
func (t zenburn) Background() color.Color        { return lipgloss.Color("#3f3f3f") }
func (t zenburn) BackgroundPanel() color.Color   { return lipgloss.Color("#4f4f4f") }
func (t zenburn) BackgroundElement() color.Color { return lipgloss.Color("#5f5f5f") }

// Text - cream
func (t zenburn) Foreground() color.Color      { return lipgloss.Color("#dcdccc") }
func (t zenburn) ForegroundMuted() color.Color { return lipgloss.Color("#9f9f9f") }
func (t zenburn) ForegroundDim() color.Color   { return lipgloss.Color("#5f5f5f") }

// Semantic colors - muted
func (t zenburn) Primary() color.Color { return lipgloss.Color("#8cd0d3") } // Soft blue
func (t zenburn) Success() color.Color { return lipgloss.Color("#7f9f7f") } // Soft green
func (t zenburn) Warning() color.Color { return lipgloss.Color("#f0dfaf") } // Soft yellow
func (t zenburn) Error() color.Color   { return lipgloss.Color("#cc9393") } // Soft red

// AccentBorder - magenta
func (t zenburn) AccentBorder() color.Color { return lipgloss.Color("#dc8cc3") }
