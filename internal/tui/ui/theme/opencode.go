package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// opencode implements the OpenCode-inspired theme
// Dark background with warm gray undertones (smoke palette)
// Minimal borders, uses background color separation
// Left accent borders for cards
type opencode struct{}

// OpenCode returns the OpenCode-inspired theme
func OpenCode() Theme { return opencode{} }

func (t opencode) Name() string { return "opencode" }

// Backgrounds - warm dark grays (smoke palette from OpenCode)
func (t opencode) Background() color.Color        { return lipgloss.Color("#0A0A0A") }
func (t opencode) BackgroundPanel() color.Color   { return lipgloss.Color("#1b1818") }
func (t opencode) BackgroundElement() color.Color { return lipgloss.Color("#252121") }

// Text - warm light grays
func (t opencode) Foreground() color.Color      { return lipgloss.Color("#f1ecec") }
func (t opencode) ForegroundMuted() color.Color { return lipgloss.Color("#b7b1b1") }
func (t opencode) ForegroundDim() color.Color   { return lipgloss.Color("#716c6b") }

// Semantic colors - based on OpenCode palette
func (t opencode) Primary() color.Color { return lipgloss.Color("#89b5ff") } // Cobalt blue
func (t opencode) Success() color.Color { return lipgloss.Color("#9dde99") } // Mint green
func (t opencode) Warning() color.Color { return lipgloss.Color("#fdd63c") } // Solaris yellow
func (t opencode) Error() color.Color   { return lipgloss.Color("#ff917b") } // Ember red/coral

// AccentBorder - blue accent for left borders
func (t opencode) AccentBorder() color.Color { return lipgloss.Color("#89b5ff") }
