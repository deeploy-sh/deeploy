package components

import (
	"image/color"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/deeploy-sh/deeploy/internal/tui/ui/styles"
)

type CardProps struct {
	Width       int
	Height      int
	Padding     []int
	Accent      bool        // Show left accent border
	AccentColor color.Color // Optional: override accent color
}

// InnerWidth berechnet die verfügbare Breite innerhalb der Card
func (p CardProps) InnerWidth() int {
	inner := p.Width
	if len(p.Padding) > 1 {
		inner -= p.Padding[1] * 2 // horizontal padding (links + rechts)
	}
	if p.Accent {
		inner -= 1 // Accent border
	}
	return inner
}

// Card creates a card style with panel background and optional left accent border
// OpenCode style: no rounded borders, separation via background color
func Card(p CardProps) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(p.Width).
		Background(styles.ColorBackgroundPanel())

	if p.Height > 0 {
		style = style.Height(p.Height)
	}

	if len(p.Padding) > 0 {
		style = style.Padding(p.Padding...)
	}

	// Left accent border (OpenCode style) - thin normal border
	if p.Accent {
		accentColor := p.AccentColor
		if accentColor == nil {
			accentColor = styles.ColorAccentBorder()
		}
		style = style.
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderLeftForeground(accentColor)
	}

	return style
}

// AccentCard creates a card with left accent border (default primary color)
func AccentCard(width int) lipgloss.Style {
	return Card(CardProps{Width: width, Padding: []int{1, 2}, Accent: true})
}

// ErrorCard creates a card with red left accent border
func ErrorCard(width int) lipgloss.Style {
	return Card(CardProps{
		Width:       width,
		Padding:     []int{1, 2},
		Accent:      true,
		AccentColor: styles.ColorError(),
	})
}

// SuccessCard creates a card with green left accent border
func SuccessCard(width int) lipgloss.Style {
	return Card(CardProps{
		Width:       width,
		Padding:     []int{1, 2},
		Accent:      true,
		AccentColor: styles.ColorSuccess(),
	})
}

// WarningCard creates a card with yellow left accent border
func WarningCard(width int) lipgloss.Style {
	return Card(CardProps{
		Width:       width,
		Padding:     []int{1, 2},
		Accent:      true,
		AccentColor: styles.ColorWarning(),
	})
}
