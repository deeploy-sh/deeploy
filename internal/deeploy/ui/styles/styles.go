package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary = lipgloss.Color("204")
	ColorError   = lipgloss.Color("9")

	FocusedStyle        = lipgloss.NewStyle().Foreground(ColorPrimary)
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	ErrorStyle          = lipgloss.NewStyle().Foreground(ColorError)
	LabelStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	AuthCard = lipgloss.NewStyle().
			Width(35).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder())

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)
