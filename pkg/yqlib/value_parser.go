package yqlib

import (
	yaml "gopkg.in/yaml.v3"
)

type ValueParser interface {
	Parse(argument string, customTag string, customStyle string, anchorName string, createAlias bool) *yaml.Node
}

type valueParser struct {
}

func NewValueParser() ValueParser {
	return &valueParser{}
}

func (v *valueParser) Parse(argument string, customTag string, customStyle string, anchorName string, createAlias bool) *yaml.Node {
	var style yaml.Style
	if customStyle == "tagged" {
		style = yaml.TaggedStyle
	} else if customStyle == "double" {
		style = yaml.DoubleQuotedStyle
	} else if customStyle == "single" {
		style = yaml.SingleQuotedStyle
	} else if customStyle == "literal" {
		style = yaml.LiteralStyle
	} else if customStyle == "folded" {
		style = yaml.FoldedStyle
	} else if customStyle == "flow" {
		style = yaml.FlowStyle
	} else if customStyle != "" {
		log.Error("Unknown style %v, ignoring", customStyle)
	}
	if argument == "[]" {
		return &yaml.Node{Tag: "!!seq", Kind: yaml.SequenceNode, Style: style}
	}

	kind := yaml.ScalarNode

	if createAlias {
		kind = yaml.AliasNode
	}

	return &yaml.Node{Value: argument, Tag: customTag, Kind: kind, Style: style, Anchor: anchorName}
}
