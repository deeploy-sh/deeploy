package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Colors - Modern, muted palette
var (
	ColorPrimary = lipgloss.Color("#FF79C6")  // Dracula Pink - Brand
	ColorSuccess = lipgloss.Color("78")  // Grün - Online, Deployed
	ColorError   = lipgloss.Color("204") // Rosa/Red - Fehler
	ColorWarning = lipgloss.Color("214") // Orange - Pending, In Progress
	ColorMuted   = lipgloss.Color("241") // Grau - Secondary text
	ColorDim     = lipgloss.Color("238") // Dunkelgrau - Disabled
)

// Text Styles
var (
	// Status
	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess)
	ErrorStyle   = lipgloss.NewStyle().Foreground(ColorError)
	WarningStyle = lipgloss.NewStyle().Foreground(ColorWarning)

	// Interactive
	FocusedStyle = lipgloss.NewStyle().Foreground(ColorPrimary)
	MutedStyle   = lipgloss.NewStyle().Foreground(ColorMuted)
	DimStyle     = lipgloss.NewStyle().Foreground(ColorDim)

	// Semantic aliases
	OnlineStyle  = SuccessStyle
	OfflineStyle = ErrorStyle

	// Legacy (für Kompatibilität)
	BlurredStyle        = DimStyle
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = MutedStyle
	CursorModeHelpStyle = MutedStyle
	LabelStyle          = MutedStyle

	// Components
	AuthCard = lipgloss.NewStyle().
			Width(35).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder())

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", DimStyle.Render("Submit"))
)
