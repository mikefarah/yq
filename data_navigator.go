package main

import (
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

func entryInSlice(context yaml.MapSlice, key interface{}) *yaml.MapItem {
	for idx := range context {
		var entry = &context[idx]
		if entry.Key == key {
			return entry
		}
	}
	return nil
}

func writeMap(context interface{}, paths []string, value interface{}) yaml.MapSlice {
	log.Debugf("writeMap for %v for %v with value %v\n", paths, context, value)

	var mapSlice yaml.MapSlice
	switch context.(type) {
	case yaml.MapSlice:
		mapSlice = context.(yaml.MapSlice)
	default:
		mapSlice = make(yaml.MapSlice, 0)
	}

	if len(paths) == 0 {
		return mapSlice
	}

	child := entryInSlice(mapSlice, paths[0])
	if child == nil {
		newChild := yaml.MapItem{Key: paths[0]}
		mapSlice = append(mapSlice, newChild)
		child = entryInSlice(mapSlice, paths[0])
		log.Debugf("\tAppended child at %v for mapSlice %v\n", paths[0], mapSlice)
	}

	log.Debugf("\tchild.Value %v\n", child.Value)

	remainingPaths := paths[1:]
	child.Value = updatedChildValue(child.Value, remainingPaths, value)
	log.Debugf("\tReturning mapSlice %v\n", mapSlice)
	return mapSlice
}

func updatedChildValue(child interface{}, remainingPaths []string, value interface{}) interface{} {
	if len(remainingPaths) == 0 {
		return value
	}

	_, nextIndexErr := strconv.ParseInt(remainingPaths[0], 10, 64)
	if nextIndexErr != nil {
		// must be a map
		return writeMap(child, remainingPaths, value)
	}

	// must be an array
	return writeArray(child, remainingPaths, value)
}

func writeArray(context interface{}, paths []string, value interface{}) []interface{} {
	log.Debugf("writeArray for %v for %v with value %v\n", paths, context, value)
	var array []interface{}
	switch context.(type) {
	case []interface{}:
		array = context.([]interface{})
	default:
		array = make([]interface{}, 1)
	}

	if len(paths) == 0 {
		return array
	}

	log.Debugf("\tarray %v\n", array)

	rawIndex := paths[0]
	index, err := strconv.ParseInt(rawIndex, 10, 64)
	if err != nil {
		die("Error accessing array: %v", err)
	}
	for index >= int64(len(array)) {
		array = append(array, nil)
	}
	currentChild := array[index]

	log.Debugf("\tcurrentChild %v\n", currentChild)

	remainingPaths := paths[1:]
	array[index] = updatedChildValue(currentChild, remainingPaths, value)
	log.Debugf("\tReturning array %v\n", array)
	return array
}

func readMap(context yaml.MapSlice, head string, tail []string) interface{} {
	if head == "*" {
		return readMapSplat(context, tail)
	}
	var value interface{}

	entry := entryInSlice(context, head)
	if entry != nil {
		value = entry.Value
	}
	return calculateValue(value, tail)
}

func readMapSplat(context yaml.MapSlice, tail []string) interface{} {
	var newArray = make([]interface{}, len(context))
	var i = 0
	for _, entry := range context {
		if len(tail) > 0 {
			newArray[i] = recurse(entry.Value, tail[0], tail[1:])
		} else {
			newArray[i] = entry.Value
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
	case yaml.MapSlice:
		return readMap(value.(yaml.MapSlice), head, tail)
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
		return recurse(value, tail[0], tail[1:])
	}
	return value
}
