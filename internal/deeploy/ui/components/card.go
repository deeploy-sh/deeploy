package components

import (
	"github.com/deeploy-sh/deeploy/internal/deeploy/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

type CardProps struct {
	Width            int
	Height           int
	Padding          []int
	BorderForeground lipgloss.TerminalColor
}

func Card(p CardProps) lipgloss.Style {
	baseStyle := lipgloss.NewStyle().
		Width(p.Width).
		Height(p.Height).
		Border(lipgloss.RoundedBorder())

	actualWidth := p.Width - baseStyle.GetHorizontalBorderSize()
	actualHeight := p.Height - baseStyle.GetVerticalBorderSize()

	return lipgloss.NewStyle().
		Width(actualWidth).
		Height(actualHeight).
		Padding(p.Padding...).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.BorderForeground)
}

func ErrorCard(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(styles.ColorError).
		Width(width).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())
}
