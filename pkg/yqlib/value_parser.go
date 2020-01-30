package yqlib

import (
	yaml "gopkg.in/yaml.v3"
)

type ValueParser interface {
	Parse(argument string, customTag string) *yaml.Node
}

type valueParser struct {
}

func NewValueParser() ValueParser {
	return &valueParser{}
}

func (v *valueParser) Parse(argument string, customTag string) *yaml.Node {
	if argument == "[]" {
		return &yaml.Node{Tag: "!!seq", Kind: yaml.SequenceNode}
	}
	return &yaml.Node{Value: argument, Tag: customTag, Kind: yaml.ScalarNode}
}
