package styles

import (
	lipgloss "charm.land/lipgloss/v2"
)

// Card width constants for consistent sizing
const (
	CardWidthSM = 40 // Delete dialogs, simple confirmations
	CardWidthMD = 55 // Standard cards, info, settings, forms
	CardWidthLG = 70 // Complex forms, detail views
)

type CardProps struct {
	Width   int
	Height  int
	Padding []int
	Accent  bool
}

// InnerWidth calculates the available width inside the card.
func (p CardProps) InnerWidth() int {
	inner := p.Width
	if len(p.Padding) > 1 {
		inner -= p.Padding[1] * 2 // horizontal padding (left + right)
	}
	if p.Accent {
		inner -= 1 // accent border
	}
	return inner
}

func (p CardProps) InnerHeight() int {
	inner := p.Height
	if len(p.Padding) > 1 {
		inner -= p.Padding[0] * 2 // vertical padding (top + bottom)
	}
	return inner
}

// Card creates a card style with panel background and optional left accent border.
func Card(p CardProps) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(p.Width).
		Background(ColorBackgroundPanel())

	if p.Height > 0 {
		style = style.Height(p.Height).MaxHeight(p.Height)
	}

	if len(p.Padding) > 0 {
		style = style.Padding(p.Padding...)
	}

	if p.Accent {
		style = style.
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderLeftForeground(ColorAccentBorder())
	}

	return style
}
