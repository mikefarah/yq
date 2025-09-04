package yqlib

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

type ColorConfig struct {
	BoolColor    color.Attribute
	NumberColor  color.Attribute
	MapKeyColor  color.Attribute
	AnchorColor  color.Attribute
	AliasColor   color.Attribute
	StringColor  color.Attribute
	CommentColor color.Attribute
}

func NewColorConfig() *ColorConfig {
	config := &ColorConfig{
		BoolColor:    color.FgHiMagenta,
		NumberColor:  color.FgHiMagenta,
		MapKeyColor:  color.FgCyan,
		AnchorColor:  color.FgHiYellow,
		AliasColor:   color.FgHiYellow,
		StringColor:  color.FgGreen,
		CommentColor: color.FgHiBlack,
	}

	if colorStr := os.Getenv("YQ_COLOR_BOOL"); colorStr != "" {
		if attr, err := parseColorAttribute(colorStr); err == nil {
			config.BoolColor = attr
		}
	}

	if colorStr := os.Getenv("YQ_COLOR_NUMBER"); colorStr != "" {
		if attr, err := parseColorAttribute(colorStr); err == nil {
			config.NumberColor = attr
		}
	}

	if colorStr := os.Getenv("YQ_COLOR_MAP_KEY"); colorStr != "" {
		if attr, err := parseColorAttribute(colorStr); err == nil {
			config.MapKeyColor = attr
		}
	}

	if colorStr := os.Getenv("YQ_COLOR_ANCHOR"); colorStr != "" {
		if attr, err := parseColorAttribute(colorStr); err == nil {
			config.AnchorColor = attr
		}
	}

	if colorStr := os.Getenv("YQ_COLOR_ALIAS"); colorStr != "" {
		if attr, err := parseColorAttribute(colorStr); err == nil {
			config.AliasColor = attr
		}
	}

	if colorStr := os.Getenv("YQ_COLOR_STRING"); colorStr != "" {
		if attr, err := parseColorAttribute(colorStr); err == nil {
			config.StringColor = attr
		}
	}

	if colorStr := os.Getenv("YQ_COLOR_COMMENT"); colorStr != "" {
		if attr, err := parseColorAttribute(colorStr); err == nil {
			config.CommentColor = attr
		}
	}

	return config
}

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
		if num, err := strconv.Atoi(colorStr); err == nil {
			return color.Attribute(num), nil
		}
		return color.Reset, fmt.Errorf("unknown color: %s", colorStr)
	}
}
