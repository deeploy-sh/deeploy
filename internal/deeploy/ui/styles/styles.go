package styles

import (
	"fmt"
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/theme"
)

// Color accessors - always use current theme
func ColorPrimary() color.Color    { return theme.Current.Primary() }
func ColorForeground() color.Color { return theme.Current.Foreground() }
func ColorSuccess() color.Color    { return theme.Current.Success() }
func ColorError() color.Color      { return theme.Current.Error() }
func ColorWarning() color.Color    { return theme.Current.Warning() }
func ColorMuted() color.Color      { return theme.Current.ForegroundMuted() }
func ColorDim() color.Color        { return theme.Current.ForegroundDim() }

// Background colors
func ColorBackground() color.Color        { return theme.Current.Background() }
func ColorBackgroundPanel() color.Color   { return theme.Current.BackgroundPanel() }
func ColorBackgroundElement() color.Color { return theme.Current.BackgroundElement() }

// Accent border color
func ColorAccentBorder() color.Color { return theme.Current.AccentBorder() }

// Style factories - call these to get fresh styles with current theme
func ForegroundStyle() lipgloss.Style { return lipgloss.NewStyle().Foreground(ColorForeground()) }
func SuccessStyle() lipgloss.Style    { return lipgloss.NewStyle().Foreground(ColorSuccess()) }
func ErrorStyle() lipgloss.Style      { return lipgloss.NewStyle().Foreground(ColorError()) }
func WarningStyle() lipgloss.Style    { return lipgloss.NewStyle().Foreground(ColorWarning()) }
func PrimaryStyle() lipgloss.Style    { return lipgloss.NewStyle().Foreground(ColorPrimary()) }
func MutedStyle() lipgloss.Style      { return lipgloss.NewStyle().Foreground(ColorMuted()) }
func DimStyle() lipgloss.Style        { return lipgloss.NewStyle().Foreground(ColorDim()) }

// Semantic aliases
func OnlineStyle() lipgloss.Style  { return SuccessStyle() }
func OfflineStyle() lipgloss.Style { return ErrorStyle() }
func FocusedStyle() lipgloss.Style { return PrimaryStyle() }
func BlurredStyle() lipgloss.Style { return DimStyle() }

// Legacy aliases (for compatibility during migration)
func CursorStyle() lipgloss.Style         { return FocusedStyle() }
func HelpStyle() lipgloss.Style           { return MutedStyle() }
func CursorModeHelpStyle() lipgloss.Style { return MutedStyle() }
func LabelStyle() lipgloss.Style          { return MutedStyle() }
func NoStyle() lipgloss.Style             { return lipgloss.NewStyle() }

// Component helpers
func FocusedButton() string { return FocusedStyle().Render("[ Submit ]") }
func BlurredButton() string { return fmt.Sprintf("[ %s ]", DimStyle().Render("Submit")) }

// AuthCard returns a card style for auth pages
// Uses panel background with left accent border
func AuthCard() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(35).
		Padding(1, 2).
		Background(ColorBackgroundPanel()).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeftForeground(ColorAccentBorder())
}
