package styles

import (
	"charm.land/bubbles/v2/help"
	lipgloss "charm.land/lipgloss/v2"
)

// NewHelpModel creates a consistently styled help model for all pages
func NewHelpModel() help.Model {
	h := help.New()

	h.Styles = help.Styles{
		ShortKey:       DimStyle.Copy(),
		ShortDesc:      MutedStyle.Copy(),
		ShortSeparator: lipgloss.NewStyle().Foreground(ColorDim),
		Ellipsis:       DimStyle.Copy(),
		FullKey:        DimStyle.Copy(),
		FullDesc:       MutedStyle.Copy(),
		FullSeparator:  lipgloss.NewStyle().Foreground(ColorDim),
	}

	return h
}
