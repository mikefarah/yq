package yqlib

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

type ColorConfig struct {
	Bool    color.Attribute
	Number  color.Attribute
	MapKey  color.Attribute
	Anchor  color.Attribute
	Alias   color.Attribute
	String  color.Attribute
	Comment color.Attribute
}

func NewColorConfig() *ColorConfig {
	config := &ColorConfig{
		Bool:    color.FgHiMagenta,
		Number:  color.FgHiMagenta,
		MapKey:  color.FgCyan,
		Anchor:  color.FgHiYellow,
		Alias:   color.FgHiYellow,
		String:  color.FgGreen,
		Comment: color.FgHiBlack,
	}

	colorMappings := map[string]*color.Attribute{
		"YQ_COLOR_BOOL":    &config.Bool,
		"YQ_COLOR_NUMBER":  &config.Number,
		"YQ_COLOR_MAP_KEY": &config.MapKey,
		"YQ_COLOR_ANCHOR":  &config.Anchor,
		"YQ_COLOR_ALIAS":   &config.Alias,
		"YQ_COLOR_STRING":  &config.String,
		"YQ_COLOR_COMMENT": &config.Comment,
	}

	for envVar, configField := range colorMappings {
		if colorStr := os.Getenv(envVar); colorStr != "" {
			if attr, err := parseColorAttribute(colorStr); err == nil {
				*configField = attr
			}
		}
	}

	return config
}

// parseColorAttribute converts a color string to a color.Attribute.
//
// Supports three types of color specifications:
// 1. Standard color names: "red", "green", "blue", "yellow", "magenta", "cyan", "white", "black"
// 2. High-intensity variants: "hi-red", "hi-green", "hi-blue", etc.
func parseColorAttribute(colorStr string) (color.Attribute, error) {
	switch colorStr {
	case "black":
		return color.FgBlack, nil
	case "red":
		return color.FgRed, nil
	case "green":
		return color.FgGreen, nil
	case "yellow":
		return color.FgYellow, nil
	case "blue":
		return color.FgBlue, nil
	case "magenta":
		return color.FgMagenta, nil
	case "cyan":
		return color.FgCyan, nil
	case "white":
		return color.FgWhite, nil
	case "hi-black":
		return color.FgHiBlack, nil
	case "hi-red":
		return color.FgHiRed, nil
	case "hi-green":
		return color.FgHiGreen, nil
	case "hi-yellow":
		return color.FgHiYellow, nil
	case "hi-blue":
		return color.FgHiBlue, nil
	case "hi-magenta":
		return color.FgHiMagenta, nil
	case "hi-cyan":
		return color.FgHiCyan, nil
	case "hi-white":
		return color.FgHiWhite, nil
	default:
		return color.Reset, fmt.Errorf("unknown color: %s", colorStr)
	}
}
