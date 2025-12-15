package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// deeploy implements the Deeploy theme
// Deep indigo-black backgrounds with electric violet accents
// A unique combination that stands out from other purple themes
type deeploy struct{}

// Deeploy returns the Deeploy theme
func Deeploy() Theme { return deeploy{} }

func (t deeploy) Name() string { return "deeploy" }

// Backgrounds - deep indigo-black
func (t deeploy) Background() color.Color        { return lipgloss.Color("#08070d") }
func (t deeploy) BackgroundPanel() color.Color   { return lipgloss.Color("#12101a") }
func (t deeploy) BackgroundElement() color.Color { return lipgloss.Color("#1c1928") }

// Text - slightly bluish-white for better contrast with violet
func (t deeploy) Foreground() color.Color      { return lipgloss.Color("#e8e6f2") }
func (t deeploy) ForegroundMuted() color.Color { return lipgloss.Color("#9490b0") }
func (t deeploy) ForegroundDim() color.Color   { return lipgloss.Color("#5a5673") }

// Semantic colors - electric violet palette
func (t deeploy) Primary() color.Color { return lipgloss.Color("#a855f7") } // Electric Violet
func (t deeploy) Success() color.Color { return lipgloss.Color("#4ade80") } // Fresh Green
func (t deeploy) Warning() color.Color { return lipgloss.Color("#fbbf24") } // Amber Gold
func (t deeploy) Error() color.Color   { return lipgloss.Color("#f87171") } // Soft Red

// AccentBorder - light violet for left borders
func (t deeploy) AccentBorder() color.Color { return lipgloss.Color("#c084fc") }
