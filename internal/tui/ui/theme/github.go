package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// github implements the GitHub Dark theme
// The official GitHub dark color scheme
type github struct{}

// GitHub returns the GitHub Dark theme
func GitHub() Theme { return github{} }

func (t github) Name() string { return "github" }

// Backgrounds
func (t github) Background() color.Color        { return lipgloss.Color("#0d1117") }
func (t github) BackgroundPanel() color.Color   { return lipgloss.Color("#010409") }
func (t github) BackgroundElement() color.Color { return lipgloss.Color("#161b22") }

// Text
func (t github) Foreground() color.Color      { return lipgloss.Color("#c9d1d9") }
func (t github) ForegroundMuted() color.Color { return lipgloss.Color("#8b949e") }
func (t github) ForegroundDim() color.Color   { return lipgloss.Color("#30363d") }

// Semantic colors
func (t github) Primary() color.Color { return lipgloss.Color("#58a6ff") } // Blue
func (t github) Success() color.Color { return lipgloss.Color("#3fb950") } // Green
func (t github) Warning() color.Color { return lipgloss.Color("#e3b341") } // Yellow
func (t github) Error() color.Color   { return lipgloss.Color("#f85149") } // Red

// AccentBorder - purple
func (t github) AccentBorder() color.Color { return lipgloss.Color("#bc8cff") }
