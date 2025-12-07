package theme

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
)

// Theme defines the color and style contract for the TUI
type Theme interface {
	Name() string

	// Backgrounds
	Background() color.Color        // Main app background (darkest)
	BackgroundPanel() color.Color   // Cards, elevated surfaces (slightly lighter)
	BackgroundElement() color.Color // Inputs, buttons, interactive elements

	// Text
	Foreground() color.Color      // Primary text
	ForegroundMuted() color.Color // Secondary text
	ForegroundDim() color.Color   // Disabled/hints

	// Semantic colors
	Primary() color.Color // Brand/accent color
	Success() color.Color
	Warning() color.Color
	Error() color.Color

	// Accent border - used for left-side card accents
	AccentBorder() color.Color
}

// Color helper - shorthand for lipgloss.Color
func Color(s string) color.Color {
	return lipgloss.Color(s)
}

// Current is the active theme - defaults to OpenCode
var Current Theme = Dracula()

// Available themes registry
var Available = map[string]Theme{
	"opencode": OpenCode(),
	"dracula":  Dracula(),
}

// SetTheme switches the active theme by name
// Returns true if theme was found and set, false otherwise
func SetTheme(name string) bool {
	if t, ok := Available[name]; ok {
		Current = t
		return true
	}
	return false
}

// ThemeNames returns a list of available theme names
func ThemeNames() []string {
	names := make([]string, 0, len(Available))
	for name := range Available {
		names = append(names, name)
	}
	return names
}
