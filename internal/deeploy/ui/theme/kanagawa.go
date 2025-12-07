package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// kanagawa implements the Kanagawa theme
// Inspired by the famous painting by Katsushika Hokusai
type kanagawa struct{}

// Kanagawa returns the Kanagawa theme
func Kanagawa() Theme { return kanagawa{} }

func (t kanagawa) Name() string { return "kanagawa" }

// Backgrounds - sumi ink
func (t kanagawa) Background() color.Color        { return lipgloss.Color("#1F1F28") }
func (t kanagawa) BackgroundPanel() color.Color   { return lipgloss.Color("#2A2A37") }
func (t kanagawa) BackgroundElement() color.Color { return lipgloss.Color("#363646") }

// Text - fuji white
func (t kanagawa) Foreground() color.Color      { return lipgloss.Color("#DCD7BA") }
func (t kanagawa) ForegroundMuted() color.Color { return lipgloss.Color("#727169") }
func (t kanagawa) ForegroundDim() color.Color   { return lipgloss.Color("#54546D") }

// Semantic colors
func (t kanagawa) Primary() color.Color { return lipgloss.Color("#7E9CD8") } // Crystal blue
func (t kanagawa) Success() color.Color { return lipgloss.Color("#98BB6C") } // Lotus green
func (t kanagawa) Warning() color.Color { return lipgloss.Color("#D7A657") } // Ronin yellow
func (t kanagawa) Error() color.Color   { return lipgloss.Color("#E82424") } // Dragon red

// AccentBorder - oni violet
func (t kanagawa) AccentBorder() color.Color { return lipgloss.Color("#957FB8") }
