package ui

import "github.com/charmbracelet/lipgloss"

// Enterprise/SaaS aesthetic colors
const (
	ColorBg        = lipgloss.Color("#0A0A0A") // Very dark gray/black
	ColorText      = lipgloss.Color("#EDEDED") // Off-white
	ColorSubtext   = lipgloss.Color("#A1A1AA") // Muted gray
	ColorBorder    = lipgloss.Color("#27272A") // Subtle border gray
	ColorAccent    = lipgloss.Color("#3B82F6") // Professional Blue (Tailwind Blue 500)
	ColorSuccess   = lipgloss.Color("#10B981") // Emerald Green
	ColorWarning   = lipgloss.Color("#F59E0B") // Amber
	ColorError     = lipgloss.Color("#EF4444") // Red
	ColorHighlight = lipgloss.Color("#1E293B") // Slate highlight
)

var (
	// Layout styles
	AppStyle = lipgloss.NewStyle().
			Margin(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorBorder).
			PaddingBottom(1).
			MarginBottom(1).
			Bold(true)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			MarginTop(1)

	// Typography
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	// Components
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2).
			MarginRight(1).
			MarginBottom(1)

	MetricLabelStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Bold(true)

	MetricValueStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true).
			MarginTop(1)

	MetricUnitStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			MarginLeft(1)

	// Status styles
	SpinnerStyle = lipgloss.NewStyle().Foreground(ColorAccent)
	ErrorStyle   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	
	// Menu styles
	ItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(ColorSubtext)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(ColorAccent).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(ColorAccent)
)
