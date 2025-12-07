package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// nord implements the Nord theme
// An arctic, north-bluish color palette
type nord struct{}

// Nord returns the Nord theme
func Nord() Theme { return nord{} }

func (t nord) Name() string { return "nord" }

// Backgrounds - polar night
func (t nord) Background() color.Color        { return lipgloss.Color("#2E3440") }
func (t nord) BackgroundPanel() color.Color   { return lipgloss.Color("#3B4252") }
func (t nord) BackgroundElement() color.Color { return lipgloss.Color("#434C5E") }

// Text - snow storm
func (t nord) Foreground() color.Color      { return lipgloss.Color("#ECEFF4") }
func (t nord) ForegroundMuted() color.Color { return lipgloss.Color("#8B95A7") }
func (t nord) ForegroundDim() color.Color   { return lipgloss.Color("#4C566A") }

// Semantic colors - frost & aurora
func (t nord) Primary() color.Color { return lipgloss.Color("#88C0D0") } // Frost cyan
func (t nord) Success() color.Color { return lipgloss.Color("#A3BE8C") } // Aurora green
func (t nord) Warning() color.Color { return lipgloss.Color("#EBCB8B") } // Aurora yellow
func (t nord) Error() color.Color   { return lipgloss.Color("#BF616A") } // Aurora red

// AccentBorder - frost blue
func (t nord) AccentBorder() color.Color { return lipgloss.Color("#81A1C1") }
