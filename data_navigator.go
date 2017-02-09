package main

import (
	// "fmt"
	"strconv"
)

func write(context map[interface{}]interface{}, head string, tail []string, value interface{}) {
	if len(tail) == 0 {
		context[head] = value
	} else {
		// e.g. if updating a.b.c, we need to get the 'b', this could be a map or an array
		var parent = readMap(context, head, tail[0:len(tail)-1])
		switch parent.(type) {
		case map[interface{}]interface{}:
			toUpdate := parent.(map[interface{}]interface{})
			// b is a map, update the key 'c' to the supplied value
			key := (tail[len(tail)-1])
			toUpdate[key] = value
		case []interface{}:
			toUpdate := parent.([]interface{})
			// b is an array, update it at index 'c' to the supplied value
			rawIndex := (tail[len(tail)-1])
			index, err := strconv.ParseInt(rawIndex, 10, 64)
			if err != nil {
				die("Error accessing array: %v", err)
			}
			toUpdate[index] = value
		}

	}
}

func readMap(context map[interface{}]interface{}, head string, tail []string) interface{} {
	if head == "*" {
		return readMapSplat(context, tail)
	}
	value := context[head]
	return calculateValue(value, tail)
}

func readMapSplat(context map[interface{}]interface{}, tail []string) interface{} {
	var newArray = make([]interface{}, len(context))
	var i = 0
	for _, value := range context {
		if len(tail) > 0 {
			newArray[i] = recurse(value, tail[0], tail[1:len(tail)])
		} else {
			newArray[i] = value
		}
		i++
	}
	return newArray
}

func recurse(value interface{}, head string, tail []string) interface{} {
	switch value.(type) {
	case []interface{}:
		if head == "*" {
			return readArraySplat(value.([]interface{}), tail)
		}
		index, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			die("Error accessing array: %v", err)
		}
		return readArray(value.([]interface{}), index, tail)
	case map[interface{}]interface{}:
		return readMap(value.(map[interface{}]interface{}), head, tail)
	default:
		return nil
	}
}

func readArray(array []interface{}, head int64, tail []string) interface{} {
	if head >= int64(len(array)) {
		return nil
	}

	value := array[head]

	return calculateValue(value, tail)
}

func readArraySplat(array []interface{}, tail []string) interface{} {
	var newArray = make([]interface{}, len(array))
	for index, value := range array {
		newArray[index] = calculateValue(value, tail)
	}
	return newArray
}

func calculateValue(value interface{}, tail []string) interface{} {
	if len(tail) > 0 {
		return recurse(value, tail[0], tail[1:len(tail)])
	}
	return value
}
