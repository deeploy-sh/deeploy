package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// dracula implements the classic Dracula color theme
// Kept for backwards compatibility and as an alternative option
type dracula struct{}

// Dracula returns the Dracula theme
func Dracula() Theme { return dracula{} }

func (t dracula) Name() string { return "dracula" }

// Backgrounds - Dracula dark purples
func (t dracula) Background() color.Color        { return lipgloss.Color("#282a36") }
func (t dracula) BackgroundPanel() color.Color   { return lipgloss.Color("#343746") }
func (t dracula) BackgroundElement() color.Color { return lipgloss.Color("#44475a") }

// Text
func (t dracula) Foreground() color.Color      { return lipgloss.Color("#f8f8f2") }
func (t dracula) ForegroundMuted() color.Color { return lipgloss.Color("#6272a4") }
func (t dracula) ForegroundDim() color.Color   { return lipgloss.Color("#44475a") }

// Semantic colors - classic Dracula
func (t dracula) Primary() color.Color { return lipgloss.Color("#ff79c6") } // Pink
func (t dracula) Success() color.Color { return lipgloss.Color("#50fa7b") } // Green
func (t dracula) Warning() color.Color { return lipgloss.Color("#ffb86c") } // Orange
func (t dracula) Error() color.Color   { return lipgloss.Color("#ff5555") } // Red

// AccentBorder - Dracula purple
func (t dracula) AccentBorder() color.Color { return lipgloss.Color("#bd93f9") }
