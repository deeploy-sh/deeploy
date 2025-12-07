package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// synthwave84 implements the Synthwave '84 theme
// Retro-futuristic neon colors inspired by 80s aesthetics
type synthwave84 struct{}

// Synthwave84 returns the Synthwave '84 theme
func Synthwave84() Theme { return synthwave84{} }

func (t synthwave84) Name() string { return "synthwave84" }

// Backgrounds - deep purple
func (t synthwave84) Background() color.Color        { return lipgloss.Color("#262335") }
func (t synthwave84) BackgroundPanel() color.Color   { return lipgloss.Color("#1e1a29") }
func (t synthwave84) BackgroundElement() color.Color { return lipgloss.Color("#2a2139") }

// Text
func (t synthwave84) Foreground() color.Color      { return lipgloss.Color("#ffffff") }
func (t synthwave84) ForegroundMuted() color.Color { return lipgloss.Color("#848bbd") }
func (t synthwave84) ForegroundDim() color.Color   { return lipgloss.Color("#495495") }

// Semantic colors - neon
func (t synthwave84) Primary() color.Color { return lipgloss.Color("#36f9f6") } // Cyan neon
func (t synthwave84) Success() color.Color { return lipgloss.Color("#72f1b8") } // Green neon
func (t synthwave84) Warning() color.Color { return lipgloss.Color("#fede5d") } // Yellow neon
func (t synthwave84) Error() color.Color   { return lipgloss.Color("#fe4450") } // Red neon

// AccentBorder - pink neon
func (t synthwave84) AccentBorder() color.Color { return lipgloss.Color("#ff7edb") }
