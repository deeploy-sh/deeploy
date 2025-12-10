package styles

import (
	lipgloss "charm.land/lipgloss/v2"
)

type CardSize int

const (
	CardSmall CardSize = iota
	CardMedium
	CardLarge
)

var cardWidths = map[CardSize]int{
	CardSmall:  40,
	CardMedium: 55,
	CardLarge:  75,
}

const cardPadding = 2

func Card(size CardSize, accent bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(cardWidths[size]).
		Padding(1, cardPadding).
		Background(ColorBackgroundPanel())

	if accent {
		style = style.
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderLeftForeground(ColorAccentBorder())
	}

	return style
}

func CardInner(size CardSize) int {
	// Subtract padding (both sides) and accent border (1 char)
	return cardWidths[size] - cardPadding*2 - 1
}
