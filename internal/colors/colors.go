// Package colors provides color management for the application.
package colors

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/context"
)

// ANSI color constants
const (
	ColorBlack ANSIColor = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

type ANSIColor int

// String returns the ANSI color code
func (c ANSIColor) String() string {
	return ansiColorMap[c]
}

var ansiColorMap = map[ANSIColor]string{
	ColorBlack:         "0",
	ColorRed:           "1",
	ColorGreen:         "2",
	ColorYellow:        "3",
	ColorBlue:          "4",
	ColorMagenta:       "5",
	ColorCyan:          "6",
	ColorWhite:         "7",
	ColorBrightBlack:   "8",
	ColorBrightRed:     "9",
	ColorBrightGreen:   "10",
	ColorBrightYellow:  "11",
	ColorBrightBlue:    "12",
	ColorBrightMagenta: "13",
	ColorBrightCyan:    "14",
	ColorBrightWhite:   "15",
}

func Yellow() lipgloss.Color {
	cfg := context.GetConfig()
	if cfg.Colors.Yellow.IsAssigned() {
		return lipgloss.Color(cfg.Colors.Yellow)
	}

	return lipgloss.Color(ColorBrightYellow.String())
}

func Green() lipgloss.Color {
	cfg := context.GetConfig()
	if cfg.Colors.Green.IsAssigned() {
		return lipgloss.Color(cfg.Colors.Green)
	}

	return lipgloss.Color(ColorBrightGreen.String())
}

func Red() lipgloss.Color {
	cfg := context.GetConfig()
	if cfg.Colors.Red.IsAssigned() {
		return lipgloss.Color(cfg.Colors.Red)
	}

	return lipgloss.Color(ColorBrightRed.String())
}

func Blue() lipgloss.Color {
	cfg := context.GetConfig()
	if cfg.Colors.Blue.IsAssigned() {
		return lipgloss.Color(cfg.Colors.Blue)
	}

	return lipgloss.Color(ColorBrightBlue.String())
}

func Primary() lipgloss.Color {
	cfg := context.GetConfig()
	if cfg.Colors.Primary.IsAssigned() {
		return lipgloss.Color(cfg.Colors.Primary)
	}

	return Blue()
}
