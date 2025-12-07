package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// gruvbox implements the Gruvbox Dark theme
// A retro groove color scheme with earthy tones
type gruvbox struct{}

// Gruvbox returns the Gruvbox Dark theme
func Gruvbox() Theme { return gruvbox{} }

func (t gruvbox) Name() string { return "gruvbox" }

// Backgrounds - warm earthy darks
func (t gruvbox) Background() color.Color        { return lipgloss.Color("#282828") }
func (t gruvbox) BackgroundPanel() color.Color   { return lipgloss.Color("#3c3836") }
func (t gruvbox) BackgroundElement() color.Color { return lipgloss.Color("#504945") }

// Text
func (t gruvbox) Foreground() color.Color      { return lipgloss.Color("#ebdbb2") }
func (t gruvbox) ForegroundMuted() color.Color { return lipgloss.Color("#928374") }
func (t gruvbox) ForegroundDim() color.Color   { return lipgloss.Color("#665c54") }

// Semantic colors
func (t gruvbox) Primary() color.Color { return lipgloss.Color("#83a598") } // Blue
func (t gruvbox) Success() color.Color { return lipgloss.Color("#b8bb26") } // Green
func (t gruvbox) Warning() color.Color { return lipgloss.Color("#fabd2f") } // Yellow
func (t gruvbox) Error() color.Color   { return lipgloss.Color("#fb4934") } // Red

// AccentBorder - aqua accent
func (t gruvbox) AccentBorder() color.Color { return lipgloss.Color("#8ec07c") }
