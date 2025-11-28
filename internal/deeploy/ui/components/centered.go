package components

import "github.com/charmbracelet/lipgloss"

func Centered(w, h int, content string) string {
	return lipgloss.Place(
		w,
		h,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
