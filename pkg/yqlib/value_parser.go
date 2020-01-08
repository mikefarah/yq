package yqlib

import (
	"strconv"

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
	var err interface{}
	var tag = customTag

	if tag == "" {
		_, err = strconv.ParseBool(argument)
		if err == nil {
			tag = "!!bool"
		}
		_, err = strconv.ParseFloat(argument, 64)
		if err == nil {
			tag = "!!float"
		}
		_, err = strconv.ParseInt(argument, 10, 64)
		if err == nil {
			tag = "!!int"
		}

		if argument == "null" {
			tag = "!!null"
		}
		if argument == "[]" {
			return &yaml.Node{Tag: "!!seq", Kind: yaml.SequenceNode}
		}
	}
	log.Debugf("parsed value '%v', tag: '%v'", argument, tag)
	return &yaml.Node{Value: argument, Tag: tag, Kind: yaml.ScalarNode}
}
