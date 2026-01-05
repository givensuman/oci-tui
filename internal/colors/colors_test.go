package colors

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/givensuman/containertui/internal/config"
	"github.com/givensuman/containertui/internal/context"
)

func TestANSIColors(t *testing.T) {
	// Test that ANSI color constants work
	if ColorBlack.String() != "0" {
		t.Errorf("expected ColorBlack to be '0', got %s", ColorBlack.String())
	}
	if ColorBrightYellow.String() != "11" {
		t.Errorf("expected ColorBrightYellow to be '11', got %s", ColorBrightYellow.String())
	}
}

func TestColorFunctions(t *testing.T) {
	// Set up default config
	context.SetConfig(config.DefaultConfig())

	// Test default colors
	yellow := Yellow()
	if yellow == "" {
		t.Error("Yellow() returned empty color")
	}

	green := Green()
	if green == "" {
		t.Error("Green() returned empty color")
	}

	red := Red()
	if red == "" {
		t.Error("Red() returned empty color")
	}

	blue := Blue()
	if blue == "" {
		t.Error("Blue() returned empty color")
	}

	primary := Primary()
	if primary == "" {
		t.Error("Primary() returned empty color")
	}

	// Test that Primary defaults to Blue
	if primary != blue {
		t.Error("Primary() should default to Blue()")
	}
}

func TestColorOverrides(t *testing.T) {
	// Set up config with color overrides
	cfg := config.DefaultConfig()
	cfg.Colors.Yellow = "#FFFF00"
	cfg.Colors.Green = "#00FF00"
	context.SetConfig(cfg)

	// Test overridden colors
	yellow := Yellow()
	if yellow != lipgloss.Color("#FFFF00") {
		t.Errorf("expected Yellow to be '#FFFF00', got %s", yellow)
	}

	green := Green()
	if green != lipgloss.Color("#00FF00") {
		t.Errorf("expected Green to be '#00FF00', got %s", green)
	}

	// Test that non-overridden colors still use defaults
	blue := Blue()
	if blue == "" {
		t.Error("Blue() should still return default color when not overridden")
	}
}
