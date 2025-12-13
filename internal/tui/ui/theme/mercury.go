package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// mercury implements the Mercury theme
// Clean purple-focused design
type mercury struct{}

// Mercury returns the Mercury theme
func Mercury() Theme { return mercury{} }

func (t mercury) Name() string { return "mercury" }

// Backgrounds - deep neutral
func (t mercury) Background() color.Color        { return lipgloss.Color("#171721") }
func (t mercury) BackgroundPanel() color.Color   { return lipgloss.Color("#10101a") }
func (t mercury) BackgroundElement() color.Color { return lipgloss.Color("#272735") }

// Text
func (t mercury) Foreground() color.Color      { return lipgloss.Color("#dddde5") }
func (t mercury) ForegroundMuted() color.Color { return lipgloss.Color("#9d9da8") }
func (t mercury) ForegroundDim() color.Color   { return lipgloss.Color("#363644") }

// Semantic colors
func (t mercury) Primary() color.Color { return lipgloss.Color("#8da4f5") } // Purple-blue
func (t mercury) Success() color.Color { return lipgloss.Color("#77c599") } // Green
func (t mercury) Warning() color.Color { return lipgloss.Color("#fc9b6f") } // Orange
func (t mercury) Error() color.Color   { return lipgloss.Color("#fc92b4") } // Pink/red

// AccentBorder - purple
func (t mercury) AccentBorder() color.Color { return lipgloss.Color("#5266eb") }
