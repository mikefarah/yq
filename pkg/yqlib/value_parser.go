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
		intValue, intErr := strconv.ParseInt(argument, 10, 64)
		floatValue, floatErr := strconv.ParseFloat(argument, 64)
		if intErr == nil && floatErr == nil {
			if int64(floatValue) == intValue {
				return intValue
			}
			return floatValue
		} else if floatErr == nil {
			// In case cannot parse the int due to large precision
			return floatValue
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
