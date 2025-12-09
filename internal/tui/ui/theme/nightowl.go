package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// nightowl implements the Night Owl theme
// A theme for night owls, fine-tuned for low-light conditions
type nightowl struct{}

// NightOwl returns the Night Owl theme
func NightOwl() Theme { return nightowl{} }

func (t nightowl) Name() string { return "nightowl" }

// Backgrounds - deep navy blue
func (t nightowl) Background() color.Color        { return lipgloss.Color("#011627") }
func (t nightowl) BackgroundPanel() color.Color   { return lipgloss.Color("#0b253a") }
func (t nightowl) BackgroundElement() color.Color { return lipgloss.Color("#0b253a") }

// Text
func (t nightowl) Foreground() color.Color      { return lipgloss.Color("#d6deeb") }
func (t nightowl) ForegroundMuted() color.Color { return lipgloss.Color("#5f7e97") }
func (t nightowl) ForegroundDim() color.Color   { return lipgloss.Color("#0b253a") }

// Semantic colors
func (t nightowl) Primary() color.Color { return lipgloss.Color("#82AAFF") } // Blue
func (t nightowl) Success() color.Color { return lipgloss.Color("#c5e478") } // Green
func (t nightowl) Warning() color.Color { return lipgloss.Color("#ecc48d") } // Yellow
func (t nightowl) Error() color.Color   { return lipgloss.Color("#EF5350") } // Red

// AccentBorder - purple
func (t nightowl) AccentBorder() color.Color { return lipgloss.Color("#c792ea") }
