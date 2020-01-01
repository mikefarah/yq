package yqlib

import "strconv"
import "strings"

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
	var isDot = strings.Contains(argument, ".")
	if !inQuotes {
		if isDot {
			value, err = strconv.ParseFloat(argument, 64)
			if err == nil {
				return value
			}
		} else {
			value, err = strconv.ParseInt(argument, 10, 64)
			if err == nil {
				return value
			}
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
