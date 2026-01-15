package colors

import (
	"fmt"
	"strings"

	"github.com/givensuman/containertui/internal/config"
)

// ParseColors parses color overrides from a slice of strings
// Format: ["primary=#b4befe'", "warning=#f9e2af", "success=#a6e3a1"]
func ParseColors(colorStrings []string) (*config.ThemeConfig, error) {
	if len(colorStrings) == 0 {
		return &config.ThemeConfig{}, nil
	}

	colorConfig := &config.ThemeConfig{}
	allPairs := []string{}

	// Collect all pairs from all strings
	for _, colorString := range colorStrings {
		colorString = strings.TrimSpace(colorString)
		if colorString == "" {
			continue
		}
		// Split each string by commas to handle cases where users might still use commas
		pairs := strings.Split(colorString, ",")
		allPairs = append(allPairs, pairs...)
	}

	for _, pair := range allPairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid color format: %s (expected key=value)", pair)
		}

		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		// Validate that value doesn't contain another '='
		if strings.Contains(value, "=") {
			return nil, fmt.Errorf("invalid color value: %s (values cannot contain '=')", value)
		}

		switch key {
		case "primary":
			colorConfig.Primary = config.ConfigString(value)
		case "border":
			colorConfig.Border = config.ConfigString(value)
		case "text":
			colorConfig.Text = config.ConfigString(value)
		case "muted":
			colorConfig.Muted = config.ConfigString(value)
		case "selected":
			colorConfig.Selected = config.ConfigString(value)
		case "success":
			colorConfig.Success = config.ConfigString(value)
		case "warning":
			colorConfig.Warning = config.ConfigString(value)
		case "error":
			colorConfig.Error = config.ConfigString(value)
		default:
			return nil, fmt.Errorf("unknown color key: %s", key)
		}
	}

	return colorConfig, nil
}
