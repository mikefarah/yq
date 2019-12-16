package yqlib

import (
	"strconv"

	logging "gopkg.in/op/go-logging.v1"
	yaml "gopkg.in/yaml.v3"
)

type ValueParser interface {
	Parse(argument string, customTag string) *yaml.Node
}

type valueParser struct {
	log *logging.Logger
}

func NewValueParser(l *logging.Logger) ValueParser {
	return &valueParser{log: l}
}

func (v *valueParser) Parse(argument string, customTag string) *yaml.Node {
	var err interface{}
	var tag = customTag

	var inQuotes = len(argument) > 0 && argument[0] == '"'
	if tag == "" && !inQuotes {

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
	v.log.Debugf("parsed value '%v', tag: '%v'", argument, tag)
	return &yaml.Node{Value: argument, Tag: tag, Kind: yaml.ScalarNode}
}
