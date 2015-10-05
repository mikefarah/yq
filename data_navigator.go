package main

import (
	// "fmt"
	"log"
	"strconv"
)

func write(context map[interface{}]interface{}, head string, tail []string, value interface{}) {
	if len(tail) == 0 {
		context[head] = value
	} else {
		// e.g. if updating a.b.c, we need to get the 'b' map...
		toUpdate := readMap(context, head, tail[0:len(tail)-1]).(map[interface{}]interface{})
		//  and then set the 'c' key.
		key := (tail[len(tail)-1])
		toUpdate[key] = value
	}
}

func readMap(context map[interface{}]interface{}, head string, tail []string) interface{} {
	value := context[head]
	if len(tail) > 0 {
		return recurse(value, tail[0], tail[1:len(tail)])
	}
	return value
}

func recurse(value interface{}, head string, tail []string) interface{} {
	switch value.(type) {
	case []interface{}:
		if head == "*" {
			return readArraySplat(value.([]interface{}), tail)
		}
		index, err := strconv.ParseInt(head, 10, 64)
		if err != nil {
			log.Fatalf("Error accessing array: %v", err)
		}
		return readArray(value.([]interface{}), index, tail)
	case map[interface{}]interface{}:
		return readMap(value.(map[interface{}]interface{}), head, tail)
	default:
		return nil
	}
}

func readArray(array []interface{}, head int64, tail []string) interface{} {
	if head > int64(len(array)) {
		return nil
	}

	value := array[head]
	if len(tail) > 0 {
		return recurse(value, tail[0], tail[1:len(tail)])
	}
	return value
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
