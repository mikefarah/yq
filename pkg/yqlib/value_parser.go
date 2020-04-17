package yqlib

import (
	yaml "gopkg.in/yaml.v3"
)

type ValueParser interface {
	Parse(argument string, customTag string, customStyle string) *yaml.Node
}

type valueParser struct {
}

func NewValueParser() ValueParser {
	return &valueParser{}
}

func (v *valueParser) Parse(argument string, customTag string, customStyle string) *yaml.Node {
	var style yaml.Style
	if customStyle == "tagged" {
		style = yaml.TaggedStyle
	} else if customStyle == "doubleQuoted" {
		style = yaml.DoubleQuotedStyle
	} else if customStyle == "singleQuoted" {
		style = yaml.SingleQuotedStyle
	} else if customStyle == "literal" {
		style = yaml.LiteralStyle
	} else if customStyle == "folded" {
		style = yaml.FoldedStyle
	}

	if argument == "[]" {
		return &yaml.Node{Tag: "!!seq", Kind: yaml.SequenceNode, Style: style}
	}
	return &yaml.Node{Value: argument, Tag: customTag, Kind: yaml.ScalarNode, Style: style}
}
