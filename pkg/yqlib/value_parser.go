package yqlib

import (
	"strconv"
)

type ValueParser interface {
	ParseValue(argument string) interface{}
}

type valueParser struct{}

func NewValueParser() ValueParser {
	return &valueParser{}
}

func (v *valueParser) ParseValue(argument string) interface{} {
	var value, err interface{}
	var inQuotes = len(argument) > 0 && argument[0] == '"'
	if !inQuotes {
		value, err = strconv.ParseFloat(argument, 64)
		if err == nil {
			return value
		}
		value, err = strconv.ParseBool(argument)
		if err == nil {
			return value
		}
		if argument == "[]" {
			return make([]interface{}, 0)
		}
		return argument
	}
	return argument[1 : len(argument)-1]
}
