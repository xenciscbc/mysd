package output

import "github.com/charmbracelet/lipgloss"

// Color constants using lipgloss hex colors.
var (
	ColorSuccess = lipgloss.Color("#22c55e") // green
	ColorError   = lipgloss.Color("#ef4444") // red
	ColorWarning = lipgloss.Color("#f59e0b") // amber
	ColorInfo    = lipgloss.Color("#3b82f6") // blue
	ColorMuted   = lipgloss.Color("#6b7280") // gray
)

// Style definitions for each output type.
var (
	StyleSuccess = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	StyleError   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	StyleWarning = lipgloss.NewStyle().Foreground(ColorWarning)
	StyleInfo    = lipgloss.NewStyle().Foreground(ColorInfo)
	StyleMuted   = lipgloss.NewStyle().Foreground(ColorMuted)
	StyleHeader  = lipgloss.NewStyle().Bold(true).Underline(true)
)
