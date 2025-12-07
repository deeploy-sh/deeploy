package styles

import (
	"charm.land/bubbles/v2/help"
	lipgloss "charm.land/lipgloss/v2"
)

// NewHelpModel creates a consistently styled help model for all pages
func NewHelpModel() help.Model {
	h := help.New()

	h.Styles = help.Styles{
		ShortKey:       DimStyle(),
		ShortDesc:      MutedStyle(),
		ShortSeparator: lipgloss.NewStyle().Foreground(ColorDim()),
		Ellipsis:       DimStyle(),
		FullKey:        DimStyle(),
		FullDesc:       MutedStyle(),
		FullSeparator:  lipgloss.NewStyle().Foreground(ColorDim()),
	}

	return h
}
