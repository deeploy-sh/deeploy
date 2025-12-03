package components

import lipgloss "charm.land/lipgloss/v2"

func Centered(w, h int, content string) string {
	return lipgloss.Place(
		w,
		h,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
