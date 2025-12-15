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

// Current is the active theme - defaults to Deeploy
var Current Theme = Deeploy()

// Available themes registry (26 themes)
var Available = map[string]Theme{
	"aura":        Aura(),
	"ayu":         Ayu(),
	"catppuccin":  Catppuccin(),
	"cobalt2":     Cobalt2(),
	"dracula":     Dracula(),
	"everforest":  Everforest(),
	"flexoki":     Flexoki(),
	"github":      GitHub(),
	"gruvbox":     Gruvbox(),
	"kanagawa":    Kanagawa(),
	"material":    Material(),
	"matrix":      Matrix(),
	"mercury":     Mercury(),
	"monokai":     Monokai(),
	"nightowl":    NightOwl(),
	"nord":        Nord(),
	"one-dark":    OneDark(),
	"deeploy":     Deeploy(),
	"palenight":   Palenight(),
	"rose-pine":   RosePine(),
	"solarized":   Solarized(),
	"synthwave84": Synthwave84(),
	"tokyonight":  TokyoNight(),
	"vercel":      Vercel(),
	"vesper":      Vesper(),
	"zenburn":     Zenburn(),
}

// SetTheme switches the active theme by name
// Returns true if theme was found and set, false otherwise
func SetTheme(name string) bool {
	t, ok := Available[name]
	if ok {
		Current = t
		return true
	}
	return false
}

// ThemeNames returns an alphabetically sorted list of available theme names
func ThemeNames() []string {
	return []string{
		"aura",
		"ayu",
		"catppuccin",
		"cobalt2",
		"deeploy",
		"dracula",
		"everforest",
		"flexoki",
		"github",
		"gruvbox",
		"kanagawa",
		"material",
		"matrix",
		"mercury",
		"monokai",
		"nightowl",
		"nord",
		"one-dark",
		"palenight",
		"rose-pine",
		"solarized",
		"synthwave84",
		"tokyonight",
		"vercel",
		"vesper",
		"zenburn",
	}
}
